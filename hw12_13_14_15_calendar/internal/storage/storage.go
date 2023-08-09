package storage

import (
	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	gomemory "hw12_13_14_15_calendar/internal/storage/gomemory"
	pgsqldtb "hw12_13_14_15_calendar/internal/storage/pgsqldtb"
)

var (
	GOMEMORY_STORAGE = "gomemory"
	POSTGRES_STORAGE = "pgsqldtb"
)

func NewStorageByType(storageType string, a ...any) interfaces.Storager {
	switch storageType {
	case GOMEMORY_STORAGE:
		return gomemory.NewStorage()
	case POSTGRES_STORAGE:
		return pgsqldtb.NewStorage(a[0].(string))
	default:
		return gomemory.NewStorage()
	}
}
