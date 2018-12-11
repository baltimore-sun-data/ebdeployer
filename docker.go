package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

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

	return *data.ProxyEndpoint, string(userpass[0]), string(userpass[1]), nil
}

func makeDockerTag(ecr, repo, image string, t time.Time) string {
	return fmt.Sprintf("%s/%s:%s-%s", ecr, repo, image, t.Format("2006-01-02-1504"))
}

func dockerLogin(endpoint, user, password string) error {
	cmd := exec.Command("docker", "login", "-u", user, "--password-stdin", endpoint)
	cmd.Stdin = strings.NewReader(password)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func dockerBuild(tag, file string) error {
	cmd := exec.Command("docker", "build", "-t", tag, file)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func dockerPush(tag string) error {
	cmd := exec.Command("docker", "push", tag)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
