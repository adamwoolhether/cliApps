//go:build !inmemory && !containers

package cmd

import (
	"github.com/spf13/viper"

	"github.com/adamwoolhether/cliApps/distributing/pomo/pomodoro"
	"github.com/adamwoolhether/cliApps/distributing/pomo/pomodoro/repository"
)

func getRepo() (pomodoro.Repository, error) {
	return repository.NewSQLite3Repo(viper.GetString("db"))
}
