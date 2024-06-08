package consumer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"select-course/demo4/src/constant/variable"
	"select-course/demo4/src/models"
	"select-course/demo4/src/models/mqm"
	"select-course/demo4/src/storage/database"
	"select-course/demo4/src/storage/mq"
	"select-course/demo4/src/utils/local"
	"select-course/demo4/src/utils/logger"
	"time"
)

var SelectConsumer *Select

type Select struct {
	channel *amqp.Channel
}

func InitSelectListener() error {

	SelectConsumer = &Select{
		channel: mq.Client,
	}
	err := SelectConsumer.Declare()
	if err != nil {
		fmt.Println(err)
		return err
	}
	SelectConsumer.Consumer()
	return nil

}

func (s *Select) Declare() error {
	var err error
	err = s.channel.ExchangeDeclare(
		variable.DeadExchange, variable.DeadKind,
		true, false, false, false, nil,
	)
	if err != nil {
		return err
	}
	_, err = s.channel.QueueDeclare(
		variable.DeadQueue, true,
		false, false, false, nil,
	)

	err = s.channel.ExchangeDeclare(variable.SelectExchange, variable.SelectKind,
		true, false, false, false, nil,
	)
	if err != nil {
		return err
	}
	err = s.channel.QueueBind(variable.DeadQueue, variable.DeadRoutingKey,
		variable.DeadExchange, false, nil)
	if err != nil {
		return err
	}
	_, err = s.channel.QueueDeclare(variable.SelectQueue, true,
		false, false, false, amqp.Table{
			"x-dead-letter-exchange":    variable.DeadExchange,
			"x-dead-letter-routing-key": variable.DeadRoutingKey,
		},
	)
	if err != nil {
		return err
	}
	err = s.channel.QueueBind(variable.SelectQueue, variable.SelectRoutingKey,
		variable.SelectExchange, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Select) Consumer() {
	results, err := SelectConsumer.channel.Consume(
		variable.SelectQueue,
		variable.SelectRoutingKey,
		false, // 关闭自动应答
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Logger.Error("消息接收失败", err)
		return
	}

	go func() {
		for res := range results {
			var msg *mqm.CourseReq
			err := json.Unmarshal(res.Body, &msg)
			if err != nil {
				logger.Logger.Error("消息反序列化失败", err)
				res.Reject(false)
				continue
			}

			err = database.Client.Transaction(func(tx *gorm.DB) error {
				switch msg.Type {
				case mqm.SelectType:
					if err := updateCourseCapacityAndUserCourse(tx, msg, true); err != nil {
						return err
					}
				case mqm.BackType:
					if err := updateCourseCapacityAndUserCourse(tx, msg, false); err != nil {
						return err
					}
				default:
					return fmt.Errorf("未知的消息类型: %s", msg.Type)
				}
				return nil
			})

			if err != nil {
				logger.Logger.Error("事务失败", err)
				res.Reject(false)
				continue
			}

			// 消息确认
			if err := res.Ack(false); err != nil {
				logger.Logger.Error("消息确认失败", err)
			}
		}
	}()
}

func updateCourseCapacityAndUserCourse(tx *gorm.DB, msg *mqm.CourseReq, selectAction bool) error {
	//1. 更新课程容量
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
	// 2. 获取课程时间并且计算offset
	offset, err := local.CalOffset(msg.CourseID)
	if err != nil {
		logger.Logger.Debug("计算offset", err)
		return err
	}
	// 3. 获取用户flag，更新flag
	var user models.User
	if err := database.Client.Where("id=?", msg.UserID).First(&user).Error; err != nil {
		logger.Logger.Debug("获取用户flag", err)
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

	// 4. 创建/更新选课记录
	var userCourse models.UserCourse
	if err := tx.Clauses(clause.Locking{Strength: "SHARE"}).
		Where("user_id=? and course_id=?", msg.UserID, msg.CourseID).
		First(&userCourse).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		userCourse = models.UserCourse{
			UserID:    msg.UserID,
			CourseID:  msg.CourseID,
			CreatedAt: msg.CreatedAt,
			IsDeleted: !selectAction,
		}

		if err := tx.Create(&userCourse).Error; err != nil {
			logger.Logger.Debug("创建/更新选课记录", err)
			return err
		}
	} else {
		userCourse.CreatedAt = msg.CreatedAt
		userCourse.IsDeleted = !selectAction
		if err := tx.Model(&userCourse).
			Where("user_id=? and course_id=?", msg.UserID, msg.CourseID).
			Updates(userCourse).Error; err != nil {
			logger.Logger.Debug("创建/更新选课记录", err)
			return err
		}
	}

	return nil
}

func (s *Select) Product(msg *mqm.CourseReq) {
	// 毫秒 记录每条消息的时间,确保加入到死信队列后期执行消费的顺序
	msg.CreatedAt = time.Now().UnixMilli()
	bytes, err := json.Marshal(msg)
	if err != nil {
		logger.Logger.Error("消息序列化失败", err)
		return
	}
	err = s.channel.Publish(
		variable.SelectExchange,
		variable.SelectRoutingKey,
		false,
		false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        bytes,
		})
	if err != nil {
		logger.Logger.Error("消息发送失败", err)
		return
	}
}
