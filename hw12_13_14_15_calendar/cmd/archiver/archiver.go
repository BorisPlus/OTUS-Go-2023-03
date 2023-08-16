package main

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	models "hw12_13_14_15_calendar/internal/models"
)

type Archiver struct {
	source                   models.RabbitMQSource
	sourceRabbitMQConnection *amqp.Connection
	archiver                 interfaces.Transmitter
	logger                   interfaces.Logger
}

func NewArchiver(
	source models.RabbitMQSource,
	archiver interfaces.Transmitter,
	logger interfaces.Logger,
) *Archiver {
	return &Archiver{
		source:   source,
		archiver: archiver,
		logger:   logger,
	}
}

func (a *Archiver) Stop() error {
	errSource := a.sourceRabbitMQConnection.Close()
	if errSource != nil {
		a.logger.Error(errSource.Error())
	}
	a.logger.Info("Archiver.Stop()")
	return errSource
}

func (a *Archiver) Start(ctx context.Context) error {
	a.logger.Info("Archiver.Start()")
	var err error
	a.sourceRabbitMQConnection, err = amqp.Dial(a.source.DSN)
	if err != nil {
		a.logger.Error(err.Error())
		return err
	}
	defer func() {
		a.Stop()
	}()
	sourceChannel, err := a.sourceRabbitMQConnection.Channel()
	if err != nil {
		a.logger.Error(err.Error())
		return err
	}
	defer func() {
		a.logger.Info("Channel close.")
		_ = sourceChannel.Close()
	}()
	queue, err := sourceChannel.QueueDeclarePassive(a.source.QueueName, true, false, false, false, nil)
	if err != nil {
		a.logger.Error(err.Error())
		return err
	}
	err = sourceChannel.Qos(1, 0, false)
	if err != nil {
		a.logger.Error(err.Error())
		return err
	}
	messageChannel, err := sourceChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		a.logger.Error(err.Error())
		return err
	}
	var event models.Event
	for {
		select {
		case d := <-messageChannel:
			data := d.Body
			json.Unmarshal(data, &event)
			a.logger.Info("Received a message: %s", data)
			if event.StartAt.Before(time.Now()) {
				if err := d.Ack(false); err != nil {
					a.logger.Info("Error acknowledging message : %s", err)
					return err
				} else {
					a.logger.Info("Archiver.Transmit()")
					err = a.archiver.Transmit(ctx, data)
					if err != nil {
						return err
					}
				}
			} else {
				d.Reject(true)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
