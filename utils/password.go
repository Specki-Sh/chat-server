package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

// TODO: make it env variable
const (
	salt = "saltkey"
)

func HashPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	hashedPassword := hash.Sum([]byte(salt))

	return hex.EncodeToString(hashedPassword)
}
