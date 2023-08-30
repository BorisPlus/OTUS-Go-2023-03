package archiver

import (
	"context"
	"encoding/json"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	models "hw12_13_14_15_calendar/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

type NoticesTarget struct {
	target     models.RabbitMQTarget
	connection *amqp.Connection
	logger     interfaces.Logger
}

func NewEventsTarget(
	target models.RabbitMQTarget,
	logger interfaces.Logger,
) *NoticesTarget {
	return &NoticesTarget{
		target: target,
		logger: logger,
	}
}

func (t *NoticesTarget) Connect(ctx context.Context) error {
	_ = ctx // TODO: usage
	var err error
	t.connection, err = amqp.Dial(t.target.DSN)
	if err != nil {
		t.logger.Error(err.Error())
		return err
	}
	return nil
}

func (t *NoticesTarget) Disconnect(ctx context.Context) error {
	_ = ctx // TODO: usage
	if t.connection.IsClosed() {
		return nil
	}
	err := t.connection.Close()
	if err != nil {
		t.logger.Error(err.Error())
		return err
	}
	return nil
}

func (t *NoticesTarget) Put(ctx context.Context, notice *models.Notice) error {
	channel, err := t.connection.Channel()
	if err != nil {
		t.logger.Error(err.Error())
		return err
	}
	err = channel.ExchangeDeclarePassive(t.target.ExchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		t.logger.Error(err.Error())
		return err
	}
	data, err := json.Marshal(notice)
	if err != nil {
		t.logger.Error(err.Error())
		return err
	}
	message := amqp.Publishing{
		Body: data,
	}
	err = channel.PublishWithContext(ctx, t.target.ExchangeName, t.target.RoutingKey, false, false, message)
	if err != nil {
		t.logger.Error(err.Error())
		return err
	}
	t.logger.Debug("DONE. Send to exchange %q with routing key %q\n",
		t.target.ExchangeName, t.target.RoutingKey)
	return nil
}
