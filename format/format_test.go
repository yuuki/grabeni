package format

import (
	"bytes"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"

	"github.com/yuuki/grabeni/aws/model"
)

func TestPrintENI(t *testing.T) {
	w := new(bytes.Buffer)
	PrintENI(w, nil)

	expected := "ID\tNAME\tSTATUS\tPRIVATE DNS NAMEPRIVATE IP\tAZ\tDEVICE INDEX\tINSTANCE ID\tINSTANCE NAME\n"
	assert.Equal(t, expected, string(w.Bytes()))

	w = new(bytes.Buffer)
	eni := model.NewENI(&ec2.NetworkInterface{})
	PrintENI(w, eni)

	expected = "ID\tNAME\tSTATUS\tPRIVATE DNS NAMEPRIVATE IP\tAZ\tDEVICE INDEX\tINSTANCE ID\tINSTANCE NAME\n\t\t\t\t\t\t\t\t-1\t\t\t\t\n"
	assert.Equal(t, expected, string(w.Bytes()))
}

func TestPrintENIs(t *testing.T) {
	w := new(bytes.Buffer)
	PrintENIs(w, nil)

	expected := "ID\tNAME\tSTATUS\tPRIVATE DNS NAMEPRIVATE IP\tAZ\tDEVICE INDEX\tINSTANCE ID\tINSTANCE NAME\n"
	assert.Equal(t, expected, string(w.Bytes()))

	w = new(bytes.Buffer)
	enis := make([]*model.ENI, 2)
	enis[0] = model.NewENI(&ec2.NetworkInterface{})

	i := model.NewInstance(&ec2.Instance{
		InstanceId: aws.String("i-1000000"),
		Tags: []*ec2.Tag{{
			Key:   aws.String("Name"),
			Value: aws.String("grabeni001"),
		}},
	})
	enis[0].SetInstance(i)
	PrintENIs(w, enis)

	expected = "ID\tNAME\tSTATUS\tPRIVATE DNS NAMEPRIVATE IP\tAZ\tDEVICE INDEX\tINSTANCE ID\tINSTANCE NAME\n\t\t\t\t\t\t\t\t-1\t\ti-1000000\tgrabeni001\n"
	assert.Equal(t, expected, string(w.Bytes()))
}
