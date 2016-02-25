package static

import "io"

type Page struct {
	Path string
	Func PageFunc
}

type PageFunc func(w io.Writer, path string) error
