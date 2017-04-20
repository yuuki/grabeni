package commands

import (
	"os"

	"github.com/urfave/cli"

	"github.com/yuuki/grabeni/aws"
	"github.com/yuuki/grabeni/format"
)

var CommandArgList = ""
var CommandList = cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "List ENIs",
	Action:  fatalOnError(doList),
}

func doList(c *cli.Context) error {
	enis, err := aws.NewENIClient().WithLogWriter(os.Stdout).DescribeENIs()
	if err != nil {
		return err
	}
	if enis == nil {
		return nil
	}

	format.PrintENIs(os.Stdout, enis)

	return nil
}
