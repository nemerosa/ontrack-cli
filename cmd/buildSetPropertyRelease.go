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
	"yontrack/client"
	"yontrack/config"
	"yontrack/utils"

	"github.com/spf13/cobra"
)

// buildSetPropertyReleaseCmd represents the buildSetPropertyRelease command
var buildSetPropertyReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Sets the release property on a build",
	Long: `Sets the release property on a build.
	
For example:

yontrack build set-property --project PROJECT --branch BRANCH --build BUILD release RC-1
`,
	Args: cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		project, branch, build, err := utils.GetProjectBranchBuildFlags(cmd, false, true)
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
}
