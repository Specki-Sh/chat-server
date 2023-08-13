package service

import (
	"context"
	"fmt"

	"chat-server/internal/domain/entity"
	u "chat-server/internal/domain/use_case"
	"chat-server/utils"
)

func NewAuthService(userUseCase u.UserUseCase, tokenUseCase u.TokenUseCase) u.AuthUseCase {
	return &authService{
		userUseCase:  userUseCase,
		tokenUseCase: tokenUseCase,
	}
}

type authService struct {
	userUseCase  u.UserUseCase
	tokenUseCase u.TokenUseCase
}

func (a *authService) Authenticate(req *entity.SignInReq) (*entity.SignInRes, error) {
	hashPassword := utils.HashPassword(req.Password)
	user, err := a.userUseCase.GetByEmailAndPassword(req.Email, hashPassword)
	if err != nil {
		return nil, fmt.Errorf("authService.Authenticate: %w", err)
	}

	tokenPair, err := a.tokenUseCase.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("authService.Authenticate: %w", err)
	}

	return &entity.SignInRes{TokenPair: *tokenPair, ID: user.ID, Username: user.Username}, nil
}

func (a *authService) Logout(ctx context.Context, refreshToken string) error {
	if err := a.tokenUseCase.InvalidateRefreshToken(ctx, refreshToken); err != nil {
		return fmt.Errorf("authService.Authenticate: %w", err)
	}
	return nil
}

func (a *authService) RefreshTokenPair(
	ctx context.Context,
	req *entity.RefreshTokenReq,
) (*entity.RefreshTokenRes, error) {
	if err := a.tokenUseCase.ValidateRefreshToken(ctx, req.RefreshToken); err != nil {
		return nil, fmt.Errorf("authService.RefreshTokenPair: %w", err)
	}

	accessToken, err := a.tokenUseCase.GenerateAccessToken(req.ID, req.Username)
	if err != nil {
		return nil, fmt.Errorf("authService.RefreshTokenPair: %w", err)
	}

	tokenPair := entity.TokenPair{RefreshToken: req.RefreshToken, AccessToken: accessToken}

	return &entity.RefreshTokenRes{TokenPair: tokenPair, ID: req.ID, Username: req.Username}, nil
}
