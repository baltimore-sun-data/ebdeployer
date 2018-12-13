package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

func getDockerLogin() (endpoint, user, password string, err error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return "", "", "", err
	}

	svc := ecr.New(cfg)
	input := &ecr.GetAuthorizationTokenInput{}

	req := svc.GetAuthorizationTokenRequest(input)
	result, err := req.Send()
	if err != nil {
		return "", "", "", err
	}

	data := result.AuthorizationData[0]
	b, err := base64.StdEncoding.DecodeString(*data.AuthorizationToken)
	if err != nil {
		return "", "", "", err
	}
	userpass := bytes.SplitN(b, []byte(":"), 2)
	if len(userpass) != 2 {
		return "", "", "", fmt.Errorf("bad auth token %q", b)
	}

	endpoint = strings.TrimPrefix(*data.ProxyEndpoint, "https://")
	user, password = string(userpass[0]), string(userpass[1])
	return
}

func makeDockerTag(ecr, repo, image, dateTag string) string {
	return fmt.Sprintf("%s/%s:%s-%s", ecr, repo, image, dateTag)
}

func dockerLogin(endpoint, user, password string) error {
	return subprocess(password, "docker", "login", "-u", user, "--password-stdin", "https://"+endpoint)
}

func dockerBuild(dir, file string) error {
	return subprocess("", "docker-compose", "-f", file, "--project-directory", dir, "build")
}

func dockerPush(dir, file string) error {
	return subprocess("", "docker-compose", "-f", file, "--project-directory", dir, "push")
}
