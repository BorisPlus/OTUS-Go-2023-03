package storage

import (
	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	gomemory "hw12_13_14_15_calendar/internal/storage/gomemory"
	pgsqldtb "hw12_13_14_15_calendar/internal/storage/pgsqldtb"
)

var (
	GoMemoryStorage = "gomemory"
	PostgresStorage = "pgsqldtb"
)

func NewStorageByType(storageType string, a ...any) interfaces.Storager {
	switch storageType {
	case GoMemoryStorage:
		return gomemory.NewStorage()
	case PostgresStorage:
		return pgsqldtb.NewStorage(a[0].(string))
	default:
		return gomemory.NewStorage()
	}
}
