package aws

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMetaData(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.String() != "/test" {
			t.Error("request URL should be /test but :", req.URL.String())
		}

		res.Header()["Content-Type"] = []string{"text/plain"}
		res.Header()["Accept-Range"] = []string{"bytes"}
		fmt.Fprint(res, "dummy")
	}))
	defer ts.Close()

	metadataURL = ts.URL + "/"

	result, err := GetMetaData("test")

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if string(result) != "dummy" {
		t.Error("result should be 'dummy'")
	}
}

func TestGetRegion(t *testing.T) {
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

	result, err := GetRegion()

	if err != nil {
		t.Error("err shoud be nil but: ", err)
	}

	if string(result) != "ap-northeast-1" {
		t.Error("result should be 'ap-northeast-1'")
	}
}

func TestGetInstanceID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.String() != "/instance-id" {
			t.Error("request URL should be /instance-id but :", req.URL.String())
		}

		res.Header()["Content-Type"] = []string{"text/plain"}
		res.Header()["Accept-Range"] = []string{"bytes"}
		fmt.Fprint(res, "i-abc10000")
	}))
	defer ts.Close()

	metadataURL = ts.URL + "/"

	result, err := GetInstanceID()

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if string(result) != "i-abc10000" {
		t.Error("result should be 'i-abc10000'")
	}
}

func TestIsInAws(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.String() != "/instance-id" {
			t.Error("request URL should be /instance-id but :", req.URL.String())
		}

		res.Header()["Content-Type"] = []string{"text/plain"}
		res.Header()["Accept-Range"] = []string{"bytes"}
		fmt.Fprint(res, "i-abc10000")
	}))
	defer ts.Close()

	metadataURL = ts.URL + "/"

	isInAws := IsInAws()

	if isInAws != true {
		t.Error("isInAws should be true")
	}
}
