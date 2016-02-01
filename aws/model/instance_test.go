package model

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
)

func TestInstanceID(t *testing.T) {
	i := NewInstance(&ec2.Instance{
		InstanceId: aws.String("i-1000000"),
	})

	assert.Equal(t, i.InstanceID(), "i-1000000", "InstanceID should be i-1000000")
}
