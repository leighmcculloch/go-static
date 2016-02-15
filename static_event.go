package static

type Event struct {
	Action string
	Path   string
	Error  error
}

type EventHandler func(event Event)
