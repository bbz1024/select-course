package consumer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/avast/retry-go"
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
	//	 发布确认
	confirmMsg chan amqp.Confirmation
	// return
	returnMsg chan amqp.Return
}

func InitSelectListener() error {
	channel, err := mq.Client.Channel()
	SelectConsumer = &Select{
		channel: channel,
	}
	err = SelectConsumer.Declare()
	if err != nil {
		return err
	}
	return nil
}
func (s *Select) Close() error {
	return s.channel.Close()
}
func (s *Select) Declare() error {

	var err error
	// 死信队列
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
	// 发布确认 异步
	err = s.channel.Confirm(false)
	if err != nil {
		return err
	}
	s.confirmMsg = s.channel.NotifyPublish(make(chan amqp.Confirmation))
	s.returnMsg = s.channel.NotifyReturn(make(chan amqp.Return))
	go s.ListenConfirm() // 无法到达交换机、或者消息路由成功
	go s.ListenReturns() // 无法路由到正确队列或者路由键规则匹配
	return nil
}
func (s *Select) ListenConfirm() {
	for msg := range s.confirmMsg {
		// 无法达到交换机
		fmt.Println(msg)
	}
}

func (s *Select) ListenReturns() {
	// 无法路由到队列
	for msg := range s.returnMsg {
		// 重新发送，尝试3次
		err := retry.Do(func() error {
			err := s.channel.Publish(
				variable.SelectExchange,
				variable.SelectRoutingKey,
				true,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        msg.Body,
				},
			)
			return err
		}, retry.Attempts(3), retry.Delay(time.Millisecond*100))
		if err != nil {
			logger.Logger.Error("消息重新发送失败", err)
			// 丢入到死信队列，进行补偿操作
			err := s.channel.Publish(
				variable.DeadExchange,
				variable.DeadRoutingKey,
				true,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        msg.Body,
				},
			)
			if err != nil {
				logger.Logger.Error("消息发送失败", err)

			}
		}
	}
}
func (s *Select) Consumer() error {
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
		return err
	}

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
				if err := updateCourseCapacity(tx, msg, true); err != nil {
					return err
				}
				if err := updateUserFlag(tx, msg, true); err != nil {
					return err
				}
				if err := updateUserCourseState(tx, msg, true); err != nil {
					return err
				}
			case mqm.BackType:
				if err := updateCourseCapacity(tx, msg, false); err != nil {
					return err
				}
				if err := updateUserFlag(tx, msg, false); err != nil {
					return err
				}
				if err := updateUserCourseState(tx, msg, false); err != nil {
					return err
				}

			default:
				return fmt.Errorf("未知的消息类型: %d", msg.Type)
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
	return nil

}

// 更新课程容量
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

// 更新用户flag
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

// 创建用户选课记录
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
			IsDeleted: !selectAction,
		}).Error; err != nil {
			logger.Logger.Info("创建选课记录失败", err)
			return err
		}
		return nil
	}
	// 存在，还是判断msg是否在创建时间之前，不在的话，不更新。
	if err := tx.Model(&models.UserCourse{}).
		Where("user_id=? and course_id=? and created_at < ?", msg.UserID, msg.CourseID, msg.CreatedAt).
		Update("is_deleted", !selectAction).
		Update("created_at", msg.CreatedAt).Error; err != nil {
		logger.Logger.Info("更新选课记录失败", err)
		return err
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
		true,
		false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        bytes,
		})
	if err != nil {
		logger.Logger.Error("消息发送失败", err)
		return
	}
}
