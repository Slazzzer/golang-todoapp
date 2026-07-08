package core_auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

type claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func NewTokenManager(config Config) *TokenManager {
	return &TokenManager{
		secret: []byte(config.Secret),
		issuer: config.Issuer,
		ttl:    config.TTL,
	}
}

func (m *TokenManager) Issue(userID int) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userID),
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	})

	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}

	return signed, nil
}

func (m *TokenManager) ParseUserID(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return m.secret, nil
		},
		jwt.WithIssuer(m.issuer),
	)
	if err != nil {
		return 0, fmt.Errorf("parse jwt: %w", err)
	}

	claims, ok := token.Claims.(*claims)
	if !ok || !token.Valid || claims.UserID <= 0 {
		return 0, fmt.Errorf("invalid jwt claims")
	}

	return claims.UserID, nil
}
