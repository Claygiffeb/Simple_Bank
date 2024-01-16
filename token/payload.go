package token

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("the token is expired")
	ErrInvalidToken = errors.New("the token is invalid")
)

type Payload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

// NewPayload creates a new tolen payload with the given username and duration
func NewPayLoad(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:       tokenID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return payload, nil
}

// For payload to become a jwt claim, it must implement the Valid() method
