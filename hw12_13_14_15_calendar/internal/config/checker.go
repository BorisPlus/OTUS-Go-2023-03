package config

import "hw12_13_14_15_calendar/internal/models"

type CheckerConfig struct {
	HTTP     HTTPConfig
	Storage  StorageConfig
	Sended   models.RabbitMQSource
	Archived models.RabbitMQSource
}

func NewCheckerConfig() *CheckerConfig {
	return &CheckerConfig{}
}
