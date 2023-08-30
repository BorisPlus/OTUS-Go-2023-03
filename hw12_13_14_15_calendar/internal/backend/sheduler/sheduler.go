package sheduler

import (
	"context"

	transmitter "hw12_13_14_15_calendar/internal/backend/transmitter"
	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	models "hw12_13_14_15_calendar/internal/models"
	common "hw12_13_14_15_calendar/internal/server/rpc/common"
	rpcapi "hw12_13_14_15_calendar/internal/server/rpc/rpcapi"
)

func HashFunc(event *rpcapi.Event) string {
	return string(event.PK)
}

func NewSheduler(
	source *EventsSource,
	target *NoticesTarget,
	logger interfaces.Logger,
	timeoutSec int64,
) *transmitter.Transmitter[*rpcapi.Event, models.Notice] {
	Transmitter := transmitter.NewTransmitter[*rpcapi.Event, models.Notice](
		source,
		target,
		*transmitter.NewSet[*rpcapi.Event](HashFunc),
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
	return Transmitter
}
