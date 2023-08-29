package sender

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"hw12_13_14_15_calendar/internal/interfaces"
	"hw12_13_14_15_calendar/internal/models"
)

type NoticesSource struct {
	source     models.RabbitMQSourceWithGetback
	connection *amqp.Connection
	logger     interfaces.Logger
}

func NewNoticesSource(
	source models.RabbitMQSourceWithGetback,
	logger interfaces.Logger,
) *NoticesSource {
	return &NoticesSource{
		source: source,
		logger: logger,
	}
}

func (s *NoticesSource) Connect(ctx context.Context) error {
	_ = ctx // TODO: usage
	var err error
	s.connection, err = amqp.Dial(s.source.DSN)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	return nil
}

func (s *NoticesSource) Disconnect(ctx context.Context) error {
	_ = ctx // TODO: usage
	if s.connection.IsClosed() {
		return nil
	}
	err := s.connection.Close()
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	return nil
}

func (s *NoticesSource) Confirm(ctx context.Context, candidate *amqp.Delivery) error {
	_ = ctx // TODO: usage
	return candidate.Ack(false)
}

func (s *NoticesSource) DataChannel(ctx context.Context) (<-chan amqp.Delivery, error) {
	_ = ctx // TODO: usage
	channel, err := s.connection.Channel()
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	queue, err := channel.QueueDeclarePassive(s.source.QueueName, true, false, false, false, nil)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	messageChannel, err := channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	return messageChannel, nil
}

func (s *NoticesSource) Getback(ctx context.Context, candidate *amqp.Delivery) error {
	_ = ctx
	data := candidate.Body
	s.logger.Info("Getback data %s", data)
	channel, err := s.connection.Channel()
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	defer channel.Close()
	err = channel.ExchangeDeclarePassive("exch_events", "direct", true, false, false, false, nil)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	message := amqp.Publishing{
		Body: data,
	}
	err = channel.PublishWithContext(ctx, s.source.GetbackToExchangeName, s.source.GetbackKey, false, false, message)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	s.logger.Info("DONE. Send to getback-exchange %q with getback-routing key %q",
		s.source.GetbackToExchangeName, s.source.GetbackKey)
	return nil
}
