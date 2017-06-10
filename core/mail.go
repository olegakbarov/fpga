package core

import "fmt"
import "sync"

type (
	Mailer interface {
		SendWelcomeMail(to string, activationURL string) error
		SendPasswordResetLink(to string, resetLink string) error
	}

	mail struct {
		Mailer
		ms MailSender
	}
)

var (
	mailOnce     sync.Once
	mailInstance *mail
)

func (f *factory) NewMail() Mailer {
	mailOnce.Do(func() {
		mailInstance = &mail{ms: f.ms}
	})
	return mailInstance
}

func (m *mail) SendWelcomeMail(to string, activationURL string) error {
	subject := "Welcome"
	body := fmt.Sprintf("Welcome! Please click here to activate your account. %s", activationURL)

	return m.ms.Send([]string{to}, subject, []byte(body))
}

func (m *mail) SendPasswordResetLink(to string, resetLink string) error {
	subject := "Password reset"
	body := fmt.Sprintf("Please click below link to reset your password <br/> %s", resetLink)

	return m.ms.Send([]string{to}, subject, []byte(body))
}
