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
	if err := database.Client.
		Preload("Category").
		Preload("Schedule").
		Find(&courseList).Error; err != nil {
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

	// 2. 数据库操作
	if err := database.Client.WithContext(context.Background()).Transaction(func(tx *gorm.DB) error {
		// 2.1 用户是否已选该门课程（这里可以不用判断因为如果用户选择了该门课程就会存在时间冲突
		if err := tx.Clauses(clause.Locking{Strength: "SHARE"}).
			Where("id=?", req.UserID).
			First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Logger.WithContext(ctx).Info("用户不存在", err)
				resp.Fail(ctx, code.NotFound, code.UserNotFound, code.UserNotFoundMsg)
			}
		}

		// 2.2 检查课程是否存在和库存
		var course models.Course
		if err := tx.Clauses(clause.Locking{Strength: "SHARE"}).
			Model(models.Course{}).
			Preload("Schedule").
			Where("id=? and capacity > 0", req.CourseID).
			First(&course).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Logger.WithContext(ctx).Info("课程不存在或库存不足", err)
				resp.Fail(ctx, code.NotFound, code.CourseNotFound, code.CourseNotFoundMsg)
				return nil //无需回滚，业务逻辑错误
			}
			return err
		}
		// 2.3 检查用户是否已选该课程或时间冲突
		offset := int(course.Schedule.Week)*3 + int(course.Schedule.Duration) - 1
		if user.Flag.TestBit(offset) {
			logger.Logger.WithContext(ctx).Info("用户已选该课程或时间冲突")
			resp.Fail(ctx, code.ParamErr, code.CourseTimeConflict, code.CourseTimeConflictMsg)
			return errors.New("用户已选该课程或时间冲突")
		}

		// 2.4 扣减课程库存
		if err := tx.Model(&course).
			Update("capacity", gorm.Expr("capacity - 1")).Error; err != nil {
			logger.Logger.WithContext(ctx).Info("更新课程容量失败", err)
			resp.DBError(ctx)
			return err
		}
		// 2.5 创建选课记录
		if err := tx.Create(&models.UserCourse{
			UserID:   req.UserID,
			CourseID: req.CourseID,
		}).Error; err != nil {
			logger.Logger.WithContext(ctx).Info("创建选课记录失败", err)
			resp.DBError(ctx)
			return err
		}
		// 2.6 更新用户选课记录
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

func BackCourse(ctx *gin.Context) {
	// 1. 参数校验
	var req request.SelectCourseReq
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Logger.WithContext(ctx).Info("参数校验失败", err)
		resp.ParamErr(ctx)
		return
	}
	// 2. 数据库操作
	if err := database.Client.WithContext(context.Background()).Transaction(func(tx *gorm.DB) error {
		// 2.1 用户是否选择了该课程
		var userCourse models.UserCourse
		if err := tx.Clauses(clause.Locking{Strength: "SHARE"}).
			Where("user_id=? and course_id=?", req.UserID, req.CourseID).
			First(&userCourse).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Logger.WithContext(ctx).Info("用户未选该课程", err)
				resp.Fail(ctx, code.ParamErr, code.CourseNotSelected, code.CourseNotSelectedMsg)
				return err
			}
			resp.DBError(ctx)
			return err
		}
		// 2.2 获取课程信息，获取schedule的week和duration 计算出offset
		var course models.Course
		if err := tx.Model(models.Course{}).
			Preload("Schedule").
			Where("id=?", req.CourseID).
			First(&course).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Logger.WithContext(ctx).Info("课程不存在", err)
				resp.Fail(ctx, code.NotFound, code.CourseNotFound, code.CourseNotFoundMsg)
				return err
			}
			return err
		}
		var user models.User
		//	2.3 这里必须加锁，因为要保证flag字段并发安全
		if err := tx.
			Clauses(clause.Locking{Strength: "SHARE"}).
			Where("id=?", req.UserID).
			First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Logger.WithContext(ctx).Info("用户不存在", err)
				resp.Fail(ctx, code.NotFound, code.UserNotFound, code.UserNotFoundMsg)
				return err
			}
			resp.DBError(ctx)
			return err
		}
		// 2.4 删除选课记录
		if err := tx.Where("user_id=? and course_id=?", req.UserID, req.CourseID).
			Delete(&models.UserCourse{}).Error; err != nil {
			logger.Logger.WithContext(ctx).Info("删除选课记录失败", err)
			resp.DBError(ctx)
			return err
		}
		// 2.5 课程容量+1
		if err := tx.Model(&course).
			Update("capacity", gorm.Expr("capacity + 1")).Error; err != nil {
			logger.Logger.WithContext(ctx).Info("更新课程容量失败", err)
			resp.DBError(ctx)
			return err
		}
		// 2.6 更新用户选课记录
		offset := int(course.Schedule.Week)*3 + int(course.Schedule.Duration) - 1
		user.Flag.ClearBit(offset)
		if err := tx.Save(&user).Error; err != nil {
			logger.Logger.WithContext(ctx).Info("更新用户选课记录失败", err)
			resp.DBError(ctx)
			return err
		}
		return nil
	}); err != nil {
		logger.Logger.WithContext(ctx).Info("事务回滚", err)
		return
	}
	// 事务成功，响应成功
	resp.Success(ctx, nil)
}

func MyCourseList(ctx *gin.Context) {
	var req request.UserReq
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Logger.Info("参数校验失败", err)
		resp.ParamErr(ctx)
		return
	}
	var userCourse []*models.UserCourse
	if err := database.Client.
		Preload("Course").
		Preload("Course.Category").
		Preload("Course.Schedule").
		Where("user_id=?", req.UserID).
		Find(&userCourse).
		Error; err != nil {
		logger.Logger.WithContext(ctx).Info("获取用户选课记录失败", err)
		resp.DBError(ctx)
		return
	}
	var course []*models.Course
	for _, v := range userCourse {
		course = append(course, v.Course)
	}
	resp.Success(ctx, course)
}
