package static_test

import (
	"testing"

	"github.com/leighmcculloch/static"
)

func TestDefaultOptions(t *testing.T) {
	options := static.DefaultOptions
	expected := static.Options{OutputDir: "build", Concurrency: 50, DirFilename: "index.html"}

	t.Logf("DefaultOptions => %#v", options)
	if options != expected {
		t.Errorf("DefaultOptions => %#v, want %#v", options, expected)
	}
}
