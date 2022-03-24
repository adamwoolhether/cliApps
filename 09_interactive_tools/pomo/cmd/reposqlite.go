//go:build !inmemory

package cmd

import (
	"github.com/spf13/viper"
	
	"github.com/adamwoolhether/cliApps/09_interactive_tools/pomo/pomodoro"
	"github.com/adamwoolhether/cliApps/09_interactive_tools/pomo/pomodoro/repository"
)

func getRepo() (pomodoro.Repository, error) {
	return repository.NewSQLite3Repo(viper.GetString("db"))
}
