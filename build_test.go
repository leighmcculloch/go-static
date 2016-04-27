package static

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func ExampleBuild() {
	handler := http.NewServeMux()
	paths := []string{}

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", path.Base(r.URL.Path))
	})

	paths = append(paths, "/world")

	Build(DefaultOptions, handler, paths, func(e Event) {
		fmt.Println(e)
	})

	// Output:
	// Action: build, StatusCode: 200, Path: /world
}

func ExampleBuildSingle() {
	handler := http.NewServeMux()

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", path.Base(r.URL.Path))
	})

	status, err := BuildSingle(DefaultOptions, handler, "/world")
	fmt.Printf("Built: /world, StatusCode: %d, Error: %v", status, err)

	// Output:
	// Built: /world, StatusCode: 200, Error: <nil>
}

func TestBuildSingle(t *testing.T) {
	t.Log("When a Handler is defined to respond to /* and response with Hello <path>!")
	handler := http.NewServeMux()
	handler.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", path.Base(r.URL.Path))
	})

	t.Log("And Options are defined with defaults and an OutputDir that does not exist.")
	options := DefaultOptions
	tempDir, _ := ioutil.TempDir("", "")
	options.OutputDir = filepath.Join(tempDir, "build")
	t.Logf("OutputDir => %s", options.OutputDir)

	t.Log("And the path to build is /hello/world.")
	path := "/hello/world"

	t.Log("Expect BuildSingle to create the output path and write a file world with contents Hello world! And return no error.")
	expectedStatus := 200
	expectedOutputFilePath := filepath.Join(options.OutputDir, "hello", "world")
	expectedOutputFileContents := "Hello world!"

	status, err := BuildSingle(options, handler, path)
	t.Logf("BuildSingle(%#v) => %v, %v", path, status, err)
	if status != expectedStatus || err != nil {
		t.Errorf("BuildSingle(%#v) => %v, %v, expected %v, nil", path, status, err, expectedStatus)
	}

	outputFileContents, err := ioutil.ReadFile(expectedOutputFilePath)
	if err != nil {
		t.Fatalf("Expected %s to exist with the output but got error when opening: %v", expectedOutputFilePath, err)
	}

	t.Logf("Contents of %s => %s", expectedOutputFilePath, outputFileContents)
	if string(outputFileContents) != expectedOutputFileContents {
		t.Fatalf(`Contents of %s => %s, expected %s`, expectedOutputFilePath, outputFileContents, expectedOutputFileContents)
	}
}

func TestBuildSingleErrors(t *testing.T) {
	t.Log("When a Handler is defined to respond to /* and response with Hello <path>!")
	handler := http.NewServeMux()
	handler.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", path.Base(r.URL.Path))
	})

	t.Log("And Options are defined with defaults and an OutputDir.")
	options := DefaultOptions
	tempDir, _ := ioutil.TempDir("", "")
	options.OutputDir = filepath.Join(tempDir, "build")
	t.Logf("OutputDir => %s", options.OutputDir)

	t.Log("And a file exists at the same path as the OutputDir (a problem).")
	f, _ := os.Create(options.OutputDir)
	defer f.Close()

	t.Log("And the path to build is /hello/world.")
	path := "/hello/world"

	t.Log("Expect BuildSingle to error with not a directory error.")

	expectedStatus := 0
	expectedErrString := "not a directory"
	status, err := BuildSingle(options, handler, path)
	if err != nil && strings.Contains(err.Error(), expectedErrString) {
		t.Logf("BuildSingle(%#v) => %v, %v", path, status, err)
	} else {
		t.Errorf("BuildSingle(%#v) => %v, %v, expected %v and a %s error", path, err, expectedStatus, expectedErrString)
	}
}

func TestBuild(t *testing.T) {
	t.Log("When a Handler is defined to respond to /* and response with Hello <path>!")
	handler := http.NewServeMux()
	handler.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", path.Base(r.URL.Path))
	})

	t.Log("And Options are defined with defaults and an OutputDir that does not exist.")
	options := DefaultOptions
	tempDir, _ := ioutil.TempDir("", "")
	options.OutputDir = filepath.Join(tempDir, "build")
	t.Logf("OutputDir: %s", options.OutputDir)

	t.Log("And there are multiple paths to build.")
	paths := []string{
		"/hello/go",
		"/hello/world",
		"/hello/universe",
		"/bye",
	}

	t.Log("Expect Build to create the output path and write a file for each path with contents Hello <file>! And send one event per path to the EventHandler containing no error.")
	expected := []struct {
		OutputFilePath     string
		OutputFileContents string
		Event              Event
	}{
		{
			filepath.Join(options.OutputDir, "hello", "go"),
			"Hello go!",
			Event{"build", 200, "/hello/go", nil},
		},
		{
			filepath.Join(options.OutputDir, "hello", "world"),
			"Hello world!",
			Event{"build", 200, "/hello/world", nil},
		},
		{
			filepath.Join(options.OutputDir, "hello", "universe"),
			"Hello universe!",
			Event{"build", 200, "/hello/universe", nil},
		},
		{
			filepath.Join(options.OutputDir, "bye"),
			"404 page not found\n",
			Event{"build", 404, "/bye", nil},
		},
	}

	expectedNumberOfEvents := len(paths)
	numberOfEvents := 0
	events := make(chan Event, expectedNumberOfEvents)
	Build(options, handler, paths, func(e Event) {
		select {
		case events <- e:
			t.Logf("Event received => %#v", e)
		default:
			t.Errorf("Additional event received => %#v", numberOfEvents, expectedNumberOfEvents)
		}
		numberOfEvents += 1
	})
	close(events)
	if numberOfEvents != expectedNumberOfEvents {
		t.Errorf("Number of events received => %d, expected %d", numberOfEvents, expectedNumberOfEvents)
	}

	eventsSeen := make(map[Event]bool)
	for event := range events {
		if eventsSeen[event] {
			t.Errorf("Event received => %#v, multiple times but was expected once.", event)
		}
		eventsSeen[event] = true
	}

	for _, expect := range expected {
		if !eventsSeen[expect.Event] {
			t.Errorf("Event not received => %#v, but was expected once.", expect.Event)
		}

		outputFileContents, err := ioutil.ReadFile(expect.OutputFilePath)
		if err != nil {
			t.Fatalf("Error opening output file => %#v, expected to exist.", err)
		}

		if string(outputFileContents) == expect.OutputFileContents {
			t.Logf("Contents of %s => %s", expect.OutputFilePath, outputFileContents)
		} else {
			t.Fatalf(`Contents of %s => %s, expected %s`, expect.OutputFilePath, outputFileContents, expect.OutputFileContents)
		}
	}
}

func TestBuildErrors(t *testing.T) {
	t.Log("When a Handler is defined to respond to /* and response with Hello <path>!")
	handler := http.NewServeMux()
	handler.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", path.Base(r.URL.Path))
	})

	t.Log("And Options are defined with defaults and an OutputDir that does not exist.")
	options := DefaultOptions
	tempDir, _ := ioutil.TempDir("", "")
	options.OutputDir = filepath.Join(tempDir, "build")
	t.Logf("OutputDir: %s", options.OutputDir)

	t.Log("And a file exists at the same path as the OutputDir (a problem).")
	f, _ := os.Create(options.OutputDir)
	defer f.Close()

	t.Log("And there are multiple paths to build.")
	paths := []string{
		"/hello/go",
		"/hello/world",
		"/hello/universe",
	}

	t.Log("Expect Build to send one event per path to the EventHandler containing a not a directory error.")

	expectedNumberOfEvents := len(paths)
	numberOfEvents := 0
	events := make(chan Event, expectedNumberOfEvents)
	Build(options, handler, paths, func(e Event) {
		select {
		case events <- e:
			t.Logf("Event received => %#v", e)
		default:
			t.Errorf("Additional event received => %#v", numberOfEvents, expectedNumberOfEvents)
		}
		numberOfEvents += 1
	})
	close(events)
	if numberOfEvents != expectedNumberOfEvents {
		t.Errorf("Number of events received => %d, expected %d", numberOfEvents, expectedNumberOfEvents)
	}

	eventsSeen := make(map[string]*Event)
	for event := range events {
		storeEvent := event
		if eventsSeen[event.Path] != nil {
			t.Errorf("Event received => %#v, multiple times but was expected once.", event)
		}
		eventsSeen[event.Path] = &storeEvent
	}

	for _, path := range paths {
		event := eventsSeen[path]
		if event == nil {
			t.Errorf("Event not received for path => %s, but was expected once.", path)
		}

		expectedAction := "build"
		if event.Action != expectedAction {
			t.Errorf("Event for %s Action => %s, expected %s.", path, event.Action, expectedAction)
		}

		if event.Path != path {
			t.Errorf("Event for %s Path => %s, expected %s.", path, event.Path, path)
		}

		expectedErrString := "not a directory"
		if event.Error != nil && strings.Contains(event.Error.Error(), expectedErrString) {
			t.Logf("Build(%#v) => %v", path, event.Error)
		} else {
			t.Errorf("Build(%#v) => %v, expected a %s error", path, event.Error, expectedErrString)
		}
	}
}

func TestBuildWithNilEventHandler(t *testing.T) {
	t.Log("When a Handler is defined to respond to /* and response with Hello <path>!")
	handler := http.NewServeMux()
	handler.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", path.Base(r.URL.Path))
	})

	t.Log("And Options are defined with defaults and an OutputDir that does not exist.")
	options := DefaultOptions
	tempDir, _ := ioutil.TempDir("", "")
	options.OutputDir = filepath.Join(tempDir, "build")
	t.Logf("OutputDir: %s", options.OutputDir)

	t.Log("And there are multiple paths to build.")
	paths := []string{
		"/hello/go",
		"/hello/world",
		"/hello/universe",
	}

	t.Log("Expect Build to create the output path and write a file for each path with contents Hello <file>! And send one event per path to the EventHandler containing no error.")
	expected := []struct {
		OutputFilePath     string
		OutputFileContents string
	}{
		{
			filepath.Join(options.OutputDir, "hello", "go"),
			"Hello go!",
		},
		{
			filepath.Join(options.OutputDir, "hello", "world"),
			"Hello world!",
		},
		{
			filepath.Join(options.OutputDir, "hello", "universe"),
			"Hello universe!",
		},
	}

	Build(options, handler, paths, nil)

	for _, expect := range expected {
		outputFileContents, err := ioutil.ReadFile(expect.OutputFilePath)
		if err != nil {
			t.Fatalf("Error opening output file => %#v, expected to exist.", err)
		}

		if string(outputFileContents) == expect.OutputFileContents {
			t.Logf("Contents of %s => %s", expect.OutputFilePath, outputFileContents)
		} else {
			t.Fatalf(`Contents of %s => %s, expected %s`, expect.OutputFilePath, outputFileContents, expect.OutputFileContents)
		}
	}
}
