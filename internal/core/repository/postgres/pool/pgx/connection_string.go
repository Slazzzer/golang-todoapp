package pgx

import (
	"net"
	"net/url"
)

func buildConnectionString(config Config) string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(config.User, config.Password),
		Host:   net.JoinHostPort(config.Host, config.Port),
		Path:   "/" + config.Database,
	}

	query := u.Query()
	query.Set("sslmode", config.SSLMode)
	u.RawQuery = query.Encode()

	return u.String()
}
