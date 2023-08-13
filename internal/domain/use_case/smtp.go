package use_case

import "chat-server/internal/domain/entity"

type SMTPUseCase interface {
	Send(mail *entity.Mail) error
}
