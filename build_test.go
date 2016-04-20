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
