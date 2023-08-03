package use_case

import (
	"context"

	"chat-server/internal/domain/entity"
)

type AuthUseCase interface {
	Authenticate(req *entity.SignInReq) (*entity.SignInRes, error)
	Logout(ctx context.Context, refreshToken string) error
	RefreshTokenPair(ctx context.Context, req *entity.RefreshTokenReq) (*entity.RefreshTokenRes, error)
}
