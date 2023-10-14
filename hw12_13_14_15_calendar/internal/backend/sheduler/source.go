package sheduler

import (
	"context"

	//unlint:gci
	"hw12_13_14_15_calendar/internal/interfaces"
	"hw12_13_14_15_calendar/internal/server/rpc/client"
	rpcapi "hw12_13_14_15_calendar/internal/server/rpc/rpcapi"
)

type EventsSource struct {
	dsn       string
	rpcClient client.Client
	logger    interfaces.Logger
}

func NewEventsSource(
	dsn string,
	logger interfaces.Logger,
) *EventsSource {
	return &EventsSource{
		dsn:    dsn,
		logger: logger,
	}
}

func (s *EventsSource) Connect(ctx context.Context) error {
	_ = ctx // TODO: usage
	err := s.rpcClient.Connect(s.dsn)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	return nil
}

func (s *EventsSource) Disconnect(ctx context.Context) error {
	_ = ctx // TODO: usage
	err := s.rpcClient.Close()
	if err != nil {
		// s.logger.Error(err.Error())
		return err
	}
	return nil
}

func (s *EventsSource) Confirm(ctx context.Context, event **rpcapi.Event) error {
	(*event).Sheduled = true
	_, err := s.rpcClient.UpdateEvent(ctx, *(event))
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	return nil
}

func (s *EventsSource) Getback(ctx context.Context, event **rpcapi.Event) error {
	(*event).Sheduled = false
	_, err := s.rpcClient.UpdateEvent(ctx, *(event))
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	return nil
}

func (s *EventsSource) DataChannel(ctx context.Context) (<-chan *rpcapi.Event, error) {
	events, err := s.rpcClient.ListNotSheduledEvents(ctx)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	eventsChan := make(chan *rpcapi.Event, len(events))
	for _, event := range events {
		eventsChan <- event
	}
	close(eventsChan)
	return eventsChan, nil
}
