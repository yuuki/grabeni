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

func TestInterfaceID(t *testing.T) {
	eni := NewENI(&ec2.NetworkInterface{
		NetworkInterfaceId: aws.String("eni-2222222"),
	})

	assert.Equal(t, eni.InterfaceID(), "eni-2222222")
}

func TestPrivateDnsName(t *testing.T) {
	eni := NewENI(&ec2.NetworkInterface{
		NetworkInterfaceId: aws.String("eni-2222222"),
		PrivateDnsName: aws.String("ip-10-0-0-100.ap-northeast-1.compute.internal"),
	})

	assert.Equal(t, eni.PrivateDnsName(), "ip-10-0-0-100.ap-northeast-1.compute.internal")
}

func TestPrivateIpAddress(t *testing.T) {
	eni := NewENI(&ec2.NetworkInterface{
		NetworkInterfaceId: aws.String("eni-2222222"),
		PrivateIpAddress: aws.String("10.0.0.100"),
	})

	assert.Equal(t, eni.PrivateIpAddress(), "10.0.0.100")
}

func TestStatus(t *testing.T) {
	eni := NewENI(&ec2.NetworkInterface{
		NetworkInterfaceId: aws.String("eni-2222222"),
		Status: aws.String("in-use"),
	})

	assert.Equal(t, eni.Status(), "in-use")
}

func TestAttachmentID(t *testing.T) {
	eni := NewENI(&ec2.NetworkInterface{
		NetworkInterfaceId: aws.String("eni-2222222"),
		Attachment: &ec2.NetworkInterfaceAttachment{
			AttachmentId: aws.String("eni-attach-11111111"),
		},
	})

	assert.Equal(t, eni.AttachmentID(), "eni-attach-11111111")

	eni = NewENI(&ec2.NetworkInterface{
		NetworkInterfaceId: aws.String("eni-2222222"),
		Attachment: nil,
	})

	assert.Equal(t, eni.AttachmentID(), "")

	eni = NewENI(&ec2.NetworkInterface{
		NetworkInterfaceId: aws.String("eni-2222222"),
		Attachment: &ec2.NetworkInterfaceAttachment{
			AttachmentId: nil,
		},
	})

	assert.Equal(t, eni.AttachmentID(), "")
}

func TestAttachedDeviceIndex(t *testing.T) {
	eni := NewENI(&ec2.NetworkInterface{
		NetworkInterfaceId: aws.String("eni-2222222"),
		Attachment: &ec2.NetworkInterfaceAttachment{
			DeviceIndex: aws.Int64(0),
		},
	})

	assert.Equal(t, eni.AttachedDeviceIndex(), int64(0))

	eni = NewENI(&ec2.NetworkInterface{
		NetworkInterfaceId: aws.String("eni-2222222"),
		Attachment: nil,
	})

	assert.Equal(t, eni.AttachedDeviceIndex(), int64(-1))

	eni = NewENI(&ec2.NetworkInterface{
		NetworkInterfaceId: aws.String("eni-2222222"),
		Attachment: &ec2.NetworkInterfaceAttachment{
			DeviceIndex: nil,
		},
	})

	assert.Equal(t, eni.AttachedDeviceIndex(), int64(-1))
}

func TestAttachedStatus(t *testing.T) {
	eni := NewENI(&ec2.NetworkInterface{
		NetworkInterfaceId: aws.String("eni-2222222"),
		Attachment: &ec2.NetworkInterfaceAttachment{
			Status: aws.String("in-use"),
		},
	})

	assert.Equal(t, eni.AttachedStatus(), "in-use")

	eni = NewENI(&ec2.NetworkInterface{
		NetworkInterfaceId: aws.String("eni-2222222"),
		Attachment: nil,
	})

	assert.Equal(t, eni.AttachedStatus(), "")

	eni = NewENI(&ec2.NetworkInterface{
		NetworkInterfaceId: aws.String("eni-2222222"),
		Attachment: &ec2.NetworkInterfaceAttachment{
			Status: nil,
		},
	})

	assert.Equal(t, eni.AttachedStatus(), "")
}

func TestAttachedInstanceID(t *testing.T) {
	eni := NewENI(&ec2.NetworkInterface{
		NetworkInterfaceId: aws.String("eni-2222222"),
		Attachment: &ec2.NetworkInterfaceAttachment{
			InstanceId: aws.String("i-1111111"),
		},
	})

	assert.Equal(t, eni.AttachedInstanceID(), "i-1111111")

	eni = NewENI(&ec2.NetworkInterface{
		NetworkInterfaceId: aws.String("eni-2222222"),
		Attachment: nil,
	})

	assert.Equal(t, eni.AttachedInstanceID(), "")

	eni = NewENI(&ec2.NetworkInterface{
		NetworkInterfaceId: aws.String("eni-2222222"),
		Attachment: &ec2.NetworkInterfaceAttachment{
			InstanceId: nil,
		},
	})

	assert.Equal(t, eni.AttachedInstanceID(), "")
}

func TestAttachedInstance(t *testing.T) {
	i := NewInstance(&ec2.Instance{
		InstanceId: aws.String("i-1000000"),
	})
	eni := NewENI(&ec2.NetworkInterface{})

	eni.SetInstance(i)

	assert.Equal(t, *eni.AttachedInstance().InstanceId, "i-1000000")
}
