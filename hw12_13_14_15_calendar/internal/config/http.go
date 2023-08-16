package config

import (
	"time"
)

type HTTPConfig struct {
	Host              string
	Port              uint16
	ReadTimeout       time.Duration // TODO: time.Duration
	ReadHeaderTimeout time.Duration // TODO: time.Duration
	WriteTimeout      time.Duration // TODO: time.Duration
	MaxHeaderBytes    int
}
