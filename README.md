# EBdeployer [![GoDoc](https://godoc.org/github.com/baltimore-sun-data/ebdeployer?status.svg)](https://godoc.org/github.com/baltimore-sun-data/ebdeployer) [![Go Report Card](https://goreportcard.com/badge/github.com/baltimore-sun-data/ebdeployer)](https://goreportcard.com/report/github.com/baltimore-sun-data/ebdeployer)

EBdeployer is a wrapper for deploying a Docker multi-image build to AWS Elastic Beanstalk.

Elastic Beanstalk is designed to make it simple to spin up new servers, but unfortunately, the awsebcli tool does not currently automated the steps involved in creating a Docker multi-image build. This tool fixes that problem.

## Installation

First install [Go](http://golang.org).

If you just want to install the binary to your current directory and don't care about the source code, run

```bash
GOBIN="$(pwd)" GOPATH="$(mktemp -d)" go get github.com/baltimore-sun-data/ebdeployer
```
