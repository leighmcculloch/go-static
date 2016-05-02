package static

// buildError is an error that occurred during build, and wraps the cause error.
type buildError struct {
	message string
	cause   error
}

// Error returns the error and it's cause as a string.
func (e buildError) Error() string {
	s := e.message
	if e.cause != nil {
		s += ": " + e.cause.Error()
	}
	return s
}
