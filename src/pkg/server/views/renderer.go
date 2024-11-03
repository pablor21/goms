package views

import (
	"context"
	"html/template"
	"io"
)

type ViewRenderer interface {
	RenderToWriter(ctx context.Context, writer io.Writer, template string, data interface{}) error
	RenderToByte(ctx context.Context, template string, data interface{}) ([]byte, error)
	Clone() (ViewRenderer, error)
	LoadTemplates(templatesDir string) error
	GetTemplates() *template.Template
}
