package users

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"select-course/demo3/src/constant/code"
	"select-course/demo3/src/models"
	"select-course/demo3/src/models/request"
	"select-course/demo3/src/storage/database"
	"select-course/demo3/src/utils/logger"
	"select-course/demo3/src/utils/resp"
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
	// First 查询不存在时会直接抛出异常，Find查询并不会，查询不到时返回的是空结果
	if err := database.Client.First(&user, req.UserID).Error; err != nil {
		// 不存在
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Fail(ctx, code.NotFound, code.UserNotFound, code.UserNotFoundMsg)
			return
		}
		logger.Logger.Error("查询用户失败", err)
		resp.DBError(ctx)
		return
	}
	//3. 返回用户信息
	resp.Success(ctx, user)
}
