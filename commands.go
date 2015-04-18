package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"

	"github.com/yuuki1/grabeni/aws"
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
	Name:  "status",
	Usage: "Show ENI status",
	Description: `
    Show the information of the ENI identified with <eni-id>.
`,
	Action: doStatus,
}

var commandList = cli.Command{
	Name:  "list",
	Usage: "Show all ENIs",
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
	"list":  {"", ""},
	"grab": {"", "[--deviceindex | -d <device index> (default: 1)] [--instanceid | -i <instance id>] <eni-id>"},
	"attach": {"", "[--deviceindex | -d <device index> (default: 1)] [--instanceid | -i <instance id>] <eni-id>"},
	"detach":  {"", "<eni-id>"},
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

func awsCli(c *cli.Context) *aws.Client {
	client, err := aws.NewClient(
		c.GlobalString("region"),
		c.GlobalString("accesskey"),
		c.GlobalString("secretkey"),
	)
	DieIf(err)
	return client
}

func fetchInstanceIDIfEmpty(c *cli.Context) string {
	if instanceID := c.String("instanceid"); instanceID != "" {
		return instanceID
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

	eni, err := awsCli(c).DescribeENIByID(eniID)
	DieIf(err)
	if eni == nil {
		os.Exit(0)
	}

	fmt.Println(StatusHeader())
	fmt.Println(StatusRow(eni))
}

func doList(c *cli.Context) {
	enis, err := awsCli(c).DescribeENIs()
	DieIf(err)
	if enis == nil {
		os.Exit(0)
	}

	fmt.Println(StatusHeader())
	for _, eni := range enis {
		fmt.Println(StatusRow(eni))
	}
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

	err := awsCli(c).GrabENI(&aws.GrabENIParam{
		InterfaceID: eniID,
		InstanceID: instanceID,
		DeviceIndex: c.Int("deviceindex"),
	}, &aws.RetryParam{
		TimeoutSec: int64(c.Int("timeout")),
		IntervalSec: int64(c.Int("interval")),
	})
	DieIf(err)

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

	instanceID := fetchInstanceIDIfEmpty(c)

	err := awsCli(c).AttachENIWithRetry(&aws.AttachENIParam{
		InterfaceID: eniID,
		InstanceID: instanceID,
		DeviceIndex: c.Int("deviceindex"),
	}, &aws.RetryParam{
		TimeoutSec: int64(c.Int("timeout")),
		IntervalSec: int64(c.Int("interval")),
	})
	DieIf(err)

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

	err := awsCli(c).DetachENIWithRetry(&aws.DetachENIParam{
		InterfaceID: eniID,
	}, &aws.RetryParam{
		TimeoutSec: int64(c.Int("timeout")),
		IntervalSec: int64(c.Int("interval")),
	})
	DieIf(err)

	Logf("detached", "eni %s detached", eniID)
}
