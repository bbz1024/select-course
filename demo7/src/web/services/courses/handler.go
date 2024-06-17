package courses

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"select-course/demo7/src/constant/code"
	"select-course/demo7/src/constant/services"
	"select-course/demo7/src/models/request"
	"select-course/demo7/src/rpc/course"
	"select-course/demo7/src/utils/bloom"
	"select-course/demo7/src/utils/grpc"
	logger2 "select-course/demo7/src/utils/logger"
	"select-course/demo7/src/utils/resp"
	"select-course/demo7/src/utils/tracing"
)

var (
	logger       *zap.Logger
	courseClient course.CourseServiceClient
)

func New() {
	conn := grpc.Connect(context.Background(), services.CourseRpcServerName)
	courseClient = course.NewCourseServiceClient(conn)
	logger = logger2.LogService(services.CourseRpcServerName)
}
func GetCourseList(ctx *gin.Context) {
	// tracing
	span := tracing.StartSpan(ctx, "GetCourseList")
	defer span.Finish()

	// call rpc
	courses, err := courseClient.GetAllCourses(ctx, nil)
	if err != nil {
		logger.Error("GetCourseList", zap.Error(err))
		tracing.RecordError(span, err)
		ctx.Render(http.StatusOK, resp.CustomJSON{
			Data:    courses,
			Context: ctx,
		})
		return
	}
	ctx.Render(http.StatusOK, resp.CustomJSON{
		Data:    courses,
		Context: ctx,
	})
}
func MyCourseList(ctx *gin.Context) {
	// tracing
	span := tracing.StartSpan(ctx, "MyCourseList")
	defer span.Finish()

	// valid
	// valid param
	var req request.UserReq
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Info("GetUserHandler", zap.Error(err))
		resp.ParamErr(ctx)
		return
	}
	logFiled := []zap.Field{zap.Int64("user_id", req.UserID)}
	// bloom filter
	if !bloom.UserBloom.TestString(fmt.Sprintf("%d", req.UserID)) {
		logger.Info("user not found", logFiled...)
		resp.Fail(ctx, code.UserNotFound, code.UserNotFoundMsg)
		return
	}
	// call rpc
	courses, err := courseClient.GetMyCourses(ctx, &course.GetMyCoursesRequest{UserId: req.UserID})
	if err != nil {
		logFiled = append(logFiled, zap.Error(err))
		logger.Error("GetMyCourses", logFiled...)
		tracing.RecordError(span, err)
	}
	ctx.Render(
		http.StatusOK, resp.CustomJSON{
			Data:    courses,
			Context: ctx,
		},
	)
}

func SelectCourse(ctx *gin.Context) {
	// tracing
	span := tracing.StartSpan(ctx, "SelectCourse")
	defer span.Finish()

	// valid param
	var req request.SelectCourseReq
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Info("SelectCourse", zap.Error(err))
		resp.ParamErr(ctx)
		return
	}
	logFiled := []zap.Field{zap.Uint("course_id", req.CourseID), zap.Uint("user_id", req.UserID)}
	if !bloom.UserBloom.TestString(fmt.Sprintf("%d", req.UserID)) {
		logger.Info("user not found", logFiled...)
		resp.Fail(ctx, code.UserNotFound, code.UserNotFoundMsg)
		return
	}
	if !bloom.CourseBloom.TestString(fmt.Sprintf("%d", req.CourseID)) {
		logger.Info("course not found", logFiled...)
		resp.Fail(ctx, code.CourseNotFound, code.CourseNotFoundMsg)
		return
	}

	// call rpc
	res, err := courseClient.SelectCourse(
		ctx, &course.CourseOptRequest{
			UserId: int64(req.UserID), CourseId: int64(req.CourseID),
		},
	)
	if err != nil {
		logFiled = append(logFiled, zap.Error(err))
		logger.Error("SelectCourse", logFiled...)
		tracing.RecordError(span, err)
	}
	ctx.Render(http.StatusOK, resp.CustomJSON{
		Data:    res,
		Context: ctx,
	})
}

func BackCourse(ctx *gin.Context) {
	// tracing
	span := tracing.StartSpan(ctx, "BackCourse")
	defer span.Finish()

	// valid param
	var req request.SelectCourseReq
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Info("BackCourse", zap.Error(err))
		resp.ParamErr(ctx)
		return
	}
	logFiled := []zap.Field{zap.Uint("course_id", req.CourseID), zap.Uint("user_id", req.UserID)}
	if !bloom.UserBloom.TestString(fmt.Sprintf("%d", req.UserID)) {
		logger.Info("user not found", logFiled...)
		resp.Fail(ctx, code.UserNotFound, code.UserNotFoundMsg)
		return
	}
	if !bloom.CourseBloom.TestString(fmt.Sprintf("%d", req.CourseID)) {
		logger.Info("course not found", logFiled...)
		resp.Fail(ctx, code.CourseNotFound, code.CourseNotFoundMsg)
		return
	}
	// call rpc
	res, err := courseClient.BackCourse(
		ctx, &course.CourseOptRequest{
			UserId: int64(req.UserID), CourseId: int64(req.CourseID),
		},
	)
	if err != nil {
		logFiled = append(logFiled, zap.Error(err))
		logger.Error("BackCourse", logFiled...)
		tracing.RecordError(span, err)
	}
	ctx.Render(http.StatusOK, resp.CustomJSON{
		Data:    res,
		Context: ctx,
	})
}
