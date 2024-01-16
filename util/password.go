package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// compute the hash string of the bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Hashed Password Failed: %w", err)
	}
	return string(hashedpassword), nil
}

// Check if the password is the password
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
