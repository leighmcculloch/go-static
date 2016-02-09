package static

import (
	"fmt"
	"net/http"
	"strings"
)

func (s *Static) ListenAndServe(addr string, eventHandler EventHandler) error {
	fileServer := http.FileServer(http.Dir(s.BuildDir))
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		var filePath string
		switch {
		case strings.HasSuffix(path, "/"):
			filePath = fmt.Sprintf("%sindex.html", path)
		case !strings.HasSuffix(path, ".html"):
			filePath = fmt.Sprintf("%s.html", path)
		}
		err := s.handleRequest(w, filePath, true)
		if err == errNotFound {
			fileServer.ServeHTTP(w, r)
			eventHandler(Event{Action: "file", Path: path})
		} else {
			eventHandler(Event{Action: "build", Path: filePath, Error: err})
		}
	})
	return http.ListenAndServe(addr, mux)
}
