package mock

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math/rand"
	"select-course/demo4/src/constant/keys"
	"select-course/demo4/src/constant/variable"
	"select-course/demo4/src/models"
	"select-course/demo4/src/models/mqm"
	"select-course/demo4/src/storage/cache"
	"select-course/demo4/src/storage/database"
	"select-course/demo4/src/storage/mq"
	"select-course/demo4/src/utils/local"
	"select-course/demo4/src/utils/logger"
	"testing"
)

func TestInsertSchedule(t *testing.T) {
	var schedule []models.Schedule
	week := 5
	duration := 3
	for i := 0; i < week; i++ {
		for j := 1; j <= duration; j++ {
			schedule = append(schedule, models.Schedule{
				Week:     models.Week(i),
				Duration: models.Duration(j),
			})
		}
	}
	fmt.Println(schedule)
	database.Client.Create(&schedule)

}
func TestInsertCourseCategory(t *testing.T) {
	var courseCategory []models.CourseCategory
	for i := 1; i < 6; i++ {
		courseCategory = append(courseCategory, models.CourseCategory{
			Name: fmt.Sprintf("分类%d", i),
			BaseModel: models.BaseModel{
				ID: uint(i),
			},
		})
	}
	database.Client.Create(&courseCategory)
}
func TestInsertCourse(t *testing.T) {
	// 星期一 上午 08:10 ~ 11:50 | 星期二 下午 14:10 ~ 16:50 | 星期三 晚上 18:50 ~ 21:20
	var course []models.Course
	for i := 1; i < 11; i++ {
		course = append(course, models.Course{
			BaseModel: models.BaseModel{
				ID: uint(i),
			},
			Title:      fmt.Sprintf("课程%d", i),
			CategoryID: uint(rand.Intn(5) + 1),
			ScheduleID: uint(rand.Intn(15) + 1),
			Capacity:   10,
		})
	}
	database.Client.Create(&course)

}

// 预热到redis
func TestPreheatMysql2Redis(t *testing.T) {
	var course []models.Course
	if err := database.Client.
		Preload("Schedule").
		Preload("Category").
		Find(&course).Error; err != nil {
		t.Error(err)
	}
	for _, v := range course {
		if err := cache.RDB.HSet(
			context.Background(),
			fmt.Sprintf(keys.CourseHsetKey, v.ID),
			keys.CourseCategoryIDKey, v.CategoryID,
			keys.CourseScheduleDurationKey, uint(v.Schedule.Duration),
			keys.CourseScheduleWeekKey, uint(v.Schedule.Week),
			keys.CourseCapacityKey, v.Capacity,
		).Err(); err != nil {
			t.Error(err)
		}
	}
}

// 手动补偿操作
func TestHandlerDeadQueue(t *testing.T) {
	channel, err := mq.Client.Channel()
	//接收消息
	results, err := channel.Consume(
		variable.DeadQueue,
		variable.DeadRoutingKey,
		false, // 关闭自动应答
		false, //
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Logger.Error("消息接收失败", err)
		return
	}
	// 获取死信队列

	for res := range results {
		var msg *mqm.CourseReq
		var err error
		err = json.Unmarshal(res.Body, &msg)
		if err != nil {
			if err := res.Nack(false, true); err != nil {
				logger.Logger.Error("消息拒绝失败", err)
			}
			logger.Logger.Error("消息反序列化失败", err)
			continue
		}
		err = database.Client.Transaction(func(tx *gorm.DB) error {
			if err := updateCourseCapacity(tx, msg, msg.Type == mqm.SelectType); err != nil {
				logger.Logger.Info("更新课程容量失败", err)
				return err
			}
			if err := updateUserCourseState(tx, msg, msg.Type == mqm.SelectType); err != nil {
				logger.Logger.Info("更新用户课程状态失败", err)
				return err
			}
			if err := updateUserFlag(tx, msg, msg.Type == mqm.SelectType); err != nil {
				logger.Logger.Info("更新课程时间失败", err)
				return err
			}
			return nil // 成功，无错误返回
		})
		if err != nil {
			logger.Logger.Error("事务处理失败", err)
			// 放回队列
			err := res.Nack(false, true)
			if err != nil {
				logger.Logger.Error("消息拒绝失败", err)
			}
			continue
		}
		if err := res.Ack(false); err != nil {
			logger.Logger.Error("消息确认失败", err)
		}
	}
}
func updateCourseCapacity(tx *gorm.DB, msg *mqm.CourseReq, selectAction bool) error {
	capacityOp := gorm.Expr("capacity - 1")
	if !selectAction {
		capacityOp = gorm.Expr("capacity + 1")
	}

	if err := tx.Model(&models.Course{}).
		Where("id=?", msg.CourseID).
		Update("capacity", capacityOp).Error; err != nil {
		logger.Logger.Debug("更新课程容量", err)
		return err
	}
	return nil
}
func updateUserFlag(tx *gorm.DB, msg *mqm.CourseReq, selectAction bool) error {

	// 2. 获取课程时间并且计算offset
	offset, err := local.CalOffset(msg.CourseID)
	if err != nil {
		logger.Logger.Debug("计算offset", err)
		return err
	}
	// 3. 获取用户flag，更新flag
	var user models.User
	if err := database.Client.Where("id=?", msg.UserID).First(&user).Error; err != nil {
		logger.Logger.Debug("获取用户信息", err)
		return fmt.Errorf("用户ID: %d 不存在", msg.UserID)
	}

	if selectAction {
		user.Flag.SetBit(offset)
	} else {
		user.Flag.ClearBit(offset)
	}
	if err := tx.Save(&user).Error; err != nil {
		logger.Logger.Debug("更新用户flag", err)
		return err
	}
	return nil
}
func updateUserCourseState(tx *gorm.DB, msg *mqm.CourseReq, selectAction bool) error {
	var userCourse models.UserCourse
	if err := tx.Clauses(clause.Locking{Strength: "SHARE"}).
		Where("user_id=? and course_id=? ", msg.UserID, msg.CourseID).
		First(&userCourse).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		// 不存在
		if err := tx.Create(&models.UserCourse{
			UserID:    msg.UserID,
			CourseID:  msg.CourseID,
			CreatedAt: msg.CreatedAt, // 创建时记录创建时间
			UpdatedAt: msg.CreatedAt,
			IsDeleted: !selectAction,
		}).Error; err != nil {
			logger.Logger.Info("创建选课记录失败", err)
			return err
		}
		return nil
	}
	// 存在，判断是否在创建时间之后，在的话，不更新。（失效）
	if msg.CreatedAt < userCourse.UpdatedAt {
		return nil
	}
	// 存在，还是判断msg是否在创建时间之前，不在的话，不更新。
	if err := tx.Model(&models.UserCourse{}).
		Where("user_id=? and course_id=?", msg.UserID, msg.CourseID).
		Update("is_deleted", !selectAction).
		Update("updated_at", msg.CreatedAt).Error; err != nil {
		logger.Logger.Info("更新选课记录失败", err)
		return err
	}
	return nil
}

// 测试redis中课程总数
func TestTotalRedisCourse(t *testing.T) {
	user := 100
	total := 0
	for i := 1; i <= user; i++ {
		v := cache.RDB.SCard(context.Background(), fmt.Sprintf(keys.UserCourseSetKey, i)).Val()

		total += int(v)
	}
	fmt.Println(total)
}
