package users

import (
	"github.com/gin-gonic/gin"
	"select-course/demo1/src/constant/code"
	"select-course/demo1/src/models"
	"select-course/demo1/src/models/request"
	"select-course/demo1/src/storage/database"
	"select-course/demo1/src/utils/logger"
	"select-course/demo1/src/utils/resp"
)

func GetUserHandler(ctx *gin.Context) {
	//1. 参数校验
	var req request.UserReq
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Logger.Info("参数校验失败", err)
		resp.ParamErr(ctx)
		return
	}
	//2. 判断用户是否存在
	var user models.User
	if err := database.Client.Where(
		"user_name = ? and password = ?", req.UserName, req.Password,
	).Find(&user).Error; err != nil {
		resp.DBError(ctx)
		return
	}
	// 不存在
	if user.ID == 0 {
		resp.Fail(ctx, code.NotFound, code.UserNotFound, code.UserNotFoundMsg)
		return
	}
	//3. 返回用户信息
	resp.Success(ctx, user)

}
