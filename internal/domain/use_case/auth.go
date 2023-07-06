package use_case

import (
	"chat-server/internal/domain/entity"
	"context"
)

type AuthUseCase interface {
	Authenticate(req *entity.SignInReq) (*entity.SignInRes, error)
	Logout(ctx context.Context, refreshToken string) error
	RefreshTokenPair(ctx context.Context, req *entity.RefreshTokenReq) (*entity.RefreshTokenRes, error)
}
