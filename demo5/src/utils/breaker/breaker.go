package breaker

import (
	"errors"
	"fmt"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
)

var (
	BreakError = errors.New("触发熔断")
)

// ErrorCountRules 熔断规则: 错误次数
var ErrorCountRules = []*circuitbreaker.Rule{
	{
		//executeLuaScript
		Resource:                     "executeLuaScript",
		Strategy:                     circuitbreaker.ErrorCount,
		RetryTimeoutMs:               3000,
		MinRequestAmount:             10,
		StatIntervalMs:               5000,
		StatSlidingWindowBucketCount: 10,
		Threshold:                    5, // 5次错误 10次请求其中有5次错误，触发熔断
	},
}

// StateChangeTestListener 熔断器
type StateChangeTestListener struct {
}

func (s *StateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {

	fmt.Println("熔断关闭")
}

func (s *StateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {

	fmt.Println("opening")
}

func (s *StateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Println("熔断恢复")
}
