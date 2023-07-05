package use_case

import (
	"chat-server/internal/domain/entity"
	"crypto/rsa"
)

type TokenUseCase interface {
	GenerateTokenPair(user *entity.User) (*entity.TokenPair, error)
	ParseToken(tokenString string, key *rsa.PublicKey) (entity.ID, entity.NonEmptyString, error)
}
