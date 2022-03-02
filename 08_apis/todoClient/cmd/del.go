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
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// delCmd represents the del command
var delCmd = &cobra.Command{
	Use:          "del <id>",
	Short:        "Deletes and item from the list",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiRoot := viper.GetString("api-root")
		
		return delAction(os.Stdout, apiRoot, args[0])
	},
}

func init() {
	rootCmd.AddCommand(delCmd)
}

func delAction(out io.Writer, apiRoot, arg string) error {
	id, err := strconv.Atoi(arg)
	if err != nil {
		return fmt.Errorf("%w: item id must be a number", ErrNotNumber)
	}
	
	if err := deleteItem(apiRoot, id); err != nil {
		return err
	}
	
	return printDelete(out, id)
}

func printDelete(out io.Writer, id int) error {
	_, err := fmt.Fprintf(out, "Item number %d deleted.\n", id)
	
	return err
}
