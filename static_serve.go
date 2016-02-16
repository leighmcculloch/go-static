package static

import (
	"fmt"
	"net/http"
	"strings"
)

func (s *Static) ListenAndServe(addr string, ev EventHandler) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.serverHandleRequest(w, r, ev)
	})
	return http.ListenAndServe(addr, mux)
}

func (s *Static) serverHandleRequest(w http.ResponseWriter, r *http.Request, ev EventHandler) {
	path := r.URL.Path

	filePath := path
	if strings.HasSuffix(path, "/") {
		filePath = fmt.Sprintf("%sindex.html", path)
	}

	err := s.handleRequest(w, filePath, true)

	if err == errNotFound {
		fileServer := http.FileServer(http.Dir(s.BuildDir))
		fileServer.ServeHTTP(w, r)
		ev(Event{Action: "file", Path: path})
	} else {
		ev(Event{Action: "build", Path: filePath, Error: err})
	}
}
