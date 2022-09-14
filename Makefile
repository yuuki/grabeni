BIN = grabeni
PKGS = $$(go list ./... | grep -v vendor)

all: clean build test

test: gen
	go test -v $(PKGS)

gen:
	go install github.com/vektra/mockery/v2@latest
	mockery --all --dir vendor/github.com/aws/aws-sdk-go/service/ec2/ec2iface --print | perl -pe 's/^package mocks/package aws/' > aws/mock_ec2api.go

build: gen
	go build -o $(BIN) ./cmd/...

lint:
	go vet
	golint

patch:
	./script/release.sh patch

minor:
	./script/release.sh minor

clean:
	go clean

.PHONY: test build lint deps testdeps clean

