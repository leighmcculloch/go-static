package static

import (
	"errors"
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		event    Event
		expected string
	}{
		{Event{Action: "action", StatusCode: 200, Path: "/path", OutputPath: "/output-path/path"}, "Action: action, StatusCode: 200, Path: /path, OutputPath: /output-path/path"},
		{Event{Action: "action", StatusCode: 404, Path: "/path", OutputPath: "/output-path/path"}, "Action: action, StatusCode: 404, Path: /path, OutputPath: /output-path/path"},
		{Event{Action: "action", StatusCode: 200, Path: "/path", OutputPath: "/output-path/path", Error: errors.New("error")}, "Action: action, StatusCode: 200, Path: /path, OutputPath: /output-path/path, Error: error"},
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
