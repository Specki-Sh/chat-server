package service

import (
	"chat-server/internal/domain/entity"
	u "chat-server/internal/domain/use_case"
	"chat-server/utils"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	tokenTTL   = 12 * time.Hour
	signingKey = "signingKey12345"
)

type tokenClaims struct {
	jwt.StandardClaims
	ID       int    `json:"user_id"`
	UserName string `json:"user_name"`
}

func NewAuthService(userUseCase u.UserUseCase) *AuthService {
	return &AuthService{userUseCase: userUseCase}
}

type AuthService struct {
	userUseCase u.UserUseCase
}

func (a *AuthService) GenerateToken(req *entity.SignInReq) (*entity.SignInRes, error) {
	password := utils.HashPassword(req.Password)
	user, err := a.userUseCase.GetByEmailAndPassword(req.Email, password)
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
		user.Username,
	})

	ss, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return &entity.SignInRes{}, err
	}

	return &entity.SignInRes{AccessToken: ss, ID: user.ID, Username: user.Username}, nil
}

func (a *AuthService) ParseToken(accessToken string) (int, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, "", errors.New("token claims are not of type *tokenClaims")
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return 0, "", errors.New("token expired")
	}

	return claims.ID, claims.UserName, nil
}
