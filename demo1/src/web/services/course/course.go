package course

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"select-course/demo1/src/constant/code"
	"select-course/demo1/src/models"
	"select-course/demo1/src/models/request"
	"select-course/demo1/src/storage/database"
	"select-course/demo1/src/utils/logger"
	"select-course/demo1/src/utils/resp"
	"sync"
)

func GetCourseList(ctx *gin.Context) {
	var courseList []*models.Course
	if err := database.Client.Find(&courseList).Error; err != nil {
		resp.DBError(ctx)
		return
	}
	resp.Success(ctx, courseList)
}

var CourseSelect sync.Mutex

func SelectCourse(ctx *gin.Context) {
	// 1. 校验参数
	var req request.SelectCourseReq
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Logger.Info("参数校验失败", err)
		resp.ParamErr(ctx)
		return
	}
	// 加全局锁
	CourseSelect.Lock()
	defer CourseSelect.Unlock()

	// 2. 校验操作
	// 2.1 课程是否存在
	var course models.Course
	if err := database.Client.First(&course, req.CourseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Info("课程不存在", err)
			resp.Fail(ctx, code.NotFound, code.CourseNotFound, code.CourseNotFoundMsg)
			return
		}
		logger.Logger.Error("课程查询失败", err)
		resp.DBError(ctx)
		return
	}
	// 3. 用户是否已经选择了该门课程 || 选的课程上课时间存在冲突
	// 获取用户已经选择的所有课程
	var userCourses []*models.UserCourse
	if err := database.Client.Where("user_id = ?", req.UserID).Preload("Course").Find(&userCourses).Error; err != nil {
		logger.Logger.Info("查询用户课程失败", err)
		resp.DBError(ctx)
		return
	}
	for _, userCourse := range userCourses {
		if userCourse.CourseID == uint(req.CourseID) {
			logger.Logger.Info("用户已经选择该门课程")
			resp.Fail(ctx, code.NotFound, code.CourseSelected, code.CourseSelectedMsg)
			return
		}
		// 3.1 课程上课时间存在冲突
		if userCourse.Course.Week == course.Week && userCourse.Course.Duration == course.Duration {
			logger.Logger.Info("课程上课时间存在冲突")
			resp.Fail(ctx, code.NotFound, code.CourseTimeConflict, code.CourseTimeConflictMsg)
			return
		}
	}
	if course.Capacity == 0 {
		logger.Logger.Info("课程已满")
		resp.Fail(ctx, code.NotFound, code.CourseFull, code.CourseFullMsg)
		return
	}
	if err := database.Client.Model(&course).Update("capacity", course.Capacity-1).Error; err != nil {
		logger.Logger.Info("更新课程容量失败", err)
		resp.DBError(ctx)
		return
	}
	// 4. 为用户创建选课记录
	if err := database.Client.Create(&models.UserCourse{
		UserID:   uint(req.UserID),
		CourseID: uint(req.CourseID),
	}).Error; err != nil {
		logger.Logger.Info("创建选课记录失败", err)
		resp.DBError(ctx)
		return
	}
	resp.Success(ctx, nil)

	/*
		以上存在并发安全问题，由于查询，扣减与创建选课记录并不是一气呵成的（原子性），不能保证在校验过程中，有其他的线程进行对课程的库存进行扣将
		方法1：
			通过加锁操作使可能存在的并发进行串行化
		方法2：
			通过mysql提供的事务机制
		细粒度对比：方法1 > 方法2
	*/

}
