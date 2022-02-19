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
	
	"github.com/adamwoolhether/cliApps/07_cobra_cli/pScan/scan"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run a port scan on the hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile := viper.GetString("hosts-file")
		
		ports, err := cmd.Flags().GetIntSlice("ports")
		if err != nil {
			return err
		}
		return scanAction(os.Stdout, hostsFile, ports)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().IntSliceP("ports", "p", []int{22, 80, 443}, "ports to scan")
}

// scanAction
func scanAction(out io.Writer, hostsFile string, ports []int) error {
	hl := &scan.HostsList{}
	
	if err := hl.Load(hostsFile); err != nil {
		return err
	}
	
	results := scan.Run(hl, ports)
	
	return printResults(out, results)
}

func printResults(out io.Writer, results []scan.Result) error {
	message := ""
	
	for _, r := range results {
		message += fmt.Sprintf("%s:", r.Host)
		
		if r.NotFound {
			message += fmt.Sprintf(" Host not found\n\n")
			continue
		}
		
		message += fmt.Sprintln()
		
		for _, p := range r.PortStates {
			message += fmt.Sprintf("\t%d: %s\n", p.Port, p.Open)
		}
		message += fmt.Sprintln()
	}
	
	_, err := fmt.Fprint(out, message)
	return err
}
