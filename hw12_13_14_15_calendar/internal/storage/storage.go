package storage

import (
	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	gomemory "hw12_13_14_15_calendar/internal/storage/gomemory"
	pgsqldtb "hw12_13_14_15_calendar/internal/storage/pgsqldtb"
)

func NewStorageByType(storageType string, a... any) interfaces.Storager {
	switch storageType {
	case "gomemory":
		return gomemory.NewStorage()
	case "pgsqldtb":
		return pgsqldtb.NewStorage(a[0].(string))
	default:
		return gomemory.NewStorage()
	}
}
