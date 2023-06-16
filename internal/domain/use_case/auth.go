package use_case

import (
	"chat-server/internal/domain/entity"
)

type AuthUseCase interface {
	GenerateToken(req *entity.SignInReq) (*entity.SignInRes, error)
	ParseToken(accessToken string) (int, error)
}
