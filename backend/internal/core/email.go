package core

import (
	"app/config"
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

func NewEmailSMTPEngine(cfg config.Config) *EmailSMTPEngine {
	return &EmailSMTPEngine{
		User:     cfg.SMTP.User,
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Password: cfg.SMTP.Password,
	}
}

func (esmtpe *EmailSMTPEngine) SendMail(subject string, text []byte, to []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("No Reply <%s>", esmtpe.User)
	e.ReplyTo = []string{esmtpe.User}
	e.To = to
	e.Subject = subject
	e.Text = text

	addr := fmt.Sprintf("%s:%s", esmtpe.Host, esmtpe.Port)
	auth := smtp.PlainAuth("", esmtpe.User, esmtpe.Password, esmtpe.Host)

	var err error
	if esmtpe.Port == "465" {
		err = e.SendWithTLS(addr, auth, &tls.Config{ServerName: esmtpe.Host})
	} else {
		err = e.Send(addr, auth)
	}
	if err != nil {
		return err
	}
	return nil
}
