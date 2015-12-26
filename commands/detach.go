package commands

import (
	"errors"

	"github.com/codegangsta/cli"

	"github.com/yuuki1/grabeni/aws"
	"github.com/yuuki1/grabeni/log"
)

var CommandArgDetach = "[--timeout TIMEOUT] [--interval INTERVAL] ENI_ID"
var CommandDetach = cli.Command{
	Name:   "detach",
	Usage:  "Detach ENI",
	Action: fatalOnError(doDetach),
	Flags: []cli.Flag{
		cli.IntFlag{Name: "t, timeout", Value: 10, Usage: "each attach and detach API request timeout seconds"},
		cli.IntFlag{Name: "i, interval", Value: 2, Usage: "each attach and detach API request polling interval seconds"},
	},
}

func doDetach(c *cli.Context) error {
	if len(c.Args()) < 1 {
		cli.ShowCommandHelp(c, "detach")
		return errors.New("ENI_ID required")
	}

	eniID := c.Args().Get(0)

	eni, err := aws.NewENIClient().DetachENIWithRetry(&aws.DetachENIParam{
		InterfaceID: eniID,
	}, &aws.RetryParam{
		TimeoutSec:  int64(c.Int("timeout")),
		IntervalSec: int64(c.Int("interval")),
	})
	if err != nil {
		return err
	}
	if eni == nil {
		log.Infof("detached: eni %s already detached", eniID)
		return nil
	}

	log.Infof("detached: eni %s detached", eniID)

	return nil
}
