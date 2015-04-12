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
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "accesskey", Value: "", Usage: "AWS ACCESS KEY ID"},
		cli.StringFlag{Name: "secretkey", Value: "", Usage: "AWS SECRET ACCESS KEY"},
		cli.StringFlag{Name: "region", Value: "", Usage: "AWS Region (eg. ap-northeast-1)"},
	}
	app.Commands = Commands

	app.Run(os.Args)
}
