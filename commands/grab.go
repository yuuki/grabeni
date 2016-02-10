package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/Songmu/prompter"

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

	if !prompter.YN("Grab following ENI.\n  "+eniID+"\nAre you sure?", true) {
		log.Infof("Grabbing is canceled")
		return nil
	}

	var instanceID string
	if instanceID = c.String("instanceid"); instanceID == "" {
		var err error
		instanceID, err = aws.NewMetaDataClient().GetInstanceID()
		if err != nil {
			return err
		}
	}

	awscli := aws.NewENIClient().WithLogWriter(os.Stdout)

	// Check instance id existence
	instance, err := awscli.DescribeInstanceByID(instanceID)
	if err != nil {
		return err
	}
	if instance == nil {
		return fmt.Errorf("No such instance %s", instanceID)
	}

	eni, err := awscli.GrabENI(&aws.GrabENIParam{
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
		log.Infof("%s already attached to instance %s", eniID, instanceID)
		return nil
	}

	log.Infof("%s attached to instance %s", eniID, instanceID)

	return nil
}
