package common

import (
	models "hw12_13_14_15_calendar/internal/models"
	calendarrpcapi "hw12_13_14_15_calendar/internal/server/rpc/rpcapi"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func Event2PBEvent(event *models.Event) *calendarrpcapi.Event {
	pbEvent := new(calendarrpcapi.Event)
	pbEvent.PK = int32(event.PK)
	pbEvent.Title = event.Title
	pbEvent.Description = event.Description
	pbEvent.StartAt = timestamppb.New(event.StartAt)
	pbEvent.Duration = int32(event.Duration)
	pbEvent.Owner = event.Owner
	pbEvent.NotifyEarly = int32(event.NotifyEarly)
	pbEvent.Sheduled = event.Sheduled
	return pbEvent
}

func PBEvent2Event(pbEvent *calendarrpcapi.Event) *models.Event {
	event := new(models.Event)
	event.PK = int(pbEvent.PK)
	event.Title = pbEvent.Title
	event.Description = pbEvent.Description
	event.StartAt = pbEvent.StartAt.AsTime()
	event.Duration = int(pbEvent.Duration)
	event.Owner = pbEvent.Owner
	event.NotifyEarly = int(pbEvent.NotifyEarly)
	event.Sheduled = pbEvent.Sheduled
	return event
}
