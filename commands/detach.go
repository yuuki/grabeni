package commands

import (
	"errors"
	"os"

	"github.com/Songmu/prompter"
	"github.com/codegangsta/cli"

	"github.com/yuuki1/grabeni/aws"
	"github.com/yuuki1/grabeni/log"
)

var CommandArgDetach = "[--max-attempts MAX_ATTEMPTS] [--interval INTERVAL] ENI_ID"
var CommandDetach = cli.Command{
	Name:   "detach",
	Usage:  "Detach ENI",
	Action: fatalOnError(doDetach),
	Flags: []cli.Flag{
		cli.IntFlag{Name: "n, max-attempts", Value: 10, Usage: "the maximum number of attempts to poll the change of ENI status (default: 10)"},
		cli.IntFlag{Name: "i, interval", Value: 2, Usage: "the interval in seconds to poll the change of ENI status (default: 2)"},
	},
}

func doDetach(c *cli.Context) error {
	if len(c.Args()) < 1 {
		cli.ShowCommandHelp(c, "detach")
		return errors.New("ENI_ID required")
	}

	eniID := c.Args().Get(0)

	if !prompter.YN("Detach following ENI.\n  "+eniID+"\nAre you sure?", true) {
		log.Infof("detachment is canceled")
		return nil
	}

	eni, err := aws.NewENIClient().WithLogWriter(os.Stdout).DetachENIWithWaiter(&aws.DetachENIParam{
		InterfaceID: eniID,
	}, &aws.WaiterParam{
		MaxAttempts: c.Int("max-attempts"),
		IntervalSec: c.Int("interval"),
	})
	if err != nil {
		return err
	}
	if eni == nil {
		log.Infof("%s already detached", eniID)
		return nil
	}

	log.Infof("%s detached", eniID)

	return nil
}
