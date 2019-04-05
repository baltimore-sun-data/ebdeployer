package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
)

func createEBEnv(repo, cfg string, date time.Time) (env string, err error) {
	env = repo + "-" + date.Format("Jan2-1504")
	env = strings.ToLower(env)

	if err = subprocess("", "eb", "create", env, "--cfg", cfg); err != nil {
		return "", err
	}
	if err = subprocess("", "eb", "use", env); err != nil {
		return "", err
	}
	return env, nil
}

func showEBEnvInfo(env string) error {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return err
	}

	svc := elasticbeanstalk.New(cfg)
	req := svc.DescribeEnvironmentsRequest(&elasticbeanstalk.DescribeEnvironmentsInput{
		EnvironmentNames: []string{env},
	})
	result, err := req.Send()
	if err != nil {
		return err
	}

	if len(result.Environments) != 1 || result.Environments[0].CNAME == nil {
		return fmt.Errorf("missing data in AWS EB result: %v", result)
	}
	cname := *result.Environments[0].CNAME
	log.Printf("Environment CNAME: %s\n", cname)
	ips, err := net.LookupHost(cname)
	if err != nil {
		return err
	}
	for _, ip := range ips {
		log.Printf("Environment IP: %s\n", ip)
	}
	return nil
}
