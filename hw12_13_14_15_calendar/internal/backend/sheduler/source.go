package sheduler

import (
	"context"

	"hw12_13_14_15_calendar/internal/interfaces"
	"hw12_13_14_15_calendar/internal/server/rpc/client"
	"hw12_13_14_15_calendar/internal/server/rpc/rpcapi"
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

func (self *EventsSource) Connect(ctx context.Context) error {
	_ = ctx // TODO: usage
	err := self.rpcClient.Connect(self.dsn)
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	return nil
}

func (self *EventsSource) Disconnect(ctx context.Context) error {
	_ = ctx // TODO: usage
	err := self.rpcClient.Close()
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	return nil
}

func (self *EventsSource) Confirm(ctx context.Context, event **rpcapi.Event) error {
	(*event).Sheduled = true
	_, err := self.rpcClient.UpdateEvent(ctx, *(event))
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	return nil
}

func (self *EventsSource) Getback(ctx context.Context, event **rpcapi.Event) error {
	(*event).Sheduled = false
	_, err := self.rpcClient.UpdateEvent(ctx, *(event))
	if err != nil {
		self.logger.Error(err.Error())
		return err
	}
	return nil
}

func (self *EventsSource) DataChannel(ctx context.Context) (<-chan *rpcapi.Event, error) {
	events, err := self.rpcClient.ListNotSheduledEvents(ctx)
	if err != nil {
		self.logger.Error(err.Error())
		return nil, err
	}
	eventsChan := make(chan *rpcapi.Event, len(events))
	for _, event := range events {
		eventsChan <- event
	}
	close(eventsChan)
	return eventsChan, nil
}
