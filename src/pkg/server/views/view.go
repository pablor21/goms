package views

import (
	"context"
	"fmt"
	"html/template"
	"io"
)

type View interface {
	SetTemplate(name string)
	Template() string
	SetLayout(layout string)
	Layout() string
	SetRenderer(renderer ViewRenderer)
	Renderer() ViewRenderer
	ExecuteTemplate(template string, data interface{}) (template.HTML, error)
	SafeAttribute(attr string, value interface{}) template.HTMLAttr
	SetContext(ctx context.Context)
	Context() context.Context
	RenderToByte() ([]byte, error)
	RenderToWriter(writer io.Writer) error
	Data() interface{}
	SetData(data interface{})
}

type HtmlView struct {
	template string
	layout   string
	renderer ViewRenderer
	ctx      context.Context
	data     interface{}
}

func NewHtmlView() *HtmlView {
	return &HtmlView{
		template: "/index",
		layout:   "/layouts/default",
	}
}

func (v *HtmlView) SetData(data interface{}) {
	v.data = data
}

func (v *HtmlView) Data() interface{} {
	return v.data
}

func (v *HtmlView) SetTemplate(name string) {
	v.template = name
}

func (v *HtmlView) Template() string {
	return v.template
}

func (v *HtmlView) SetLayout(layout string) {
	v.layout = layout
}

func (v *HtmlView) Layout() string {
	return v.layout
}

func (v *HtmlView) SetRenderer(renderer ViewRenderer) {
	v.renderer = renderer
}

func (v *HtmlView) Renderer() ViewRenderer {
	return v.renderer
}

func (v *HtmlView) ExecuteTemplate(t string, data interface{}) (template.HTML, error) {
	res, err := v.Renderer().RenderToByte(v.ctx, t, data)
	if err != nil {
		return "", err
	}

	return template.HTML(res), nil
}

func (v *HtmlView) SafeAttribute(attr string, value interface{}) template.HTMLAttr {
	return template.HTMLAttr(fmt.Sprintf(`%v="%v"`, attr, value))
}

func (v *HtmlView) SetContext(ctx context.Context) {
	v.ctx = ctx
}

func (v *HtmlView) Context() context.Context {
	return v.ctx
}

func (v *HtmlView) RenderToByte() ([]byte, error) {
	toRender := v.Layout()
	if toRender == "" {
		toRender = v.Template()
	}
	return v.Renderer().RenderToByte(v.ctx, toRender, v)
}

func (v *HtmlView) RenderToWriter(writer io.Writer) error {
	toRender := v.Layout()
	if toRender == "" {
		toRender = v.Template()
	}

	return v.Renderer().RenderToWriter(v.ctx, writer, toRender, v)
}
