package main

import (
	"context"

	"hw12_13_14_15_calendar/internal/interfaces"
	"hw12_13_14_15_calendar/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventsSource struct {
	source     models.RabbitMQSource
	connection *amqp.Connection
	channel    *amqp.Channel
	logger     interfaces.Logger
}

func NewEventsSource(
	source models.RabbitMQSource,
	logger interfaces.Logger,
) *EventsSource {
	return &EventsSource{
		source: source,
		logger: logger,
	}
}

func (self *EventsSource) Disconnect(ctx context.Context) error {
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

func (self *EventsSource) Connect(ctx context.Context) error {
	_ = ctx // TODO: usage
	var err error
	self.connection, err = amqp.Dial(self.source.DSN)
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

func (self *EventsSource) AcknowledgeEvent(ctx context.Context, candidate amqp.Delivery) error {
	_ = ctx // TODO: usage
	return candidate.Ack(false)
}

func (self *EventsSource) GetEvents(ctx context.Context) (<-chan amqp.Delivery, error) {
	queue, err := self.channel.QueueDeclarePassive(self.source.QueueName, true, false, false, false, nil)
	if err != nil {
		self.logger.Error(err.Error())
		return nil, err
	}
	err = self.channel.Qos(1, 0, false)
	if err != nil {
		self.logger.Error(err.Error())
		return nil, err
	}
	messageChannel, err := self.channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		self.logger.Error(err.Error())
		return nil, err
	}
	return messageChannel, nil
}
