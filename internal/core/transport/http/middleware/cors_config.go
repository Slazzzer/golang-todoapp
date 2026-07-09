package core_http_middleware

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// CORSConfig настройки CORS из переменных окружения.
// ALLOWED_ORIGINS — список через запятую: полный origin (http://host:port),
// hostname (80.249.145.151 — любой порт) или "null" для file://.
type CORSConfig struct {
	AllowedOrigins string `envconfig:"ALLOWED_ORIGINS" default:"http://localhost:5050,http://127.0.0.1:5050,null,80.249.145.151,http://80.249.145.151:5050"`
}

type corsChecker struct {
	exactOrigins map[string]struct{}
	allowedHosts map[string]struct{}
}

func NewCORSConfig() (CORSConfig, error) {
	var cfg CORSConfig

	if err := envconfig.Process("", &cfg); err != nil {
		return CORSConfig{}, fmt.Errorf("process envconfig: %w", err)
	}

	return cfg, nil
}

func NewCORSConfigMust() CORSConfig {
	cfg, err := NewCORSConfig()
	if err != nil {
		panic(fmt.Errorf("get cors config: %w", err))
	}

	return cfg
}

func (c CORSConfig) checker() corsChecker {
	exact := make(map[string]struct{})
	hosts := make(map[string]struct{})

	for _, part := range strings.Split(c.AllowedOrigins, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if part == "null" {
			exact["null"] = struct{}{}
			continue
		}

		if strings.Contains(part, "://") {
			exact[part] = struct{}{}
			if u, err := url.Parse(part); err == nil && u.Hostname() != "" {
				hosts[u.Hostname()] = struct{}{}
			}
			continue
		}

		hosts[part] = struct{}{}
	}

	return corsChecker{
		exactOrigins: exact,
		allowedHosts: hosts,
	}
}

func (c corsChecker) isOriginAllowed(origin string) bool {
	if origin == "" {
		return true
	}

	if _, ok := c.exactOrigins[origin]; ok {
		return true
	}

	if origin == "null" {
		_, ok := c.exactOrigins["null"]
		return ok
	}

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	_, ok := c.allowedHosts[u.Hostname()]
	return ok
}
