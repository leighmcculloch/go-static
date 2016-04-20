package static_test

import (
	"."
	"testing"
)

func TestDefaultOptions(t *testing.T) {
	options := static.DefaultOptions()
	expected := static.Options{OutputDir: "build", Concurrency: 50, DirFilename: "index.html"}

	t.Logf("static.DefaultOptions() => %#v", options)
	if options != expected {
		t.Errorf("static.DefaultOptions() => %#v, want %#v", options, expected)
	}
}
