package main

import (
	"context"

	"hw12_13_14_15_calendar/internal/interfaces"
	"hw12_13_14_15_calendar/internal/server/rpc/rpcapi"
	"hw12_13_14_15_calendar/internal/transmitter"
)

type Sheduler struct {
	Transmitter transmitter.Transmitter[*rpcapi.Event]
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
	target *EventsTarget,
	logger interfaces.Logger,
	timeoutSec int64,
) *Sheduler {
	Transmitter := transmitter.NewTransmitter[*rpcapi.Event](
		source,
		target,
		logger,
		timeoutSec,
	)
	Transmitter.Transmit = func(ctx context.Context, candidate *rpcapi.Event) error {
		err := Transmitter.Target.PutEvent(ctx, candidate)
		if err != nil {
			Transmitter.Logger.Error(err.Error())
			return err
		}
		candidate.Sheduled = true
		err = Transmitter.Source.AcknowledgeEvent(ctx, candidate)
		if err != nil {
			Transmitter.Logger.Error(err.Error())
			return err
		}
		return nil
	}
	return &Sheduler{
		Transmitter: *Transmitter,
		Logger:      logger,
	}
}
