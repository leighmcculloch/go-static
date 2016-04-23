package static_test

import (
	"."
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"testing"
)

func ExampleBuild() {
	handler := http.NewServeMux()
	paths := []string{}

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", path.Base(r.URL.Path))
	})

	paths = append(paths, "/world")

	options := static.DefaultOptions()
	static.Build(options, handler, paths, func(e static.Event) {
		fmt.Println(e)
	})

	// Output:
	// Action: build, Path: /world
}

func ExampleBuildSingle() {
	handler := http.NewServeMux()

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", path.Base(r.URL.Path))
	})

	options := static.DefaultOptions()
	err := static.BuildSingle(options, handler, "/world")
	fmt.Println("Built: /world, Error:", err)

	// Output:
	// Built: /world, Error: <nil>
}

func TestBuildSingle(t *testing.T) {
	t.Log("When a Handler is defined to respond to /* and response with Hello <path>!")
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", path.Base(r.URL.Path))
	})

	t.Log("And Options are defined with defaults and an OutputDir that does not exist.")
	options := static.DefaultOptions()
	tempDir, _ := ioutil.TempDir("", "")
	options.OutputDir = filepath.Join(tempDir, "build")
	t.Logf("OutputDir => %s", options.OutputDir)

	t.Log("And the path to build is /world.")
	path := "/world"

	t.Log("Expect BuildSingle to create the OutputDir, write a file world with contents Hello world! And return no error.")
	expectedOutputFilePath := filepath.Join(options.OutputDir, "world")
	expectedOutputFileContents := "Hello world!"

	err := static.BuildSingle(options, handler, path)
	t.Logf("BuildSingle(..., %#v) => %v", path, err)
	if err != nil {
		t.Errorf("BuildSingle(..., %#v) => %v, expected nil", path, err)
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

func TestBuild(t *testing.T) {
	t.Log("When a Handler is defined to respond to /* and response with Hello <path>!")
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", path.Base(r.URL.Path))
	})

	t.Log("And Options are defined with defaults and an OutputDir that does not exist.")
	options := static.DefaultOptions()
	tempDir, _ := ioutil.TempDir("", "")
	options.OutputDir = filepath.Join(tempDir, "build")
	t.Logf("OutputDir: %s", options.OutputDir)

	t.Log("And there are multiple paths to build.")
	paths := []string{"/go", "/world", "/universe"}

	t.Log("Expect Build to create the OutputDir, write a file for each path with contents Hello <file>! And send one event per path to the EventHandler containing no error.")
	expected := []struct {
		OutputFilePath     string
		OutputFileContents string
		Event              static.Event
	}{
		{
			filepath.Join(options.OutputDir, "go"),
			"Hello go!",
			static.Event{"build", "/go", nil},
		},
		{
			filepath.Join(options.OutputDir, "world"),
			"Hello world!",
			static.Event{"build", "/world", nil},
		},
		{
			filepath.Join(options.OutputDir, "universe"),
			"Hello universe!",
			static.Event{"build", "/universe", nil},
		},
	}

	expectedNumberOfEvents := len(paths)
	numberOfEvents := 0
	events := make(chan static.Event, expectedNumberOfEvents)
	static.Build(options, handler, paths, func(e static.Event) {
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

	eventsSeen := make(map[static.Event]bool)
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

func TestBuildWithNilEventHandler(t *testing.T) {
	t.Log("When a Handler is defined to respond to /* and response with Hello <path>!")
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", path.Base(r.URL.Path))
	})

	t.Log("And Options are defined with defaults and an OutputDir that does not exist.")
	options := static.DefaultOptions()
	tempDir, _ := ioutil.TempDir("", "")
	options.OutputDir = filepath.Join(tempDir, "build")
	t.Logf("OutputDir: %s", options.OutputDir)

	t.Log("And there are multiple paths to build.")
	paths := []string{"/go", "/world", "/universe"}

	t.Log("Expect Build to create the OutputDir, write a file for each path with contents Hello <file>!")
	expected := []struct {
		OutputFilePath     string
		OutputFileContents string
	}{
		{filepath.Join(options.OutputDir, "go"), "Hello go!"},
		{filepath.Join(options.OutputDir, "world"), "Hello world!"},
		{filepath.Join(options.OutputDir, "universe"), "Hello universe!"},
	}

	static.Build(options, handler, paths, nil)

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
