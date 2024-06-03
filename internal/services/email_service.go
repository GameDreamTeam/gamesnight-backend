package services

import (
	"fmt"
	"net/smtp"
)

type EmailService struct {
	auth smtp.Auth
	from string
}

var es *EmailService

func NewEmailService(username, password, host string, port int) {
	auth := smtp.PlainAuth("", username, password, host)
	es = &EmailService{
		auth: auth,
		from: username,
	}
}

func GetEmailService() *EmailService {
	return es
}

func (es *EmailService) SendEmail(to, subject, body string) error {
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", to, subject, body))
	addr := fmt.Sprintf("%s:%d", "smtp.gmail.com", 587)

	return smtp.SendMail(addr, es.auth, es.from, []string{to}, msg)
}
