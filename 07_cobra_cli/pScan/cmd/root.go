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
// cobra init --pkg-name GO_MODULE
// cobra add hosts
// cobra add list -p hostsCmd
// cobra add completion
//
// PSCAN_HOSTS_FILE=newFile.hosts

package cmd

import (
	"fmt"
	"os"
	"strings"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "pScan",
	Version: "0.1",
	Short:   "Fast TCP port scanner",
	Long: `pScan - short for Port Scanner - executes TCP port scan
on a list of hosts.

pScan allows you to add, list, and delete hosts from the list.

pScan executes a port scan on specified TCP ports. You can customize the
target ports using a command line flag.`,
	// Run: func(cmd *cobra.Command, args []string) { },
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
	
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pScan.yaml)")
	rootCmd.PersistentFlags().StringP("hosts-file", "f", "pScan.hosts", "pScan hosts file")
	
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("PSCAN")
	err := viper.BindPFlag("hosts-file", rootCmd.PersistentFlags().Lookup("hosts-file"))
	if err != nil {
		return
	}
	
	versionTemplate := `{{printf "%s: %s - version %s\n" .Name .Short .Version}}`
	rootCmd.SetVersionTemplate(versionTemplate)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		
		viper.AddConfigPath(home)
		viper.SetConfigName(".pScan")
	}
	
	viper.AutomaticEnv()
	
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
