package format

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/service/ec2"
)

const header = "ENI ID\tPrivate DNS Name\tPrivate IP\tInstance ID\tDevice index\tStatus\tName\tAZ"

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
		var networkInterfaceId string
		var privateDnsName string
		var privateIpAddress string
		var instanceID string
		var deviceIndex int64 = -1
		var name string
		var status string
		var az string

		// avoid to occur panic bacause of invalid memory address or nil pointer dereference
		if eni.NetworkInterfaceId != nil {
			networkInterfaceId = *eni.NetworkInterfaceId
		}
		if eni.PrivateDnsName != nil {
			privateDnsName = *eni.PrivateDnsName
		}
		if eni.PrivateIpAddress != nil {
			privateIpAddress = *eni.PrivateIpAddress
		}
		if eni.Status != nil {
			status = *eni.Status
		}
		if eni.AvailabilityZone != nil {
			az = *eni.AvailabilityZone
		}

		if eni.Attachment == nil {
			instanceID, deviceIndex = "", -1
		} else {
			// managed services such as Amazon RDS don't have InstanceId.
			if eni.Attachment.InstanceId != nil {
				instanceID = *eni.Attachment.InstanceId
			}

			if eni.Attachment.DeviceIndex != nil {
				deviceIndex = *eni.Attachment.DeviceIndex
			}
		}

		if len(eni.TagSet) > 0 {
			for _, tag := range eni.TagSet {
				if *tag.Key == "Name" {
					name = *tag.Value
					break
				}
			}
		}

		fmt.Fprintln(tw, fmt.Sprintf("%s\t%s\t%s\t%s\t%d\t%s\t%s\t%s",
			networkInterfaceId,
			privateDnsName,
			privateIpAddress,
			instanceID,
			deviceIndex,
			status,
			name,
			az,
		))
	}

	tw.Flush()
}
