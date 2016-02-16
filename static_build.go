package static

import (
	"fmt"
	"os"
	"sync"
)

func (s Static) Build(eventHandler EventHandler) error {
	var wg sync.WaitGroup
	var wgEventHandling sync.WaitGroup

	paths := make(chan string)
	buildEvents := make(chan Event)

	for i := 0; i < s.BuildConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.buildPaths(paths, buildEvents)
		}()
	}

	wgEventHandling.Add(1)
	go func() {
		defer wgEventHandling.Done()
		for event := range buildEvents {
			eventHandler(event)
		}
	}()

	for path := range s.pages {
		paths <- path
	}

	close(paths)

	wg.Wait()

	close(buildEvents)

	wgEventHandling.Wait()

	return nil
}

func (s Static) buildPaths(paths <-chan string, buildEvents chan<- Event) {
	for path := range paths {
		err := s.buildPage(path)
		buildEvents <- Event{Action: "build", Path: path, Error: err}
	}
}

func (s Static) buildPage(path string) error {
	fp := fmt.Sprintf("%s%s", s.BuildDir, path)
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()
	err = s.handleRequest(f, path, true)
	return err
}
