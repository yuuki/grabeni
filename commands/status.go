package commands

import (
	"errors"
	"os"

	"github.com/urfave/cli"

	"github.com/yuuki/grabeni/aws"
	"github.com/yuuki/grabeni/format"
)

var CommandArgStatus = "ENI_ID"
var CommandStatus = cli.Command{
	Name:    "status",
	Aliases: []string{"st"},
	Usage:   "Show ENI status",
	Action:  fatalOnError(doStatus),
}

func doStatus(c *cli.Context) error {
	if len(c.Args()) < 1 {
		cli.ShowCommandHelp(c, "status")
		return errors.New("ENI_ID required")
	}

	eniID := c.Args().Get(0)

	eni, err := aws.NewENIClient().WithLogWriter(os.Stdout).DescribeENIByID(eniID)
	if err != nil {
		return err
	}
	if eni == nil {
		return nil
	}

	format.PrintENI(os.Stdout, eni)

	return nil
}
