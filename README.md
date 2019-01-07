# EBdeployer [![GoDoc](https://godoc.org/github.com/baltimore-sun-data/ebdeployer?status.svg)](https://godoc.org/github.com/baltimore-sun-data/ebdeployer) [![Go Report Card](https://goreportcard.com/badge/github.com/baltimore-sun-data/ebdeployer)](https://goreportcard.com/report/github.com/baltimore-sun-data/ebdeployer)

EBdeployer is a wrapper for deploying a Docker multi-image build to AWS Elastic Beanstalk.

Elastic Beanstalk is designed to make it simple to spin up new servers, but unfortunately, the awsebcli tool does not currently automate the steps involved in creating a Docker multi-image build. This tool fixes that problem.

## Installation

First install [Go](http://golang.org).

If you just want to install the binary to your current directory and don't care about the source code, run

```bash
GOBIN="$(pwd)" GOPATH="$(mktemp -d)" go get github.com/baltimore-sun-data/ebdeployer
```

Other requirements:

- A valid AWS key (run `aws configure` to set this with the AWS CLI)
- awsebcli should be installed as `eb` (on Mac `brew install awsebcli`)
- Docker/docker-compose

## CLI options
```
$ ebdeployer -h
Usage of ebdeployer:
  -cfg string
        name of the eb config
  -compose path
        path to the docker-compose file (default "docker-compose.yaml")
  -deploy
        run eb create env after setting up Dockerrun file (default true)
  -dir path
        path of Docker project directory (default ".")
  -dockerrun file
        json file to use as base template for Dockerrun.aws.json (default Dockerrun.base.json)
  -repo string
        name of the repo
  -secrets file
        json file to read secrets out of (default secrets.json)
```
