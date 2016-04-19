package static

import (
	"bytes"
	"net/http"
)

type responseBuffer struct {
	bytes.Buffer
	header http.Header
}

func newResponseBuffer() responseBuffer {
	return responseBuffer{}
}

func (rc *responseBuffer) Header() http.Header {
	if rc.header == nil {
		rc.header = make(map[string][]string)
	}
	return rc.header
}

func (rc *responseBuffer) WriteHeader(int) {
}
