package aws

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/ec2"
)

// ENIClient is interface that implements only methods about ENI of the all methods that ec2.EC2 has
// ENIClient helps to mock ec2.EC2 class for unit testing
type ENIClient interface {
	DescribeNetworkInterfaces(*ec2.DescribeNetworkInterfacesInput) (*ec2.DescribeNetworkInterfacesOutput, error)
	AttachNetworkInterface(*ec2.AttachNetworkInterfaceInput) (*ec2.AttachNetworkInterfaceOutput, error)
	DetachNetworkInterface(*ec2.DetachNetworkInterfaceInput) (*ec2.DetachNetworkInterfaceOutput, error)
}

type Client struct {
	ENI	ENIClient
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

func NewClient(region string, accessKey string, secretKey string) (*Client, error) {
	config := aws.Config{}
	envRegion := os.Getenv("AWS_REGION")

	if region == "" && envRegion == "" {
		var err error

		config.Region, err = GetRegion()
		if err != nil {
			return nil, err
		}
	} else if region == "" && envRegion != "" {
		config.Region = envRegion
	} else {
		config.Region = region
	}

	if accessKey != "" && secretKey != "" {
		config.Credentials = aws.Creds(accessKey, secretKey, "")
	}

	return &Client{ec2.New(&config)}, nil
}

func (cli *Client) DescribeENIByID(InterfaceID string) (*ec2.NetworkInterface, error) {
	params := &ec2.DescribeNetworkInterfacesInput{
		NetworkInterfaceIDs: []*string{
			aws.String(InterfaceID),
		},
	}
	resp, err := cli.ENI.DescribeNetworkInterfaces(params)
	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		return nil, errors.New(awserr.Error())
	} else if err != nil {
		// A non-service error occurred.
		return nil, err
	}

	if len(resp.NetworkInterfaces) < 1 {
		return nil, nil
	}

	return resp.NetworkInterfaces[0], nil
}

func (cli *Client) DescribeENIs() ([]*ec2.NetworkInterface, error) {
	resp, err := cli.ENI.DescribeNetworkInterfaces(nil)
	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		return nil, errors.New(awserr.Error())
	} else if err != nil {
		// A non-service error occurred.
		return nil, err
	}

	if len(resp.NetworkInterfaces) < 1 {
		return nil, nil
	}

	return resp.NetworkInterfaces, nil
}

func (cli *Client) AttachENI(param *AttachENIParam) (*ec2.NetworkInterface, error) {
	eni, err := cli.DescribeENIByID(param.InterfaceID)
	if err != nil {
		return nil, err
	}

	// Do nothing if the target ENI already attached with the target instance
	if eni.Attachment != nil {
		if *eni.Attachment.InstanceID == param.InstanceID {
			return nil, nil
		}
	}

	input := &ec2.AttachNetworkInterfaceInput{
		NetworkInterfaceID: aws.String(param.InterfaceID),
		InstanceID:         aws.String(param.InstanceID),
		DeviceIndex:        aws.Long(int64(param.DeviceIndex)),
	}
	_, err = cli.ENI.AttachNetworkInterface(input)
	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		return nil, errors.New(awserr.Error())
	} else if err != nil {
		// A non-service error occurred.
		return nil, err
	}

	return eni, nil
}

func (cli *Client) AttachENIWithRetry(param *AttachENIParam, retryParam *RetryParam) (*ec2.NetworkInterface, error) {
	if err := validateRetryParam(retryParam); err != nil {
		return nil, err
	}

	if eni, err := cli.AttachENI(param); eni == nil || err != nil {
		return nil, err
	}

	// Retry until attach event completed or timeout
	for {
		select {
		case <-time.After(time.Duration(retryParam.TimeoutSec) * time.Second):
			return nil, fmt.Errorf("timeout occured. %d seconds elapsed.", retryParam.TimeoutSec)
		case <-time.Tick(time.Duration(retryParam.IntervalSec) * time.Second):
			eni, err := cli.DescribeENIByID(param.InterfaceID)
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

func (cli *Client) DetachENIByAttachmentID(attachmentID string) error {
	params := &ec2.DetachNetworkInterfaceInput{
		AttachmentID: aws.String(attachmentID),
		Force:        aws.Boolean(false),
	}
	_, err := cli.ENI.DetachNetworkInterface(params)
	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		return errors.New(awserr.Error())
	} else if err != nil {
		// A non-service error occurred.
		return err
	}

	return nil
}

func (cli *Client) DetachENI(param *DetachENIParam) (*ec2.NetworkInterface, error) {
	eni, err := cli.DescribeENIByID(param.InterfaceID)
	if err != nil {
		return nil, err
	}

	if *eni.Status == "available" {
		// already detached
		return nil, nil
	}

	if err := cli.DetachENIByAttachmentID(*eni.Attachment.AttachmentID); err != nil {
		return nil, err
	}

	return eni, nil
}

func (cli *Client) DetachENIWithRetry(param *DetachENIParam, retryParam *RetryParam) (*ec2.NetworkInterface, error) {
	if err := validateRetryParam(retryParam); err != nil {
		return nil, err
	}

	if eni, err := cli.DetachENI(param); eni == nil || err != nil {
		return nil, err
	}

	// Retry until detach event completed or timeout
	for {
		select {
		case <-time.After(time.Duration(retryParam.TimeoutSec) * time.Second):
			return nil, fmt.Errorf("timeout occured. %d seconds elapsed.", retryParam.TimeoutSec)
		case <-time.Tick(time.Duration(retryParam.IntervalSec) * time.Second):
			eni, err := cli.DescribeENIByID(param.InterfaceID)
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

func (cli *Client) GrabENI(param *GrabENIParam, retryParam *RetryParam) (*ec2.NetworkInterface, error) {
	eni, err := cli.DescribeENIByID(param.InterfaceID)
	if err != nil {
		return nil, err
	}

	// Skip detaching if the target ENI has still not attached with the other instance
	if *eni.Status == "in-use" && eni.Attachment != nil && *eni.Attachment.Status == "attached" {
		// Do nothing if the target ENI already attached with the target instance
		if *eni.Attachment.InstanceID == param.InstanceID {
			return nil, nil
		}

		if _, err := cli.DetachENIWithRetry(&DetachENIParam{InterfaceID: param.InterfaceID}, retryParam); err != nil {
			return nil, err
		}
	}

	p := AttachENIParam(*param)
	if eni, err = cli.AttachENIWithRetry(&p, retryParam); err != nil {
		return nil, err
	}

	return eni, nil
}
