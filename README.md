grabeni
=======
[![Latest Version](http://img.shields.io/github/release/yuuki/grabeni.svg?style=flat-square)][release]
[![Build Status](http://img.shields.io/travis/yuuki/grabeni.svg?style=flat-square)][travis]
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]
[![Go Report Card](https://goreportcard.com/badge/github.com/yuuki/grabeni)][goreport]
[![License](http://img.shields.io/:license-mit-blue.svg)][license]

[release]: https://github.com/yuuki/grabeni/releases
[travis]: http://travis-ci.org/yuuki/grabeni
[godocs]: http://godoc.org/github.com/yuuki/grabeni
[goreport]: https://goreportcard.com/report/github.com/yuuki/grabeni
[license]: http://doge.mit-license.org

`grabeni` is an operation-friendly tool to grab an EC2 Network Interface (ENI) from an other EC2 instance.

Detaching and attaching (grabbing) ENI is a common way to realize VIP (Virtual IP address) in EC2 with Heartbeat, Keepalived, MHA, etc.
`grabeni` provides command line interface for handling ENI.

The list of `grabeni`'s features below.

- Listing ENI information.
- Attacing the specified ENI to the specified instance.
- Detaching the specified ENI.
- Grabbing (Attaching and Detaching) the specified ENI to the specified instance.
- timeout/retry for requesting AWS API.

## Requirements
```
ec2:DescribeInstances
ec2:DescribeNetworkInterfaces
ec2:AttachNetworkInterface
ec2:DetachNetworkInterface
```

## Usage

```bash
$ export AWS_ACCESS_KEY_ID='...'
$ export AWS_SECRET_ACCESS_KEY='...'
$ export AWS_REGION='us-east-1'
$ grabeni grab eni-xxxxxx --instanceid i-xxxxxxd # attach eni-xxxxxx to EC2 instance where grabeni runs if instanceid option is skipped
```

See also `grabeni --help`.

## Example

```bash
$ grabeni ls
ID            NAME    STATUS      PRIVATE DNS NAME                              PRIVATE IP  AZ              DEVICE INDEX    INSTANCE ID INSTANCE NAME
eni-00000000  eni01   in-use      ip-10-0-0-100.ap-northeast-1.compute.internal 10.0.0.100  ap-northeast-1b 0   i-00000000  instance01
eni-11111111  eni02   available   ip-10-0-0-10.ap-northeast-1.compute.internal	10.0.0.10   ap-northeast-1c -1
eni-22222222  eni03   avaolable   ip-10-0-0-11.ap-northeast-1.compute.internal	10.0.0.11   ap-northeast-1c 1

$ grabeni status eni-2222222
ID            NAME    STATUS      PRIVATE DNS NAME                              PRIVATE IP  AZ              DEVICE INDEX    INSTANCE ID INSTANCE NAME
eni-22222222  eni03   avaolable   ip-10-0-0-11.ap-northeast-1.compute.internal	10.0.0.11   ap-northeast-1c 1

$ grabeni grab eni-2222222
--> Detaching:    eni-2222222
--> Attaching:    eni-2222222
eni eni-2222222 attached to instance i-xxxxxx
```

## Installation

### Homebrew
```bash
$ brew tap yuuki/grabeni
$ brew install grabeni
```

### Download binary from GitHub Releases
[Releases・yuuki/grabeni - GitHub](https://github.com/yuuki/grabeni/releases)

### Build from source
```bash
 $ go get github.com/yuuki/grabeni
 $ go install github.com/yuuki/grabeni/...
```

## Roadmap

- `attach`, `detach`, `grab`: Show ENI information before execution
- `list`: Filter option
- Add `check` command to check an availability zone
- `attach`, `detach`, `grab`: dryrun option

## Contribution

1. Fork ([https://github.com/y_uuki/grabeni/fork](https://github.com/y_uuki/grabeni/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create new Pull Request

## License

[The MIT License](./LICENSE).

## Author

[y_uuki](https://github.com/yuuki)
