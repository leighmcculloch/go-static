package templates

import (
	"crypto/sha1"
	"fmt"
	"html/template"
	"io"
)

type Templates struct {
	TemplateFuncs template.FuncMap
	tmpls         map[string]*template.Template
}

func New(templateFuncs template.FuncMap) Templates {
	return Templates{
		TemplateFuncs: templateFuncs,
		tmpls:         make(map[string]*template.Template),
	}
}

func (t Templates) Render(w io.Writer, data interface{}, tmplPaths []string, tmpl string, ignoreCache bool) error {
	tmpls, err := t.getTmpl(tmplPaths, ignoreCache)
	if err != nil {
		return err
	}
	return tmpls.ExecuteTemplate(w, tmpl, data)
}

func (t Templates) getTmpl(tmplPaths []string, cache bool) (*template.Template, error) {
	h := hash(tmplPaths)

	tmpl := t.tmpls[h]
	if cache && tmpl != nil {
		return tmpl, nil
	}

	tmpl, err := template.New("all").Funcs(t.TemplateFuncs).ParseFiles(tmplPaths...)
	if err != nil {
		return nil, err
	}

	if cache {
		t.tmpls[h] = tmpl
	}

	return tmpl, nil
}

func hash(tmplPaths []string) string {
	h := sha1.New()
	for _, t := range tmplPaths {
		io.WriteString(h, t)
	}
	return fmt.Sprintf("% x", h.Sum(nil))
}
