package course

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"select-course/demo2/src/constant/code"
	"select-course/demo2/src/constant/keys"
	"select-course/demo2/src/models"
	"select-course/demo2/src/models/request"
	"select-course/demo2/src/storage/cache"
	"select-course/demo2/src/storage/database"
	"select-course/demo2/src/utils/bloom"
	"select-course/demo2/src/utils/logger"
	"select-course/demo2/src/utils/resp"
	"strconv"
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

	// 2. 校验操作
	// 2.1 判断课程是否存在
	if !bloom.CourseBloom.TestString(fmt.Sprintf("%d", req.CourseID)) {
		resp.Fail(ctx, code.NotFound, code.CourseNotFound, code.CourseNotFoundMsg)
		return
	}
	// 2.2 判断用户是否存在 这里可以防止缓存穿透
	if !bloom.UserBloom.TestString(fmt.Sprintf("%d", req.UserID)) {
		resp.Fail(ctx, code.NotFound, code.UserNotFound, code.UserNotFoundMsg)
		return
	}
	// 2.3. 用户是否已经选择该门课程
	if exist := cache.RDB.HExists(ctx, keys.UserCourseSetKey, fmt.Sprintf("%d", req.UserID)).Val(); exist {
		logger.Logger.WithContext(ctx).Info("用户已经选择该门课程")
		resp.Fail(ctx, code.Fail, code.CourseSelected, code.CourseSelectedMsg)
		return
	}
	// 2.4. 获取用户flag
	var user models.User
	if err := database.Client.Clauses(clause.Locking{Strength: "SHARE"}).
		Where("id=?", req.UserID).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.WithContext(ctx).Info("用户不存在", err)
			resp.Fail(ctx, code.NotFound, code.UserNotFound, code.UserNotFoundMsg)
			return
		}
	}

	// 2.5. 是否存在课程时间冲突
	key := fmt.Sprintf(keys.CourseHsetKey, req.CourseID)
	weekStr := cache.RDB.HGet(ctx, key, keys.CourseScheduleWeekKey).Val()
	durationStr := cache.RDB.HGet(ctx, key, keys.CourseScheduleDurationKey).Val()
	week, _ := strconv.Atoi(weekStr)
	duration, _ := strconv.Atoi(durationStr)
	offset := week*3 + duration - 1
	if user.Flag.TestBit(offset) {
		logger.Logger.WithContext(ctx).Info("用户已选该课程或时间冲突")
		resp.Fail(ctx, code.ParamErr, code.CourseTimeConflict, code.CourseTimeConflictMsg)
		return
	}

	// 2.6. 扣减课程容量与创建课程操作
	var err error
	var capacityStr string
	if capacityStr, err = cache.RDB.HGet(ctx, key, keys.CourseCapacityKey).Result(); err != nil {
		logger.Logger.WithContext(ctx).Info("获取课程容量失败", err)
	}
	if capacity, _ := strconv.Atoi(capacityStr); capacity <= 0 {
		logger.Logger.WithContext(ctx).Info("课程容量不足", err)
		resp.Fail(ctx, code.Fail, code.CourseFull, code.CourseFullMsg)
		return
	}
	// 3. 创建操作
	// 3.1 扣减课程容量
	if err := cache.RDB.HIncrBy(ctx, key, keys.CourseCapacityKey, -1).Err(); err != nil {
		logger.Logger.WithContext(ctx).Info("扣减课程容量失败", err)
		resp.DBError(ctx)
	}
	// 3.2 加入到用户课程集合
	if err := cache.RDB.SAdd(ctx, fmt.Sprintf(keys.UserCourseSetKey, req.UserID), req.CourseID).Err(); err != nil {
		logger.Logger.WithContext(ctx).Info("加入到用户课程集合失败", err)
		resp.DBError(ctx)
	}

	// 异步写入到数据库
	go Async2Mysql(ctx, &user, req.CourseID, offset)

	// 事务成功，响应成功
	resp.Success(ctx, nil)
}
func Async2Mysql(ctx *gin.Context, user *models.User, courseID uint, offset int) {
	// 2. 数据库操作
	if err := database.Client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 2.4 扣减课程库存
		if err := tx.Model(&models.Course{}).
			Where("id=?", courseID).
			Update("capacity", gorm.Expr("capacity - 1")).Error; err != nil {
			logger.Logger.WithContext(ctx).Info("更新课程容量失败", err)
			resp.DBError(ctx)
			return err
		}
		// 2.5 创建选课记录
		if err := tx.Create(&models.UserCourse{
			UserID:   user.ID,
			CourseID: courseID,
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
