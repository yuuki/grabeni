package aws

import (
	"errors"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/ec2"
	"github.com/codegangsta/cli"
)

type Client struct {
	*ec2.EC2
}

func NewClient(c *cli.Context) (*Client, error) {
	config := aws.DefaultConfig

	config.Region = c.GlobalString("region")
	if config.Region == "" {
		region, err := GetRegion()
		if err != nil {
			return nil, err
		}
		config.Region = region
	}

	if c.GlobalString("accesskey") != "" && c.GlobalString("secretkey") != "" {
		config.Credentials = aws.Creds(c.GlobalString("accesskey"), c.GlobalString("secretkey"), "")
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
