package app

import (
	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	models "hw12_13_14_15_calendar/internal/models"
)

type App struct {
	logger  interfaces.Logger
	storage interfaces.Storager
}

func NewApp(logger interfaces.Logger, storage interfaces.Storager) *App {
	return &App{logger, storage}
}

func (a *App) CreateEvent(event *models.Event) (*models.Event, error) {
	err := a.storage.Connect()
	if err != nil {
		return nil, err
	}
	defer a.storage.Close()
	// TODO: in args `ctx context.Context`
	return a.storage.CreateEvent(event)
}

func (a *App) ReadEvent(pk int) (*models.Event, error) {
	err := a.storage.Connect()
	if err != nil {
		return nil, err
	}
	defer a.storage.Close()
	return a.storage.ReadEvent(pk)
}

func (a *App) UpdateEvent(e *models.Event) (*models.Event, error) {
	err := a.storage.Connect()
	if err != nil {
		return nil, err
	}
	defer a.storage.Close()
	return a.storage.UpdateEvent(e)
}

func (a *App) DeleteEvent(e *models.Event) (*models.Event, error) {
	err := a.storage.Connect()
	if err != nil {
		return nil, err
	}
	defer a.storage.Close()
	return a.storage.DeleteEvent(e)
}

func (a *App) ListEvents() ([]models.Event, error) {
	err := a.storage.Connect()
	if err != nil {
		return nil, err
	}
	defer a.storage.Close()
	return a.storage.ListEvents()
}

func (a *App) ListNotSheduledEvents() ([]models.Event, error) {
	err := a.storage.Connect()
	if err != nil {
		return nil, err
	}
	defer a.storage.Close()
	return a.storage.ListNotSheduledEvents()
}
