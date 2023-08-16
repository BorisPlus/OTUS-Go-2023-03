package rmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	models "hw12_13_14_15_calendar/internal/models"
)

type Notifier struct {
	notifyTo models.RabbitMQTarget
	logger   interfaces.Logger
}

func NewNotifier(
	notifyTo models.RabbitMQTarget,
	logger interfaces.Logger,
) *Notifier {
	return &Notifier{
		notifyTo: notifyTo,
		logger:   logger,
	}
}

func (n *Notifier) Transmit(ctx context.Context, data []byte) error {
	var err error
	conn, err := amqp.Dial(n.notifyTo.DSN)
	if err != nil {
		n.logger.Error(err.Error())
		return err
	}
	defer conn.Close()
	channel, err := conn.Channel()
	if err != nil {
		n.logger.Error(err.Error())
		return err
	}
	defer channel.Close()
	err = channel.ExchangeDeclarePassive(n.notifyTo.ExchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		n.logger.Error(err.Error())
		return err
	}
	message := amqp.Publishing{
		Body: data,
	}
	err = channel.PublishWithContext(ctx, n.notifyTo.ExchangeName, n.notifyTo.RoutingKey, false, false, message)
	if err != nil {
		n.logger.Error(err.Error())
		return err
	}
	return nil
}
