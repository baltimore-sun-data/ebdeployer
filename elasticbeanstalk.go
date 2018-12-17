package main

import (
	"strings"
	"time"
)

func createEBEnv(repo, cfg string, date time.Time) error {
	env := repo + "-" + date.Format("Jan2-1504")
	env = strings.ToLower(env)
	err := subprocess("", "eb", "create", env, "--cfg", cfg)
	if err != nil {
		return err
	}
	return subprocess("", "eb", "use", env)
}
