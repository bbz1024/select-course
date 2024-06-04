package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
	"select-course/demo3/src/constant/variable"
	"select-course/demo3/src/models"
	"select-course/demo3/src/models/mqm"
	"select-course/demo3/src/storage/database"
	"select-course/demo3/src/storage/mq"
	"select-course/demo3/src/utils/logger"
)

var BackConsumer *Back

type Back struct {
	channel *amqp.Channel
}

func InitBackListener() error {

	BackConsumer = &Back{
		channel: mq.Client,
	}
	err := BackConsumer.Declare()
	if err != nil {
		fmt.Println(err)
		return err
	}
	BackConsumer.Consumer()
	return nil

}

func (s *Back) Declare() error {
	err := s.channel.ExchangeDeclare(variable.BackExchange, variable.BackKind,
		true, false, false, false, nil,
	)
	if err != nil {
		return err
	}
	_, err = s.channel.QueueDeclare(variable.BackQueue, true,
		false, false, false, nil,
	)
	if err != nil {
		return err
	}
	// 将队列绑定到交换机上
	err = s.channel.QueueBind(variable.BackQueue, variable.BackRoutingKey,
		variable.BackExchange, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Back) Consumer() {
	//接收消息
	results, err := BackConsumer.channel.Consume(
		variable.BackQueue,
		variable.BackRoutingKey,
		false, // 关闭自动应答
		false, //
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Logger.Error("消息接收失败", err)
	}
	//启用后台协程处理消息
	go func() {
		for res := range results {
			var msg *mqm.CourseReq
			err := json.Unmarshal(res.Body, &msg)
			if err != nil {
				logger.Logger.Error("消息反序列化失败", err)
				continue
			}

			// 扣减库存操作
			if err := database.Client.Transaction(func(tx *gorm.DB) error {
				if err := tx.Model(&models.Course{}).
					Where("id=?", msg.CourseID).
					Update("capacity", gorm.Expr("capacity + 1")).Error; err != nil {
					logger.Logger.Info("更新课程容量失败", err)
					return err
				}
				if err := tx.Where("user_id=? and course_id=?", msg.UserID, msg.CourseID).Delete(&models.UserCourse{}).Error; err != nil {
					logger.Logger.Info("删除选课记录失败", err)
					return err
				}

				return nil
			}); err != nil {
				if err := res.Nack(false, true); err != nil {
					logger.Logger.Error("消息确认失败", err)
					return
				}
				logger.Logger.Info("事务回滚", err)
				return
			}
			if err := res.Ack(false); err != nil {
				logger.Logger.Error("消息确认失败", err)
				return
			}
		}
		logger.Logger.Info("消息接收协程退出")
	}()
}
func (s *Back) Product(msg *mqm.CourseReq) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		logger.Logger.Error("消息序列化失败", err)
		return
	}
	err = s.channel.Publish(
		variable.BackExchange,
		variable.BackRoutingKey,
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
