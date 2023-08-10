package service

import (
	"fmt"
	"net/smtp"

	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
)

func NewSMTPService(config *SMTPConfig) use_case.SMTPUseCase {
	return &smtpService{
		config: config,
	}
}

type SMTPConfig struct {
	Email    string
	Password string
	Host     string
	Port     string
}

type smtpService struct {
	config *SMTPConfig
}

func (s *smtpService) Send(mail *entity.Mail) error {
	from := s.config.Email
	password := s.config.Password
	toList := []string{mail.To}
	host := s.config.Host
	port := s.config.Port

	// Формирование заголовка письма
	header := make(map[string]string)
	header["From"] = from
	header["To"] = mail.To
	header["Subject"] = mail.Subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	// Формирование тела письма
	body := ""
	for k, v := range header {
		body += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	body += "\r\n" + mail.Body

	// Отправка письма
	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(host+":"+port, auth, from, toList, []byte(body))
	if err != nil {
		return fmt.Errorf("smtpService.Send: %w", err)
	}
	return nil
}
