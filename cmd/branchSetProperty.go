/*
Copyright © 2021 Damien Coraboeuf <damien.coraboeuf@nemerosa.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// branchSetPropertyCmd represents the branchSetProperty command
var branchSetPropertyCmd = &cobra.Command{
	Use:   "set-property",
	Short: "Sets a property on a branch",
	Long: `Sets a property on a branch.
	
This can be used for setting generic properties, using their full qualified class name and some JSON content:
	
	yontrack branch set-property --project PROJECT --branch BRANCH generic --property "net.nemerosa.ontrack.extension.git.property.GitBranchConfigurationPropertyType" --value '{branch:"main"}'

Some specific commands are also available for the most common property types. The example below does exactly the
same update than the one just above:

	yontrack branch set-property --project PROJECT --branch BRANCH git --git-branch main
`,
	// Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	branchCmd.AddCommand(branchSetPropertyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// branchSetPropertyCmd.PersistentFlags().String("foo", "", "A help for foo")
	branchSetPropertyCmd.PersistentFlags().StringP("project", "p", "", "Name of the project")
	branchSetPropertyCmd.PersistentFlags().StringP("branch", "b", "", "Name of the branch")

	branchSetPropertyCmd.MarkPersistentFlagRequired("project")
	branchSetPropertyCmd.MarkPersistentFlagRequired("branch")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// branchSetPropertyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
