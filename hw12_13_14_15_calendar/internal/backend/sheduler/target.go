package sheduler

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	models "hw12_13_14_15_calendar/internal/models"
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

func (s *NoticesTarget) Connect(ctx context.Context) error {
	_ = ctx // TODO
	var err error
	s.connection, err = amqp.Dial(s.target.DSN)
	if err != nil {
		s.logger.Error(err.Error())
		return nil
	}
	return nil
}

func (s *NoticesTarget) Disconnect(ctx context.Context) error {
	_ = ctx // TODO: usage
	err := s.connection.Close()
	if s.connection.IsClosed() {
		return nil
	}
	if err != nil {
		// s.logger.Error(err.Error())
		return err
	}
	return nil
}

func (s *NoticesTarget) Put(ctx context.Context, notice *models.Notice) error {
	channel, err := s.connection.Channel()
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	defer channel.Close()
	err = channel.ExchangeDeclarePassive(s.target.ExchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	data, err := json.Marshal(notice)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	message := amqp.Publishing{
		Body: data,
	}
	err = channel.PublishWithContext(ctx, s.target.ExchangeName, s.target.RoutingKey, false, false, message)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	return nil
}
