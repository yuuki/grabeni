package aws

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type ENIClient struct {
	svc ec2iface.EC2API
}

type AttachENIParam struct {
	InterfaceID string
	InstanceID  string
	DeviceIndex int
}

type DetachENIParam struct {
	InterfaceID string
}

type GrabENIParam AttachENIParam

type RetryParam struct {
	TimeoutSec  int64
	IntervalSec int64
}

func validateRetryParam(param *RetryParam) error {
	if param == nil {
		return fmt.Errorf("RetryParam require")
	}

	if param.TimeoutSec <= 0 {
		return fmt.Errorf("invalid timeout (%d) seconds", param.TimeoutSec)
	}
	if param.IntervalSec <= 0 {
		return fmt.Errorf("invalid interval (%d) seconds", param.IntervalSec)
	}
	if param.TimeoutSec < param.IntervalSec {
		return fmt.Errorf(
			"interval (%d) should be less than timeout (%d) seconds",
			param.IntervalSec,
			param.TimeoutSec,
		)
	}

	return nil
}

func NewENIClient() *ENIClient {
	session := session.New()
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region, _ = NewMetaDataClientFromSession(session).GetRegion()
	}
	svc := ec2.New(session, &aws.Config{Region: aws.String(region)})
	return &ENIClient{svc: svc}
}

func (c *ENIClient) DescribeENIByID(InterfaceID string) (*ec2.NetworkInterface, error) {
	params := &ec2.DescribeNetworkInterfacesInput{
		NetworkInterfaceIds: []*string{
			aws.String(InterfaceID),
		},
	}
	resp, err := c.svc.DescribeNetworkInterfaces(params)
	if err != nil {
		return nil, err
	}

	if len(resp.NetworkInterfaces) < 1 {
		return nil, nil
	}

	return resp.NetworkInterfaces[0], nil
}

func (c *ENIClient) DescribeENIs() ([]*ec2.NetworkInterface, error) {
	resp, err := c.svc.DescribeNetworkInterfaces(nil)
	if err != nil {
		return nil, err
	}

	if len(resp.NetworkInterfaces) < 1 {
		return nil, nil
	}

	return resp.NetworkInterfaces, nil
}

func (c *ENIClient) AttachENI(param *AttachENIParam) (*ec2.NetworkInterface, error) {
	eni, err := c.DescribeENIByID(param.InterfaceID)
	if err != nil {
		return nil, err
	}

	// Do nothing if the target ENI already attached with the target instance
	if eni.Attachment != nil {
		if *eni.Attachment.InstanceId == param.InstanceID {
			return nil, nil
		}
	}

	input := &ec2.AttachNetworkInterfaceInput{
		NetworkInterfaceId: aws.String(param.InterfaceID),
		InstanceId:         aws.String(param.InstanceID),
		DeviceIndex:        aws.Int64(int64(param.DeviceIndex)),
	}
	_, err = c.svc.AttachNetworkInterface(input)
	if err != nil {
		return nil, err
	}

	return eni, nil
}

func (c *ENIClient) AttachENIWithRetry(param *AttachENIParam, retryParam *RetryParam) (*ec2.NetworkInterface, error) {
	if err := validateRetryParam(retryParam); err != nil {
		return nil, err
	}

	if eni, err := c.AttachENI(param); eni == nil || err != nil {
		return nil, err
	}

	// Retry until attach event completed or timeout
	for {
		select {
		case <-time.After(time.Duration(retryParam.TimeoutSec) * time.Second):
			return nil, fmt.Errorf("timeout occured. %d seconds elapsed.", retryParam.TimeoutSec)
		case <-time.Tick(time.Duration(retryParam.IntervalSec) * time.Second):
			eni, err := c.DescribeENIByID(param.InterfaceID)
			if err != nil {
				return nil, err
			}
			if *eni.Status == "in-use" && eni.Attachment != nil && *eni.Attachment.Status == "attached" {
				return eni, nil // detach completed
			}
		}
	}

	return nil, nil
}

func (c *ENIClient) DetachENIByAttachmentID(attachmentID string) error {
	params := &ec2.DetachNetworkInterfaceInput{
		AttachmentId: aws.String(attachmentID),
		Force:        aws.Bool(false),
	}
	_, err := c.svc.DetachNetworkInterface(params)
	if err != nil {
		return err
	}

	return nil
}

func (c *ENIClient) DetachENI(param *DetachENIParam) (*ec2.NetworkInterface, error) {
	eni, err := c.DescribeENIByID(param.InterfaceID)
	if err != nil {
		return nil, err
	}

	if *eni.Status == "available" {
		// already detached
		return nil, nil
	}

	if err := c.DetachENIByAttachmentID(*eni.Attachment.AttachmentId); err != nil {
		return nil, err
	}

	return eni, nil
}

func (c *ENIClient) DetachENIWithRetry(param *DetachENIParam, retryParam *RetryParam) (*ec2.NetworkInterface, error) {
	if err := validateRetryParam(retryParam); err != nil {
		return nil, err
	}

	if eni, err := c.DetachENI(param); eni == nil || err != nil {
		return nil, err
	}

	// Retry until detach event completed or timeout
	for {
		select {
		case <-time.After(time.Duration(retryParam.TimeoutSec) * time.Second):
			return nil, fmt.Errorf("timeout occured. %d seconds elapsed.", retryParam.TimeoutSec)
		case <-time.Tick(time.Duration(retryParam.IntervalSec) * time.Second):
			eni, err := c.DescribeENIByID(param.InterfaceID)
			if err != nil {
				return nil, err
			}
			if *eni.Status == "available" {
				return eni, nil // detach completed
			}
		}
	}

	return nil, nil
}

func (c *ENIClient) GrabENI(param *GrabENIParam, retryParam *RetryParam) (*ec2.NetworkInterface, error) {
	eni, err := c.DescribeENIByID(param.InterfaceID)
	if err != nil {
		return nil, err
	}

	// Skip detaching if the target ENI has still not attached with the other instance
	if *eni.Status == "in-use" && eni.Attachment != nil && *eni.Attachment.Status == "attached" {
		// Do nothing if the target ENI already attached with the target instance
		if *eni.Attachment.InstanceId == param.InstanceID {
			return nil, nil
		}

		if _, err := c.DetachENIWithRetry(&DetachENIParam{InterfaceID: param.InterfaceID}, retryParam); err != nil {
			return nil, err
		}
	}

	p := AttachENIParam(*param)
	if eni, err = c.AttachENIWithRetry(&p, retryParam); err != nil {
		return nil, err
	}

	return eni, nil
}
