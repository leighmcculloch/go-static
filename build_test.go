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
	tempDir, err := ioutil.TempDir("", "test_build_single")
	if err != nil {
		t.Fatalf("Error creating temp dir for OutputDir: %v", err)
	}
	t.Logf("Created temp dir for OutputDir: %s", tempDir)
	options.OutputDir = filepath.Join(tempDir, "build")
	t.Logf("OutputDir: %s", options.OutputDir)

	path := "/world"
	expectedOutputFilePath := filepath.Join(options.OutputDir, "world")
	expectedOutputFileContents := "Hello /world!"

	err = static.BuildSingle(options, handler, path)
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
