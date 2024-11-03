package mailer

import (
	"context"
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/pablor21/goms/pkg/logger"
	"github.com/pablor21/goms/pkg/server/views"
)

var htmlRenderer views.ViewRenderer

//go:embed all:templates
var templateFiles embed.FS
var templatesFolder = "templates"

type HtmlRenderer struct {
	templates *template.Template
}

func NewHtmlRenderer() views.ViewRenderer {
	return &HtmlRenderer{}
}

func GetHtmlRenderer() views.ViewRenderer {
	if htmlRenderer == nil {
		htmlRenderer = NewHtmlRenderer()
	}
	return htmlRenderer
}

func (r *HtmlRenderer) RenderToWriter(ctx context.Context, writer io.Writer, name string, data interface{}) error {
	return r.templates.ExecuteTemplate(writer, name, data)
}

func (r *HtmlRenderer) RenderToByte(ctx context.Context, name string, data interface{}) (res []byte, err error) {
	err = r.LoadTemplates(templatesFolder)
	if err != nil {
		return
	}
	buf := new(strings.Builder)
	err = r.templates.ExecuteTemplate(buf, name, data)
	if err != nil {
		return nil, err
	}
	return []byte(buf.String()), nil
}

func (r *HtmlRenderer) Clone() (views.ViewRenderer, error) {
	clone, err := r.templates.Clone()
	if err != nil {
		return nil, err
	}
	return &HtmlRenderer{
		templates: clone,
	}, nil
}

func (r *HtmlRenderer) GetFuncMap() template.FuncMap {
	return template.FuncMap{
		// "safeHTML": func(s string) template.HTML {
		// 	return template.HTML(s)
		// },
		"nl2br": func(s string) template.HTML {
			return template.HTML(strings.Replace(s, "\n", "<br>", -1))
		},
		"time":      time.Now,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		"json": func(v interface{}) (string, error) {
			b, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			return string(b), nil
		},
		"html": func(s string) template.HTML {
			return template.HTML(s)
		},
		"defaultValue": func(value, def interface{}) interface{} {
			if value == nil {
				return def
			}
			return value
		},
	}
}

func (r *HtmlRenderer) LoadTemplates(templatesDir string) (err error) {
	if r.templates != nil {
		return
	}

	cleanRoot := filepath.Clean(templatesDir)
	pfx := len(cleanRoot) + 1
	root := template.New("/").Funcs(r.GetFuncMap())
	err = fs.WalkDir(templateFiles, cleanRoot, func(path string, d fs.DirEntry, e error) (err error) {
		if e != nil {
			return e
		}
		if !d.IsDir() {
			if strings.HasSuffix(path, ".html") {
				b, e2 := fs.ReadFile(templateFiles, path)
				if e2 != nil {
					return e2
				}
				name := "/" + path[pfx:]
				name = strings.TrimSuffix(name, ".html")
				logger.Debug().Str("path", path).Str("name", name).Msg("Loading template")
				_, err := root.Parse(string(b))
				if err != nil {
					return err
				}

			}
			r.templates = root
		}

		return nil
	})
	r.templates = root
	return
}

func (r *HtmlRenderer) GetTemplates() *template.Template {
	return r.templates
}

func (r *HtmlRenderer) SetTemplates(templates *template.Template) {
	r.templates = templates
}

func (r *HtmlRenderer) AddTemplate(name string, tmpl *template.Template) {
	r.templates.AddParseTree(name, tmpl.Tree)
}
