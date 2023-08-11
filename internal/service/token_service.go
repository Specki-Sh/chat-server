package service

import (
	"context"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"

	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
)

type KeyPair struct {
	PrivKey *rsa.PrivateKey
	PubKey  *rsa.PublicKey
}

type tokenService struct {
	TokenRepository   use_case.TokenStorage
	AccessKeys        *KeyPair
	RefreshKeys       *KeyPair
	AccessExpiration  *time.Duration
	RefreshExpiration *time.Duration
}

type TSConfig struct {
	TokenRepository   use_case.TokenStorage
	AccessKeys        *KeyPair
	RefreshKeys       *KeyPair
	AccessExpiration  *time.Duration
	RefreshExpiration *time.Duration
}

func NewTokenService(c *TSConfig) use_case.TokenUseCase {
	return &tokenService{
		TokenRepository:   c.TokenRepository,
		AccessKeys:        c.AccessKeys,
		RefreshKeys:       c.RefreshKeys,
		AccessExpiration:  c.AccessExpiration,
		RefreshExpiration: c.RefreshExpiration,
	}
}

type tokenClaims struct {
	jwt.StandardClaims
	ID       entity.ID             `json:"user_id"`
	UserName entity.NonEmptyString `json:"user_name"`
}

func (ts *tokenService) GenerateTokenPair(user *entity.User) (*entity.TokenPair, error) {
	accessToken, err := ts.generateToken(user.ID, user.Username, ts.AccessKeys.PrivKey, ts.AccessExpiration)
	if err != nil {
		return nil, fmt.Errorf("tokenService.GenerateTokenPair: %w", err)
	}

	refreshToken, err := ts.generateToken(user.ID, user.Username, ts.RefreshKeys.PrivKey, ts.RefreshExpiration)
	if err != nil {
		return nil, fmt.Errorf("tokenService.GenerateTokenPair: %w", err)
	}

	tokenPair := &entity.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokenPair, nil
}

func (ts *tokenService) ParseAccessToken(accessToken string) (entity.ID, entity.NonEmptyString, error) {
	claims, err := ts.parseToken(accessToken, ts.AccessKeys.PubKey)
	if err != nil {
		return 0, "", fmt.Errorf("tokenService.ParseAccessToken: %w", err)
	}
	return claims.ID, claims.UserName, nil
}

func (ts *tokenService) ParseRefreshToken(tokenString string) (entity.ID, entity.NonEmptyString, error) {
	claims, err := ts.parseToken(tokenString, ts.RefreshKeys.PubKey)
	if err != nil {
		return 0, "", fmt.Errorf("tokenService.ParseRefreshToken: %w", err)
	}
	return claims.ID, claims.UserName, nil
}

func (ts *tokenService) GenerateAccessToken(userID entity.ID, username entity.NonEmptyString) (string, error) {
	accessToken, err := ts.generateToken(userID, username, ts.AccessKeys.PrivKey, ts.AccessExpiration)
	if err != nil {
		return "", fmt.Errorf("tokenService.GenerateAccessToken: %w", err)
	}
	return accessToken, nil
}

func (ts *tokenService) ValidateRefreshToken(ctx context.Context, refreshToken string) error {
	exists, err := ts.TokenRepository.InvalidRefreshTokenExists(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("tokenService.ValidateRefreshToken: %w", err)
	}
	if exists {
		return fmt.Errorf("tokenService.ValidateRefreshToken: %w", fmt.Errorf("refresh token is invalid: %s", refreshToken))
	}

	if _, err := ts.parseToken(refreshToken, ts.RefreshKeys.PubKey); err != nil {
		return fmt.Errorf("tokenService.ValidateRefreshToken: %w", err)
	}

	return nil
}

func (ts *tokenService) InvalidateRefreshToken(ctx context.Context, refreshToken string) error {
	claims, err := ts.parseToken(refreshToken, ts.RefreshKeys.PubKey)
	if err != nil {
		return fmt.Errorf("tokenService.InvalidateRefreshToken: %w", err)
	}
	expiresIn := time.Until(time.Unix(claims.ExpiresAt, 0))
	if err := ts.TokenRepository.SetInvalidRefreshToken(ctx, claims.ID, refreshToken, expiresIn); err != nil {
		return fmt.Errorf("tokenService.InvalidateRefreshToken: %w", err)
	}
	return nil
}

func (ts *tokenService) parseToken(tokenString string, key *rsa.PublicKey) (*tokenClaims, error) {
	claims := &tokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("tokenService.parseToken: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("tokenService.parseToken: %w", fmt.Errorf("invalid jwt token: %s", tokenString))
	}

	return claims, nil
}

func (ts *tokenService) generateToken(userID entity.ID, userName entity.NonEmptyString, key *rsa.PrivateKey, expiresIn *time.Duration) (string, error) {
	expirationTime := time.Now().Add(*expiresIn)

	claims := &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
		ID:       userID,
		UserName: userName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("tokenService.generateToken: %w", err)
	}

	return tokenString, nil
}
