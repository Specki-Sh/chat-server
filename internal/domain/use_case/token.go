package use_case

import (
	"context"

	"chat-server/internal/domain/entity"
)

type TokenUseCase interface {
	GenerateTokenPair(user *entity.User) (*entity.TokenPair, error)
	ParseAccessToken(accessToken string) (entity.ID, entity.NonEmptyString, error)
	ParseRefreshToken(tokenString string) (entity.ID, entity.NonEmptyString, error)
	GenerateAccessToken(userID entity.ID, username entity.NonEmptyString) (string, error)
	ValidateRefreshToken(ctx context.Context, refreshToken string) error
	InvalidateRefreshToken(ctx context.Context, refreshToken string) error
}
