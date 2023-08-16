package config

import "hw12_13_14_15_calendar/internal/models"

type RMQOpsConfig struct {
	RabbitMQ models.RabbitMQNode
}

func NewRMQOpsConfig() *RMQOpsConfig {
	return &RMQOpsConfig{}
}
