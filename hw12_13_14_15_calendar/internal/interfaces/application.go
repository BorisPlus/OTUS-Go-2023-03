package interfaces

import (
	models "hw12_13_14_15_calendar/internal/models"
)

type Applicationer interface {
	CreateEvent(event *models.Event) (*models.Event, error)
	ReadEvent(pk int) (*models.Event, error)
	UpdateEvent(event *models.Event) (*models.Event, error)
	DeleteEvent(event *models.Event) (*models.Event, error)
	ListEvents() ([]models.Event, error)
	ListNotSheduledEvents() ([]models.Event, error)
}
