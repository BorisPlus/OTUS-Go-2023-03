package transmitter

import (
	"context"
)

type Source[FROM Item] interface {
	Streamer
	DataChannel(context.Context) (<-chan FROM, error)
	Confirm(context.Context, *FROM) error
	Getback(context.Context, *FROM) error
}
