package courses

import (
	"context"
	"errors"
	"fmt"
	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v3"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	grpc2 "select-course/demo7/src/utils/grpc"
	"select-course/demo8/src/constant/code"
	"select-course/demo8/src/constant/config"
	"select-course/demo8/src/constant/services"
	"select-course/demo8/src/models/request"
	"select-course/demo8/src/rpc/course"
	"select-course/demo8/src/utils/bloom"
	logger2 "select-course/demo8/src/utils/logger"
	"select-course/demo8/src/utils/resp"
	"select-course/demo8/src/utils/tracing"
)

var (
	logger       *zap.Logger
	courseClient course.CourseServiceClient
)

func New() {

	conn := grpc2.Connect(context.Background(), services.CourseRpcServerName)
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
	// call grpc by transaction
	courseGrpcAddr := fmt.Sprintf("%s%s", config.EnvCfg.BaseHost, services.CourseRpcServerAddr)
	dtmAddr := fmt.Sprintf("%s:%d", config.EnvCfg.DtmHost, config.EnvCfg.DtmPort)
	logger.Info("transaction start")
	var res = &course.CourseOptResponse{}
	err := dtmgrpc.TccGlobalTransaction(dtmAddr, shortuuid.New(), func(tcc *dtmgrpc.TccGrpc) error {
		r := &emptypb.Empty{}
		if err := tcc.CallBranch(
			&course.CourseOptRequest{
				UserId:   int64(req.UserID),
				CourseId: int64(req.CourseID),
			},
			courseGrpcAddr+"/course.CourseService/TryTryDeductCourse",
			courseGrpcAddr+"/course.CourseService/TryConfirmDeductCourse",
			courseGrpcAddr+"/course.CourseService/TryCancelDeductCourse",
			r,
		); err != nil {
			return err
		}
		fmt.Println(r)

		if err := tcc.CallBranch(&course.EnQueueCourseRequest{
			CreateAt: res.CreateAt,
			UserId:   int64(req.UserID),
			CourseId: int64(req.CourseID),
			IsSelect: true,
		},
			courseGrpcAddr+"/course.CourseService/TryTryEnqueueMessage",
			courseGrpcAddr+"/course.CourseService/TryConfirmEnqueueMessage",
			courseGrpcAddr+"/course.CourseService/TryCancelEnqueueMessage",
			res,
		); err != nil {
			return err
		}
		return nil
	})

	fmt.Println(err)
	if err != nil && !errors.Is(err, code.AbortedError) {
		logFiled = append(logFiled, zap.Error(err))
		logger.Error("tcc error", logFiled...)
		tracing.RecordError(span, err)
	}
	logger.Info("transaction end")
	ctx.Render(http.StatusOK, resp.CustomJSON{
		Data:    nil,
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
	if err != nil && !errors.Is(err, code.AbortedError) {
		logFiled = append(logFiled, zap.Error(err))

		logger.Error("BackCourse", logFiled...)
		tracing.RecordError(span, err)
	}
	ctx.Render(http.StatusOK, resp.CustomJSON{
		Data:    res,
		Context: ctx,
	})
}
