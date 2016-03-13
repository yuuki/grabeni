package commands

import (
	"fmt"

	"github.com/codegangsta/cli"

	"github.com/yuuki/grabeni/log"
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
			log.Error(fmt.Sprintf("error: %s", err.Error()))
		}
	}
}
