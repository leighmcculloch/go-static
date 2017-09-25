package static_test

import (
	"errors"
	"testing"

	"4d63.com/static"
)

func TestString(t *testing.T) {
	tests := []struct {
		event    static.Event
		expected string
	}{
		{static.Event{Action: "action", Path: "/path", StatusCode: 200, OutputPath: "/output-path/path"}, "Action: action, Path: /path, StatusCode: 200, OutputPath: /output-path/path"},
		{static.Event{Action: "action", Path: "/path", StatusCode: 404, OutputPath: "/output-path/path"}, "Action: action, Path: /path, StatusCode: 404, OutputPath: /output-path/path"},
		{static.Event{Action: "action", Path: "/path", StatusCode: 200, OutputPath: "/output-path/path", Error: errors.New("error")}, "Action: action, Path: /path, StatusCode: 200, OutputPath: /output-path/path, Error: error"},
	}

	for _, test := range tests {
		s := test.event.String()
		if s == test.expected {
			t.Logf("%#v.String() => %v", test.event, s)
		} else {
			t.Errorf("%#v.String() => %v, want %v", test.event, s, test.expected)
		}
	}
}
