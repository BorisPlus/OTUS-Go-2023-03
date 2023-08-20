package main

import (
	"context"
	"encoding/json"

	"hw12_13_14_15_calendar/internal/interfaces"
	"hw12_13_14_15_calendar/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventsTarget struct {
	target     models.RabbitMQTarget
	connection *amqp.Connection
	channel    *amqp.Channel
	logger     interfaces.Logger
}

func NewEventsTarget(
	target models.RabbitMQTarget,
	logger interfaces.Logger,
) *EventsTarget {
	return &EventsTarget{
		target: target,
		logger: logger,
	}
}

func (self *EventsTarget) Connect(ctx context.Context) error {
	_ = ctx // TODO
	var err error
	self.connection, err = amqp.Dial(self.target.DSN)
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	self.channel, err = self.connection.Channel()
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	return nil
}

func (self *EventsTarget) Disconnect(ctx context.Context) error {
	_ = ctx // TODO: usage
	err := self.channel.Close()
	if err != nil {
		self.logger.Error(err.Error())
		// return err
	}
	err = self.connection.Close()
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	return nil
}

func (self *EventsTarget) PutEvent(ctx context.Context, event amqp.Delivery) error {
	err := self.channel.ExchangeDeclarePassive(self.target.ExchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	data, err := json.Marshal(event)
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	message := amqp.Publishing{
		Body: data,
	}
	err = self.channel.PublishWithContext(ctx, self.target.ExchangeName, self.target.RoutingKey, false, false, message)
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	return nil
}
