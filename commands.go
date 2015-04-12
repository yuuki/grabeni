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
	Usage: "",
	Description: `
`,
	Action: doGrab,
}

var commandAttach = cli.Command{
	Name:  "attach",
	Usage: "",
	Description: `
`,
	Action: doAttach,
}

var commandDetach = cli.Command{
	Name:  "detach",
	Usage: "",
	Description: `
`,
	Action: doDetach,
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
}

func doAttach(c *cli.Context) {
}

func doDetach(c *cli.Context) {
}
