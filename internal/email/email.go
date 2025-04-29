package email

import (
	"github.com/labovector/vecsys-api/infrastructure/config"
	"gopkg.in/gomail.v2"
)

type EmailDialer struct {
	*gomail.Dialer
	SenderName string
}

func NewEmailDialer(config *config.EmailConfig) *EmailDialer {
	return &EmailDialer{
		Dialer:     gomail.NewDialer(config.Host, config.Port, config.AuthEmail, config.AuthPassword),
		SenderName: config.SenderName,
	}
}
