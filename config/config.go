package config

import (
	"encoding/json"
	"time"
)

type Config struct {
	Endpoint        string
	PollingInterval time.Duration
	VirtualChains   []uint32
	DB              string
}

var defaultConfig = Config{
	PollingInterval: 2 * time.Second,
	Endpoint:        "http://localhost:8080",
	DB:              "./data",
}

func Parse(input []byte) (*Config, error) {
	cfg := &Config{}
	err := json.Unmarshal(input, cfg)

	if cfg.PollingInterval == 0 {
		cfg.PollingInterval = defaultConfig.PollingInterval
	}

	if cfg.Endpoint == "" {
		cfg.Endpoint = defaultConfig.Endpoint
	}

	if cfg.DB == "" {
		cfg.DB = defaultConfig.DB
	}

	return cfg, err
}
