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

type Transmitter[FROM Item, TO Item] struct {
	Source         Source[FROM]
	Target         Target[TO]
	Transmit       func(ctx context.Context, candidate FROM) (bool, error)
	Logger         interfaces.Logger
	LoopTimeoutSec int64
}

func NewTransmitter[FROM Item, TO Item](
	Source Source[FROM],
	Target Target[TO],
	Logger interfaces.Logger,
	LoopTimeoutSec int64,
) *Transmitter[FROM, TO] {
	return &Transmitter[FROM, TO]{
		Source:         Source,
		Target:         Target,
		Logger:         Logger,
		LoopTimeoutSec: LoopTimeoutSec,
	}
}

func (t *Transmitter[FROM, TO]) Stop(ctx context.Context) error {
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

func (t *Transmitter[FROM, TO]) Start(ctx context.Context) error {
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
	for {
		t.Logger.Info("Retry")
		eventsChan, err := t.Source.DataChannel(ctx)
		if err != nil {
			t.Logger.Error(err.Error())
			return err
		}
		breakMe := false
		for !breakMe {
			select {
			case candidate, ok := <-eventsChan:
				if !ok {
					breakMe = true
					t.Logger.Warning("No candidates")
					break
				}
				t.Logger.Info("Transmit candidate begin")
				indicator, err := t.Transmit(ctx, candidate)
				if err != nil {
					t.Logger.Error(err.Error())
					return err
				}
				t.Logger.Info("Transmit candidate done with status %v", indicator)
			case <-ctx.Done():
				return nil
			}
		}
		SleepWithContext(ctx, time.Duration(t.LoopTimeoutSec)*time.Second)
	}
}
