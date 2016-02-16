package static

import (
	"errors"
	"html/template"
	"io"
	"strings"

	"github.com/leighmcculloch/static/templates"
)

type Static struct {
	// The directory where files will be written when building.
	BuildDir string
	// The number of files that will be built concurrently, default 50.
	BuildConcurrency int
	// The port served from when in serve mode, default 4567.
	ServerPort int

	// Functions available to templates.
	TemplateFuncs template.FuncMap

	// Pages registered
	pages map[string]*Page

	// Templates registered
	templates templates.Templates
}

// Create a new Static with defaults.
func New() Static {
	defaultTemplateFuncs := template.FuncMap{
		"UnsafeHTML": unsafeHTML,
		"ToLower":    strings.ToLower,
		"ToUpper":    strings.ToUpper,
	}
	return Static{
		BuildDir:         "build",
		BuildConcurrency: 50,
		ServerPort:       4567,
		TemplateFuncs:    defaultTemplateFuncs,
		pages:            make(map[string]*Page),
		templates:        templates.New(defaultTemplateFuncs),
	}
}

// Register a page with a relative path and function to call when the page is served or built.
func (s Static) Page(path string, pageFunc PageFunc) {
	s.pages[path] = &Page{
		Path: path,
		Func: pageFunc,
	}
}

var errNotFound = errors.New("No handler for path")

func (s Static) handleRequest(w io.Writer, path string, ignoreCache bool) error {
	p := s.getPageForPath(path)
	if p == nil {
		return errNotFound
	}
	data, tmplPaths, tmpl := p.Func(p.Path)
	return s.templates.Render(w, data, tmplPaths, tmpl, ignoreCache)
}

func (s Static) getPageForPath(path string) *Page {
	return s.pages[path]
}

func unsafeHTML(s string) template.HTML {
	return template.HTML(s)
}
