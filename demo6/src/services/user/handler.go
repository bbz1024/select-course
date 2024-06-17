package main

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"select-course/demo6/src/constant/code"
	"select-course/demo6/src/constant/services"
	"select-course/demo6/src/models"
	"select-course/demo6/src/rpc/user"
	"select-course/demo6/src/storage/database"
	"select-course/demo6/src/utils/logger"
	"select-course/demo6/src/utils/tracing"
)

type User struct {
	user.UnimplementedUserServiceServer
}

var (
	Logger *zap.Logger
)

func (u *User) New() {
	// init rpc client
	Logger = logger.LogService(services.UserRpcServerName)

}
func (u *User) GetUserInfo(ctx context.Context, req *user.UserRequest) (*user.UserResponse, error) {
	// tracing
	span := tracing.StartSpan(ctx, "GetUserInfo")
	defer span.Finish()
	// get userinfo
	logField := []zap.Field{zap.Int64("user_id", req.UserId)}
	var userInfo models.User
	if err := database.Client.First(&userInfo, req.UserId).Error; err != nil {
		tracing.RecordError(span, err)
		logField = append(logField, zap.Error(err))
		Logger.Error("GetUserInfo", logField...)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user.UserResponse{
				StatusCode: code.UserNotFound,
				StatusMsg:  code.UserNotFoundMsg,
			}, err
		}
		return &user.UserResponse{
			StatusCode: code.DBError,
			StatusMsg:  code.DBErrorMsg,
		}, err
	}
	// build res
	res := &user.UserResponse{
		StatusCode: int32(code.Success),
		StatusMsg:  code.SuccessMsg,
		UserId:     int64(userInfo.ID),
		UserName:   userInfo.UserName,
		Password:   userInfo.Password,
	}
	return res, nil
}
