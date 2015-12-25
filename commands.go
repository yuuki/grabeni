package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"

	"github.com/yuuki1/grabeni/aws"
	"github.com/yuuki1/grabeni/format"
	. "github.com/yuuki1/grabeni/log"
)

var Commands = []cli.Command{
	commandStatus,
	commandList,
	commandGrab,
	commandAttach,
	commandDetach,
}

var commandStatus = cli.Command{
	Name:    "status",
	Aliases: []string{"st"},
	Usage:   "Show ENI status",
	Description: `
    Show the information of the ENI identified with <eni-id>.
`,
	Action: doStatus,
}

var commandList = cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "Show all ENIs",
	Description: `
    List the ENIs owned by your account.
`,
	Action: doList,
}

var commandGrab = cli.Command{
	Name:  "grab",
	Usage: "Detach and attach ENI",
	Description: `
    Detach the ENI identified with <eni-id> whether the eni has already attached or not. And then attach the ENI to EC2 instance identified with --instanceid option.
`,
	Action: doGrab,
	Flags: []cli.Flag{
		cli.IntFlag{Name: "d, deviceindex", Value: 1, Usage: "device index number"},
		cli.StringFlag{Name: "I, instanceid", Usage: "attach-targeted instance id"},
		cli.IntFlag{Name: "t, timeout", Value: 10, Usage: "each attach and detach API request timeout seconds"},
		cli.IntFlag{Name: "i, interval", Value: 2, Usage: "each attach and detach API request polling interval seconds"},
	},
}

var commandAttach = cli.Command{
	Name:  "attach",
	Usage: "Attach ENI",
	Description: `
    Just attach the ENI identified with <eni-id> to EC2 instance identified with --instanceid option.
`,
	Action: doAttach,
	Flags: []cli.Flag{
		cli.IntFlag{Name: "d, deviceindex", Value: 1, Usage: "device index number"},
		cli.StringFlag{Name: "I, instanceid", Usage: "attach-targeted instance id"},
		cli.IntFlag{Name: "t, timeout", Value: 10, Usage: "each attach and detach API request timeout seconds"},
		cli.IntFlag{Name: "i, interval", Value: 2, Usage: "each attach and detach API request polling interval seconds"},
	},
}

var commandDetach = cli.Command{
	Name:  "detach",
	Usage: "Detach ENI",
	Description: `
    Just detach the ENI identified with <eni-id>.
`,
	Action: doDetach,
	Flags: []cli.Flag{
		cli.IntFlag{Name: "t, timeout", Value: 10, Usage: "each attach and detach API request timeout seconds"},
		cli.IntFlag{Name: "i, interval", Value: 2, Usage: "each attach and detach API request polling interval seconds"},
	},
}

type commandDoc struct {
	Parent    string
	Arguments string
}

var commandDocs = map[string]commandDoc{
	"status": {"", "<eni-id>"},
	"list":   {"", ""},
	"grab":   {"", "[--deviceindex | -d <device index> (default: 1)] [--instanceid | -i <instance id>] <eni-id>"},
	"attach": {"", "[--deviceindex | -d <device index> (default: 1)] [--instanceid | -i <instance id>] <eni-id>"},
	"detach": {"", "<eni-id>"},
}

// Makes template conditionals to generate per-command documents.
func mkCommandsTemplate(genTemplate func(commandDoc) string) string {
	template := "{{if false}}"
	for _, command := range append(Commands) {
		template = template + fmt.Sprintf("{{else if (eq .Name %q)}}%s", command.Name, genTemplate(commandDocs[command.Name]))
	}
	return template + "{{end}}"
}

func init() {
	argsTemplate := mkCommandsTemplate(func(doc commandDoc) string { return doc.Arguments })
	parentTemplate := mkCommandsTemplate(func(doc commandDoc) string { return string(strings.TrimLeft(doc.Parent+" ", " ")) })

	cli.CommandHelpTemplate = `NAME:
    {{.Name}} - {{.Usage}}

USAGE:
    grabeni ` + parentTemplate + `{{.Name}} ` + argsTemplate + `
{{if (len .Description)}}
DESCRIPTION: {{.Description}}
{{end}}{{if (len .Flags)}}
OPTIONS:
    {{range .Flags}}{{.}}
    {{end}}
{{end}}`
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

	eni, err := aws.NewENIClient().DescribeENIByID(eniID)
	if err != nil {
		Logf("error", err.Error())
		os.Exit(1)
	}
	if eni == nil {
		os.Exit(0)
	}

	format.PrintENI(os.Stdout, eni)
}

func doList(c *cli.Context) {
	enis, err := aws.NewENIClient().DescribeENIs()
	if err != nil {
		Logf("error", err.Error())
		os.Exit(1)
	}
	if enis == nil {
		os.Exit(0)
	}

	format.PrintENIs(os.Stdout, enis)
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

	var instanceID string
	if instanceID = c.String("instanceid"); instanceID == "" {
		var err error
		instanceID, err = aws.NewMetaDataClient().GetInstanceID()
		if err != nil {
			Logf("error", err.Error())
			os.Exit(1)
		}
	}

	eni, err := aws.NewENIClient().GrabENI(&aws.GrabENIParam{
		InterfaceID: eniID,
		InstanceID:  instanceID,
		DeviceIndex: c.Int("deviceindex"),
	}, &aws.RetryParam{
		TimeoutSec:  int64(c.Int("timeout")),
		IntervalSec: int64(c.Int("interval")),
	})
	if err != nil {
		Logf("error", err.Error())
		os.Exit(1)
	}
	if eni == nil {
		Logf("attached", "eni %s already attached to instance %s", eniID, instanceID)
		os.Exit(0)
	}

	Logf("grabbed", "eni %s attached to instance %s", eniID, instanceID)
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

	var instanceID string
	if instanceID = c.String("instanceid"); instanceID == "" {
		var err error
		instanceID, err = aws.NewMetaDataClient().GetInstanceID()
		if err != nil {
			Logf("error", err.Error())
			os.Exit(1)
		}
	}

	eni, err := aws.NewENIClient().AttachENIWithRetry(&aws.AttachENIParam{
		InterfaceID: eniID,
		InstanceID:  instanceID,
		DeviceIndex: c.Int("deviceindex"),
	}, &aws.RetryParam{
		TimeoutSec:  int64(c.Int("timeout")),
		IntervalSec: int64(c.Int("interval")),
	})
	if err != nil {
		Logf("error", err.Error())
		os.Exit(1)
	}
	if eni == nil {
		Logf("attached", "eni %s already attached to instance %s", eniID, instanceID)
		os.Exit(0)
	}

	Logf("attached", "eni %s attached to instance %s", eniID, instanceID)
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

	eni, err := aws.NewENIClient().DetachENIWithRetry(&aws.DetachENIParam{
		InterfaceID: eniID,
	}, &aws.RetryParam{
		TimeoutSec:  int64(c.Int("timeout")),
		IntervalSec: int64(c.Int("interval")),
	})
	if err != nil {
		Logf("error", err.Error())
		os.Exit(1)
	}
	if eni == nil {
		Logf("detached", "eni %s already detached", eniID)
		os.Exit(0)
	}

	Logf("detached", "eni %s detached", eniID)
}
