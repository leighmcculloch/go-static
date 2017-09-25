// +build !windows

package static_test

import (
	"fmt"
	"net/http"
	"path"

	"4d63.com/static"
)

func ExampleBuild() {
	handler := http.NewServeMux()
	paths := []string{}

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", path.Base(r.URL.Path))
	})

	paths = append(paths, "/world")

	static.Build(static.DefaultOptions, handler, paths, func(e static.Event) {
		fmt.Println(e)
	})

	// Output:
	// Action: build, Path: /world, StatusCode: 200, OutputPath: build/world
}

func ExampleBuildSingle() {
	handler := http.NewServeMux()

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", path.Base(r.URL.Path))
	})

	statusCode, outputPath, err := static.BuildSingle(static.DefaultOptions, handler, "/world")
	fmt.Printf("Built: /world, StatusCode: %d, OutputPath: %v, Error: %v", statusCode, outputPath, err)

	// Output:
	// Built: /world, StatusCode: 200, OutputPath: build/world, Error: <nil>
}
