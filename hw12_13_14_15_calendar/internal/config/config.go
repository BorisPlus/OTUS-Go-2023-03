package config

import (
	"time"

	logger "hw12_13_14_15_calendar/internal/logger"
)

type RPCConfig struct {
	Host string
	Port uint16
}

type HTTPConfig struct {
	Host              string
	Port              uint16
	ReadTimeout       time.Duration // TODO: time.Duration
	ReadHeaderTimeout time.Duration // TODO: time.Duration
	WriteTimeout      time.Duration // TODO: time.Duration
	MaxHeaderBytes    int
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
	RPC     RPCConfig
	Storage StorageConfig
	Log     LogConfig
}

func NewConfig() *Config {
	return &Config{}
}
