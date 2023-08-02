package gomemory

import (
	"testing"
	"time"

	models "hw12_13_14_15_calendar/internal/models"
)

// CRUD ethalon.
var event = models.Event{
	Title:       "Title",
	StartAt:     time.Now(),
	Duration:    3_600,
	Description: "no description",
	Owner:       "admin",
	NotifyEarly: 600,
}

func TestStorage(t *testing.T) {
	storage := NewStorage()
	err := storage.Connect()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err := storage.Close()
		if err != nil {
			t.Error(err)
		}
	}()
	// CREATE
	firstEvent, _ := storage.CreateEvent(&event)
	if firstEvent.PK != 1 {
		t.Error("event.PK must be auto set after insert to 1")
	}
	secondEvent, _ := storage.CreateEvent(&event)
	if secondEvent.PK != 2 {
		t.Error("event.PK must be auto set after insert to 2")
	}
	// DELETE
	storage.DeleteEvent(secondEvent)
	// CREATE
	thirdEvent, _ := storage.CreateEvent(&event)
	if thirdEvent.PK != 3 {
		t.Error("event.PK must be auto set after insert to 3")
	}
	// READ (all)
	events, err := storage.ListEvents()
	if err != nil {
		t.Error(err)
	}
	if len(events) != 2 {
		t.Error("events list must be 2")
	}
	// UPDATE
	thirdEvent.Description = "thirdEvent"
	_, err = storage.UpdateEvent(thirdEvent)
	if err != nil {
		t.Error(err)
	}
	// READ
	eventCopyFromStorageAfterUpdate, err := storage.ReadEvent(thirdEvent.PK)
	if err != nil {
		t.Error(err)
	}
	if event.Description != eventCopyFromStorageAfterUpdate.Description {
		t.Error("event was not updated")
	}
}
