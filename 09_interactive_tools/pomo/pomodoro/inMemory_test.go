//go:build inmemory

package pomodoro_test

import (
	"testing"
	
	"github.com/adamwoolhether/cliApps/09_interactive_tools/pomo/pomodoro"
	"github.com/adamwoolhether/cliApps/09_interactive_tools/pomo/pomodoro/repository"
)

// getRepo returns the repo instance and a cleanup func.
func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()
	
	return repository.NewInMemoryRepo(), func() {}
}
