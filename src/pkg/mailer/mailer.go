package mailer

import (
	"context"
	"crypto/tls"

	"github.com/pablor21/goms/app/config"
	mailer_config "github.com/pablor21/goms/pkg/mailer/config"
	"github.com/wneessen/go-mail"
)

type MailAddress struct {
	Email string
	Name  string
}

type MailMessage struct {
	To          []MailAddress
	Subject     string
	HtmlBody    []byte
	PlainBody   []byte
	Attachments []string
}

type Mailer interface {
	SendTemplate(ctx context.Context, template MailView, message MailMessage) error
	SendMessage(ctx context.Context, message MailMessage) error
	SendMessageFrom(ctx context.Context, from MailAddress, message MailMessage) error
}

type MailerService struct {
	config mailer_config.MailerConnectionConfig
}

func NewMailer(config mailer_config.MailerConnectionConfig) Mailer {
	return &MailerService{
		config: config,
	}
}

func GetMailer(name string) Mailer {
	return NewMailer(config.GetConfig().Mailer[name])
}

func GetDefaultMailer() Mailer {
	return GetMailer("default")
}

func (ms *MailerService) Close(sender *mail.Client) error {
	return sender.Close()
}

func (ms *MailerService) Open(ctx context.Context) (sender *mail.Client, err error) {
	sender, err = mail.NewClient(
		ms.config.Host,
		mail.WithSMTPAuth(mail.SMTPAuthType(ms.config.SmtpAuthType)),
		mail.WithUsername(ms.config.Username),
		mail.WithPassword(ms.config.Password),
		mail.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
		mail.WithPort(ms.config.Port))
	if err != nil {
		return
	}
	err = sender.DialWithContext(ctx)
	return
}

func (ms *MailerService) SendMessage(ctx context.Context, message MailMessage) error {
	return ms.SendMessageFrom(ctx, MailAddress{Email: ms.config.From}, message)
}

func (ms *MailerService) SendTemplate(ctx context.Context, template MailView, message MailMessage) error {
	str, err := template.RenderToByte()
	if err != nil {
		return err
	}
	message.HtmlBody = str

	return ms.SendMessage(ctx, message)
}

func (ms *MailerService) SendMessageFrom(ctx context.Context, from MailAddress, message MailMessage) error {
	sender, err := ms.Open(ctx)
	if err != nil {
		return err
	}
	defer ms.Close(sender)
	messages := make([]*mail.Msg, 0)
	for _, to := range message.To {
		gMessage := mail.NewMsg()
		gMessage.From(from.Email)
		gMessage.To(to.Email)
		gMessage.Subject(message.Subject)
		gMessage.SetBodyString("text/html", string(message.HtmlBody))
		gMessage.AddAlternativeString("text/plain", string(message.PlainBody))
		for _, attachment := range message.Attachments {
			gMessage.AttachFile(attachment)
		}
		messages = append(messages, gMessage)
	}
	err = sender.Send(messages...)
	if err != nil {
		return err
	}
	return nil
}
