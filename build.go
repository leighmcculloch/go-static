package static

import (
	"net/http"
	"os"
	pathutil "path"
	"strings"
	"sync"
)

// Build the paths. Uses the http.Handler to get the response for each path, and writes that response to a file with it's respective path in the OutputDir specified in the Options. Does so concurrently as defined in the Options, and calls the EventHandler for every path with an Event that states that the path was built and if an error occurred.
func Build(o Options, h http.Handler, paths []string, eh EventHandler) {
	var wg sync.WaitGroup

	pathsChan := make(chan string)

	for i := 0; i < o.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buildWorker(o, h, pathsChan, eh)
		}()
	}

	for _, path := range paths {
		pathsChan <- path
	}

	close(pathsChan)

	wg.Wait()
}

func buildWorker(o Options, h http.Handler, paths <-chan string, eh EventHandler) {
	for path := range paths {
		status, err := BuildSingle(o, h, path)
		if eh != nil {
			eh(Event{Action: "build", StatusCode: status, Path: path, Error: err})
		}
	}
}

// Build a single path. Uses the http.Handler to get the response for each path, and writes that response to a file with it's respective path in the OutputDir specified in the Options. Returns the HTTP status code return and an error if one occurs.
func BuildSingle(o Options, h http.Handler, path string) (int, error) {
	fp := o.OutputDir + path
	if strings.HasSuffix(fp, "/") {
		fp += o.DirFilename
	}

	dp := pathutil.Dir(fp)
	_, err := os.Stat(dp)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dp, 0755)
	}
	if err != nil {
		return 0, err
	}

	f, err := os.Create(fp)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	r, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return 0, err
	}
	rw := newResponseWriter(f)
	h.ServeHTTP(&rw, r)

	return rw.Status(), nil
}
