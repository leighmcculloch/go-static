package static_test

import (
	"."
	"testing"
)

func TestNewOptions(t *testing.T) {
	options := static.NewOptions()
	expected := static.Options{OutputDir: "build", Concurrency: 50, DirFilename: "index.html"}

	t.Logf("static.NewOptions() => %#v", options)
	if options != expected {
		t.Errorf("static.NewOptions() => %#v, want %#v", options, expected)
	}
}
