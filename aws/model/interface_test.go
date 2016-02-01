package model

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
)

func TestSetInstance(t *testing.T) {
	i := NewInstance(&ec2.Instance{
		InstanceId: aws.String("i-1000000"),
	})
	eni := NewENI(&ec2.NetworkInterface{})

	eni.SetInstance(i)

	assert.Equal(t, eni.instance, i)
}

