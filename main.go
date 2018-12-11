package main // import "github.com/baltimore-sun-data/ebdeployer"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/carlmjohnson/flagext"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Fatal error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	now := time.Now()
	repo := flag.String("repo", "", "name of the repo")
	cfg := flag.String("cfg", "", "name of the eb config")
	secretsFile := flagext.File("secrets.json")
	flag.Var(secretsFile, "secrets", "json `file` to read secrets out of")
	baseFile := flagext.File("Dockerrun.base.json")
	flag.Var(baseFile, "dockerrun", "json `file` to use as base template for Dockerrun.aws.json")
	flag.Parse()

	if *repo == "" {
		return fmt.Errorf("must set repo name with -repo")
	}
	if *cfg == "" {
		return fmt.Errorf("must set cfg name with -cfg")
	}

	log.Println("Get docker login from AWS")
	ecr, user, password, err := getDockerLogin()
	if err != nil {
		return err
	}

	log.Println("Docker login")
	if err = dockerLogin(ecr, user, password); err != nil {
		return err
	}

	log.Printf("Read %q", secretsFile)
	var secrets map[string]map[string]string
	if err = readJSON(secretsFile, &secrets); err != nil {
		return err
	}

	log.Printf("Read %q", baseFile)
	var base Dockerrun
	if err = readJSON(baseFile, &base); err != nil {
		return err
	}

	for _, cd := range base.ContainerDefinitions {
		// Fix repo/tag
		cd.Image = makeDockerTag(ecr, *repo, cd.Image, now)
		log.Printf("Docker build %s from %q", cd.Image, cd.Dockerfile)
		if err = dockerBuild(cd.Image, cd.Dockerfile); err != nil {
			return err
		}
		// AWS doesn't expect a Dockerfile field, so drop it
		cd.Dockerfile = ""
		log.Printf("Docker push %s", cd.Image)
		if err = dockerPush(cd.Image); err != nil {
			return err
		}
		// Add in secrets
		for name, val := range secrets[cd.Name] {
			cd.Environment = append(cd.Environment, EnvPair{name, val})
		}
	}

	log.Println("Write Dockerrun.aws.json")
	if err = writeJSON("Dockerrun.aws.json", &base); err != nil {
		return err
	}

	log.Println("Create EB environment")
	return createEBEnv(*repo, *cfg, now)
}

func subprocess(stdin string, name string, args ...string) error {
	log.Printf("Running %q", strings.Join(append([]string{name}, args...), " "))
	cmd := exec.Command(name, args...)
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
