package interfaces

import "context"

type Transmitter interface {
	Transmit(context.Context, []byte) error
}
