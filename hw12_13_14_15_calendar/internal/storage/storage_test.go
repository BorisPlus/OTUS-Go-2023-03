package storage

import (
	"testing"
	"time"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	models "hw12_13_14_15_calendar/internal/models"
	gomemory "hw12_13_14_15_calendar/internal/storage/gomemory"
	pgsqldtb "hw12_13_14_15_calendar/internal/storage/pgsqldtb"
)

var postgresDsn = "user=hw12user password=hw12user host='127.0.0.1' database=hw12calendar search_path=hw12calendar"

var _ = pgsqldtb.NewStorage(postgresDsn) // TODO: need mock

var testCases = []struct {
	storager interfaces.Storager
}{
	{
		pgsqldtb.NewStorage(postgresDsn),
	},
	{
		gomemory.NewStorage(),
	},
}

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
	for _, testCase := range testCases { //
		err := testCase.storager.Connect()
		if err != nil {
			t.Error(err)
		}
		defer func() {
			err := testCase.storager.Close()
			if err != nil {
				t.Error(err)
			}
		}()
		// CREATE
		testCase.storager.CreateEvent(&event)
		if event.PK == 0 {
			t.Error("event.PK must be auto set after insert")
		}
		// READ
		eventCopyFromStorageAfterInsert, err := testCase.storager.ReadEvent(event.PK)
		if err != nil {
			t.Error(err)
		}
		if eventCopyFromStorageAfterInsert == nil {
			t.Error("event not found in database after insert")
		}
		// READ (all)
		events, err := testCase.storager.ListEvents()
		if err != nil {
			t.Error(err)
		}
		if len(events) == 0 {
			t.Error("events list must be not empty")
		}
		// UPDATE
		event.Description = "new description"
		err = testCase.storager.UpdateEvent(&event)
		if err != nil {
			t.Error(err)
		}
		eventCopyFromStorageAfterUpdate, err := testCase.storager.ReadEvent(event.PK)
		if err != nil {
			t.Error(err)
		}
		if event.Description != eventCopyFromStorageAfterUpdate.Description {
			t.Error("event was not updated")
		}
		// DELETE
		testCase.storager.DeleteEvent(&event)
		checkExist, err := testCase.storager.ReadEvent(event.PK)
		if err != nil {
			t.Error(err)
		}
		if checkExist != nil {
			t.Error("event was not deleted")
		}
	}
}
