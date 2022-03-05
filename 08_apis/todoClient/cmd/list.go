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
	"text/tabwriter"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:          "list",
	Short:        "A brief description of your command",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiRoot := viper.GetString("api-root")
		active, err := cmd.Flags().GetBool("active")
		if err != nil {
			return err
		}
		
		return listAction(os.Stdout, apiRoot, active)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("active", "a", false, "Display only incomplete tasks")
}

func listAction(out io.Writer, apiRoot string, active bool) error {
	items, err := getAll(apiRoot)
	if err != nil {
		return err
	}
	
	return printAll(out, items, active)
}

func printAll(out io.Writer, items []item, active bool) error {
	w := tabwriter.NewWriter(out, 3, 2, 0, ' ', 0)
	
	for k, v := range items {
		if v.Done && active {
			continue
		}
		
		done := "-"
		if v.Done {
			done = "X"
		}
		fmt.Fprintf(w, "%s\t%d\t%s\t\n", done, k+1, v.Task)
	}
	
	return w.Flush()
}
