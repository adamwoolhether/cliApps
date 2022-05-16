//go:build inmemory || containers

package cmd

import (
	"github.com/adamwoolhether/cliApps/distributing/pomo/pomodoro"
	"github.com/adamwoolhether/cliApps/distributing/pomo/pomodoro/repository"
)

// getRepo returns an instance of pomodoro.Repository.
func getRepo() (pomodoro.Repository, error) {
	return repository.NewInMemoryRepo(), nil
}
