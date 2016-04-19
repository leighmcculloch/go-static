package static

import (
	"io"
	"net/http"
)

type responseWriter struct {
	writer io.Writer
	header http.Header
}

func newResponseWriter(w io.Writer) responseWriter {
	return responseWriter{writer: w}
}

func (rc *responseWriter) Header() http.Header {
	if rc.header == nil {
		rc.header = make(map[string][]string)
	}
	return rc.header
}

func (rc *responseWriter) WriteHeader(int) {
}

func (rc *responseWriter) Write(p []byte) (n int, err error) {
	return rc.writer.Write(p)
}
