FROM golang:1.5.2

RUN go get  github.com/golang/lint/golint \
            golang.org/x/tools/cmd/vet \
            github.com/tools/godep \
	          github.com/axw/gocov/gocov \
            github.com/laher/goxc

ENV USER root
WORKDIR /go/src/github.com/yuuki1/grabeni

ADD . /go/src/github.com/yuuki1/grabeni
