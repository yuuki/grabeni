package model

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

type ENI struct {
	iface    *ec2.NetworkInterface
	instance *Instance
}

func NewENI(networkInterface *ec2.NetworkInterface) *ENI {
	return &ENI{iface: networkInterface, instance: nil}
}

func (e *ENI) SetInstance(i *Instance) {
	e.instance = i
}

func (e *ENI) InterfaceID() string {
	if e.iface.NetworkInterfaceId != nil {
		return *e.iface.NetworkInterfaceId
	}
	return ""
}

func (e *ENI) PrivateDnsName() string {
	if e.iface.PrivateDnsName != nil {
		return *e.iface.PrivateDnsName
	}
	return ""
}

func (e *ENI) PrivateIpAddress() string {
	if e.iface.PrivateIpAddress != nil {
		return *e.iface.PrivateIpAddress
	}
	return ""
}

func (e *ENI) Status() string {
	if e.iface.Status != nil {
		return *e.iface.Status
	}
	return ""
}

func (e *ENI) AttachmentID() string {
	if e.iface.Attachment != nil && e.iface.Attachment.AttachmentId != nil {
		return *e.iface.Attachment.AttachmentId
	}
	return ""
}

func (e *ENI) AttachedDeviceIndex() int64 {
	if e.iface.Attachment != nil && e.iface.Attachment.DeviceIndex != nil {
		return *e.iface.Attachment.DeviceIndex
	}
	return -1
}

func (e *ENI) AttachedStatus() string {
	if e.iface.Attachment != nil && e.iface.Attachment.Status != nil {
		return *e.iface.Attachment.Status
	}
	return ""
}

func (e *ENI) AttachedInstanceID() string {
	if e.iface.Attachment != nil && e.iface.Attachment.InstanceId != nil {
		return *e.iface.Attachment.InstanceId
	}
	return ""
}

func (e *ENI) AttachedInstance() *Instance {
	return e.instance
}

func (e *ENI) AvailabilityZone() string {
	if e.iface.AvailabilityZone != nil {
		return *e.iface.AvailabilityZone
	}
	return ""
}

func (e *ENI) Name() string {
	if len(e.iface.TagSet) > 0 {
		for _, tag := range e.iface.TagSet {
			if *tag.Key == "Name" {
				return *tag.Value
			}
		}
	}
	return ""
}
