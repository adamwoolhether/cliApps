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
	"text/tabwriter"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var viewCmd = &cobra.Command{
	Use:          "view <id>",
	Short:        "View details about a single item",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiRoot := viper.GetString("api-root")
		
		return viewAction(os.Stdout, apiRoot, args[0])
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}

func viewAction(out io.Writer, apiRoot, arg string) error {
	id, err := strconv.Atoi(arg)
	if err != nil {
		return fmt.Errorf("%w: item id must be a number", ErrNotNumber)
	}
	
	i, err := getOne(apiRoot, id)
	if err != nil {
		return err
	}
	
	return printOne(out, i)
}

func printOne(out io.Writer, i item) error {
	w := tabwriter.NewWriter(out, 14, 2, 0, ' ', 0)
	fmt.Fprintf(w, "Task:\t%s\n", i.Task)
	fmt.Fprintf(w, "Created at:\t%s\n", i.CreatedAt.Format(timeFormat))
	if i.Done {
		fmt.Fprintf(w, "Completed:\t%s\n", "Yes")
		fmt.Fprintf(w, "Completed At:\t%s\n", i.CompletedAt.Format(timeFormat))
	}
	
	fmt.Fprintf(w, "Completed:\t%s\n", "No")
	
	return w.Flush()
}
