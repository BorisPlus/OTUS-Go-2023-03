package config

import "hw12_13_14_15_calendar/internal/models"

type DatasetsCount struct {
	Send, Archive, Defer int
}

type CheckerConfig struct {
	HTTP     HTTPClient
	Storage  StorageConfig
	Sended   models.RabbitMQSource
	Archived models.RabbitMQSource
	Log      LogConfig
	Counts   DatasetsCount
}

func NewCheckerConfig() *CheckerConfig {
	return &CheckerConfig{}
}
