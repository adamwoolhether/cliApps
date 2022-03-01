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

/*
curl -L -X POST -d '{"task": "Task 1"}' -H 'Content-Type: application/json' http://localhost:8080/todo
curl -L -X POST -d '{"task": "Task 2"}' -H 'Content-Type: application/json' http://localhost:8080/todo

./todoServer -f /tmp/testtodoclient01.json
*/
package cmd

import (
	"os"
	"strings"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "todoClient",
	Short: "A brief description of your application",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// cobra.OnInitialize(initConfig)
	
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.todoClient.yaml")
	rootCmd.PersistentFlags().String("api-root", "http://localhost:8080", "Todo APR URL")
	
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("TODO")
	
	viper.BindPFlag("api-root", rootCmd.PersistentFlags().Lookup("api-root"))
}

// func initConfig() {
//
// }
