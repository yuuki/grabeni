BIN = grabeni

all: clean build test

test: testdeps
	go test -v ./...

build: deps
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

