package static

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (s *Static) ListenAndServe(addr string) error {
	fileServer := http.FileServer(http.Dir(s.BuildDir))
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case path == "/":
			path = "/index.html"
		case !strings.HasSuffix(path, ".html"):
			path = fmt.Sprintf("%s.html", path)
		}
		err := s.handleRequest(w, path, true)
		if err != nil {
			if err == errNotFound {
				fileServer.ServeHTTP(w, r)
			} else {
				log.Fatal(err)
			}
		}
	})
	return http.ListenAndServe(addr, mux)
}
