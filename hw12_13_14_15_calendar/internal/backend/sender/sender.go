package sender

import (
	"context"
	"encoding/json"
	"time"

	"hw12_13_14_15_calendar/internal/backend/archiver"
	"hw12_13_14_15_calendar/internal/backend/transmitter"
	"hw12_13_14_15_calendar/internal/interfaces"
	"hw12_13_14_15_calendar/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewSender(
	source *archiver.NoticesSource,
	target *NoticesTarget,
	logger interfaces.Logger,
	timeoutSec int64,
) *transmitter.Transmitter[amqp.Delivery, models.Notice] {
	Transmitter := transmitter.NewTransmitter[amqp.Delivery, models.Notice](
		source,
		target,
		*transmitter.NewSet[amqp.Delivery](archiver.HashFunc),
		logger,
		timeoutSec,
	)
	Transmitter.Transmit = func(ctx context.Context, candidate amqp.Delivery) (bool, error) {
		data := candidate.Body
		var notice models.Notice
		json.Unmarshal(data, &notice)
		Transmitter.Logger.Info("Received a message: %s", data)
		now := time.Now()
		if now.After(notice.StartAt.Add(-time.Second*time.Duration(notice.NotifyEarly))) && now.Before(notice.StartAt) {
			Transmitter.Logger.Info("Must be send: %d", notice.PK)
			time.Sleep(2 * time.Second)
			if err := Transmitter.Source.Confirm(ctx, &candidate); err != nil {
				return false, err
			}
			if err := Transmitter.Target.Put(ctx, &notice); err != nil {
				Transmitter.Source.Getback(ctx, &candidate)
				return false, err
			}
			return true, nil
		}
		Transmitter.Logger.Info("Must be getback PK:%d", notice.PK)
		if err := Transmitter.Source.Getback(ctx, &candidate); err != nil {
			return false, err
		}
		if err := Transmitter.Source.Confirm(ctx, &candidate); err != nil {
			return false, err
		}
		return false, nil
	}
	return Transmitter
}
