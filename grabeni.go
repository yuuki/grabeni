package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "grabeni"
	app.Version = Version
	app.Usage = "AWS ENI grabbing tool in Go"
	app.Author = "y_uuki"
	app.Commands = Commands

	app.Run(os.Args)
}
