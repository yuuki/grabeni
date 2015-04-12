package aws

// https://github.com/pdalinis/aws-sdk-go/blob/153_Metadata/aws/metaData.go

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	metadataURL = "http://169.254.169.254/latest/meta-data/"
)

func GetMetaData(path string) (contents []byte, err error) {
	url := metadataURL + path

	req, _ := http.NewRequest("GET", url, nil)
	client := http.Client{
		Timeout: time.Millisecond * 500,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("Code %d returned for url %s", resp.StatusCode, url)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return []byte(body), nil
}

// Gets the Region that the instance is running in.
func GetRegion() (string, error) {
	path := "placement/availability-zone"
	resp, err := GetMetaData(path)
	if err != nil {
		return "", err
	}

	az := string(resp)

	//returns us-west-2a, just return us-west-2
	return string(az[:len(az)-1]), nil
}

func GetInstanceID() (string, error) {
	path := "instance-id"
	resp, err := GetMetaData(path)
	if err != nil {
		return "", err
	}

	return string(resp), nil
}

func IsInAws() bool {
	path := "instance-id"
	_, err := GetMetaData(path)

	if err != nil {
		return false
	}

	return true
}
