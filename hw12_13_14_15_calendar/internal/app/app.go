package app

import (
	"context"

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

func (a *App) CreateEvent(ctx context.Context, title string) error {
	_ = ctx
	event := models.Event{}
	event.Title = title
	return a.storage.CreateEvent(&event)
}

// TODO: see realizations in storage_test.go
