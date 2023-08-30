package sender

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"hw12_13_14_15_calendar/internal/backend/transmitter"
	"hw12_13_14_15_calendar/internal/interfaces"
	"hw12_13_14_15_calendar/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

func HashFunc(candidate amqp.Delivery) string { // TODO: err
	var notice models.Notice
	json.Unmarshal(candidate.Body, &notice) // TODO: err
	return strconv.Itoa(notice.PK)
}

type Sender struct {
	Transmitter transmitter.Transmitter[amqp.Delivery, models.Notice]
	Logger      interfaces.Logger
}

func (s *Sender) Start(ctx context.Context) error {
	return s.Transmitter.Start(ctx)
}

func (s *Sender) Stop(ctx context.Context) error {
	return s.Transmitter.Stop(ctx)
}

func NewSender(
	source *NoticesSource,
	target *NoticesTarget,
	logger interfaces.Logger,
	timeoutSec int64,
) *Sender {
	Transmitter := transmitter.NewTransmitter[amqp.Delivery, models.Notice](
		source,
		target,
		transmitter.NewSet[amqp.Delivery](HashFunc),
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
	return &Sender{
		Transmitter: *Transmitter,
		Logger:      logger,
	}
}
