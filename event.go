package static

import (
	"fmt"
)

// Event represents action has taken place for a path in the build process, and includes an error if an error occurred while the action took place.
type Event struct {
	// The action taken place on the path.
	Action Action
	// The path the action has taken place on.
	Path string
	// The HTTP status code returned when the path was built.
	StatusCode int
	// The output path where the output of the action was written to.
	OutputPath string
	// An error if an error occurred while performing the action, otherwise nil.
	Error error
}

// Action is something taken place, captured in an Event.
type Action string

const (
	// BUILD is the building of a path.
	BUILD Action = "build"
)

// A simple string representation of an Event in the format:
//	 Action: build, Path: <path>, StatusCode: 200|404|etc, OutputPath: <output-path>
// And when the Event has an error:
//	 Action: build, Path: <path>, StatusCode: 200|404|etc, OutputPath: <output-path>, Error: <error>
func (e Event) String() string {
	s := fmt.Sprintf("Action: %s, Path: %s, StatusCode: %d, OutputPath: %s", e.Action, e.Path, e.StatusCode, e.OutputPath)
	if e.Error != nil {
		s += fmt.Sprintf(", Error: %v", e.Error)
	}
	return s
}
