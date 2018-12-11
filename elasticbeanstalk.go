package main

import (
	"os"
	"os/exec"
	"time"
)

func createEBEnv(repo, cfg string, date time.Time) error {
	env := repo + "-" + date.Format("jan2-1504")
	cmd := exec.Command("eb", "create", env, "--cfg", cfg)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("eb", "use", env)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
