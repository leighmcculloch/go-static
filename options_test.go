package static

import (
	"testing"
)

func TestDefaultOptions(t *testing.T) {
	options := DefaultOptions()
	expected := Options{OutputDir: "build", Concurrency: 50, DirFilename: "index.html"}

	t.Logf("DefaultOptions() => %#v", options)
	if options != expected {
		t.Errorf("DefaultOptions() => %#v, want %#v", options, expected)
	}
}
