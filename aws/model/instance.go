package model

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Instance ec2.Instance

func NewInstance(i *ec2.Instance) *Instance {
	instance := Instance(*i)
	return &instance
}

func (i *Instance) InstanceID() string {
	return *i.InstanceId
}

func (i *Instance) Name() string {
	if len(i.Tags) > 0 {
		for _, tag := range i.Tags {
			if *tag.Key == "Name" {
				return *tag.Value
			}
		}
	}
	return ""
}

