package commands

import (
	"github.com/codegangsta/cli"

	. "github.com/yuuki1/grabeni/log"
)

var Commands = []cli.Command{
	CommandStatus,
	CommandList,
	CommandAttach,
	CommandDetach,
	CommandGrab,
}

func fatalOnError(command func(context *cli.Context) error) func(context *cli.Context) {
	return func(context *cli.Context) {
		if err := command(context); err != nil {
			Logf("error", err.Error())
		}
	}
}
