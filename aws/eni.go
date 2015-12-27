package aws

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	"github.com/yuuki1/grabeni/aws/model"
)

type ENIClient struct {
	svc       ec2iface.EC2API
	logger    *log.Logger
	logWriter io.Writer
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

type WaiterParam struct {
	MaxAttempts int
	IntervalSec int
}

func validateWaitUntilParam(p *WaiterParam) error {
	if p == nil {
		return fmt.Errorf("WaitUntilParam require")
	}

	if p.MaxAttempts <= 0 {
		return fmt.Errorf("invalid max attempts (%d)", p.MaxAttempts)
	}
	if p.IntervalSec <= 0 {
		return fmt.Errorf("invalid interval (%d) seconds", p.IntervalSec)
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

	f, _ := os.Open(os.DevNull)
	l := log.New(f, "", 0)

	return &ENIClient{svc: svc, logger: l}
}

func (c *ENIClient) WithLogWriter(w io.Writer) *ENIClient {
	c.logger.SetOutput(w)
	c.logWriter = w
	return c
}

func (c *ENIClient) DescribeENIByID(InterfaceID string) (*model.ENI, error) {
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

	eni := model.NewENI(resp.NetworkInterfaces[0])
	instance, err := c.DescribeInstanceByID(eni.AttachedInstanceID())
	if err != nil {
		return nil, err
	}
	eni.SetInstance(instance)

	return eni, nil
}

func (c *ENIClient) DescribeENIs() ([]*model.ENI, error) {
	resp, err := c.svc.DescribeNetworkInterfaces(nil)
	if err != nil {
		return nil, err
	}

	if len(resp.NetworkInterfaces) < 1 {
		return nil, nil
	}

	enis := make([]*model.ENI, 0)
	for _, iface := range resp.NetworkInterfaces {
		enis = append(enis, model.NewENI(iface))
	}

	instanceIDs := make([]string, 0)
	for _, eni := range enis {
		if id := eni.AttachedInstanceID(); id != "" {
			instanceIDs = append(instanceIDs, id)
		}
	}

	instances, err := c.DescribeInstancesByIDs(instanceIDs)
	if err != nil {
		return nil, err
	}
	if len(instances) < 1 {
		return enis, nil
	}
	//TODO make hashmap

	for _, eni := range enis {
		for _, instance := range instances {
			if eni.AttachedInstanceID() == instance.InstanceID() {
				eni.SetInstance(instance)
			}
		}
	}

	return enis, nil
}

func (c *ENIClient) AttachENI(param *AttachENIParam) (*model.ENI, error) {
	eni, err := c.DescribeENIByID(param.InterfaceID)
	if err != nil {
		return nil, err
	}

	// Do nothing if the target ENI already attached with the target instance
	if eni.AttachedInstanceID() == param.InstanceID {
		return nil, nil
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

func (c *ENIClient) AttachENIWithWaiter(p *AttachENIParam, wp *WaiterParam) (*model.ENI, error) {
	if err := validateWaitUntilParam(wp); err != nil {
		return nil, err
	}

	if eni, err := c.AttachENI(p); eni == nil || err != nil {
		return nil, err
	}

	c.logger.Printf("--> Attaching: %15s\n", p.InterfaceID)

	// Wait until attach event completed or timeout
	for i := 0; i < wp.MaxAttempts; i++ {
		fmt.Fprint(c.logWriter, ".") // use fmt.Fprint because standard log package always newline

		eni, err := c.DescribeENIByID(p.InterfaceID)
		if err != nil {
			return nil, err
		}
		if eni.Status() == "in-use" && eni.AttachedStatus() == "attached" {
			c.logger.Println()
			c.logger.Printf("--> Attached: %15s\n", p.InterfaceID)
			return eni, nil // detach completed
		}

		time.Sleep(time.Duration(wp.IntervalSec) * time.Second)
	}

	return nil, fmt.Errorf("attach %s error: over %d polling attempts", p.InterfaceID, wp.MaxAttempts)
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

func (c *ENIClient) DetachENI(param *DetachENIParam) (*model.ENI, error) {
	eni, err := c.DescribeENIByID(param.InterfaceID)
	if err != nil {
		return nil, err
	}

	if eni.Status() == "available" {
		// already detached
		return nil, nil
	}

	if err := c.DetachENIByAttachmentID(eni.AttachmentID()); err != nil {
		return nil, err
	}

	return eni, nil
}

func (c *ENIClient) DetachENIWithWaiter(p *DetachENIParam, wp *WaiterParam) (*model.ENI, error) {
	if err := validateWaitUntilParam(wp); err != nil {
		return nil, err
	}

	if eni, err := c.DetachENI(p); eni == nil || err != nil {
		return nil, err
	}

	c.logger.Printf("--> Detaching: %15s\n", p.InterfaceID)

	// Wait until detach event completed or timeout
	for i := 0; i < wp.MaxAttempts; i++ {
		fmt.Fprint(c.logWriter, ".") // use fmt.Fprint because standard log package always newline

		eni, err := c.DescribeENIByID(p.InterfaceID)
		if err != nil {
			return nil, err
		}
		if eni.Status() == "available" {
			c.logger.Println()
			c.logger.Printf("--> Detached: %15s\n", eni.InterfaceID())
			return eni, nil // detach completed
		}

		time.Sleep(time.Duration(wp.IntervalSec) * time.Second)
	}

	return nil, fmt.Errorf("detach %s error: over %d polling attempts", p.InterfaceID, wp.MaxAttempts)
}

func (c *ENIClient) GrabENI(p *GrabENIParam, wp *WaiterParam) (*model.ENI, error) {
	eni, err := c.DescribeENIByID(p.InterfaceID)
	if err != nil {
		return nil, err
	}

	// Skip detaching if the target ENI has still not attached with the other instance
	if eni.Status() == "in-use" && eni.AttachedStatus() == "attached" {
		// Do nothing if the target ENI already attached with the target instance
		if eni.AttachedInstanceID() == p.InstanceID {
			return nil, nil
		}

		if _, err := c.DetachENIWithWaiter(&DetachENIParam{InterfaceID: eni.InterfaceID()}, wp); err != nil {
			return nil, err
		}
	}

	param := AttachENIParam(*p)
	if eni, err = c.AttachENIWithWaiter(&param, wp); err != nil {
		return nil, err
	}

	return eni, nil
}

func (c *ENIClient) DescribeInstanceByID(instanceID string) (*model.Instance, error) {
	p := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
		MaxResults:  aws.Int64(1),
	}
	resp, err := c.svc.DescribeInstances(p)
	if err != nil {
		return nil, err
	}

	if len(resp.Reservations) < 1 {
		return nil, nil // Not found
	}

	instances := resp.Reservations[0].Instances

	if len(instances) < 1 {
		return nil, nil // Not found
	}

	return model.NewInstance(instances[0]), nil
}

func (c *ENIClient) DescribeInstancesByIDs(instanceIDs []string) ([]*model.Instance, error) {
	p := &ec2.DescribeInstancesInput{
		InstanceIds: aws.StringSlice(instanceIDs),
	}
	resp, err := c.svc.DescribeInstances(p)
	if err != nil {
		return nil, err
	}

	instances := make([]*model.Instance, 0)

	for _, r := range resp.Reservations {
		for _, i := range r.Instances {
			instances = append(instances, model.NewInstance(i))
		}
	}

	return instances, nil
}
