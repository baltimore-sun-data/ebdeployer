package main

import (
	"fmt"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
)

func repoHash() (string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", err
	}
	head, err := repo.Head()
	if err != nil {
		return "", err
	}

	return head.Hash().String(), nil
}

func repoPushTag(name string) error {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return err
	}
	head, err := repo.Head()
	if err != nil {
		return err
	}

	// Ignore errors. Race-y, but can't be helped.
	_ = repo.DeleteTag(name)

	newRef, err := repo.CreateTag(name, head.Hash(), nil)
	if err != nil {
		return err
	}
	rs := config.RefSpec(fmt.Sprintf("%s:%s",
		newRef.Name(), newRef.Name(),
	))
	err = repo.Push(&git.PushOptions{
		RefSpecs: []config.RefSpec{rs},
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	return nil
}
