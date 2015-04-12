package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"

	"github.com/y-uuki/grabeni/aws"
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
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "n, nametag", Usage: "ENI Tag Name"},
	},
}

var commandGrab = cli.Command{
	Name:  "grab",
	Usage: "Detach and attach ENI",
	Description: `
`,
	Action: doGrab,
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "n, nametag", Usage: "ENI Tag Name"},
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
		cli.BoolFlag{Name: "n, nametag", Usage: "ENI Tag Name"},
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
		cli.BoolFlag{Name: "n, nametag", Usage: "ENI Tag Name"},
		cli.StringFlag{Name: "i, instanceid", Usage: "Instance Id"},
	},
}

func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Println(v...)
	}
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
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
	if err != nil {
		assert(err)
		os.Exit(1)
	}

	eni, err := cli.DescribeENIByID(eniID)
	if err != nil {
		assert(err)
		os.Exit(1)
	}
	if eni == nil {
		os.Exit(1)
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

	instanceID := c.String("instanceid")
	if instanceID == "" {
		cli.ShowCommandHelp(c, "grab")
		os.Exit(1)
	}

	deviceIndex := c.Int("deviceindex")

	cli, err := aws.NewClient(c)
	if err != nil {
		assert(err)
		os.Exit(1)
	}

	err = cli.GrabENI(eniID, instanceID, deviceIndex)
	if err != nil {
		assert(err)
		os.Exit(1)
	}
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

	instanceID := c.String("instanceid")
	if instanceID == "" {
		cli.ShowCommandHelp(c, "attach")
		os.Exit(1)
	}

	deviceIndex := c.Int("deviceindex")

	cli, err := aws.NewClient(c)
	if err != nil {
		assert(err)
		os.Exit(1)
	}

	err = cli.AttachENI(eniID, instanceID, deviceIndex)
	if err != nil {
		assert(err)
		os.Exit(1)
	}
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

	instanceID := c.String("instanceid")
	if instanceID == "" {
		cli.ShowCommandHelp(c, "detach")
		os.Exit(1)
	}

	cli, err := aws.NewClient(c)
	if err != nil {
		assert(err)
		os.Exit(1)
	}

	err = cli.DetachENI(eniID, instanceID)
	if err != nil {
		assert(err)
		os.Exit(1)
	}
}
