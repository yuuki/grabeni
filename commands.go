package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"

	"github.com/y-uuki/grabeni/aws"
	. "github.com/y-uuki/grabeni/log"
)

var Commands = []cli.Command{
	commandStatus,
	commandGrab,
	commandAttach,
	commandDetach,
}

var commandStatus = cli.Command{
	Name:  "status",
	Usage: "Show ENI status",
	Description: `
`,
	Action: doStatus,
}

var commandGrab = cli.Command{
	Name:  "grab",
	Usage: "Detach and attach ENI",
	Description: `
`,
	Action: doGrab,
	Flags: []cli.Flag{
		cli.IntFlag{Name: "d, deviceindex", Value: 1, Usage: "Device Index Number"},
		cli.StringFlag{Name: "i, instanceid", Usage: "Instance Id"},
	},
}

var commandAttach = cli.Command{
	Name:  "attach",
	Usage: "",
	Description: `
`,
	Action: doAttach,
	Flags: []cli.Flag{
		cli.IntFlag{Name: "d, deviceindex", Value: 1, Usage: "Device Index Number"},
		cli.StringFlag{Name: "i, instanceid", Usage: "Instance Id"},
	},
}

var commandDetach = cli.Command{
	Name:  "detach",
	Usage: "",
	Description: `
`,
	Action: doDetach,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "i, instanceid", Usage: "Instance Id"},
	},
}

func fetchInstanceIDIfEmpty(c *cli.Context) string {
	if instanceID := c.String("instanceid"); instanceID != "" {
		return c.String(instanceID)
	}

	instanceID, err := aws.GetInstanceID()
	if err != nil {
		cli.ShowCommandHelp(c, c.Command.Name)
		DieIf(err)
	}

	return instanceID
}

func doStatus(c *cli.Context) {
	if len(c.Args()) < 1 {
		cli.ShowCommandHelp(c, "status")
		os.Exit(1)
	}

	eniID := c.Args()[0]
	if eniID == "" {
		cli.ShowCommandHelp(c, "status")
		os.Exit(1)
	}

	cli, err := aws.NewClient(c)
	DieIf(err)

	eni, err := cli.DescribeENIByID(eniID)
	DieIf(err)
	if eni == nil {
		os.Exit(0)
	}

	name := ""
	if len(eni.TagSet) > 0 {
		for _, tag := range eni.TagSet {
			if *tag.Key == "Name" {
				name = *tag.Value
				break
			}
		}
	}

	fmt.Println("Name\tNetworkInterfaceID\tPrivateDNSName\tPrivateIPAddress\tInstanceID\tDeviceIndex\tStatus")
	fmt.Printf("%s\t%s\t%s\t%s\t%s\t%d\t%s\t\n",
		name,
		*eni.NetworkInterfaceID,
		*eni.PrivateDNSName,
		*eni.PrivateIPAddress,
		*eni.Attachment.InstanceID,
		*eni.Attachment.DeviceIndex,
		*eni.Status,
	)
}

func doGrab(c *cli.Context) {
	if len(c.Args()) < 1 {
		cli.ShowCommandHelp(c, "grab")
		os.Exit(1)
	}

	eniID := c.Args()[0]
	if eniID == "" {
		cli.ShowCommandHelp(c, "grab")
		os.Exit(1)
	}

	instanceID := fetchInstanceIDIfEmpty(c)
	deviceIndex := c.Int("deviceindex")

	cli, err := aws.NewClient(c)
	DieIf(err)

	err = cli.GrabENI(eniID, instanceID, deviceIndex)
	DieIf(err)
}

func doAttach(c *cli.Context) {
	if len(c.Args()) < 1 {
		cli.ShowCommandHelp(c, "attach")
		os.Exit(1)
	}

	eniID := c.Args()[0]
	if eniID == "" {
		cli.ShowCommandHelp(c, "attach")
		os.Exit(1)
	}

	instanceID := fetchInstanceIDIfEmpty(c)
	deviceIndex := c.Int("deviceindex")

	cli, err := aws.NewClient(c)
	DieIf(err)

	err = cli.AttachENI(eniID, instanceID, deviceIndex)
	DieIf(err)
}

func doDetach(c *cli.Context) {
	if len(c.Args()) < 1 {
		cli.ShowCommandHelp(c, "detach")
		os.Exit(1)
	}

	eniID := c.Args()[0]
	if eniID == "" {
		cli.ShowCommandHelp(c, "detach")
		os.Exit(1)
	}

	instanceID := fetchInstanceIDIfEmpty(c)

	cli, err := aws.NewClient(c)
	DieIf(err)

	err = cli.DetachENI(eniID, instanceID)
	DieIf(err)
}
