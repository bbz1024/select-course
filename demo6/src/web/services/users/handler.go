package users

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"select-course/demo6/src/constant/code"
	"select-course/demo6/src/constant/services"
	"select-course/demo6/src/models/request"
	"select-course/demo6/src/rpc/user"
	"select-course/demo6/src/utils/bloom"
	grpc2 "select-course/demo6/src/utils/grpc"
	logger2 "select-course/demo6/src/utils/logger"
	"select-course/demo6/src/utils/resp"
	"select-course/demo6/src/utils/tracing"
)

var (
	userClient user.UserServiceClient
	logger     *zap.Logger
)

func New() {
	logger = logger2.LogService(services.UserRpcServerName)
	conn := grpc2.Connect(context.Background(), services.UserRpcServerName)
	userClient = user.NewUserServiceClient(conn)
}

func GetUserHandler(ctx *gin.Context) {
	// tracing
	span := tracing.StartSpan(ctx, "GetUserHandler")
	defer span.Finish()

	// valid param
	var req request.UserReq
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Info("GetUserHandler", zap.Error(err))
		resp.ParamErr(ctx)
		return
	}
	logField := []zap.Field{zap.Int64("user_id", req.UserID)}
	// bloom filter
	if !bloom.UserBloom.TestString(fmt.Sprintf("%d", req.UserID)) {
		logger.Info("user not found", logField...)
		resp.Fail(ctx, code.UserNotFound, code.UserNotFoundMsg)
		return
	}
	// call rpc
	var res *user.UserResponse
	var err error
	if res, err = userClient.GetUserInfo(ctx, &user.UserRequest{
		UserId: req.UserID,
	}); err != nil {
		logField = append(logField, zap.Error(err))
		logger.Warn("call userClient.GetUserInfo failed", logField...)
		tracing.RecordError(span, err)
		ctx.Render(http.StatusOK, resp.CustomJSON{
			Context: ctx,
			Data:    res,
		})
		return
	}
	ctx.Render(http.StatusOK, resp.CustomJSON{
		Context: ctx,
		Data:    res,
	})
}
