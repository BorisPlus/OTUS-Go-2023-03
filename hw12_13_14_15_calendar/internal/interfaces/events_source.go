package interfaces

import (
	"context"
)

type EventsSource[T Event] interface {
	Streamer
	GetEvents(context.Context) (<-chan T, error)
	AcknowledgeEvent(context.Context, T) error
}
