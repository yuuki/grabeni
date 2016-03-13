package format

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/yuuki/grabeni/aws/model"
)

const header = "ID\tNAME\tSTATUS\tPRIVATE DNS NAME\tPRIVATE IP\tAZ\tDEVICE INDEX\tINSTANCE ID\tINSTANCE NAME"

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
		if eni == nil {
			continue
		}

		instance := eni.AttachedInstance()
		var instanceID, instanceName string
		if instance != nil {
			instanceID, instanceName = instance.InstanceID(), instance.Name()
		}

		fmt.Fprintln(tw, fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%d\t%s\t%s",
			eni.InterfaceID(),
			eni.Name(),
			eni.Status(),
			eni.PrivateDnsName(),
			eni.PrivateIpAddress(),
			eni.AvailabilityZone(),
			eni.AttachedDeviceIndex(),
			instanceID,
			instanceName,
		))
	}

	tw.Flush()
}
