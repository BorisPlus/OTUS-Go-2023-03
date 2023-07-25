package config

import (
	logger "hw12_13_14_15_calendar/internal/logger"
)

type HTTPConfig struct {
	Host string
	Port uint16
}

type StorageConfig struct {
	Type string
	DSN  string
}

type LogConfig struct {
	Level logger.LogLevel
}

type Config struct {
	HTTP    HTTPConfig
	Storage StorageConfig
	Log     LogConfig
}

func NewConfig() *Config {
	return &Config{}
}
