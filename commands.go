package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

var Commands = []cli.Command{
	commandStatus,
	commandGrab,
	commandAttach,
	commandDetach,
}

var commandStatus = cli.Command{
	Name:  "status",
	Usage: "",
	Description: `
`,
	Action: doStatus,
}

var commandGrab = cli.Command{
	Name:  "grab",
	Usage: "",
	Description: `
`,
	Action: doGrab,
}

var commandAttach = cli.Command{
	Name:  "attach",
	Usage: "",
	Description: `
`,
	Action: doAttach,
}

var commandDetach = cli.Command{
	Name:  "detach",
	Usage: "",
	Description: `
`,
	Action: doDetach,
}

func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Println(v...)
	}
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func doStatus(c *cli.Context) {
}

func doGrab(c *cli.Context) {
}

func doAttach(c *cli.Context) {
}

func doDetach(c *cli.Context) {
}
