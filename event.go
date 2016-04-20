package static

import (
	"fmt"
)

// A build event representing that an action has taken place for a path in the build process, and includes an error if an error occurred while the action took place.
type Event struct {
	// The action taken place on the path.
	Action string
	// The path the action has taken place on.
	Path string
	// An error if an error occurred while performing the action, otherwise nil.
	Error error
}

// A simple string representation of an Event in the format:
//	 Action: <action>, Path: <path>
// And when the Event has an error:
//	 Action: <action>, Path: <path>, Error: <error>
func (e Event) String() string {
	s := fmt.Sprintf("Action: %s, Path: %s", e.Action, e.Path)
	if e.Error != nil {
		s += fmt.Sprintf(", Error: %v", e.Error)
	}
	return s
}
