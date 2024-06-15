package users

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"select-course/demo5/src/constant/code"
	"select-course/demo5/src/constant/services"
	"select-course/demo5/src/models/request"
	"select-course/demo5/src/rpc/user"
	grpc2 "select-course/demo5/src/utils/grpc"
	logger2 "select-course/demo5/src/utils/logger"
	"select-course/demo5/src/utils/resp"
	"select-course/demo5/src/utils/tracing"
)

var (
	userClient user.UserServiceClient
	logger     *zap.Logger
)

func New() {
	conn := grpc2.Connect(context.Background(), services.UserRpcServerName)
	userClient = user.NewUserServiceClient(conn)
	logger = logger2.LogService(services.UserRpcServerName)
}
func GetUserHandler(ctx *gin.Context) {
	span := tracing.StartSpan(ctx, "GetUserHandler")
	defer span.Finish()
	var req request.UserReq
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Info("GetUserHandler", zap.Error(err))
		resp.ParamErr(ctx)
		return
	}

	userInfo, err := userClient.GetUserInfo(ctx, &user.UserRequest{
		UserId: req.UserID,
	})
	if err != nil {
		logger.Warn("call userClient.GetUserInfo failed", zap.Error(err))
		tracing.RecordError(span, err)
		resp.Fail(ctx, code.Fail, code.Fail, code.FailMsg)
		return
	}
	resp.Success(ctx, userInfo)
}
