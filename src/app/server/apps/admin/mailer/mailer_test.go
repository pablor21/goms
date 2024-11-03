package mailer_test

import (
	"context"
	"testing"

	"github.com/pablor21/goms/app/config"
	admin_mailer "github.com/pablor21/goms/app/server/apps/admin/mailer"
	"github.com/pablor21/goms/pkg/mailer"
)

func TestHtmlRenderer_RenderToByte(t *testing.T) {
	config.InitConfig([]string{"config.yml"})
	view := mailer.NewMailHtmlView()
	view.SetTemplate("/forgot-password")
	view.SetRenderer(admin_mailer.GetHtmlRenderer())
	view.SetData(struct {
		Code string
	}{
		Code: "1111",
	})

	mailMessage := mailer.MailMessage{
		To: []mailer.MailAddress{
			{Name: "John Doe", Email: "john@doe.com"},
		},
		Subject: "One Time Password",
	}
	err := mailer.GetDefaultMailer().SendTemplate(context.Background(), view, mailMessage)
	if err != nil {
		t.Error(err)
	}
	t.Log("Mail sent")
}
