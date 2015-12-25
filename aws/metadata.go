package aws

import (
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

type MetaDataClient struct {
	svc *ec2metadata.EC2Metadata
}

func NewMetaDataClient() *MetaDataClient {
	return &MetaDataClient{svc: ec2metadata.New(session.New())}
}

func NewMetaDataClientFromSession(s *session.Session) *MetaDataClient {
	return &MetaDataClient{svc: ec2metadata.New(s)}
}

func (c *MetaDataClient) GetInstanceID() (string, error) {
	return c.svc.GetMetadata("instance-id")
}

func (c *MetaDataClient) GetRegion() (string, error) {
	return c.svc.Region()
}
