package static

import (
	"net/http"
	"os"
	"sync"
)

func Build(h http.Handler, paths []string, eh EventHandler) error {
	return BuildOptions(NewOptions(), h, paths, eh)
}

func BuildOptions(o Options, h http.Handler, paths []string, eh EventHandler) error {
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

	return nil
}

func buildPaths(o Options, h http.Handler, paths <-chan string, eh EventHandler) {
	for path := range paths {
		err := buildPath(o, h, path)
		eh(Event{Action: "build", Path: path, Error: err})
	}
}

func buildPath(o Options, h http.Handler, path string) error {
	fp := o.OutputDir + path
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	r, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}
	rw := newResponseBuffer()
	h.ServeHTTP(rw, r)
	rw.WriteTo(f)

	return nil
}
