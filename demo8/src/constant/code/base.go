package code

import (
	"github.com/dtm-labs/client/dtmcli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	Success  int = 200
	ParamErr     = 400
	NotFound     = 404
	Fail         = 500

	// users

	UserNotFound = 1000

	// course

	CourseNotFound      = 2001
	CourseFull          = 2002
	CourseSelected      = 2003
	CourseTimeConflict  = 2004
	CourseNotSelected   = 2005
	CourseCalOffsetFail = 2006
	CourseLuaRunErr     = 2007

	// mysql

	DBError = 3000

	// 熔断机制触发

	CircuitBreakerTrigger = 4000
	LimitTrigger          = 5000
)

var AbortedError = status.New(codes.Aborted, dtmcli.ResultFailure).Err()

const (
	// base

	SuccessMsg  string = "success"
	ParamErrMsg        = "param error"
	NotFoundMsg        = "not found"
	FailMsg            = "fail"

	// users

	UserNotFoundMsg = "用户不存在"

	// course

	CourseNotFoundMsg      = "课程不存在"
	CourseSelectedMsg      = "用户已经选择该门课程"
	CourseTimeConflictMsg  = "课程上课时间存在冲突"
	CourseFullMsg          = "课程已满"
	CourseNotSelectedMsg   = "用户未选择该门课程"
	CourseCalOffsetFailMsg = "计算课程上课时间失败"
	CourseLuaRunErrMsg     = "lua脚本执行失败"

	// db

	DBErrorMsg = "查询失败"

	CircuitBreakerTriggerMsg = "系统繁忙，请稍后再试"
	LimitTriggerMsg          = "请求太频繁，请稍后再试"
)
