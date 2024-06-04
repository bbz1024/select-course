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
	err := s.channel.ExchangeDeclare(variable.SelectExchange, variable.SelectKind,
		true, false, false, false, nil,
	)
	if err != nil {
		return err
	}
	_, err = s.channel.QueueDeclare(variable.SelectQueue, true,
		false, false, false, nil,
	)
	if err != nil {
		return err
	}
	// 将队列绑定到交换机上
	err = s.channel.QueueBind(variable.SelectQueue, variable.SelectRoutingKey,
		variable.SelectExchange, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Select) Consumer() {
	//接收消息
	results, err := SelectConsumer.channel.Consume(
		variable.SelectQueue,
		variable.SelectRoutingKey,
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
				// 2.4 扣减课程库存
				if err := tx.Model(&models.Course{}).
					Where("id=?", msg.CourseID).
					Update("capacity", gorm.Expr("capacity - 1")).Error; err != nil {
					logger.Logger.Info("更新课程容量失败", err)
					return err
				}
				// 2.5 创建选课记录
				if err := tx.Create(&models.UserCourse{
					UserID:   msg.UserID,
					CourseID: msg.CourseID,
				}).Error; err != nil {
					logger.Logger.Info("创建选课记录失败", err)
					return err
				}
				return nil // 成功，无错误返回
			}); err != nil {
				err := res.Nack(false, true)
				if err != nil {
					logger.Logger.Error("消息确认失败", err)
					return
				}
				logger.Logger.Info("事务回滚", err)
				return
			}
			err = res.Ack(false)
			if err != nil {
				logger.Logger.Error("消息确认失败", err)
				return
			}
		}
	}()
}
func (s *Select) Product(msg *mqm.CourseReq) {
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
