package mailer

import (
	"io"

	"github.com/pablor21/goms/pkg/server/views"
)

type MailView interface {
	views.View
}

type MailHtmlView struct {
	*views.HtmlView
	Preheader  string
	FooterText string
	HeroImage  string
}

func NewMailHtmlView() *MailHtmlView {
	ret := &MailHtmlView{
		HtmlView:   views.NewHtmlView(),
		Preheader:  "",
		FooterText: "Hello, Brands.com",
		HeroImage:  "",
	}
	return ret
}

func (v *MailHtmlView) SetPreheader(preheader string) {
	v.Preheader = preheader
}

func (v *MailHtmlView) GetPreheader() string {
	return v.Preheader
}

func (v *MailHtmlView) SetFooterText(footerText string) {
	v.FooterText = footerText
}

func (v *MailHtmlView) GetFooterText() string {
	return v.FooterText
}
func (v *MailHtmlView) RenderToByte() ([]byte, error) {
	toRender := v.Layout()
	if toRender == "" {
		toRender = v.Template()
	}
	return v.Renderer().RenderToByte(v.Context(), toRender, v)
}

func (v *MailHtmlView) RenderToWriter(writer io.Writer) error {
	toRender := v.Layout()
	if toRender == "" {
		toRender = v.Template()
	}

	return v.Renderer().RenderToWriter(v.Context(), writer, toRender, v)
}
