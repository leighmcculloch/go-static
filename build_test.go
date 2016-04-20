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

	paths = append(paths, "/")
	paths = append(paths, "/world")
	paths = append(paths, "/go")

	options := static.DefaultOptions()
	static.Build(options, handler, paths, func(e static.Event) {
		fmt.Println(e)
	})

	// Output:
	// Action: build, Path: /
	// Action: build, Path: /world
	// Action: build, Path: /go
}

func ExampleBuildSingle() {
	handler := http.NewServeMux()

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", r.URL.Path)
	})

	options := static.DefaultOptions()

	var err error

	err = static.BuildSingle(options, handler, "/")
	fmt.Println("Built: /, Error:", err)

	err = static.BuildSingle(options, handler, "/world")
	fmt.Println("Built: /world, Error:", err)

	err = static.BuildSingle(options, handler, "/go")
	fmt.Println("Built: /go, Error:", err)

	// Output:
	// Built: /, Error: <nil>
	// Built: /world, Error: <nil>
	// Built: /go, Error: <nil>
}
