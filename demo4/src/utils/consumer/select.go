package consumer

import (
	"context"
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
	"select-course/demo4/src/storage/cache"
	"select-course/demo4/src/storage/database"
	"select-course/demo4/src/storage/mq"
	"select-course/demo4/src/utils/local"
	"select-course/demo4/src/utils/logger"
	"sync"
	"sync/atomic"
	"time"
)

// TestMsgLose 测试消息是否被丢失
const TestMsgLose = true

var SelectConsumer *Select

// UnconfirmedMessage 定义未确认消息结构体
type UnconfirmedMessage struct {
	Body []byte
}

type Select struct {
	channel    *amqp.Channel
	confirmMsg chan amqp.Confirmation
	// returnMsg  chan amqp.Return
	// 使用sync.Map记录未确认消息
	unconfirmedMessages sync.Map
	cnt                 atomic.Int64
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
	//s.returnMsg = s.channel.NotifyReturn(make(chan amqp.Return))
	go s.ListenConfirm() // 无法到达交换机
	//go s.ListenReturns() // 无法路由到正确队列或者路由键规则匹配，但是消息已经在交换机上持久化了的
	return nil
}
func (s *Select) ListenConfirm() {
	for msg := range s.confirmMsg {
		if TestMsgLose {
			cache.RDB.HIncrBy(context.Background(), "record", "total", 1)
		}
		if !msg.Ack {
			if TestMsgLose {
				cache.RDB.HIncrBy(context.Background(), "record", "ack-fail", 1)
			}
			val, ok := s.unconfirmedMessages.Load(msg.DeliveryTag)
			if !ok {
				logger.Logger.Error("消息确认失败", msg)
				continue
			}
			data := val.(UnconfirmedMessage)
			// 重新发送，尝试3次
			var cnt int
			if err := retry.Do(func() error {
				if err := s.PushDeadQueue(data.Body); err != nil {
					cnt++
					logger.Logger.Error("消息发送失败%d次", cnt)
					return err
				}
				return nil
			}, retry.Attempts(3), retry.Delay(time.Millisecond*100)); err != nil {
				logger.Logger.Error("消息重新发送失败", err)
				// 丢入到死信队列，进行补偿操作
				if err := s.PushDeadQueue(data.Body); err != nil {
					logger.Logger.Error("消息发送失败", err)
				}
			}
		}
		s.unconfirmedMessages.Delete(msg.DeliveryTag)
	}
}

/*
	func (s *Select) ListenReturns() {
		// 无法路由到队列,但是消息到达交换机，如果设置了死信队列就没必要写了。且这里存在消息的重复发送
		for msg := range s.returnMsg {
			// 重新发送，尝试3次
			var cnt int
			err := retry.Do(func() error {
				if err := s.channel.Publish(
					variable.SelectExchange,
					variable.SelectRoutingKey,
					true,
					false,
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        msg.Body,
					},
				); err != nil {
					cnt++
					logger.Logger.Error("消息发送失败%d次", cnt)
					return err
				}
				return nil
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
*/
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
	//time.Sleep(time.Second * 20)
	for res := range results {
		if TestMsgLose {
			cache.RDB.HIncrBy(context.Background(), "record", "consume-total", 1)

		}
		var msg *mqm.CourseReq
		if err := json.Unmarshal(res.Body, &msg); err != nil {
			logger.Logger.Error("消息反序列化失败", err)
			if err := res.Reject(false); err != nil {
				logger.Logger.Error("消息确认失败", err)
			}
			if TestMsgLose {
				cache.RDB.HIncrBy(context.Background(), "record", "consume-fail", 1)

			}
			continue
		}
		//time.Sleep(time.Millisecond * 500)
		//fmt.Println(msg)
		if err = database.Client.Transaction(func(tx *gorm.DB) error {
			if err := updateCourseCapacity(tx, msg, msg.Type == mqm.SelectType); err != nil {
				return err
			}
			if err := updateUserFlag(tx, msg, msg.Type == mqm.SelectType); err != nil {
				return err
			}
			if err := updateUserCourseState(tx, msg, msg.Type == mqm.SelectType); err != nil {
				return err
			}

			if TestMsgLose {
				cache.RDB.HIncrBy(context.Background(), "record", "transaction-success", 1)
			}
			return nil
		}); err != nil {
			logger.Logger.Error("事务失败", err)
			if err := res.Reject(false); err != nil {

				logger.Logger.Error("消息确认失败", err)
			}
			if TestMsgLose {
				cache.RDB.HIncrBy(context.Background(), "record", "consume-fail", 1)
				cache.RDB.HIncrBy(context.Background(), "record", "transaction-fail", 1)

			}
			continue
		}

		// 消息确认
		if err := res.Ack(false); err != nil {
			logger.Logger.Error("消息确认失败", err)
		}
		if TestMsgLose {
			cache.RDB.HIncrBy(context.Background(), "record", "consume-success", 1)
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
		logger.Logger.Error("更新课程容量", err)
		return err
	}
	return nil
}

// 更新用户flag
func updateUserFlag(tx *gorm.DB, msg *mqm.CourseReq, selectAction bool) error {

	// 2. 获取课程时间并且计算offset
	offset, err := local.CalOffset(msg.CourseID)
	if err != nil {
		logger.Logger.Error("计算offset", err)
		return err
	}
	// 3. 获取用户flag，更新flag
	var user models.User
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id=?", msg.UserID).First(&user).Error; err != nil {
		logger.Logger.Error("获取用户信息", err)
		return fmt.Errorf("用户ID: %d 不存在", msg.UserID)
	}
	if selectAction {
		user.Flag.SetBit(offset)
	} else {
		user.Flag.ClearBit(offset)
	}
	if err := tx.Save(&user).Error; err != nil {
		logger.Logger.Error("更新用户flag", err)
		return err
	}
	return nil
}

// 创建用户选课记录
func updateUserCourseState(tx *gorm.DB, msg *mqm.CourseReq, selectAction bool) error {

	var userCourse models.UserCourse

	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id=? and course_id=?", msg.UserID, msg.CourseID).
		First(&userCourse).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("获取用户选课记录失败", err)
			return err
		}
		// 不存在，则创建
		if err := tx.Create(&models.UserCourse{
			UserID:    msg.UserID,
			CourseID:  msg.CourseID,
			CreatedAt: msg.CreatedAt,
			UpdatedAt: msg.CreatedAt,
			IsDeleted: !selectAction,
		}).Error; err != nil {
			logger.Logger.Error("创建选课记录失败", err)
			return err
		}
		return nil
	}
	// 消息时间大于记录时间，则更新，否则就是为失效消息（过期）
	if msg.CreatedAt > userCourse.UpdatedAt {
		// 记录已存在，检查并可能更新
		updateData := map[string]interface{}{"is_deleted": !selectAction, "updated_at": msg.CreatedAt}
		if err := tx.Model(&userCourse).
			Where("user_id=? and course_id=? ", msg.UserID, msg.CourseID).
			Updates(updateData).Error; err != nil {
			logger.Logger.Error("更新选课记录失败", err)
			return err
		}
	}
	return nil
}

func (s *Select) Product(msg *mqm.CourseReq) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		logger.Logger.Error("消息序列化失败", err)
		return
	}
	if TestMsgLose {
		cache.RDB.HIncrBy(context.Background(), "record", "produce-total", 1)

	}
	s.cnt.Add(1)
	var cnt uint8
	if err := retry.Do(func() error {
		if err = s.PushMainQueue(bytes); err != nil {
			cnt++
			logger.Logger.Warningf("消息发送失败,尝试次数: %d", cnt)
			return err
		}
		if TestMsgLose {
			cache.RDB.HIncrBy(context.Background(), "record", "produce-success", 1)
		}
		// 记录消息
		s.unconfirmedMessages.Store(s.cnt.Load(), UnconfirmedMessage{
			Body: bytes,
		})
		return nil
	}, retry.Attempts(3), retry.Delay(time.Millisecond*100)); err != nil {
		logger.Logger.Error("尝试消息发送失败", err)
		if TestMsgLose {
			cache.RDB.HIncrBy(context.Background(), "record", "produce-fail", 1)
		}
		// 死信
		if err := s.PushDeadQueue(bytes); err != nil {
			logger.Logger.Error("死信消息发送失败", err)
		}
	}
}
func (s *Select) PushDeadQueue(body []byte) error {
	return s.channel.Publish(
		variable.DeadExchange,
		variable.DeadRoutingKey,
		true,
		false, amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
}
func (s *Select) PushMainQueue(body []byte) error {
	return s.channel.Publish(
		variable.SelectExchange,
		variable.SelectRoutingKey,
		true,
		false, amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
}
