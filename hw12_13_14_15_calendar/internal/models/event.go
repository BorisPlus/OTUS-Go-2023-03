package models

import (
	"time"
)

type Event struct {
	PK          int
	Title       string    `json:"title"`
	StartAt     time.Time `json:"startat"`
	Duration    int       `json:"duration"` // TODO: time.Duration
	Description string    `json:"description"`
	Owner       string    `json:"owner"`
	NotifyEarly int       `json:"notifyearly"` // TODO: time.Duration
	Sheduled    bool      `json:"sheduled"`
}
