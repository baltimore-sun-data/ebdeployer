package main

import (
	"time"
)

func createEBEnv(repo, cfg string, date time.Time) error {
	env := repo + "-" + date.Format("jan2-1504")
	err := subprocess("", "eb", "create", env, "--cfg", cfg)
	if err != nil {
		return err
	}
	return subprocess("", "eb", "use", env)
}
