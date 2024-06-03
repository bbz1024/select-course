package course

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"select-course/demo2/src/constant/code"
	"select-course/demo2/src/constant/keys"
	"select-course/demo2/src/constant/lua"
	"select-course/demo2/src/models"
	"select-course/demo2/src/models/request"
	"select-course/demo2/src/storage/cache"
	"select-course/demo2/src/storage/database"
	"select-course/demo2/src/utils/bloom"
	"select-course/demo2/src/utils/logger"
	"select-course/demo2/src/utils/resp"
	"strconv"
	"sync"
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
	// 3. 执行操作
	if err := database.Client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 获取课程的week和duration
		var course models.Course
		if err := tx.Model(&models.Course{}).Preload("Schedule").First(&course, req.CourseID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			return err
		}
		// 2. 获取用户的flag字段
		var user models.User
		//
		if err := tx.Clauses(clause.Locking{Strength: "SHARE"}).Where("id = ?", req.UserID).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			return err
		}

		/*
		   args:
		      key1 = 用户key
		      key2 = 课程id
		      key3 = 课程key
		      key4 = capacity key
		   return:
		      0 : 执行正常
		      1 : 用户是否已经选择了
		      2 : 课程满了
		*/
		// 3. 判断是否选存在时间冲突
		offset := int(course.Schedule.Week*3) + int(course.Schedule.Duration)
		if user.Flag.TestBit(offset) {
			resp.Fail(ctx, code.Fail, code.CourseSelected, code.CourseSelectedMsg)
			return errors.New("用户已经选过该课程")
		}
		// 4. 修改用户flag
		user.Flag.SetBit(offset)
		if err := tx.Save(&user).Error; err != nil {
			return err
		}
		// 5. 执行脚本进行扣减课程和创建
		script := redis.NewScript(lua.CourseSelectLuaScript)
		var val interface{}
		var err2 error
		if val, err2 = script.Run(ctx, cache.RDB, []string{
			fmt.Sprintf(keys.UserCourseSetKey, req.UserID),
			strconv.Itoa(int(req.CourseID)),
			fmt.Sprintf(keys.CourseHsetKey, req.CourseID),
			keys.CourseCapacityKey,
		}, req.UserID, req.CourseID).Result(); err2 != nil {
			logger.Logger.WithContext(ctx).Info("执行lua脚本失败", err2)
			resp.Fail(ctx, code.Fail, code.CourseSelected, code.CourseSelectedMsg)
			return err2
		}
		switch val.(int64) {
		case lua.CourseSelectOK:
			logger.Logger.WithContext(ctx).Info("选课成功")
			go AsyncSelect2Mysql(ctx, req.UserID, req.CourseID)
		case lua.CourseSelected:
			logger.Logger.WithContext(ctx).Info("用户已经选择该门课程")
			resp.Fail(ctx, code.Fail, code.CourseSelected, code.CourseSelectedMsg)
			return errors.New("用户已经选择该门课程")
		case lua.CourseFull:
			logger.Logger.WithContext(ctx).Info("课程已满")
			resp.Fail(ctx, code.Fail, code.CourseFull, code.CourseFullMsg)
			return errors.New("课程已满")
		default:
			logger.Logger.WithContext(ctx).Info("未知错误")
			return errors.New("未知错误")
		}

		return nil
	}); err != nil {
		logger.Logger.WithContext(ctx).Info("事务失败", err)
		resp.Fail(ctx, code.Fail, code.DBError, code.DBErrorMsg)
		return
	}

	// 事务成功，响应成功
	resp.Success(ctx, nil)
}

var mutex sync.Mutex

func AsyncSelect2Mysql(ctx *gin.Context, userID uint, courseID uint) {
	mutex.Lock()
	defer mutex.Unlock()
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
			UserID:   userID,
			CourseID: courseID,
		}).Error; err != nil {
			logger.Logger.WithContext(ctx).Info("创建选课记录失败", err)
			resp.DBError(ctx)
			return err
		}
		return nil // 成功，无错误返回
	}); err != nil {
		logger.Logger.WithContext(ctx).Info("事务回滚", err)
		return
	}
}

func AsyncBack2Mysql(ctx *gin.Context, userID uint, courseID uint) {
	mutex.Lock()
	defer mutex.Unlock()
	if err := database.Client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Course{}).
			Where("id=?", courseID).
			Update("capacity", gorm.Expr("capacity + 1")).Error; err != nil {
			logger.Logger.WithContext(ctx).Info("更新课程容量失败", err)
			resp.DBError(ctx)
			return err
		}
		if err := tx.Where("user_id=? and course_id=?", userID, courseID).Delete(&models.UserCourse{}).Error; err != nil {
			logger.Logger.WithContext(ctx).Info("删除选课记录失败", err)
			resp.DBError(ctx)
			return err
		}
		return nil
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

		if !user.Flag.TestBit(int(course.Schedule.Week*3) + int(course.Schedule.Duration)) {
			logger.Logger.WithContext(ctx).Info("用户没有选过该课程")
			resp.Fail(ctx, code.Fail, code.CourseNotSelected, code.CourseNotSelectedMsg)
			return errors.New("用户没有选过该课程")
		}
		// 2.6 更新用户选课记录
		offset := int(course.Schedule.Week)*3 + int(course.Schedule.Duration)
		user.Flag.ClearBit(offset)
		if err := tx.Save(&user).Error; err != nil {
			logger.Logger.WithContext(ctx).Info("更新用户选课记录失败", err)
			resp.DBError(ctx)
			return err
		}
		script := redis.NewScript(lua.CourseBackLuaScript)
		var val interface{}
		var err2 error
		if val, err2 = script.Run(ctx, cache.RDB, []string{
			fmt.Sprintf(keys.UserCourseSetKey, req.UserID),
			strconv.Itoa(int(req.CourseID)),
			fmt.Sprintf(keys.CourseHsetKey, req.CourseID),
			keys.CourseCapacityKey,
		}, req.UserID, req.CourseID).Result(); err2 != nil {
			logger.Logger.WithContext(ctx).Info("执行lua脚本失败", err2)
			resp.Fail(ctx, code.Fail, code.CourseNotSelected, code.CourseNotSelectedMsg)
			return err2
		}

		switch val.(int64) {
		case lua.CourseBackOK:
			logger.Logger.WithContext(ctx).Info("退课成功")
			go AsyncBack2Mysql(ctx, req.UserID, req.CourseID)
			return nil
		case lua.CourseNotSelected:
			logger.Logger.WithContext(ctx).Info()
			resp.Fail(ctx, code.Fail, code.CourseNotSelected, code.CourseNotSelectedMsg)
			return errors.New("用户没有选过该课程")
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
