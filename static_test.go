package static

import (
	"testing"
)

func TestFactoryDefaultSourceDir(t *testing.T) {
	const want = "source"

	got := New().SourceDir
	if got != want {
		t.Errorf("Expected default SourceDir to be %#v, but was %#v", want, got)
	}
}

func TestFactoryDefaultBuildDir(t *testing.T) {
	const want = "build"

	got := New().BuildDir
	if got != want {
		t.Errorf("Expected default BuildDir to be %#v, but was %#v", want, got)
	}
}

func TestFactoryDefaultBuildConcurrency(t *testing.T) {
	const want = 50

	got := New().BuildConcurrency
	if got != want {
		t.Errorf("Expected default BuildConcurrency to be %#v, but was %#v", want, got)
	}
}

func TestFactoryDefaultServerPort(t *testing.T) {
	const want = 4567

	got := New().ServerPort
	if got != want {
		t.Errorf("Expected default ServerPort to be %#v, but was %#v", want, got)
	}
}
