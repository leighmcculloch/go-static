package static

import (
	"net/http"
	"os"
	"strings"
	"sync"
)

func Build(o Options, h http.Handler, paths []string, eh EventHandler) {
	var wg sync.WaitGroup

	pathsChan := make(chan string)

	for i := 0; i < o.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buildPaths(o, h, pathsChan, eh)
		}()
	}

	for _, path := range paths {
		pathsChan <- path
	}

	close(pathsChan)

	wg.Wait()
}

func buildPaths(o Options, h http.Handler, paths <-chan string, eh EventHandler) {
	for path := range paths {
		err := BuildSingle(o, h, path)
		if eh != nil {
			eh(Event{Action: "build", Path: path, Error: err})
		}
	}
}

func BuildSingle(o Options, h http.Handler, path string) error {
	_, err := os.Stat(o.OutputDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(o.OutputDir, 0755)
	}
	if err != nil {
		return err
	}

	fp := o.OutputDir + path
	if strings.HasSuffix(fp, "/") {
		fp += o.DirFilename
	}

	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	r, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}
	rw := newResponseWriter(f)
	h.ServeHTTP(&rw, r)

	return nil
}
