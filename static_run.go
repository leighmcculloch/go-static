package static

import (
	"log"
	"os"
)

const (
	RunCommandServer = "server"
	RunCommandBuild  = "build"
)

func (s *Static) Run() {
	command := RunCommandServer
	if len(os.Args) >= 2 {
		command = os.Args[1]
	}

	var err error
	switch command {
	case RunCommandBuild:
		err = s.Build()
	case RunCommandServer:
		err = s.ListenAndServe(":4567")
	}

	if err != nil {
		log.Fatal(err)
	}
}
