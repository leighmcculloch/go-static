package static

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Build the paths. Uses the http.Handler to get the response for each path, and writes that response to a file with it's respective path in the OutputDir specified in the Options. Does so concurrently as defined in the Options, and calls the EventHandler for every path with an Event that states that the path was built and if an error occurred. EventHandler may be nil.
func Build(o Options, h http.Handler, paths []string, eh EventHandler) {
	if eh == nil {
		eh = defaultEventHandler
	}

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
		statusCode, outputPath, err := BuildSingle(o, h, path)
		eh(Event{Action: "build", StatusCode: statusCode, Path: path, OutputPath: outputPath, Error: err})
	}
}

// BuildSingle builds a single path. It uses the http.Handler to get the response for each path, and writes that response to a file with it's respective path in the OutputDir specified in the Options. Returns the HTTP status code returned by the handler, the output path written to and an error if one occurs.
func BuildSingle(o Options, h http.Handler, path string) (statusCode int, outputPath string, err error) {
	pathIsDir := strings.HasSuffix(path, "/")

	filePath := filepath.FromSlash(path)
	if pathIsDir {
		filePath = filepath.Join(filePath, o.DirFilename)
	}

	outputPath = filepath.Join(o.OutputDir, filePath)
	outputDir := filepath.Dir(outputPath)
	_, err = os.Stat(outputDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, 0755)
	}
	if err != nil {
		message := fmt.Sprintf("Unable to create dir %s for path %s", outputDir, path)
		return 0, "", buildError{message, err}
	}

	f, err := os.Create(outputPath)
	if err != nil {
		message := fmt.Sprintf("Unable to create file %s for path %s", outputPath, path)
		return 0, "", buildError{message, err}
	}
	defer f.Close()

	r, err := http.NewRequest("GET", path, nil)
	if err != nil {
		message := fmt.Sprintf("Unable to create http.Request for path %s", path)
		return 0, "", buildError{message, err}
	}
	rw := newResponseWriter(f)
	h.ServeHTTP(&rw, r)

	return rw.StatusCode(), outputPath, nil
}
