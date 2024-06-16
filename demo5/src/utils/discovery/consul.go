package discovery

import (
	"context"
	"fmt"
	capi "github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"log"
	config2 "select-course/demo5/src/constant/config"
	"select-course/demo5/src/utils/logger"
	"strconv"
)

type ConsulDiscovery struct {
	prefix  string
	Address string
	client  *capi.Client
}

func NewConsulDiscovery(address string, prefix string) *ConsulDiscovery {
	cfg := capi.DefaultConfig()
	cfg.Address = address
	client, err := capi.NewClient(cfg)
	if err != nil {
		log.Printf("Connect Consul happens error: %v", err)
	}
	return &ConsulDiscovery{
		prefix:  prefix,
		Address: address,
		client:  client,
	}

}
func (c *ConsulDiscovery) Register(ctx context.Context, service Service) error {
	parsePort, err := strconv.Atoi(service.Port[1:])
	if err != nil {
		return err
	}
	reg := &capi.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", c.prefix, service.Name, parsePort),
		Name:    service.Name,
		Address: config2.EnvCfg.BaseHost,
		Port:    parsePort,
		//Http检查
		//Check: &capi.AgentServiceCheck{
		//	Interval:                       "5s",
		//	Timeout:                        "5s",
		//	DeregisterCriticalServiceAfter: "10s",
		//	HTTP:                           fmt.Sprintf("http://%s:%d/health", service.Address, port),
		//},
	}
	if err := c.client.Agent().ServiceRegister(reg); err != nil {
		return err
	}
	logger.LogService(service.Name).Debug("register service success",
		zap.String("address", fmt.Sprintf("%s:%s", config2.EnvCfg.BaseHost, service.Port)),
	)
	return nil
}

func (c *ConsulDiscovery) Deregister(ctx context.Context, name string) error {
	panic("implement me")
}

func (c *ConsulDiscovery) GetService(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("consul://%s/%s?wait=15s", c.Address, name), nil
}

var Consul = NewConsulDiscovery(
	fmt.Sprintf("%s:%s", config2.EnvCfg.ConsulHost, config2.EnvCfg.ConsulPort), config2.EnvCfg.ProjectName)
