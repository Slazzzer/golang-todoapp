package core_ratelimit

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RegisterMax    int           `envconfig:"REGISTER_MAX" default:"30"`
	RegisterWindow time.Duration `envconfig:"REGISTER_WINDOW" default:"1m"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("RATE_LIMIT", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	if config.RegisterMax <= 0 {
		return Config{}, fmt.Errorf("RATE_LIMIT_REGISTER_MAX must be positive")
	}

	if config.RegisterWindow <= 0 {
		return Config{}, fmt.Errorf("RATE_LIMIT_REGISTER_WINDOW must be positive")
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("get rate limit config: %w", err))
	}

	return config
}
