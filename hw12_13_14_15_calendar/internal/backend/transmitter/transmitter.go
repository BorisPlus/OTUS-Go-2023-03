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
		return
	case <-timer.C:
		return
	}
}

type Transmitter[FROM Item, TO Item] struct {
	Source         Source[FROM]
	Target         Target[TO]
	Processed      Set[FROM]
	Transmit       func(ctx context.Context, candidate FROM) (bool, error)
	Logger         interfaces.Logger
	LoopTimeoutSec int64
}

func NewTransmitter[FROM Item, TO Item](
	source Source[FROM],
	target Target[TO],
	processed Set[FROM],
	levelLogger interfaces.Logger,
	loopTimeoutSec int64,
) *Transmitter[FROM, TO] {
	return &Transmitter[FROM, TO]{
		Source:         source,
		Target:         target,
		Processed:      processed,
		Logger:         levelLogger,
		LoopTimeoutSec: loopTimeoutSec,
	}
}

func (t *Transmitter[FROM, TO]) connect(ctx context.Context) error {
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
	return nil
}

func (t *Transmitter[FROM, TO]) disconnect(ctx context.Context) error {
	err := t.Source.Disconnect(ctx)
	if err != nil {
		return err
	}
	err = t.Target.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (t *Transmitter[FROM, TO]) Stop(ctx context.Context) error {
	t.Logger.Info("Transmitter.Stop()")
	return t.disconnect(ctx)
}

func (t *Transmitter[FROM, TO]) Start(ctx context.Context) error {
	t.Logger.Info("Transmitter.Start()")
	for {
		err := t.connect(ctx)
		if err != nil {
			t.Logger.Error(err.Error())
			return err
		}
		t.Logger.Info("Transmitter step.")
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
					t.Logger.Debug("No candidates")
					break
				}
				t.Logger.Debug("Transmit candidate begin")
				indicator, err := t.Transmit(ctx, candidate)
				if err != nil {
					t.Logger.Error(err.Error())
					return err
				}
				t.Logger.Debug("Transmit candidate done with status %v", indicator)
				if !indicator {
					if t.Processed.has(candidate) {
						t.Processed.clear()
						breakMe = true
						break
					}
					t.Logger.Info("new candidate.")
					t.Processed.add(candidate)
				}
			case <-ctx.Done():
				return nil
			}
		}
		SleepWithContext(ctx, time.Duration(t.LoopTimeoutSec)*time.Second)
		err = t.disconnect(ctx)
		if err != nil {
			return err
		}
	}
}
