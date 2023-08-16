package config

import "hw12_13_14_15_calendar/internal/models"

type PRCClientDSN struct {
	DSN string
}

type ShedulerConfig struct {
	Source     PRCClientDSN
	Target     models.RabbitMQTarget
	Log        LogConfig
	TimeoutSec int64
}

func NewShedulerConfig() *ShedulerConfig {
	return &ShedulerConfig{}
}
