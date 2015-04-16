package main

import (
	"fmt"

	"github.com/awslabs/aws-sdk-go/service/ec2"
)

func StatusHeader() string {
	// arrange columns by multiple "\t"
	return "NetworkInterfaceID\tPrivateDNSName\t\t\t\tPrivateIPAddress\tInstanceID\tDeviceIndex\tStatus\tName"
}

func StatusRow(eni *ec2.NetworkInterface) string {
	var instanceID string
	var deviceIndex int64
	if eni.Attachment == nil {
		instanceID, deviceIndex = "", -1
	} else {
		instanceID = *eni.Attachment.InstanceID
		deviceIndex = *eni.Attachment.DeviceIndex
	}

	if instanceID == "" {
		instanceID = "\t"
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

	return fmt.Sprintf("%s\t\t%s\t%s\t%s\t%d\t\t%s\t%s\t\n",
		*eni.NetworkInterfaceID,
		*eni.PrivateDNSName,
		*eni.PrivateIPAddress,
		instanceID,
		deviceIndex,
		*eni.Status,
		name,
	)
}
