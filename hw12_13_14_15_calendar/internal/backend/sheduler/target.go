package sheduler

import (
	"context"
	"encoding/json"

	"hw12_13_14_15_calendar/internal/interfaces"
	"hw12_13_14_15_calendar/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

type NoticesTarget struct {
	target     models.RabbitMQTarget
	connection *amqp.Connection
	logger     interfaces.Logger
}

func NewNoticesTarget(
	target models.RabbitMQTarget,
	logger interfaces.Logger,
) *NoticesTarget {
	return &NoticesTarget{
		target: target,
		logger: logger,
	}
}

func (self *NoticesTarget) Connect(ctx context.Context) error {
	_ = ctx // TODO
	var err error
	self.connection, err = amqp.Dial(self.target.DSN)
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	return nil
}

func (self *NoticesTarget) Disconnect(ctx context.Context) error {
	_ = ctx // TODO: usage
	err := self.connection.Close()
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	return nil
}

func (self *NoticesTarget) Put(ctx context.Context, notice *models.Notice) error {
	channel, err := self.connection.Channel()
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	defer channel.Close()
	err = channel.ExchangeDeclarePassive(self.target.ExchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	data, err := json.Marshal(notice)
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	message := amqp.Publishing{
		Body: data,
	}
	err = channel.PublishWithContext(ctx, self.target.ExchangeName, self.target.RoutingKey, false, false, message)
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	return nil
}
