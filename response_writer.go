package static

import (
	"io"
	"net/http"
)

type responseWriter struct {
	writer    io.Writer
	header    http.Header
	status    int
	statusSet bool
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

func (rc *responseWriter) WriteHeader(code int) {
	rc.status = code
	rc.statusSet = true
}

func (rc *responseWriter) Status() int {
	return rc.status
}

func (rc *responseWriter) Write(p []byte) (n int, err error) {
	if !rc.statusSet {
		rc.status = http.StatusOK
	}
	return rc.writer.Write(p)
}
