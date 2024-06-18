package main

import (
	"context"
	"errors"
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"select-course/demo8/src/constant/code"
	"select-course/demo8/src/constant/keys"
	"select-course/demo8/src/constant/lua"
	"select-course/demo8/src/constant/services"
	"select-course/demo8/src/models"
	"select-course/demo8/src/models/mqm"
	"select-course/demo8/src/rpc/course"
	"select-course/demo8/src/storage/cache"
	"select-course/demo8/src/storage/database"
	"select-course/demo8/src/utils/breaker"
	"select-course/demo8/src/utils/consumer"
	"select-course/demo8/src/utils/local"
	"select-course/demo8/src/utils/logger"
	"select-course/demo8/src/utils/tracing"
	"strconv"
)

type Course struct {
	course.UnimplementedCourseServiceServer
}

var (
	Logger *zap.Logger
)

const TestRandoBroke = true

var AbortedError = status.New(codes.Aborted, dtmcli.ResultFailure).Err()

func (c Course) New() {
	Logger = logger.LogService(services.CourseRpcServerName)
}
func (c Course) GetAllCourses(ctx context.Context, request *course.GetAllCoursesRequest) (*course.GetAllCoursesResponse, error) {
	// tracing
	span := tracing.StartSpan(ctx, "GetAllCourses")
	defer span.Finish()

	// find all courses
	var courseList []*models.Course
	if err := database.Client.Find(&courseList).Error; err != nil {
		tracing.RecordError(span, err)
		Logger.Error("GetAllCourses opt database error", zap.Error(err))
		return nil, err
	}

	// build res
	var courseListRes []*course.Course
	for _, c := range courseList {
		courseListRes = append(courseListRes, &course.Course{
			Id:       int64(c.ID),
			Name:     c.Title,
			Capacity: int64(c.Capacity),
		})
	}
	return &course.GetAllCoursesResponse{
		Courses: courseListRes,
	}, nil
}

func (c Course) GetMyCourses(ctx context.Context, request *course.GetMyCoursesRequest) (*course.GetMyCoursesResponse, error) {
	// tracing
	span := tracing.StartSpan(ctx, "GetMyCourses")
	defer span.Finish()
	// limiter
	e, b := sentinel.Entry("GetMyCourses", sentinel.WithArgs(request.UserId))
	if b != nil {
		fmt.Println("GetMyCourses: LimitTrigger")
		// 对该用户进行限流
		return &course.GetMyCoursesResponse{
			StatusMsg:  code.LimitTriggerMsg,
			StatusCode: code.LimitTrigger,
		}, nil
	}
	defer e.Exit()
	// find user courses
	var userCourseList []*models.UserCourse
	if err := database.Client.
		Model(models.UserCourse{}).
		Preload("Course").
		Where("user_id = ?", request.UserId).
		Find(&userCourseList).Error; err != nil {
		tracing.RecordError(span, err)
		Logger.Error("GetMyCourses opt database error", zap.Error(err))
		return nil, err
	}

	// build res
	var courseListRes []*course.Course
	//for _, userCourse := range userCourseList {
	//	courseListRes = append(courseListRes, &course.Course{
	//		Id:       int64(userCourse.Course.ID),
	//		Name:     userCourse.Course.Title,
	//		Capacity: int64(userCourse.Course.Capacity),
	//	})
	//}
	return &course.GetMyCoursesResponse{
		Courses: courseListRes,
	}, nil
}

func (c Course) SelectCourse(ctx context.Context, request *course.CourseOptRequest) (*course.CourseOptResponse, error) {
	// tracing
	span := tracing.StartSpan(ctx, "SelectCourse")
	defer span.Finish()
	// get offset
	logFiled := []zap.Field{zap.Int64("course_id", request.CourseId), zap.Int64("user_id", request.UserId)}
	offset, err := local.CalOffset(uint(request.CourseId))

	if err != nil {
		logFiled = append(logFiled, zap.Error(err))
		tracing.RecordError(span, err)
		Logger.Error("SelectCourse cal offset error", logFiled...)
		return &course.CourseOptResponse{
			StatusCode: code.Fail,
			StatusMsg:  code.FailMsg,
		}, AbortedError
	}
	val, err := executeLuaScript(ctx, cache.RDB, redis.NewScript(lua.CourseSelectLuaScript), []string{
		fmt.Sprintf(keys.UserCourseSetKey, request.UserId),
		strconv.Itoa(int(request.CourseId)),
		fmt.Sprintf(keys.CourseHsetKey, request.CourseId),
		keys.CourseCapacityKey,
		fmt.Sprintf(keys.UserCourseScheduleBitMapKey, request.UserId),
		strconv.Itoa(offset),
		keys.CourseSequenceKey,
	}, request.UserId, request.CourseId)
	if err != nil {
		tracing.RecordError(span, err)
		logFiled = append(logFiled, zap.Error(err))
		if errors.Is(breaker.BreakError, err) {
			Logger.Warn("circuit breaker", logFiled...)
			return &course.CourseOptResponse{
				StatusCode: code.CircuitBreakerTrigger,
				StatusMsg:  code.CircuitBreakerTriggerMsg,
			}, nil
		}
		Logger.Error("SelectCourse execute lua script error", logFiled...)
		return &course.CourseOptResponse{
			StatusCode: code.Fail,
			StatusMsg:  code.FailMsg,
		}, AbortedError
	}
	switch {
	case val >= lua.CourseOptOK:
		Logger.Info("用户选择课程", logFiled...)
		return &course.CourseOptResponse{
			CreateAt: val,
		}, nil
	case val == lua.CourseSelected:
		Logger.Info("用户已经选择该门课程", logFiled...)
		return &course.CourseOptResponse{
			StatusCode: code.CourseSelected,
			StatusMsg:  code.CourseSelectedMsg,
		}, AbortedError
	case val == lua.CourseFull:
		Logger.Info("课程已满", logFiled...)
		return &course.CourseOptResponse{
			StatusCode: code.CourseFull,
			StatusMsg:  code.CourseFullMsg,
		}, AbortedError
	case val == lua.CourseTimeConflict:
		Logger.Info("课程时间冲突", logFiled...)
		return &course.CourseOptResponse{
			StatusCode: code.CourseTimeConflict,
			StatusMsg:  code.CourseTimeConflictMsg,
		}, AbortedError
	}

	return &course.CourseOptResponse{
		CreateAt: val,
	}, AbortedError
}

func (c Course) BackCourse(ctx context.Context, request *course.CourseOptRequest) (*course.CourseOptResponse, error) {
	// tracing
	span := tracing.StartSpan(ctx, "BackCourse")
	defer span.Finish()

	// cal offset
	offset, err := local.CalOffset(uint(request.CourseId))
	if err != nil {
		tracing.RecordError(span, err)
		Logger.Error("BackCourse cal offset error", zap.Error(err))
		return &course.CourseOptResponse{
			StatusCode: code.Fail,
			StatusMsg:  code.FailMsg,
		}, status.New(codes.Aborted, dtmcli.ResultFailure).Err()
	}

	logFiled := []zap.Field{zap.Int64("course_id", request.CourseId), zap.Int64("user_id", request.UserId)}
	val, err := executeLuaScript(ctx, cache.RDB, redis.NewScript(lua.CourseBackLuaScript), []string{
		fmt.Sprintf(keys.UserCourseSetKey, request.UserId),
		strconv.Itoa(int(request.CourseId)),
		fmt.Sprintf(keys.CourseHsetKey, request.CourseId),
		keys.CourseCapacityKey,
		fmt.Sprintf(keys.UserCourseScheduleBitMapKey, request.UserId),
		strconv.Itoa(offset),
		keys.CourseSequenceKey,
	}, request.UserId, request.CourseId)
	if err != nil {
		tracing.RecordError(span, err)
		logFiled = append(logFiled, zap.Error(err))
		if errors.Is(breaker.BreakError, err) {
			Logger.Warn("circuit breaker", logFiled...)
			return &course.CourseOptResponse{
				StatusCode: code.CircuitBreakerTrigger,
				StatusMsg:  code.CircuitBreakerTriggerMsg,
			}, nil
		}
		Logger.Error("BackCourse execute lua script error", logFiled...)
		return &course.CourseOptResponse{
			StatusCode: code.Fail,
			StatusMsg:  code.FailMsg,
		}, AbortedError
	}
	switch {
	case val >= lua.CourseOptOK:
		Logger.Info("用户退选课程", logFiled...)
		return &course.CourseOptResponse{
			CreateAt: val,
		}, nil
	case val == lua.CourseNotSelected:
		Logger.Info("用户未选择该门课程", logFiled...)
		return &course.CourseOptResponse{
			StatusCode: code.CourseNotSelected,
			StatusMsg:  code.CourseNotSelectedMsg,
		}, AbortedError
	}
	return &course.CourseOptResponse{
		CreateAt: val,
	}, AbortedError
}

// EnQueueCourse 异步消息入队
func (c Course) EnQueueCourse(ctx context.Context, request *course.EnQueueCourseRequest) (*course.CourseOptResponse, error) {
	span := tracing.StartSpan(ctx, "EnQueueCourse")
	defer span.Finish()

	opt := mqm.SelectType
	if !request.IsSelect {
		opt = mqm.BackType
	}
	if err := consumer.SelectConsumer.Product(&mqm.CourseReq{
		UserID: uint(request.UserId), CourseID: uint(request.CourseId), Type: opt,
		CreatedAt: request.CreateAt,
	}); err != nil {
		return &course.CourseOptResponse{
			StatusCode: code.Fail,
			StatusMsg:  code.FailMsg,
		}, AbortedError
	}
	return nil, nil
}

func executeLuaScript(ctx context.Context, rdb *redis.Client, script *redis.Script, keys []string, args ...interface{}) (int64, error) {
	// sentinel
	entry, errBroke := sentinel.Entry("executeLuaScript")
	// 如果errBroke非nil，则触发熔断，进行快速失败（降级，静态返回）
	if errBroke != nil {
		Logger.Warn("open circuit breaker", zap.Error(errBroke))
		return 0, breaker.BreakError
	}
	defer entry.Exit()
	val, err := script.Run(ctx, rdb, keys, args...).Result()
	if err != nil {
		sentinel.TraceError(entry, err)
		Logger.Warn("执行lua脚本失败", zap.Error(err))
		return 0, err
	}
	return val.(int64), nil
}
