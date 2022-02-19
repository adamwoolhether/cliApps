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
	"io"
	"os"
	
	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate bash completion for your command",
	Long: `To load your completions run
source <(pScan completion)

To load completions automatically on login, add this line to your .bashrc file:
$ ~/.bashrc
source <(pScan completion)
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return completionAction(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

func completionAction(out io.Writer) error {
	// return rootCmd.GenZshCompletion(out)
	return rootCmd.GenBashCompletion(out)
	// return rootCmd.GenBashCompletionV2(out, true)
}
