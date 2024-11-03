package mailer

import (
	"context"

	"github.com/pablor21/goms/app/dtos"
	"github.com/pablor21/goms/pkg/mailer"
)

func SendOtp(ctx context.Context, mailData dtos.OTPResultMailData) (err error) {
	view := mailer.NewMailHtmlView()
	view.SetTemplate("/otp-code")
	view.SetRenderer(GetHtmlRenderer())
	view.SetData(mailData)
	mailMessage := mailer.MailMessage{
		To: []mailer.MailAddress{
			{Name: mailData.User.DisplayName(), Email: mailData.User.Email},
		},
		Subject: "One Time Password",
	}
	err = mailer.GetDefaultMailer().SendTemplate(context.Background(), view, mailMessage)
	return
}
