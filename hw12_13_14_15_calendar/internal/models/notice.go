package models

import (
	"time"
)

type Notice struct {
	PK          int
	Title       string    `json:"title"`
	StartAt     time.Time `json:"startat"`
	Owner       string    `json:"owner"`
	NotifyEarly int       `json:"notifyearly"`
}

func NewNotice(e Event) Notice {
	return Notice{
		PK:          e.PK,
		Title:       e.Title,
		StartAt:     e.StartAt,
		Owner:       e.Owner,
		NotifyEarly: e.NotifyEarly,
	}
}
