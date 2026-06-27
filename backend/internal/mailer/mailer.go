package mailer

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	host     string
	port     int
	user     string
	pass     string
	from     string
	enabled  bool
}

func New(host string, port int, user, pass, from string) *Mailer {
	return &Mailer{
		host:    host,
		port:    port,
		user:    user,
		pass:    pass,
		from:    from,
		enabled: host != "" && user != "",
	}
}

func (m *Mailer) Enabled() bool { return m.enabled }

func (m *Mailer) Send(to, subject, body string) error {
	if !m.enabled {
		return nil
	}

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s\r\n", m.from, to, subject, body))

	addr := fmt.Sprintf("%s:%d", m.host, m.port)
	auth := smtp.PlainAuth("", m.user, m.pass, m.host)

	return smtp.SendMail(addr, auth, m.from, []string{to}, msg)
}

func (m *Mailer) SendVerificationCode(to, code string) error {
	subject := "Seu codigo de verificacao - Rifa Online"
	body := fmt.Sprintf("Seu codigo de verificacao e: %s\n\nEste codigo expira em 30 minutos.", code)
	return m.Send(to, subject, body)
}
