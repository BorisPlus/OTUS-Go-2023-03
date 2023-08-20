package main

import (
	"context"
	"encoding/json"

	"hw12_13_14_15_calendar/internal/interfaces"
	"hw12_13_14_15_calendar/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventsTarget struct {
	targets    models.RabbitMQMultyTarget
	connection *amqp.Connection
	channel    *amqp.Channel
	logger     interfaces.Logger
}

func NewEventsTarget(
	targets models.RabbitMQMultyTarget,
	logger interfaces.Logger,
) *EventsTarget {
	return &EventsTarget{
		targets: targets,
		logger:  logger,
	}
}

func (self *EventsTarget) Connect(ctx context.Context) error {
	_ = ctx // TODO
	var err error
	self.connection, err = amqp.Dial(self.targets.DSN)
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
	}
	err = self.connection.Close()
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	return nil
}

func (self *EventsTarget) PutEvent(ctx context.Context, event amqp.Delivery) error {
	err := self.channel.ExchangeDeclarePassive(self.targets.ExchangeName, "direct", true, false, false, false, nil)
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
	for _, key := range self.targets.RoutingKeys {
		err = self.channel.PublishWithContext(ctx, self.targets.ExchangeName, key, false, false, message)
		if err != nil {
			self.logger.Error(err.Error())
			return err
		}
	}
	return nil
}
