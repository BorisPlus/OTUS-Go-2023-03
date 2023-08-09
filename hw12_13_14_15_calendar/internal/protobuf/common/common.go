package calendar_common

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	models "hw12_13_14_15_calendar/internal/models"
	pb "hw12_13_14_15_calendar/internal/protobuf/api"
)

func Event2PBEvent(event *models.Event) *pb.Event {
	pbEvent := new(pb.Event)
	// `event::models.Event`` - переменная models.Event-типа, возвращаемого БД через App
	// Как конвертировать `pbEvent::pb.Event` в `event::models.Event` и обратно?
	pbEvent.PK = int32(event.PK)
	pbEvent.Title = event.Title
	pbEvent.Description = event.Description
	pbEvent.StartAt = timestamppb.New(event.StartAt)
	pbEvent.Duration = int32(event.Duration)
	pbEvent.Owner = event.Owner
	pbEvent.NotifyEarly = int32(event.NotifyEarly)
	return pbEvent
}

func PBEvent2Event(pbEvent *pb.Event) *models.Event {
	event := new(models.Event)
	event.PK = int(pbEvent.PK)
	event.Title = pbEvent.Title
	event.Description = pbEvent.Description
	event.StartAt = pbEvent.StartAt.AsTime()
	event.Duration = int(pbEvent.Duration)
	event.Owner = pbEvent.Owner
	event.NotifyEarly = int(pbEvent.NotifyEarly)
	return event
}
