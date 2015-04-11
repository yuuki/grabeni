package aws

import (
	"errors"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/ec2"
)

func DescribeENIByID(eniID string) (*ec2.NetworkInterface, error) {
	region, err := GetRegion()
	if err != nil {
		return nil, err
	}

	svc := ec2.New(&aws.Config{Region: region})

	params := &ec2.DescribeNetworkInterfacesInput{
		NetworkInterfaceIDs: []*string{
			aws.String(eniID),
		},
	}
	resp, err := svc.DescribeNetworkInterfaces(params)
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
