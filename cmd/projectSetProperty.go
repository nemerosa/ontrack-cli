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

// projectSetPropertyCmd represents the projectSetProperty command
var projectSetPropertyCmd = &cobra.Command{
	Use:   "set-property",
	Short: "Sets a property on a project",
	Long: `Sets a property on a project.
	
This can be used for setting generic properties, using their full qualified class name and some JSON content:
	
	yontrack project set-property --project PROJECT generic --property "net.nemerosa.ontrack.extension.github.property.GitHubProjectConfigurationPropertyType" --value '{"configuration":"GitHub","repository":"nemerosa/ontrack"}'

Some specific commands are also available for the most common property types. The example below does exactly the
same update than the one just above:

	yontrack project set-property --project PROJECT github --configuration GitHub --repository nemerosa/ontrack
`,
	// Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	projectCmd.AddCommand(projectSetPropertyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	projectSetPropertyCmd.PersistentFlags().StringP("project", "p", "", "Name of the project")
	projectSetPropertyCmd.MarkPersistentFlagRequired("project")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// projectSetPropertyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
