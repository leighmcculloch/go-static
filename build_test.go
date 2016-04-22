package static_test

import (
	"."
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
)

func ExampleBuild() {
	handler := http.NewServeMux()
	paths := []string{}

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", r.URL.Path)
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
		fmt.Fprintf(w, "Hello %s!", r.URL.Path)
	})

	options := static.DefaultOptions()
	err := static.BuildSingle(options, handler, "/world")
	fmt.Println("Built: /world, Error:", err)

	// Output:
	// Built: /world, Error: <nil>
}

func TestBuildSingle(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", r.URL.Path)
	})

	options := static.DefaultOptions()
	tempDir, _ := ioutil.TempDir("", "test_build_single")
	t.Logf("Created temp dir for OutputDir: %s", tempDir)
	options.OutputDir = filepath.Join(tempDir, "build")
	t.Logf("OutputDir: %s", options.OutputDir)

	path := "/world"
	expectedOutputFilePath := filepath.Join(options.OutputDir, "world")
	expectedOutputFileContents := "Hello /world!"

	err := static.BuildSingle(options, handler, path)
	t.Logf("BuildSingle(..., %#v) => %v", path, err)
	if err != nil {
		t.Errorf("BuildSingle(..., %#v) => %v, want nil", path, err)
	}

	outputFileContents, err := ioutil.ReadFile(expectedOutputFilePath)
	if err != nil {
		t.Fatalf("Expected %s to exist with the output but got error when opening: %v", expectedOutputFilePath, err)
	}

	t.Logf("Contents of %s => %s", expectedOutputFilePath, outputFileContents)
	if string(outputFileContents) != expectedOutputFileContents {
		t.Fatalf(`Contents of %s => %s, want %s`, expectedOutputFilePath, outputFileContents, expectedOutputFileContents)
	}
}

func TestBuild(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", r.URL.Path)
	})

	options := static.DefaultOptions()
	tempDir, _ := ioutil.TempDir("", "test_build")
	options.OutputDir = filepath.Join(tempDir, "build")
	t.Logf("OutputDir: %s", options.OutputDir)

	paths := []string{"/go", "/world", "/universe"}
	expected := []struct {
		OutputFilePath     string
		OutputFileContents string
		Event              static.Event
	}{
		{filepath.Join(options.OutputDir, "go"), "Hello /go!", static.Event{"build", "/go", nil}},
		{filepath.Join(options.OutputDir, "world"), "Hello /world!", static.Event{"build", "/world", nil}},
		{filepath.Join(options.OutputDir, "universe"), "Hello /universe!", static.Event{"build", "/universe", nil}},
	}

	expectedNumberOfEvents := len(paths)
	events := make(chan static.Event, expectedNumberOfEvents)
	static.Build(options, handler, paths, func(e static.Event) {
		select {
		case events <- e:
		default:
			t.Errorf("Event %#v was seen, but we have already seen %d events and didn't expect anymore.", e, expectedNumberOfEvents)
		}
	})
	close(events)

	eventsSeen := make(map[static.Event]bool)
	for event := range events {
		t.Logf("Event %#v was seen.", event)
		if eventsSeen[event] {
			t.Errorf("Event %#v was seen more than once, but was only expected once.", event)
		}
		eventsSeen[event] = true
	}

	for _, expect := range expected {
		if !eventsSeen[expect.Event] {
			t.Errorf("Event %#v was not seen during the build.", expect.Event)
		}

		outputFileContents, err := ioutil.ReadFile(expect.OutputFilePath)
		if err != nil {
			t.Fatalf("Expected %s to exist with the output but got error when opening: %v", expect.OutputFilePath, err)
		}

		t.Logf("Contents of %s => %s", expect.OutputFilePath, outputFileContents)
		if string(outputFileContents) != expect.OutputFileContents {
			t.Fatalf(`Contents of %s => %s, want %s`, expect.OutputFilePath, outputFileContents, expect.OutputFileContents)
		}
	}
}

func TestBuild_WithoutEventHandler(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", r.URL.Path)
	})

	options := static.DefaultOptions()
	tempDir, _ := ioutil.TempDir("", "test_build")
	t.Logf("Created temp dir for OutputDir: %s", tempDir)
	options.OutputDir = filepath.Join(tempDir, "build")
	t.Logf("OutputDir: %s", options.OutputDir)

	paths := []string{"/go", "/world", "/universe"}
	expected := []struct {
		OutputFilePath     string
		OutputFileContents string
	}{
		{filepath.Join(options.OutputDir, "go"), "Hello /go!"},
		{filepath.Join(options.OutputDir, "world"), "Hello /world!"},
		{filepath.Join(options.OutputDir, "universe"), "Hello /universe!"},
	}

	static.Build(options, handler, paths, nil)

	for _, expect := range expected {
		outputFileContents, err := ioutil.ReadFile(expect.OutputFilePath)
		if err != nil {
			t.Fatalf("Expected %s to exist with the output but got error when opening: %v", expect.OutputFilePath, err)
		}

		t.Logf("Contents of %s => %s", expect.OutputFilePath, outputFileContents)
		if string(outputFileContents) != expect.OutputFileContents {
			t.Fatalf(`Contents of %s => %s, want %s`, expect.OutputFilePath, outputFileContents, expect.OutputFileContents)
		}
	}
}
