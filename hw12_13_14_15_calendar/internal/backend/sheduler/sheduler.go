package sheduler

import (
	"context"

	"hw12_13_14_15_calendar/internal/backend/transmitter"
	"hw12_13_14_15_calendar/internal/interfaces"
	"hw12_13_14_15_calendar/internal/models"
	"hw12_13_14_15_calendar/internal/server/rpc/rpcapi"
	"hw12_13_14_15_calendar/internal/server/rpc/common"
)

type Sheduler struct {
	Transmitter transmitter.Transmitter[*rpcapi.Event, models.Notice]
	Logger      interfaces.Logger
}

func (self *Sheduler) Start(ctx context.Context) error {
	return self.Transmitter.Start(ctx)
}

func (self *Sheduler) Stop(ctx context.Context) error {
	return self.Transmitter.Stop(ctx)
}

func NewSheduler(
	source *EventsSource,
	target *NoticesTarget,
	logger interfaces.Logger,
	timeoutSec int64,
) *Sheduler {
	Transmitter := transmitter.NewTransmitter[*rpcapi.Event, models.Notice](
		source,
		target,
		logger,
		timeoutSec,
	)
	Transmitter.Transmit = func(ctx context.Context, candidate *rpcapi.Event) (bool, error) {
		notice := models.NewNotice(*common.PBEvent2Event(candidate))
		err := Transmitter.Target.Put(ctx, &notice)
		if err != nil {
			Transmitter.Logger.Error(err.Error())
			return false, err
		}
		candidate.Sheduled = true
		err = Transmitter.Source.Confirm(ctx, &candidate)
		if err == nil {
			return true, nil
		}
		subErr := Transmitter.Source.Getback(ctx, &candidate)
		if subErr != nil {
			Transmitter.Logger.Error(subErr.Error())
			return false, err
		}
		Transmitter.Logger.Error(err.Error())
		return false, err
	}
	return &Sheduler{
		Transmitter: *Transmitter,
		Logger:      logger,
	}
}