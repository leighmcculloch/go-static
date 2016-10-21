package static_test

import (
	"errors"
	"testing"

	"github.com/leighmcculloch/static"
)

func TestString(t *testing.T) {
	tests := []struct {
		event    static.Event
		expected string
	}{
		{static.Event{Action: "action", StatusCode: 200, Path: "/path"}, "Action: action, StatusCode: 200, Path: /path"},
		{static.Event{Action: "action", StatusCode: 404, Path: "/path"}, "Action: action, StatusCode: 404, Path: /path"},
		{static.Event{Action: "action", StatusCode: 200, Path: "/path", Error: errors.New("error")}, "Action: action, StatusCode: 200, Path: /path, Error: error"},
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
