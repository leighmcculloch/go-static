// +build !windows

package static

import (
	"fmt"
	"net/http"
	"path"
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
	// Action: build, Path: /world, StatusCode: 200, OutputPath: build/world
}

func ExampleBuildSingle() {
	handler := http.NewServeMux()

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", path.Base(r.URL.Path))
	})

	statusCode, outputPath, err := BuildSingle(DefaultOptions, handler, "/world")
	fmt.Printf("Built: /world, StatusCode: %d, OutputPath: %v, Error: %v", statusCode, outputPath, err)

	// Output:
	// Built: /world, StatusCode: 200, OutputPath: build/world, Error: <nil>
}
