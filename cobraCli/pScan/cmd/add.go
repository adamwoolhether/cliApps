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
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	
	"github.com/adamwoolhether/cliApps/cobraCli/pScan/scan"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Aliases:      []string{"a"},
	Args:         cobra.MinimumNArgs(1),
	Use:          "add <host1>...<hostn>",
	Short:        "Add new host(s) to list",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile := viper.GetString("hosts-file")
		
		return addAction(os.Stdout, hostsFile, args)
	},
}

func init() {
	hostsCmd.AddCommand(addCmd)
}

// addAction adds given hosts to the hosts file.
func addAction(out io.Writer, hostsFile string, args []string) error {
	hl := &scan.HostsList{}
	
	if err := hl.Load(hostsFile); err != nil {
		return err
	}
	for _, h := range args {
		if err := hl.Add(h); err != nil {
			return err
		}
		fmt.Fprintln(out, "Added host:", h)
	}
	
	return hl.Save(hostsFile)
}
