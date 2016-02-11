package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/Songmu/prompter"
	"github.com/codegangsta/cli"

	"github.com/yuuki1/grabeni/aws"
	"github.com/yuuki1/grabeni/log"
)

var CommandArgAttach = "[--instanceid INSTANCE_ID] [--deviceindex DEVICE_INDEX] [--max-attempts MAX_ATTEMPTS] [--interval INTERVAL] ENI_ID"
var CommandAttach = cli.Command{
	Name:   "attach",
	Usage:  "Attach ENI",
	Action: fatalOnError(doAttach),
	Flags: []cli.Flag{
		cli.IntFlag{Name: "d, deviceindex", Value: 1, Usage: "device index number"},
		cli.StringFlag{Name: "I, instanceid", Usage: "attach-targeted instance id"},
		cli.IntFlag{Name: "n, max-attempts", Value: 10, Usage: "the maximum number of attempts to poll the change of ENI status (default: 10)"},
		cli.IntFlag{Name: "i, interval", Value: 2, Usage: "the interval in seconds to poll the change of ENI status (default: 2)"},
	},
}

func doAttach(c *cli.Context) error {
	if len(c.Args()) < 1 {
		cli.ShowCommandHelp(c, "attach")
		return errors.New("ENI_ID required")
	}

	eniID := c.Args().Get(0)

	if !prompter.YN("Attach following ENI.\n  "+eniID+"\nAre you sure?", true) {
		log.Infof("Attachment is canceled")
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

	eni, err := awscli.AttachENIWithWaiter(&aws.AttachENIParam{
		InterfaceID: eniID,
		InstanceID:  *instance.InstanceId,
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
