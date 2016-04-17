package static

import "io"

type WriteFunc func(w io.Writer, path string) error
