package static

import (
	"errors"
	"fmt"
	"log"
	"os"
)

var (
	ErrUnknownCommand = errors.New("Unknown command. Valid commands are: 'server', 'build'.")
)

const (
	RunCommandServer = "server"
	RunCommandBuild  = "build"
)

func (s Static) Run() error {
	command := RunCommandServer
	if len(os.Args) >= 2 {
		command = os.Args[1]
	}

	switch command {
	case RunCommandBuild:
		s.Build(logEvent)
		return nil
	case RunCommandServer:
		addr := fmt.Sprintf(":%d", s.ServerPort)
		return s.ListenAndServe(addr, logEvent)
	}

	return ErrUnknownCommand
}

func logEvent(event Event) {
	var s string
	if event.Error == nil {
		s = fmt.Sprintf("%10s  %-20s", event.Action, event.Path)
	} else {
		s = fmt.Sprintf("%10s  %-20s  %v", "error", event.Path, event.Error)
	}
	log.Println(s)
}
