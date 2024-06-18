package discovery

import (
	"context"
)

type Service struct {
	Name   string
	Port   string
	Weight int
}

type Discovery interface {
	// Register register service
	Register(ctx context.Context, service Service) error
	// Deregister deregister service
	Deregister(ctx context.Context, name string) error
	// GetService get service
	GetService(ctx context.Context, name string) (string, error)
}
