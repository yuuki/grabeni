package aws

import (
	"errors"
	"time"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/ec2"
)

type Client struct {
	*ec2.EC2
}

func NewClient(region string, accessKey string, secretKey string) (*Client, error) {
	config := aws.DefaultConfig

	config.Region = region
	if config.Region == "" {
		region, err := GetRegion()
		if err != nil {
			return nil, err
		}
		config.Region = region
	}

	if accessKey != "" && secretKey != "" {
		config.Credentials = aws.Creds(accessKey, secretKey, "")
	}

	return &Client{ec2.New(config)}, nil
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

func (cli *Client) DetachENI(eniID string, instanceID string) error {
	eni, err := cli.DescribeENIByID(eniID)
	if err != nil {
		return err
	}

	return cli.DetachENIByAttachmentID(*eni.Attachment.AttachmentID)
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

		err = cli.DetachENIByAttachmentID(*eni.Attachment.AttachmentID)
		if err != nil {
			return err, false
		}

		// Sometimes detaching ENI is too slow
		time.Sleep(10 * time.Second)
	}

	err = cli.AttachENI(eniID, instanceID, deviceIndex)
	if err != nil {
		return err, false
	}

	return nil, true
}
