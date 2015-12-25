BIN = grabeni

all: clean cross test

test: testdeps
	go test -v ./...

build: deps
	go build ./...

lint: deps testdeps
	go vet
	golint

cross: deps
	goxc -tasks='xc archive' -bc 'linux,!arm windows darwin' -d . -build-ldflags "-X main.Version \"$(VERSION)\"" -resources-include='README*'
	cp -p $(PWD)/snapshot/linux_amd64/grabeni $(PWD)/snapshot/grabeni_linux_amd64
	cp -p $(PWD)/snapshot/linux_386/grabeni $(PWD)/snapshot/grabeni_linux_386
	cp -p $(PWD)/snapshot/darwin_amd64/grabeni $(PWD)/snapshot/grabeni_darwin_amd64
	cp -p $(PWD)/snapshot/darwin_386/grabeni $(PWD)/snapshot/grabeni_darwin_386
	cp -p $(PWD)/snapshot/windows_amd64/grabeni.exe $(PWD)/snapshot/grabeni_windows_amd64.exe
	cp -p $(PWD)/snapshot/windows_386/grabeni.exe $(PWD)/snapshot/grabeni_windows_386.exe

deps:
	go get -d -v .

testdeps:
	go get -d -v -t .

clean:
	go clean

.PHONY: test build cross lint deps testdeps clean

