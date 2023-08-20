package interfaces

import (
	"context"
)

type EventsTarget[T Event] interface {
	Streamer
	PutEvent(context.Context, T) error
}
