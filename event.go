package static

import (
	"fmt"
)

type Event struct {
	Action string
	Path   string
	Error  error
}

func (e Event) String() string {
	s := fmt.Sprintf("Action: %s, Path: %s", e.Action, e.Path)
	if e.Error != nil {
		s += fmt.Sprintf(", Error: %v", e.Error)
	}
	return s
}
