package build

import (
	"fmt"
	"os"
	"sync"

	"github.com/leighmcculloch/static"
)

const (
	defaultBuildDir         = "build"
	defaultBuildConcurrency = 50
)

type Builder struct {
	// The directory where files will be written when building.
	BuildDir string
	// The number of files that will be built concurrently, default 50.
	BuildConcurrency int
}

func NewBuilder() Builder {
	return Builder{
		BuildDir:         defaultBuildDir,
		BuildConcurrency: defaultBuildConcurrency,
	}
}

func (b Builder) Start(s static.Static, ev static.EventHandler) error {
	var wg sync.WaitGroup
	var wgEventHandling sync.WaitGroup

	paths := make(chan string)
	buildEvents := make(chan static.Event)

	for i := 0; i < b.BuildConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.buildPaths(s, paths, buildEvents)
		}()
	}

	wgEventHandling.Add(1)
	go func() {
		defer wgEventHandling.Done()
		for event := range buildEvents {
			ev(event)
		}
	}()

	for path := range s.Pages {
		paths <- path
	}

	close(paths)

	wg.Wait()

	close(buildEvents)

	wgEventHandling.Wait()

	return nil
}

func (b Builder) buildPaths(s static.Static, paths <-chan string, buildEvents chan<- static.Event) {
	for path := range paths {
		err := b.buildPage(s, path)
		buildEvents <- static.Event{Action: "build", Path: path, Error: err}
	}
}

func (b Builder) buildPage(s static.Static, path string) error {
	fp := fmt.Sprintf("%s%s", b.BuildDir, path)
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()
	return s.WritePage(f, path, true)
}
