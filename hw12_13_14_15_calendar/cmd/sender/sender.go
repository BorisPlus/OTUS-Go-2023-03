package main

import (
	"context"
	"encoding/json"
	"time"

	"hw12_13_14_15_calendar/internal/interfaces"
	"hw12_13_14_15_calendar/internal/server/rpc/rpcapi"
	"hw12_13_14_15_calendar/internal/transmitter"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Sender struct {
	Transmitter transmitter.Transmitter[amqp.Delivery]
	Logger      interfaces.Logger
}

func (self *Sender) Start(ctx context.Context) error {
	return self.Transmitter.Start(ctx)
}

func (self *Sender) Stop(ctx context.Context) error {
	return self.Transmitter.Stop(ctx)
}

func NewSender(
	source *EventsSource,
	target *EventsTarget,
	logger interfaces.Logger,
	timeoutSec int64,
) *Sender {
	Transmitter := transmitter.NewTransmitter[amqp.Delivery](
		source,
		target,
		logger,
		timeoutSec,
	)
	Transmitter.Transmit = func(ctx context.Context, candidate amqp.Delivery) error {
		data := candidate.Body
		var event rpcapi.Event
		json.Unmarshal(data, &event)
		Transmitter.Logger.Info("Received a message: %s", data)
		now := time.Now()
		if event.StartAt.AsTime().Add(-time.Second*time.Duration(event.Duration)).After(now) && event.StartAt.AsTime().Before(now) {
			if err := Transmitter.Source.AcknowledgeEvent(ctx, candidate); err != nil {
				Transmitter.Logger.Info("Error acknowledging message : %s", err)
				return err
			} else {
				Transmitter.Logger.Info("Archiver.Transmit()")
				if err != nil {
					return err
				}
			}
		} else {
			candidate.Reject(true)
		}
		return nil
	}
	return &Sender{
		Transmitter: *Transmitter,
		Logger:      logger,
	}
}
