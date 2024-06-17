package limiter

import (
	"github.com/alibaba/sentinel-golang/core/hotspot"
)

var LimitRules = []*hotspot.Rule{
	{
		Resource:        "GetMyCourses",
		MetricType:      hotspot.QPS,    // 基于QPS
		ControlBehavior: hotspot.Reject, //直接拒绝，
		ParamIndex:      0,              // 0:第一个参数
		Threshold:       50,             // 每秒请求数不超过100，针对某个参数
		DurationInSec:   1,              // 1:统计时间间隔
	},
}
