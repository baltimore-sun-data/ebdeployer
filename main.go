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
	compose := flag.String("compose", "docker-compose.yaml", "`path` to the docker-compose file")
	dir := flag.String("dir", ".", "`path` of Docker project directory")
	repo := flag.String("repo", "", "name of the repo")
	cfg := flag.String("cfg", "", "name of the eb config")
	secretsFile := flagext.File("secrets.json")
	flag.Var(secretsFile, "secrets", "json `file` to read secrets out of")
	baseFile := flagext.File("Dockerrun.base.json")
	flag.Var(baseFile, "dockerrun", "json `file` to use as base template for Dockerrun.aws.json")
	deploy := flag.Bool("deploy", true, "run eb create env after setting up Dockerrun file")
	flag.Parse()

	if *repo == "" {
		return fmt.Errorf("must set repo name with -repo")
	}
	if *cfg == "" {
		return fmt.Errorf("must set cfg name with -cfg")
	}

	now := time.Now()
	dateTag := os.Getenv("DATE_TAG")
	if dateTag == "" {
		dateTag = now.Format("2006-01-02-1504")
		os.Setenv("DATE_TAG", dateTag)
	}

	var err error
	appVersion := os.Getenv("APP_VERSION")
	if appVersion == "" {
		appVersion, err = repoHash()
		if err != nil {
			return err
		}
		os.Setenv("APP_VERSION", appVersion)
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

	os.Setenv("ECR_REPO", ecr)

	log.Println("Docker build")
	if err = dockerBuild(*dir, *compose); err != nil {
		return err
	}
	log.Println("Docker push")
	if err = dockerPush(*dir, *compose); err != nil {
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
		cd.Image = makeDockerTag(ecr, *repo, cd.Image, dateTag)

		// Add in secrets
		for name, val := range secrets[cd.Name] {
			cd.Environment = append(cd.Environment, EnvPair{name, val})
		}
	}

	log.Println("Write Dockerrun.aws.json")
	if err = writeJSON("Dockerrun.aws.json", &base); err != nil {
		return err
	}

	if !*deploy {
		return nil
	}

	log.Println("Create EB environment")
	envname, err := createEBEnv(*repo, *cfg, now)
	if err != nil {
		return err
	}

	log.Println("Tagging release")
	if err = repoPushTag("release/" + envname); err != nil {
		return err
	}

	log.Println("Getting EB environment description")
	return showEBEnvInfo(envname)
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
