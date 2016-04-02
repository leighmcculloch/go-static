package serve

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/leighmcculloch/static"
)

const (
	defaultStaticDir  = "build"
	defaultServerPort = 4567
)

type Server struct {
	// The directory where static resources can be found
	StaticDir string
	// The port served from when in serve mode, default 4567.
	ServerPort int
}

func NewServer() Server {
	return Server{
		StaticDir:  defaultStaticDir,
		ServerPort: defaultServerPort,
	}
}

func (v Server) Render(s static.Static, ev static.EventHandler) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		v.WritePage(w, r, s, ev)
	})
	addr := fmt.Sprintf(":%d", v.ServerPort)
	return http.ListenAndServe(addr, mux)
}

func (v Server) WritePage(w http.ResponseWriter, r *http.Request, s static.Static, ev static.EventHandler) {
	path := r.URL.Path

	filePath := path
	if strings.HasSuffix(path, "/") {
		filePath = fmt.Sprintf("%sindex.html", path)
	}

	err := s.WritePage(w, filePath, false)

	if err == static.ErrNotFound {
		fileServer := http.FileServer(http.Dir(v.StaticDir))
		fileServer.ServeHTTP(w, r)
		ev(static.Event{Action: "file", Path: path})
	} else {
		ev(static.Event{Action: "build", Path: filePath, Error: err})
	}
}
