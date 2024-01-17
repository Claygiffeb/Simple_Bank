package token

import "time"

// This is an interface for managing tokens/ JWT or sth elese
type Maker interface {

	// CreateToken creates a new token for the given username and duration
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	// VerifyToken checks if the given token is valid
	VerifyToken(token string) (*Payload, error)
}
