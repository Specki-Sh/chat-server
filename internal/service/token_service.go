package service

import (
	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type KeyPair struct {
	PrivKey *rsa.PrivateKey
	PubKey  *rsa.PublicKey
}

// ts used for injecting an implementation of TokenRepository
// for use in service methods along with keys and secrets for
// signing JWTs
type tokenService struct {
	TokenRepository   use_case.TokenStorage
	AccessKeys        *KeyPair
	RefreshKeys       *KeyPair
	RefreshSecret     string
	AccessExpiration  *time.Duration
	RefreshExpiration *time.Duration
}

// TSConfig will hold repositories that will eventually be injected into this
// this service layer
type TSConfig struct {
	TokenRepository   use_case.TokenStorage
	AccessKeys        *KeyPair
	RefreshKeys       *KeyPair
	RefreshSecret     string
	AccessExpiration  *time.Duration
	RefreshExpiration *time.Duration
}

// NewTokenService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewTokenService(c *TSConfig) use_case.TokenUseCase {
	return &tokenService{
		TokenRepository:   c.TokenRepository,
		AccessKeys:        c.AccessKeys,
		RefreshKeys:       c.RefreshKeys,
		AccessExpiration:  c.AccessExpiration,
		RefreshExpiration: c.RefreshExpiration,
	}
}

type CachedTokens struct {
	AccessUID  string `json:"access"`
	RefreshUID string `json:"refresh"`
}

func (ts *tokenService) GenerateTokenPair(user *entity.User) (*entity.TokenPair, error) {
	accessToken, err := ts.generateToken(user.ID, user.Username, ts.AccessKeys.PrivKey, ts.AccessExpiration)
	if err != nil {
		return nil, fmt.Errorf("Could not generate access token: %v\n", err)
	}

	refreshToken, err := ts.generateToken(user.ID, user.Username, ts.RefreshKeys.PrivKey, ts.RefreshExpiration)
	if err != nil {
		return nil, fmt.Errorf("Could not generate refresh token: %v\n", err)
	}

	tokenPair := &entity.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokenPair, nil
}

func (ts *tokenService) ParseToken(tokenString string, key *rsa.PublicKey) (entity.ID, entity.NonEmptyString, error) {
	claims := &tokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return 0, "", fmt.Errorf("Could not parse JWT token: %v\n", err)
	}

	if !token.Valid {
		return 0, "", fmt.Errorf("Invalid JWT token\n")
	}

	return claims.ID, claims.UserName, nil
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
		return "", fmt.Errorf("Could not generate JWT token: %v\n", err)
	}

	return tokenString, nil
}
