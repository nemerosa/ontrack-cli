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

// branchSetPropertyGenericCmd represents the branchSetPropertyGeneric command
var branchSetPropertyGenericCmd = &cobra.Command{
	Use:   "generic",
	Short: "Sets a branch property using its type and its value as JSON",
	Long: `Sets a branch property using its type and its value as JSON.

Example:

ontrack-cli branch set-property --project PROJECT --branch BRANCH generic --property "net.nemerosa.ontrack.extension.git.property.GitBranchConfigurationPropertyType" --value '{branch:"main"}'

	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, err := cmd.Flags().GetString("project")
		if err != nil {
			return err
		}

		branch, err := cmd.Flags().GetString("branch")
		if err != nil {
			return err
		}

		property, err := cmd.Flags().GetString("property")
		if err != nil {
			return err
		}

		value, err := cmd.Flags().GetString("value")
		if err != nil {
			return err
		}

		return SetProperty("branch", map[string]string{
			"project": project,
			"branch":  branch,
		}, property, value)
	},
}

func init() {
	branchSetPropertyCmd.AddCommand(branchSetPropertyGenericCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// branchSetPropertyGenericCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// branchSetPropertyGenericCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	branchSetPropertyGenericCmd.Flags().StringP("property", "t", "", "FQCN of the property")
	branchSetPropertyGenericCmd.Flags().StringP("value", "v", "", "Value of the property as a JSON string")

	branchSetPropertyGenericCmd.MarkFlagRequired("property")
	branchSetPropertyGenericCmd.MarkFlagRequired("value")
}
