package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "grabeni"
	app.Version = Version
	app.Usage = ""
	app.Author = "y_uuki"
	app.Email = "yuki.tsubo@gmail.com"
	app.Commands = Commands

	app.Run(os.Args)
}
