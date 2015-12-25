grabeni
=======
[![Latest Version](http://img.shields.io/github/release/yuuki1/grabeni.svg?style=flat-square)][release]
[![Build Status](http://img.shields.io/travis/yuuki1/grabeni.svg?style=flat-square)][travis]
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]

[release]: https://github.com/yuuki1/grabeni/releases
[travis]: http://travis-ci.org/yuuki1/grabeni
[godocs]: http://godoc.org/github.com/yuuki1/grabeni

`grabeni` is an operation-friendly tool to grab an EC2 Network Interface (ENI) from an other EC2 instance.

Detaching and attaching (grabbing) ENI is a common way to realize VIP (Virtual IP address) in EC2 with Heartbeat, Keepalived, MHA, etc.
`grabeni` provides command line interface for handling ENI.

The list of `grabeni`'s features below.

- Listing ENI information.
- Attacing the specified ENI to the specified instance.
- Detaching the specified ENI.
- Grabbing (Attaching and Detaching) the specified ENI to the specified instance.
- timeout/retry for requesting AWS API.

## Usage

```bash
$ export AWS_ACCESS_KEY_ID='...'
$ export AWS_SECRET_ACCESS_KEY='...'
$ export AWS_REGION='us-east-1'
$ grabeni grab eni-xxxxxx --instanceid i-xxxxxxd # attach eni-xxxxxx to EC2 instance where grabeni runs if instanceid option is skipped
```

See also `granebi --help`.

## Example

```bash
$ grabeni ls
NetworkInterfaceID	PrivateDNSName				PrivateIPAddress	InstanceID	DeviceIndex	Status	Name
eni-0000000	  ip-10-0-0-100.ap-northeast-1.compute.internal	10.0.0.100	   0		in-use
eni-1111111		ip-10-0-0-10.ap-northeast-1.compute.internal	10.0.0.10			-1		available	eni01
eni-2222222		ip-10-0-0-11.ap-northeast-1.compute.internal	10.0.0.11	     1		in-use	eni02

$ grabeni status eni-2222222
NetworkInterfaceID	PrivateDNSName				PrivateIPAddress	InstanceID	DeviceIndex	Status	Name
eni-2222222		ip-10-0-0-11.ap-northeast-1.compute.internal	10.0.0.11	     1		in-use	eni02

$ grabeni grab eni-2222222
grabbed: eni eni-2222222 attached to instance i-xxxxxx
```

## Installation

### Homebrew
```bash
$ brew tap yuuki1/grabeni
$ brew install grabeni
```

### Download binary from GitHub Releases
[Releases・yuuki1/grabeni - GitHub](https://github.com/yuuki1/grabeni/releases)

### Build from source
```bash
 $ go get github.com/yuuki1/grabeni
 $ go install github.com/yuuki1/grabeni
```

## Release

```bash
$ vim ./version.go
$ make cross
$ ghr -u yuuki1 -p 2 $VERSION snapshot/
```

## Roadmap

- `attach`, `detach`, `grab`: Display prompt
- `attach`, `detach`, `grab`: Show ENI information before execution
- `attach`, `detach`, `grab`: Enrich log
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

[y_uuki](https://github.com/yuuki1)
