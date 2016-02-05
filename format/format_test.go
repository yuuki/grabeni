package format

import (
	"bytes"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"

	"github.com/yuuki1/grabeni/aws/model"
)

func TestPrintENI(t *testing.T) {
	w := new(bytes.Buffer)
	PrintENI(w, nil)

	expected := "ID\tNAME\tSTATUS\tPRIVATE DNS NAMEPRIVATE IP\tAZ\tDEVICE INDEX\tINSTANCE ID\tINSTANCE NAME\n"
	assert.Equal(t, expected, string(w.Bytes()))

	eni := model.NewENI(&ec2.NetworkInterface{})
	PrintENI(w, eni)

	expected = "ID\tNAME\tSTATUS\tPRIVATE DNS NAMEPRIVATE IP\tAZ\tDEVICE INDEX\tINSTANCE ID\tINSTANCE NAME\nID\tNAME\tSTATUS\tPRIVATE DNS NAMEPRIVATE IP\tAZ\tDEVICE INDEX\tINSTANCE ID\tINSTANCE NAME\n\t\t\t\t\t\t\t\t-1\t\t\t\t\n"
	assert.Equal(t, expected, string(w.Bytes()))
}
