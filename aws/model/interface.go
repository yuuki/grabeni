package model

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

type ENI struct {
	*ec2.NetworkInterface
	instance *Instance
}

func NewENI(networkInterface *ec2.NetworkInterface) *ENI {
	return &ENI(*networkInterface)
}

func (e *ENI) SetInstance(i *Instance) {
	e.instance = Instance
}

func (e *ENI) InterfaceID() string {
	if e.NetworkInterfaceId != nil {
		return *e.NetworkInterfaceId
	}
	return ""
}

func (e *ENI) PrivateDnsName() string {
	if e.PrivateDnsName != nil {
		return *e.PrivateDnsName
	}
	return ""
}

func (e *ENI) PrivateIpAddress() string {
	if e.PrivateIpAddress != nil {
		return *e.PrivateIpAddress
	}
	return ""
}

func (e *ENI) Status() string {
	if e.Status != nil {
		return *e.Status
	}
	return ""
}

func (e *ENI) AttachedDeviceIndex() int {
	if e.Attachment != nil && e.Attachment.DeviceIndex != nil {
		return *e.Attachment.DeviceIndex
	}
	return -1
}

func (e *ENI) AttachedInstanceID() string {
	var instanceID string
	if e.Attachment != nil && e.Attachment.InstanceId != nil {
		 return *e.Attachment.InstanceId
	}
	return ""
}

func (e *ENI) AttachedInstance() string {
	return e.instance
}

func (e *ENI) AvailabilityZone() string {
	if e.AvailabilityZone != nil {
		return *e.AvailabilityZone
	}
	return ""
}

func (e *ENI) Name() string {
	if len(e.TagSet) > 0 {
		for _, tag := range e.TagSet {
			if *tag.Key == "Name" {
				return *tag.Value
			}
		}
	}
	return ""
}
