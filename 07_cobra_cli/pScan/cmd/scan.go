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
	"strconv"
	"strings"
	
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
		
		portRange, err := cmd.Flags().GetString("range")
		if err != nil {
			return err
		}
		
		var filter string
		open, err := cmd.Flags().GetBool("open")
		if err != nil {
			return err
		}
		
		closed, err := cmd.Flags().GetBool("closed")
		if err != nil {
			return err
		}
		
		if open && closed {
			fmt.Errorf("cannot filter both open and closed ports")
		}
		if open {
			filter = "open"
		}
		if closed {
			filter = "closed"
		}
		
		return scanAction(os.Stdout, filter, hostsFile, portRange, ports)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().IntSliceP("ports", "p", []int{22, 80, 443}, "ports to scan")
	scanCmd.Flags().StringP("range", "r", "", "range of ports to scan")
	scanCmd.Flags().BoolP("open", "o", false, "show only open ports")
	scanCmd.Flags().BoolP("closed", "x", false, "show only closed ports")
}

// scanAction loads the hosts file, runs the port scan, and returns printed results.
func scanAction(out io.Writer, filter, hostsFile, portRange string, ports []int) error {
	if portRange != "" {
		portRangeMembers, err := rangePorts(portRange)
		if err != nil {
			return err
		}
		
		ports = append(ports, portRangeMembers...)
	}
	
	hl := &scan.HostsList{}
	
	if err := hl.Load(hostsFile); err != nil {
		return err
	}
	
	if !validatePort(ports) {
		return scan.ErrInvalidPort
	}
	
	results := scan.Run(hl, ports)
	
	return printResults(out, filter, results)
}

func validatePort(ports []int) bool {
	for _, port := range ports {
		if port < 1 || port > 65535 {
			return false
		}
	}
	
	return true
}

func rangePorts(portRange string) ([]int, error) {
	span := strings.Split(portRange, "-")
	if len(span) < 2 {
		return nil, scan.ErrInvalidRange
	}
	
	low, err := strconv.Atoi(span[0])
	if err != nil {
		return nil, err
	}
	high, err := strconv.Atoi(span[1])
	if err != nil {
		return nil, err
	}
	
	if low < 1 || high > 65535 || low > high {
		return nil, scan.ErrInvalidRange
	}
	
	ports := make([]int, 0, high-low+1)
	
	for i := low; i < high+1; i++ {
		ports = append(ports, i)
	}
	
	return ports, nil
}

// printResults ranges over results, formats them, and writes them to stdout.
func printResults(out io.Writer, filter string, results []scan.Result) error {
	message := ""
	
	for _, r := range results {
		message += fmt.Sprintf("%s:", r.Host)
		
		if r.NotFound {
			message += fmt.Sprintf(" Host not found\n\n")
			continue
		}
		
		message += fmt.Sprintln()
		
		for _, p := range r.PortStates {
			if filter != "" && filter != p.Open.String() {
				break
			}
			message += fmt.Sprintf("\t%d: %s\n", p.Port, p.Open)
		}
		message += fmt.Sprintln()
	}
	
	_, err := fmt.Fprint(out, message)
	return err
}
