package models

import (
	"time"
)

type Event struct {
	PK          int
	Title       string
	StartAt     time.Time
	Duration    int // TODO: time.Duration
	Description string
	Owner       string
	NotifyEarly int // TODO: time.Duration
}
