package format

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/service/ec2"
)

const header = "ENI ID\tPrivate DNS Name\tPrivate IP\tInstance ID\tDevice index\tStatus\tName"

func PrintENI(w io.Writer, eni *ec2.NetworkInterface) {
	enis := make([]*ec2.NetworkInterface, 1)
	enis[0] = eni
	PrintENIs(w, enis)
}

func PrintENIs(w io.Writer, enis []*ec2.NetworkInterface) {
	// Format in tab-separated columns with a tab stop of 8.
	tw := tabwriter.NewWriter(w, 0, 8, 0, '\t', 0)

	fmt.Fprintln(tw, header)

	for _, eni := range enis {
		var instanceID string
		var deviceIndex int64

		if eni.Attachment == nil {
			instanceID, deviceIndex = "", -1
		} else {
			instanceID = *eni.Attachment.InstanceId
			deviceIndex = *eni.Attachment.DeviceIndex
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

		fmt.Fprintln(tw, fmt.Sprintf("%s\t%s\t%s\t%s\t%d\t%s\t%s",
			*eni.NetworkInterfaceId,
			*eni.PrivateDnsName,
			*eni.PrivateIpAddress,
			instanceID,
			deviceIndex,
			*eni.Status,
			name,
		))
	}

	tw.Flush()
}
