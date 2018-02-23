package config

import "time"

type Config struct {
	PollIntervalMs time.Duration `mapstructure:"poll_interval_ms"`
	HTTPPort       int           `mapstructure:"http_port"`
	Interfaces     []Interface   `mapstructure:"interfaces"`
}

type Interface struct {
	Domain string `mapstructure:"domain"`
	Name   string `mapstructure:"name"`
}

func NewConfig() *Config {
	return &Config{
		PollIntervalMs: 1000,
		HTTPPort:       8888,
	}
}
