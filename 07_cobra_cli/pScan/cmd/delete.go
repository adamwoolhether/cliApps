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
	
	"github.com/adamwoolhether/cliApps/07_cobra_cli/pScan/scan"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Aliases:      []string{"d"},
	Use:          "delete <host1>...<hostn>",
	Short:        "Delete hosts(s) from list",
	SilenceUsage: true,
	Args:         cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile, err := cmd.Flags().GetString("hosts-file")
		if err != nil {
			return err
		}
		
		return deleteAction(os.Stdout, hostsFile, args)
	},
}

func init() {
	hostsCmd.AddCommand(deleteCmd)
}

// deleteAction deletes the given hosts from the hosts file.
func deleteAction(out io.Writer, hostsFile string, args []string) error {
	hl := &scan.HostsList{}
	
	if err := hl.Load(hostsFile); err != nil {
		return err
	}
	
	for _, h := range args {
		if err := hl.Remove(h); err != nil {
			return err
		}
		fmt.Fprintln(out, "Deleted host:", h)
	}
	
	return hl.Save(hostsFile)
}
