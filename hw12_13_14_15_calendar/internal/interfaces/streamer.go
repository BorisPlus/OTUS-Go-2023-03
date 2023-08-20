package interfaces

import (
	"context"
)

type Streamer interface {
	Connect(context.Context) error
	Disconnect(context.Context) error
}
