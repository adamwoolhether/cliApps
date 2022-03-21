package cmd

import (
	"github.com/adamwoolhether/cliApps/09_interactive_tools/pomo/pomodoro"
	"github.com/adamwoolhether/cliApps/09_interactive_tools/pomo/pomodoro/repository"
)

// getRepo returns an instance of pomodoro.Repository.
func getRepo() (pomodoro.Repository, error) {
	return repository.NewInMemoryRepo(), nil
}
