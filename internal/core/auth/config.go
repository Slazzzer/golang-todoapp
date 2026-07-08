package core_auth

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Secret string        `envconfig:"SECRET" required:"true"`
	TTL    time.Duration `envconfig:"TTL"    default:"720h"`
	Issuer string        `envconfig:"ISSUER" default:"todoapp"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("JWT", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("get JWT config: %w", err))
	}

	return config
}
