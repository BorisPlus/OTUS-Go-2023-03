package transmitter

import (
	"context"
	"time"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
)

func SleepWithContext(ctx context.Context, d time.Duration) {
	timer := time.NewTimer(d)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
	case <-timer.C:
	}
}

type Transmitter[T interfaces.Event] struct {
	Source         interfaces.EventsSource[T]
	Target         interfaces.EventsTarget[T]
	Transmit       func(ctx context.Context, candidate T) error
	Logger         interfaces.Logger
	LoopTimeoutSec int64
}

func NewTransmitter[T interfaces.Event](
	Source interfaces.EventsSource[T],
	Target interfaces.EventsTarget[T],
	Logger interfaces.Logger,
	LoopTimeoutSec int64,
) *Transmitter[T] {
	return &Transmitter[T]{
		Source:         Source,
		Target:         Target,
		Logger:         Logger,
		LoopTimeoutSec: LoopTimeoutSec,
	}
}

func (t *Transmitter[T]) Stop(ctx context.Context) error {
	t.Logger.Info("Transmitter.Stop()")
	err := t.Source.Disconnect(ctx)
	if err != nil {
		t.Logger.Error(err.Error())
		return err
	}
	err = t.Target.Disconnect(ctx)
	if err != nil {
		t.Logger.Error(err.Error())
		return err
	}
	return nil
}

func (t *Transmitter[T]) Start(ctx context.Context) error {
	t.Logger.Info("Transmitter.Start()")
	err := t.Source.Connect(ctx)
	if err != nil {
		t.Logger.Error(err.Error())
		return err
	}
	err = t.Target.Connect(ctx)
	if err != nil {
		t.Logger.Error(err.Error())
		return err
	}
	eventsChan := make(<-chan T)
	for {
		select {
		case candidate := <-eventsChan:
			err := t.Transmit(ctx, candidate)
			if err != nil {
				t.Logger.Error(err.Error())
				return err
			}
		case <-ctx.Done():
			return nil
		default:
			eventsChan, err = t.Source.GetEvents(ctx)
			if err != nil {
				t.Logger.Error(err.Error())
				return err
			}
			if len(eventsChan) == 0 {
				// time.Sleep(time.Duration(t.LoopTimeoutSec) * time.Second)
				SleepWithContext(ctx, time.Duration(t.LoopTimeoutSec)*time.Second)
			}
			continue
		}
	}
}
