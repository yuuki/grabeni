package format

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/yuuki1/grabeni/aws/model"
)

const header = "ENI ID\tPrivate DNS Name\tPrivate IP\tInstance ID\tInsrance Name\tDevice index\tStatus\tName\tAZ"

func PrintENI(w io.Writer, eni *model.ENI) {
	enis := make([]*model.ENI, 1)
	enis[0] = eni
	PrintENIs(w, enis)
}

func PrintENIs(w io.Writer, enis []*model.ENI) {
	// Format in tab-separated columns with a tab stop of 8.
	tw := tabwriter.NewWriter(w, 0, 8, 0, '\t', 0)

	fmt.Fprintln(tw, header)

	for _, eni := range enis {
		instance := eni.AttachedInstance()
		var instanceID, instanceName string
		if instance != nil {
			instanceID, instanceName = instance.InstanceID(), instance.Name()
		}

		fmt.Fprintln(tw, fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%d\t%s\t%s\t%s",
			eni.InterfaceID(),
			eni.PrivateDnsName(),
			eni.PrivateIpAddress(),
			instanceID,
			instanceName,
			eni.AttachedDeviceIndex(),
			eni.Status(),
			eni.Name(),
			eni.AvailabilityZone(),
		))
	}

	tw.Flush()
}
