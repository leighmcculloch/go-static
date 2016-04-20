package static_test

import (
	"."
	"fmt"
	"net/http"
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
