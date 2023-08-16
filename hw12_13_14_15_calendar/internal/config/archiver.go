package config

import "hw12_13_14_15_calendar/internal/models"

type ArchiverConfig struct {
	Source     models.RabbitMQSource
	ArchiveTo  models.RabbitMQTarget
	Log        LogConfig
	TimeoutSec int64
}

func NewArchiverConfig() *ArchiverConfig {
	return &ArchiverConfig{}
}
