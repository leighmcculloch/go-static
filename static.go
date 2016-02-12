package static

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"path"
	"strings"
)

type Static struct {
	SourceDir        string
	BuildDir         string
	BuildConcurrency int
	ServerPort       int

	TemplateFuncs template.FuncMap

	*static
}

type static struct {
	templates map[string]*template.Template
	pages     map[string]*Page
}

func New() Static {
	return Static{
		SourceDir:        "source",
		BuildDir:         "build",
		BuildConcurrency: 50,
		ServerPort:       4567,
		TemplateFuncs: template.FuncMap{
			"UnsafeHTML": unsafeHTML,
			"ToLower":    strings.ToLower,
			"ToUpper":    strings.ToUpper,
		},
		static: &static{
			templates: make(map[string]*template.Template),
			pages:     make(map[string]*Page),
		},
	}
}

func (s Static) Page(path string, pageFunc PageFunc) {
	s.pages[path] = &Page{
		Path: path,
		Func: pageFunc,
	}
}

func (s Static) BuildPage(path string) error {
	fp := fmt.Sprintf("%s%s", s.BuildDir, path)
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()
	err = s.handleRequest(f, path, false)
	return err
}

var errNotFound = errors.New("No handler for path")

func (s Static) handleRequest(w io.Writer, path string, ignoreCache bool) error {
	page := s.getPageForPath(path)
	if page == nil {
		return errNotFound
	}
	data, tmpls, tmpl := page.Func(path)
	err := s.writeResponse(w, data, tmpls, tmpl, ignoreCache)
	return err
}

func (s Static) getPageForPath(path string) *Page {
	return s.pages[path]
}

func (s Static) expandTemplatePaths(tmpl []string) []string {
	expandedTmpl := make([]string, len(tmpl))
	for i, t := range tmpl {
		expandedTmpl[i] = s.templatePath(t)
	}
	return expandedTmpl
}

func getTemplateCacheHash(tmpl []string) string {
	h := sha1.New()
	for _, t := range tmpl {
		io.WriteString(h, t)
	}
	return fmt.Sprintf("% x", h.Sum(nil))
}

func (s Static) getTemplates(tmpl []string, ignoreCache bool) (*template.Template, error) {
	var err error
	tmpl = s.expandTemplatePaths(tmpl)
	h := getTemplateCacheHash(tmpl)
	if s.templates[h] == nil || ignoreCache {
		s.templates[h], err = template.New("all").Funcs(s.TemplateFuncs).ParseFiles(tmpl...)
	}
	return s.templates[h], err
}

func (s Static) writeResponse(w io.Writer, data interface{}, tmpls []string, tmpl string, ignoreCache bool) error {
	templates, err := s.getTemplates(tmpls, ignoreCache)
	if err != nil {
		return err
	}
	return templates.ExecuteTemplate(w, tmpl, data)
}

func (s Static) templatePath(filename string) string {
	return path.Join(s.SourceDir, filename)
}

func unsafeHTML(s string) template.HTML {
	return template.HTML(s)
}
