package course

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"select-course/demo4/src/constant/code"
	"select-course/demo4/src/constant/keys"
	"select-course/demo4/src/constant/lua"
	"select-course/demo4/src/models"
	"select-course/demo4/src/models/mqm"
	"select-course/demo4/src/models/request"
	"select-course/demo4/src/storage/cache"
	"select-course/demo4/src/storage/database"
	"select-course/demo4/src/utils/bloom"
	"select-course/demo4/src/utils/consumer"
	"select-course/demo4/src/utils/local"
	"select-course/demo4/src/utils/logger"
	"select-course/demo4/src/utils/resp"
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

func handleRequestError(ctx *gin.Context, err error) {
	logger.Logger.WithContext(ctx).Info("参数校验失败", err)
	resp.ParamErr(ctx)
	return
}

func validateAndLogError(ctx *gin.Context, err error, failCode int, failMsg string) {
	if err != nil {
		logger.Logger.WithContext(ctx).Info(failMsg, err)
		resp.Fail(ctx, failCode, failCode, failMsg)
		return
	}
}

func executeLuaScript(ctx context.Context, rdb *redis.Client, script *redis.Script, keys []string, args ...interface{}) (interface{}, error) {
	val, err := script.Run(ctx, rdb, keys, args...).Result()
	if err != nil {
		logger.Logger.WithContext(ctx).Info("执行lua脚本失败", err)
		return nil, err
	}
	return val, nil
}

func SelectCourse(ctx *gin.Context) {
	var req request.SelectCourseReq
	if err := ctx.ShouldBind(&req); err != nil {
		handleRequestError(ctx, err)
		return
	}

	if !bloom.CourseBloom.TestString(fmt.Sprintf("%d", req.CourseID)) {
		resp.Fail(ctx, code.Fail, code.CourseNotFound, code.CourseNotFoundMsg)
		return
	}
	if !bloom.UserBloom.TestString(fmt.Sprintf("%d", req.UserID)) {
		resp.Fail(ctx, code.Fail, code.UserNotFound, code.UserNotFoundMsg)
		return
	}

	offset, err := local.CalOffset(req.CourseID)
	validateAndLogError(ctx, err, code.Fail, "获取课程时间失败")

	val, err := executeLuaScript(ctx, cache.RDB, redis.NewScript(lua.CourseSelectLuaScript), []string{
		fmt.Sprintf(keys.UserCourseSetKey, req.UserID),
		strconv.Itoa(int(req.CourseID)),
		fmt.Sprintf(keys.CourseHsetKey, req.CourseID),
		keys.CourseCapacityKey,
		fmt.Sprintf(keys.UserCourseScheduleBitMapKey, req.UserID),
		strconv.Itoa(offset),
	}, req.UserID, req.CourseID)
	if err != nil {
		logger.Logger.WithContext(ctx).Info("执行lua脚本失败", err)
		resp.Fail(ctx, code.Fail, code.Fail, code.FailMsg)
		return
	}

	switch val.(int64) {
	case lua.CourseSelectOK:
		logger.Logger.WithContext(ctx).Info("选课成功")
		consumer.SelectConsumer.Product(&mqm.CourseReq{UserID: req.UserID, CourseID: req.CourseID, Type: mqm.SelectType})
		resp.Success(ctx, nil)
	case lua.CourseSelected:
		logger.Logger.WithContext(ctx).Info("用户已经选择该门课程")
		resp.Fail(ctx, code.Fail, code.CourseSelected, code.CourseSelectedMsg)
	case lua.CourseFull:
		logger.Logger.WithContext(ctx).Info("课程已满")
		resp.Fail(ctx, code.Fail, code.CourseFull, code.CourseFullMsg)
	case lua.CourseTimeConflict:
		logger.Logger.WithContext(ctx).Info("课程时间冲突")
		resp.Fail(ctx, code.Fail, code.CourseTimeConflict, code.CourseTimeConflictMsg)
	default:
		logger.Logger.WithContext(ctx).Info("未知错误")
		resp.Fail(ctx, code.Fail, code.Fail, code.FailMsg)
	}
}

func BackCourse(ctx *gin.Context) {
	var req request.SelectCourseReq
	if err := ctx.ShouldBind(&req); err != nil {
		handleRequestError(ctx, err)
		return
	}

	offset, err := local.CalOffset(req.CourseID)
	validateAndLogError(ctx, err, code.Fail, "计算offset失败")

	val, err := executeLuaScript(ctx, cache.RDB, redis.NewScript(lua.CourseBackLuaScript), []string{
		fmt.Sprintf(keys.UserCourseSetKey, req.UserID),
		strconv.Itoa(int(req.CourseID)),
		fmt.Sprintf(keys.CourseHsetKey, req.CourseID),
		keys.CourseCapacityKey,
		fmt.Sprintf(keys.UserCourseScheduleBitMapKey, req.UserID),
		strconv.Itoa(offset),
	}, req.UserID, req.CourseID)
	if err != nil {
		return
	}

	switch val.(int64) {
	case lua.CourseBackOK:
		logger.Logger.WithContext(ctx).Info("退课成功")
		consumer.SelectConsumer.Product(&mqm.CourseReq{
			UserID:   req.UserID,
			CourseID: req.CourseID,
			Type:     mqm.BackType,
		})
		resp.Success(ctx, nil)
	case lua.CourseNotSelected:
		logger.Logger.WithContext(ctx).Info("退课失败：课程未选择")
		resp.Fail(ctx, code.Fail, code.CourseNotSelected, code.CourseNotSelectedMsg)
	default:
		logger.Logger.WithContext(ctx).Info("未知错误")
		resp.Fail(ctx, code.Fail, code.Fail, code.FailMsg)
	}
}
