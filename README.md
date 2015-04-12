grabeni
=======

Grab ENI tool from an other EC2 instance.

# Usage

```
grabeni status <eniId>
grabeni grab <eniId> [-d <deviceIndex>] [-i <instanceId>] [--dryrun]
grabeni attach [<eniId>] [-d <deviceIndex>] [-i <instanceId>] [--dryrun]
grabeni detach [<eniId>] [-i <instanceId>] [--dryrun]
```

# Install

To install, use `go get`:

```bash
$ go get -d github.com/y_uuki/grabeni
```

# Contribution

1. Fork ([https://github.com/y_uuki/grabeni/fork](https://github.com/y_uuki/grabeni/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create new Pull Request
