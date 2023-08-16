package config

import "hw12_13_14_15_calendar/internal/models"

type SenderConfig struct {
	Source     models.RabbitMQSource
	SendTo     models.RabbitMQTarget
	ArchiveTo  models.RabbitMQTarget
	Log        LogConfig
	TimeoutSec int64
}

func NewSenderConfig() *SenderConfig {
	return &SenderConfig{}
}
