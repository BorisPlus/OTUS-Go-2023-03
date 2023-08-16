package main

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	models "hw12_13_14_15_calendar/internal/models"
)

type Sender struct {
	source                   models.RabbitMQSource
	sourceRabbitMQConnection *amqp.Connection
	sender                   interfaces.Transmitter
	archiver                 interfaces.Transmitter
	logger                   interfaces.Logger
}

func NewSender(
	source models.RabbitMQSource,
	sender interfaces.Transmitter,
	archiver interfaces.Transmitter,
	logger interfaces.Logger,
) *Sender {
	return &Sender{
		source:   source,
		sender:   sender,
		archiver: archiver,
		logger:   logger,
	}
}

func (s *Sender) Stop() error {
	errSource := s.sourceRabbitMQConnection.Close()
	if errSource != nil {
		s.logger.Error(errSource.Error())
	}
	s.logger.Info("Sender.Stop()")
	return errSource
}

func (s *Sender) Start(ctx context.Context) error {
	s.logger.Info("Sender.Start()")
	var err error
	s.sourceRabbitMQConnection, err = amqp.Dial(s.source.DSN)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	defer func() {
		s.Stop()
	}()
	sourceChannel, err := s.sourceRabbitMQConnection.Channel()
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	defer func() {
		s.logger.Info("Channel close.")
		_ = sourceChannel.Close()
	}()
	queue, err := sourceChannel.QueueDeclarePassive(s.source.QueueName, true, false, false, false, nil)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	err = sourceChannel.Qos(1, 0, false)
	if err != nil {
		s.logger.Error(err.Error())
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
		s.logger.Error(err.Error())
		return err
	}

	var event models.Event
	for {
		select {
		case d := <-messageChannel:
			data := d.Body
			json.Unmarshal(data, &event)
			s.logger.Info("Received a message: %s", data)
			// now := time.Now()
			_ = time.Now()
			if true {
				// if event.StartAt.Add(-time.Duration(event.Duration)*time.Second).Before(now) && now.Before(event.StartAt) {
				if err := d.Ack(false); err != nil {
					s.logger.Info("Error acknowledging message : %s", err)
					return err
				} else {
					err = s.sender.Transmit(ctx, data)
					if err != nil {
						return err
					}
					err = s.archiver.Transmit(ctx, data)
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
