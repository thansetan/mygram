package helper

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(raw string) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), 5)
	if err != nil {
		return nil, fmt.Errorf("helper.HashPassword: %w", err)
	}

	return hashed, nil
}

func IsValidPassword(hashed []byte, raw string) bool {
	err := bcrypt.CompareHashAndPassword(hashed, []byte(raw))
	return err == nil
}
