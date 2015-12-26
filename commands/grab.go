package commands

import (
	"errors"

	"github.com/codegangsta/cli"

	"github.com/yuuki1/grabeni/aws"
	"github.com/yuuki1/grabeni/log"
)

var CommandArgGrab = "[--instanceid INSTANCE_ID] [--deviceindex DEVICE_INDEX] [--max-attempts MAX_ATTEMPTS] [--interval INTERVAL] ENI_ID"
var CommandGrab = cli.Command{
	Name:   "grab",
	Usage:  "Detach and attach ENI whether the eni has already attached or not.",
	Action: fatalOnError(doGrab),
	Flags: []cli.Flag{
		cli.IntFlag{Name: "d, deviceindex", Value: 1, Usage: "device index number"},
		cli.StringFlag{Name: "I, instanceid", Usage: "attach-targeted instance id"},
		cli.IntFlag{Name: "n, max-attempts", Value: 10, Usage: "the maximum number of attempts to poll the change of ENI status (default: 10)"},
		cli.IntFlag{Name: "i, interval", Value: 2, Usage: "the interval in seconds to poll the change of ENI status (default: 2)"},
	},
}

func doGrab(c *cli.Context) error {
	if len(c.Args()) < 1 {
		cli.ShowCommandHelp(c, "grab")
		return errors.New("ENI_ID required")
	}

	eniID := c.Args().Get(0)

	var instanceID string
	if instanceID = c.String("instanceid"); instanceID == "" {
		var err error
		instanceID, err = aws.NewMetaDataClient().GetInstanceID()
		if err != nil {
			return err
		}
	}

	eni, err := aws.NewENIClient().GrabENI(&aws.GrabENIParam{
		InterfaceID: eniID,
		InstanceID:  instanceID,
		DeviceIndex: c.Int("deviceindex"),
	}, &aws.WaiterParam{
		MaxAttempts: c.Int("max-attempts"),
		IntervalSec: c.Int("interval"),
	})
	if err != nil {
		return err
	}
	if eni == nil {
		log.Infof("attached: eni %s already attached to instance %s", eniID, instanceID)
		return nil
	}

	log.Infof("grabbed: eni %s attached to instance %s", eniID, instanceID)

	return nil
}
