package adminauth

import (
	"crypto/subtle"
	"os"
)

type Config struct {
	Login        string
	Password     string
	SessionToken string
}

func NewConfig() Config {
	return Config{
		Login:        os.Getenv("ADMIN_LOGIN"),
		Password:     os.Getenv("ADMIN_PASSWORD"),
		SessionToken: os.Getenv("ADMIN_SESSION_TOKEN"),
	}
}

func (c Config) ValidateCredentials(login, password string) bool {
	if c.Login == "" || c.Password == "" || c.SessionToken == "" {
		return false
	}
	loginMatch := subtle.ConstantTimeCompare([]byte(login), []byte(c.Login)) == 1
	passwordMatch := subtle.ConstantTimeCompare([]byte(password), []byte(c.Password)) == 1
	return loginMatch && passwordMatch
}

func (c Config) ValidateSession(token string) bool {
	if c.SessionToken == "" || token == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(token), []byte(c.SessionToken)) == 1
}
