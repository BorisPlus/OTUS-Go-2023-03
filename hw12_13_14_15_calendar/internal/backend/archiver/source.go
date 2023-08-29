package archiver

import (
	"context"
	// "encoding/json"

	"hw12_13_14_15_calendar/internal/interfaces"
	"hw12_13_14_15_calendar/internal/models"

	// "hw12_13_14_15_calendar/internal/server/rpc/rpcapi"

	amqp "github.com/rabbitmq/amqp091-go"
)

type NoticesSource struct {
	source     models.RabbitMQSourceWithGetback
	connection *amqp.Connection
	logger     interfaces.Logger
}

func NewEventsSource(
	source models.RabbitMQSourceWithGetback,
	logger interfaces.Logger,
) *NoticesSource {
	return &NoticesSource{
		source: source,
		logger: logger,
	}
}

func (self *NoticesSource) Connect(ctx context.Context) error {
	_ = ctx // TODO: usage
	var err error
	self.connection, err = amqp.Dial(self.source.DSN)
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	return nil
}

func (self *NoticesSource) Disconnect(ctx context.Context) error {
	_ = ctx // TODO: usage
	err := self.connection.Close()
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	return nil
}

func (self *NoticesSource) Confirm(ctx context.Context, candidate *amqp.Delivery) error {
	_ = ctx                     // TODO: usage
	return candidate.Ack(false) // TODO: usage Nack(false, true)
}

func (self *NoticesSource) Getback(ctx context.Context, candidate *amqp.Delivery) error {
	_ = ctx
	data := candidate.Body
	self.logger.Info("Getback data %s", data)
	channel, err := self.connection.Channel()
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	defer channel.Close()
	err = channel.ExchangeDeclarePassive("exch_events", "direct", true, false, false, false, nil)
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	message := amqp.Publishing{
		Body: data,
	}
	err = channel.PublishWithContext(ctx, self.source.GetbackToExchangeName, self.source.GetbackKey, false, false, message)
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	self.logger.Info("DONE. Send to getback-exchange %q with getback-routing key %q\n", self.source.GetbackToExchangeName, self.source.GetbackKey)
	return nil
}

func (self *NoticesSource) DataChannel(ctx context.Context) (<-chan amqp.Delivery, error) {
	channel, err := self.connection.Channel()
	if err != nil {
		self.logger.Error(err.Error())
		return nil, err
	}
	queue, err := channel.QueueDeclarePassive(self.source.QueueName, true, false, false, false, nil)
	if err != nil {
		self.logger.Error(err.Error())
		return nil, err
	}
	// err = channel.Qos(1, 0, false) // 1, 0, false
	// if err != nil {
	// 	self.logger.Error(err.Error())
	// 	return nil, err
	// }
	messageChannel, err := channel.Consume(
		queue.Name,
		queue.Name,
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
