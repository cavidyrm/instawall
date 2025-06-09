package email

import (
	"fmt"
	"net/smtp"
)

type SMTPSender struct {
	From     string
	Host     string
	Port     int
	Username string
	Password string
}

func NewSMTPSender(from, host string, port int, username, password string) *SMTPSender {
	return &SMTPSender{
		From:     from,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (s *SMTPSender) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		body + "\r\n")

	return smtp.SendMail(addr, auth, s.From, []string{to}, msg)
}
