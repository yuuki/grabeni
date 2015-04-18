package aws

import (
	"errors"
	"fmt"
	"time"
	"os"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/ec2"
)

type Client struct {
	*ec2.EC2
	Timeout time.Duration
	Interval time.Duration
}

func NewClient(region string, accessKey string, secretKey string, timeoutSec int, intervalSec int) (*Client, error) {
	if timeoutSec <= 0 {
		return nil, errors.New("invalid timeout")
	}
	if intervalSec <= 0 {
		return nil, errors.New("invalid interval")
	}
	if timeoutSec < intervalSec {
		return nil, errors.New("interval should be less than timeout")
	}

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

	return &Client{ec2.New(&config), time.Duration(timeoutSec) * time.Second, time.Duration(intervalSec) * time.Second}, nil
}

func (cli *Client) DescribeENIByID(eniID string) (*ec2.NetworkInterface, error) {
	params := &ec2.DescribeNetworkInterfacesInput{
		NetworkInterfaceIDs: []*string{
			aws.String(eniID),
		},
	}
	resp, err := cli.EC2.DescribeNetworkInterfaces(params)
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
	resp, err := cli.EC2.DescribeNetworkInterfaces(nil)
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

func (cli *Client) AttachENI(eniID string, instanceID string, deviceIndex int) error {
	params := &ec2.AttachNetworkInterfaceInput{
		NetworkInterfaceID: aws.String(eniID),
		InstanceID:         aws.String(instanceID),
		DeviceIndex:        aws.Long(int64(deviceIndex)),
	}
	_, err := cli.EC2.AttachNetworkInterface(params)

	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		return errors.New(awserr.Error())
	} else if err != nil {
		// A non-service error occurred.
		return err
	}

	return nil
}

func (cli *Client) AttachENIWithRetry(eniID string, instanceID string, deviceIndex int) error {
	if err := cli.AttachENI(eniID, instanceID, deviceIndex); err != nil {
		return err
	}

	// Retry until attach event completed or timeout
	for {
		select {
		case <-time.After(cli.Timeout):
			return errors.New(fmt.Sprintf("timeout occured. %d seconds elapsed.", cli.Timeout))
		case <-time.Tick(cli.Interval):
			eni, err := cli.DescribeENIByID(eniID)
			if err != nil {
				return err
			}
			if *eni.Status == "in-use" && eni.Attachment != nil && *eni.Attachment.Status == "attached" {
				return nil // detach completed
			}
		}
	}

	return nil
}

func (cli *Client) DetachENIByAttachmentID(attachmentID string) error {
	params := &ec2.DetachNetworkInterfaceInput{
		AttachmentID: aws.String(attachmentID),
		Force:        aws.Boolean(false),
	}
	_, err := cli.EC2.DetachNetworkInterface(params)

	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		return errors.New(awserr.Error())
	} else if err != nil {
		// A non-service error occurred.
		return err
	}

	return nil
}

func (cli *Client) DetachENI(eniID string) error {
	eni, err := cli.DescribeENIByID(eniID)
	if err != nil {
		return err
	}

	if eni.Attachment == nil {
		// already detached
		return nil
	}

	return cli.DetachENIByAttachmentID(*eni.Attachment.AttachmentID)
}

func (cli *Client) DetachENIWithRetry(eniID string) error {
	if err := cli.DetachENI(eniID); err != nil {
		return err
	}

	// Retry until detach event completed or timeout
	for {
		select {
		case <-time.After(cli.Timeout):
			return errors.New(fmt.Sprintf("timeout occured. %d seconds elapsed.", cli.Timeout))
		case <-time.Tick(cli.Interval):
			eni, err := cli.DescribeENIByID(eniID)
			if err != nil {
				return err
			}
			if *eni.Status == "available" {
				return nil // detach completed
			}
		}
	}

	return nil
}

func (cli *Client) GrabENI(eniID string, instanceID string, deviceIndex int) (error, bool) {
	eni, err := cli.DescribeENIByID(eniID)
	if err != nil {
		return err, false
	}

	// Skip detaching if the target ENI has still not attached with the other instance
	if eni.Attachment != nil {
		// Do nothing if the target ENI already attached with the target instance
		if *eni.Attachment.InstanceID == instanceID {
			return nil, false
		}

		err = cli.DetachENIWithRetry(eniID)
		if err != nil {
			return err, false
		}
	}

	err = cli.AttachENIWithRetry(eniID, instanceID, deviceIndex)
	if err != nil {
		return err, false
	}

	return nil, true
}
