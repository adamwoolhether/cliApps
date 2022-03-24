/*
Copyright © 2022 Adam Woolhether

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"time"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	
	"github.com/adamwoolhether/cliApps/09_interactive_tools/pomo/app"
	"github.com/adamwoolhether/cliApps/09_interactive_tools/pomo/pomodoro"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pomo",
	Short: "Interactive Pomodoro Timer",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := getRepo()
		if err != nil {
			return err
		}
		
		config := pomodoro.NewConfig(
			repo,
			viper.GetDuration("pomo"),
			viper.GetDuration("short"),
			viper.GetDuration("long"),
		)
		
		return rootAction(os.Stdout, config)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pomo.yaml)")
	rootCmd.Flags().StringP("db", "d", "pomo.db", "Database file")
	rootCmd.Flags().DurationP("pomo", "p", 25*time.Minute, "Pomodoro duration")
	rootCmd.Flags().DurationP("short", "s", 5*time.Minute, "Short break duration")
	rootCmd.Flags().DurationP("long", "l", 15*time.Minute, "Long break duration")
	
	viper.BindPFlag("db", rootCmd.Flags().Lookup("db"))
	viper.BindPFlag("pomo", rootCmd.Flags().Lookup("pomo"))
	viper.BindPFlag("short", rootCmd.Flags().Lookup("short"))
	viper.BindPFlag("short", rootCmd.Flags().Lookup("long"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		
		viper.AddConfigPath(home)
		viper.SetConfigName(".pomo")
	}
	
	viper.AutomaticEnv()
	
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("can't use config file:", viper.ConfigFileUsed())
		
	}
}

func rootAction(out io.Writer, config *pomodoro.IntervalConfig) error {
	a, err := app.New(config)
	if err != nil {
		return err
	}
	
	return a.Run()
}
