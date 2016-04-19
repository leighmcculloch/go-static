package static_test

import (
	"."
	"errors"
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		event    static.Event
		expected string
	}{
		{static.Event{Action: "action", Path: "/path"}, "Action: action, Path: /path"},
		{static.Event{Action: "action", Path: "/path", Error: errors.New("error")}, "Action: action, Path: /path, Error: error"},
	}

	for _, test := range tests {
		s := test.event.String()
		t.Logf("%#v.String() => %v", test.event, s)
		if s != test.expected {
			t.Errorf("%#v.String() => %v, want %v", test.event, s, test.expected)
		}
	}
}
