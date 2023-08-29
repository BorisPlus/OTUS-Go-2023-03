package archiver

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"hw12_13_14_15_calendar/internal/backend/transmitter"
	"hw12_13_14_15_calendar/internal/interfaces"
	"hw12_13_14_15_calendar/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Archiver struct {
	Transmitter transmitter.Transmitter[amqp.Delivery, models.Notice]
	Logger      interfaces.Logger
}

func (self *Archiver) Start(ctx context.Context) error {
	return self.Transmitter.Start(ctx)
}

func (self *Archiver) Stop(ctx context.Context) error {
	return self.Transmitter.Stop(ctx)
}

func NewArchiver(
	source *NoticesSource,
	target *NoticesTarget,
	logger interfaces.Logger,
	timeoutSec int64,
) *Archiver {
	Transmitter := transmitter.NewTransmitter[amqp.Delivery, models.Notice](
		source,
		target,
		logger,
		timeoutSec,
	)
	Transmitter.Transmit = func(ctx context.Context, candidate amqp.Delivery) (bool, error) {
		data := candidate.Body
		var notice models.Notice
		json.Unmarshal(data, &notice)
		now := time.Now()
		if notice.StartAt.Before(now) {
			Transmitter.Logger.Info("Must be archived: %d", notice.PK)
			fmt.Println()
			fmt.Println(now)
			fmt.Println(notice.StartAt)
			if err := Transmitter.Source.Confirm(ctx, &candidate); err != nil {
				return false, err
			}
			if err := Transmitter.Target.Put(ctx, &notice); err != nil {
				Transmitter.Source.Getback(ctx, &candidate)
				return false, err
			}
			return true, nil
		}
		Transmitter.Logger.Info("Must be getback PK:%d\n", notice.PK)
		// TODO: how to make it in one transaction
		// BEGIN TRANSACTION
		if err := Transmitter.Source.Getback(ctx, &candidate); err != nil {
			return false, err
		}
		if err := Transmitter.Source.Confirm(ctx, &candidate); err != nil {
			return false, err
		}
		// END TRANSACTION
		return false, nil
	}
	return &Archiver{
		Transmitter: *Transmitter,
		Logger:      logger,
	}
}
