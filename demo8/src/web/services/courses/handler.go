package courses

import (
	"context"
	"errors"
	"fmt"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/client/dtmgrpc/dtmgimp"
	"github.com/dtm-labs/client/workflow"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"net/http"
	grpc2 "select-course/demo7/src/utils/grpc"
	"select-course/demo8/src/constant/code"
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
	// register workflow transaction
	rpcServer := grpc.NewServer(grpc.UnaryInterceptor(dtmgimp.GrpcServerLog))
	const GrpcAddr = "localhost:10001"
	const DTMAddr = "localhost:36790"
	workflow.InitGrpc(DTMAddr, GrpcAddr, rpcServer)
	var res *course.CourseOptResponse
	err := workflow.Register(services.CourseRpcServerName, func(wf *workflow.Workflow, data []byte) error {
		// 预扣减容量
		// tackle rollback
		var body course.CourseOptRequest
		if err := proto.Unmarshal(data, &body); err != nil {
			logger.Error("json unmarshal", zap.Error(err))
			tracing.RecordError(span, err)
			return err
		}
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			if _, err := courseClient.BackCourse(wf.Context, &body); err != nil {
				logger.Error("rollback", zap.Error(err))
				return err
			}
			return nil
		})
		var err error
		res, err = courseClient.SelectCourse(wf.Context, &body)
		if err != nil {
			logger.Error("SelectCourse", zap.Error(err))
			tracing.RecordError(span, err)
			return err
		}
		// 异步提交消息，是否提交成功（消息队列是否落库）
		res, err = courseClient.EnQueueCourse(wf.Context, &course.EnQueueCourseRequest{
			UserId: int64(req.UserID), CourseId: int64(req.CourseID),
			CreateAt: res.CreateAt, IsSelect: true,
		})
		if err != nil {
			logger.Error("EnQueueCourse", zap.Error(err))
			// 回滚事务
			return err
		}
		return nil
	})
	if err != nil {
		tracing.RecordError(span, err)
		logFiled = append(logFiled, zap.Error(err))
		logger.Error("SelectCourse", logFiled...)
		resp.Fail(ctx, code.Fail, code.FailMsg)
		return
	}
	data, err := proto.Marshal(&course.CourseOptRequest{
		UserId: int64(req.UserID), CourseId: int64(req.CourseID),
	})
	if err != nil {
		tracing.RecordError(span, err)
		logFiled = append(logFiled, zap.Error(err))
		logger.Error("SelectCourse", logFiled...)
		resp.Fail(ctx, code.Fail, code.FailMsg)
		return
	}
	_, err = workflow.ExecuteCtx(ctx.Request.Context(), services.CourseRpcServerName, shortuuid.New(), data)
	if err != nil {
		tracing.RecordError(span, err)
		logFiled = append(logFiled, zap.Error(err))
		// rollback
		if errors.Is(err, dtmcli.ErrFailure) {
			logger.Error("rollback", logFiled...)
			ctx.Render(http.StatusOK, resp.CustomJSON{
				Data:    res,
				Context: ctx,
			})
			return
		}
		logger.Error("SelectCourse", logFiled...)
		resp.Fail(ctx, code.Fail, code.FailMsg)
		return
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
