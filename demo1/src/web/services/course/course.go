package course

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"select-course/demo1/src/constant/code"
	"select-course/demo1/src/models"
	"select-course/demo1/src/models/request"
	"select-course/demo1/src/storage/database"
	"select-course/demo1/src/utils/logger"
	"select-course/demo1/src/utils/resp"
)

func GetCourseList(ctx *gin.Context) {
	var courseList []*models.Course
	if err := database.Client.Find(&courseList).Error; err != nil {
		resp.DBError(ctx)
		return
	}
	resp.Success(ctx, courseList)
}

func SelectCourse(ctx *gin.Context) {
	// 1. 参数校验
	var req request.SelectCourseReq
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Logger.WithContext(ctx).Info("参数校验失败", err)
		resp.ParamErr(ctx)
		return
	}
	var user models.User
	if err := database.Client.Where("id=?", req.UserID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.WithContext(ctx).Info("用户不存在", err)
			resp.Fail(ctx, code.NotFound, code.UserNotFound, code.UserNotFoundMsg)
			return
		}
	}
	// 2. 数据库操作
	// 使用数据库事务确保操作的原子性
	if err := database.Client.WithContext(context.Background()).Transaction(func(tx *gorm.DB) error {
		// 2.1 检查课程是否存在和库存
		var course models.Course
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(models.Course{}).Preload("Schedule").
			Where("id=? and capacity > 0", req.CourseID).First(&course).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Logger.WithContext(ctx).Info("课程不存在或库存不足", err)
				resp.Fail(ctx, code.NotFound, code.CourseNotFound, code.CourseNotFoundMsg)
				// return nil 无需回滚，业务逻辑错误
			}
			return err
		}
		offset := int(course.Schedule.Week)*3 + int(course.Schedule.Duration)
		// 2.2 检查用户是否已选该课程或时间冲突
		if user.Flag.TestBit(offset) {
			logger.Logger.WithContext(ctx).Info("用户已选该课程或时间冲突")
			resp.Fail(ctx, code.ParamErr, code.CourseSelected, code.CourseSelectedMsg)
			return errors.New("用户已选该课程或时间冲突")
		}
		// 2.3 扣减课程库存并创建用户选课记录
		if err := tx.Model(&course).Update("capacity", gorm.Expr("capacity - 1")).Error; err != nil {
			logger.Logger.WithContext(ctx).Info("更新课程容量失败", err)
			resp.DBError(ctx)
			return err
		}
		if err := tx.Create(&models.UserCourse{
			UserID:   req.UserID,
			CourseID: req.CourseID,
		}).Error; err != nil {
			logger.Logger.WithContext(ctx).Info("创建选课记录失败", err)
			resp.DBError(ctx)
			return err
		}
		user.Flag.SetBit(offset)
		if err := tx.Save(&user).Error; err != nil {
			logger.Logger.WithContext(ctx).Info("更新用户选课记录失败", err)
			resp.DBError(ctx)
			return err
		}
		return nil // 成功，无错误返回
	}); err != nil {
		logger.Logger.WithContext(ctx).Info("事务回滚", err)
		return
	}
	// 事务成功，响应成功
	resp.Success(ctx, nil)
}
