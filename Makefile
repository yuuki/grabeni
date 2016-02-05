BIN = grabeni

all: clean build test

test: testdeps gen
	go test -v ./...

gen:
	go get github.com/vektra/mockery/.../
	mockery -all -dir ${GOPATH}/src/github.com/aws/aws-sdk-go/service/ec2/ec2iface -print | perl -pe 's/^package mocks/package aws/' > aws/mock_ec2api.go

build: deps gen
	go build -o $(BIN) ./cmd

lint: deps testdeps
	go vet
	golint

patch: gobump
	./script/release.sh patch

minor: gobump
	./script/release.sh minor

gobump:
	go get github.com/motemen/gobump/cmd/gobump

deps:
	go get -d -v ./...

testdeps:
	go get -d -v -t ./...

clean:
	go clean

.PHONY: test build lint deps testdeps clean

