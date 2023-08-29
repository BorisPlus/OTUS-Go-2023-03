package transmitter

import (
	"context"
)

type Target[TO Item] interface {
	Streamer
	Put(context.Context, *TO) error
}
