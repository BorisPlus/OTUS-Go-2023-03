package interfaces

import (
	models "hw12_13_14_15_calendar/internal/models"
)

type Storager interface {
	Connect() error
	Close() error
	CreateEvent(*models.Event) error
	ReadEvent(int) (*models.Event, error)
	UpdateEvent(*models.Event) error
	DeleteEvent(*models.Event) error
	ListEvents() ([]models.Event, error)
}
