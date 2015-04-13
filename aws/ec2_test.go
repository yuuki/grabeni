package aws

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	// case 1
	{
		os.Clearenv()

		cli, err := NewClient("ap-northeast-1", "dummyaccess", "dummysecret")
		if err != nil {
			t.Error("err should be nil but: ", err)
		}

		config := cli.EC2.Service.Config
		if config.Region != "ap-northeast-1" {
			t.Error("Region should be 'ap-northeast-1' but: ", config.Region)
		}

		cred, err := config.Credentials.Credentials()
		if err != nil {
			t.Error("err should be nil but: ", err)
		}
		if cred.AccessKeyID != "dummyaccess" {
			t.Error("AccessKeyID should be 'dummyaccess' but: ", cred.AccessKeyID)
		}
		if cred.SecretAccessKey != "dummysecret" {
			t.Error("SecretAccessKey should be 'dummysecret' but: ", cred.SecretAccessKey)
		}
	}

	// case 2
	{

		os.Clearenv()
		os.Setenv("AWS_REGION", "us-west-2")

		cli, err := NewClient("", "dummyaccess", "dummysecret")
		if err != nil {
			t.Error("err should be nil but: ", err)
		}

		config := cli.EC2.Service.Config
		if config.Region != "us-west-2" {
			t.Error("Region should be 'us-west-2' but:", config.Region)
		}
	}

	// case 3
	{
		os.Clearenv()

		ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if req.URL.String() != "/placement/availability-zone" {
				t.Error("request URL should be /placement/availability-zone but :", req.URL.String())
			}

			res.Header()["Content-Type"] = []string{"text/plain"}
			res.Header()["Accept-Range"] = []string{"bytes"}
			fmt.Fprint(res, "ap-northeast-1a")
		}))
		defer ts.Close()

		metadataURL = ts.URL + "/"

		cli, err := NewClient("", "dummyaccess", "dummysecret")
		if err != nil {
			t.Error("err should be nil but: ", err)
		}

		config := cli.EC2.Service.Config
		if config.Region != "ap-northeast-1" {
			t.Error("Region should be 'ap-northeast-1' but:", config.Region)
		}
	}

	// case 4
	{
		os.Clearenv()
		os.Setenv("AWS_ACCESS_KEY_ID", "dummyaccess")
		os.Setenv("AWS_SECRET_ACCESS_KEY", "dummysecret")

		cli, err := NewClient("ap-northeast-1", "", "")
		if err != nil {
			t.Error("err should be nil but: ", err)
		}

		config := cli.EC2.Service.Config
		if config.Region != "ap-northeast-1" {
			t.Error("Region should be 'ap-northeast-1' but: ", config.Region)
		}

		cred, err := config.Credentials.Credentials()
		if err != nil {
			t.Error("err should be nil but: ", err)
		}
		if cred.AccessKeyID != "dummyaccess" {
			t.Error("AccessKeyID should be 'dummyaccess' but: ", cred.AccessKeyID)
		}
		if cred.SecretAccessKey != "dummysecret" {
			t.Error("SecretAccessKey should be 'dummysecret' but: ", cred.SecretAccessKey)
		}
	}

	//TODO case 5
	{
		os.Clearenv()

		// cli, err := NewClient("", "", "")
		// if err != nil {
		// 	t.Error("err should be nil but: ", err)
		// }
	}
}
