package notification

import (
	"bytes"
	"net"
	mail "net/mail"
	"net/smtp"
	"strings"

	"github.com/GGP1/comeet/event"

	"github.com/pkg/errors"
)

type mailNotifier struct {
	config MailConfig
}

// NewMailNotifier returns a notifier that sends mails.
func NewMailNotifier(config MailConfig) (Notifier, error) {
	return &mailNotifier{config: config}, nil
}

func (m *mailNotifier) Notify(event *event.Event) error {
	var auth smtp.Auth
	// TODO: use https://github.com/go-mail/mail?
	if m.config.SMTPHostAddress == "smtp-mail.outlook.com" {
		auth = LoginAuth(m.config.SenderAddress, m.config.SenderPassword)
	} else {
		auth = smtp.PlainAuth("", m.config.SenderAddress, m.config.SenderPassword, m.config.SMTPHostAddress)
	}

	from := mail.Address{Name: "", Address: m.config.SenderAddress}
	to := strings.Join(m.config.ReceiverAddresses, ", ")

	headers := make(map[string]string, 4)
	headers["From"] = from.String()
	headers["To"] = to
	headers["Subject"] = event.Title
	headers["Content-Type"] = `text/html; charset="UTF-8"`

	content := buildMailContent(headers, event.Message())

	hostAddr := net.JoinHostPort(m.config.SMTPHostAddress, m.config.SMTPHostPort)
	if err := smtp.SendMail(hostAddr, auth, from.Address, m.config.ReceiverAddresses, content.Bytes()); err != nil {
		return errors.Wrap(err, "couldn't send the mail")
	}

	return nil
}

func buildMailContent(headers map[string]string, body string) *bytes.Buffer {
	message := new(bytes.Buffer)

	for k, v := range headers {
		// key: value\r\n
		message.WriteString(k)
		message.WriteString(": ")
		message.WriteString(v)
		message.WriteString("\r\n")
	}

	// TODO: use a template?
	message.WriteString(body)

	return message
}

type loginAuth struct {
	username, password string
}

// LoginAuth returns an authentication type that is not supported by the standard smtp package.
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if !more {
		return nil, nil
	}

	switch {
	case bytes.Equal(fromServer, []byte("Username:")):
		return []byte(a.username), nil

	case bytes.Equal(fromServer, []byte("Password:")):
		return []byte(a.password), nil

	default:
		return nil, errors.Errorf("unexpected server challenge: %s", fromServer)
	}
}
