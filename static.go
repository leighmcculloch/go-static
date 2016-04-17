package static

import (
	"errors"
	"fmt"
	"io"
	"log"
)

var (
	ErrNotFound = errors.New("No page registered for path.")
)

type Static struct {
	// Pages registered
	Pages map[string]*Page
}

// Create a new Static with defaults.
func New() Static {
	return Static{
		Pages: make(map[string]*Page),
	}
}

// Register a page with a relative path and function to call when the page is served or built that will write the page.
func (s Static) AddPage(path string, writeFunc WriteFunc) {
	s.Pages[path] = &Page{
		Path:      path,
		WriteFunc: writeFunc,
	}
}

func (s Static) RenderEventHandler(r Renderer, ev EventHandler) error {
	return r.Render(s, ev)
}

func (s Static) Render(r Renderer) error {
	return s.RenderEventHandler(r, logEvent)
}

func (s Static) WritePage(w io.Writer, path string, ignoreCache bool) error {
	p := s.getPageForPath(path)
	if p == nil {
		return ErrNotFound
	}
	return p.WriteFunc(w, p.Path)
}

func (s Static) getPageForPath(path string) *Page {
	return s.Pages[path]
}

func logEvent(event Event) {
	var s string
	if event.Error == nil {
		s = fmt.Sprintf("%10s  %-20s", event.Action, event.Path)
	} else {
		s = fmt.Sprintf("%10s  %-20s  %v", "error", event.Path, event.Error)
	}
	log.Println(s)
}
