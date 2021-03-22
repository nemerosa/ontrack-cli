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
	"ontrack-cli/client"
	"ontrack-cli/config"

	"github.com/spf13/cobra"
)

// buildSetPropertyReleaseCmd represents the buildSetPropertyRelease command
var buildSetPropertyReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Sets the release property on a build",
	Long: `Sets the release property on a build.
	
For example:

ontrack-cli build set-property --project PROJECT --branch BRANCH --build BUILD release RC-1
`,
	Args: cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		project, err := cmd.Flags().GetString("project")
		if err != nil {
			return err
		}

		branch, err := cmd.Flags().GetString("branch")
		if err != nil {
			return err
		}
		branch = NormalizeBranchName(branch)

		build, err := cmd.Flags().GetString("build")
		if err != nil {
			return err
		}

		// Property value
		value := args[0]

		// Configuration
		cfg, err := config.GetSelectedConfiguration()
		if err != nil {
			return err
		}

		// Data
		var data struct {
			SetBuildReleaseProperty struct {
				Errors []struct {
					Message string
				}
			}
		}

		// Call
		if err := client.GraphQLCall(cfg, `
			mutation SetBuildReleaseProperty(
				$project: String!,
				$branch: String!,
				$build: String!,
				$release: String!
			) {
				setBuildReleaseProperty(input: {
					project: $project,
					branch: $branch,
					build: $build,
					release: $release
				}) {
					errors {
						message
					}
				}
			}
		`, map[string]interface{}{
			"project": project,
			"branch":  branch,
			"build":   build,
			"release": value,
		}, &data); err != nil {
			return err
		}

		if err := client.CheckDataErrors(data.SetBuildReleaseProperty.Errors); err != nil {
			return err
		}

		// OK
		return nil
	},
}

func init() {
	buildSetPropertyCmd.AddCommand(buildSetPropertyReleaseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildSetPropertyReleaseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildSetPropertyReleaseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
