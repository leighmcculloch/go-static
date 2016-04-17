package static

import (
	"bytes"
	"io"
	"net/http"
)

type responseBuffer struct {
	buffer *bytes.Buffer
	header http.Header
}

func newResponseBuffer() responseBuffer {
	return responseBuffer{
		buffer: new(bytes.Buffer),
		header: make(map[string][]string),
	}
}

func (rc responseBuffer) Header() http.Header {
	return rc.header
}

func (rc responseBuffer) Write(p []byte) (int, error) {
	return rc.buffer.Write(p)
}

func (rc responseBuffer) WriteHeader(int) {
}

func (rc responseBuffer) WriteTo(w io.Writer) (n int64, err error) {
	return rc.buffer.WriteTo(w)
}
