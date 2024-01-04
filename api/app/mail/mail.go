package mail

import (
	"blog/app"
	"blog/config"
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
	"net/textproto"
)

type Mail struct {
	*email.Email
	ServerAddr string
	Auth       smtp.Auth
}

func (m *Mail) Send(subject, message string) (err error) {
	m.Subject = subject
	m.HTML = []byte(message)

	app.Logger.Sugar().Debugf("%+v", m)

	return m.Email.Send(m.ServerAddr, m.Auth)
}

func (m *Mail) AppendReplyTo(replyTo ...string) *Mail {
	m.ReplyTo = append(m.ReplyTo, replyTo...)

	return m
}

func (m *Mail) AppendTo(to ...string) *Mail {
	m.To = append(m.To, to...)

	return m
}

func (m *Mail) AppendCc(cc ...string) *Mail {
	m.Cc = append(m.Cc, cc...)

	return m
}

func (m *Mail) AppendBcc(bcc ...string) *Mail {
	m.Bcc = append(m.Bcc, bcc...)

	return m
}

func New() *Mail {
	return &Mail{
		ServerAddr: fmt.Sprintf("%s:%d", config.Mail.Host, config.Mail.Port),
		Auth:       smtp.PlainAuth("", config.Mail.Username, config.Mail.Password, config.Mail.Host),
		Email: &email.Email{
			Headers: textproto.MIMEHeader{},
			From:    fmt.Sprintf("%s <%s>", config.Mail.FromMame, config.Mail.FromAddr),
		},
	}
}
