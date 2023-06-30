package utils

import (
	"chat-server/internal/domain/entity"
	"crypto/sha1"
	"encoding/hex"
)

// TODO: make it env variable
const (
	salt = "saltkey"
)

func HashPassword(password entity.Password) entity.HashPassword {
	hash := sha1.New()
	hash.Write([]byte(password))
	hashedPassword := hash.Sum([]byte(salt))

	return entity.HashPassword(hex.EncodeToString(hashedPassword))
}
