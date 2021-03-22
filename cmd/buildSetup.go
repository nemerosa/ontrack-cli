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

	client "ontrack-cli/client"
	config "ontrack-cli/config"
)

// buildSetupCmd represents the buildSetup command
var buildSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Creates a build if it does not exist",
	Long: `Creates a build if it does not exist.

For example, the following command will create the build:

    ontrack-cli build setup --project my-project --branch release/1.0 --build 1

and the same command run a second time won't do anything.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return buildSetup(cmd)
	},
}

func buildSetup(cmd *cobra.Command) error {
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

	description, err := cmd.Flags().GetString("description")
	if err != nil {
		return err
	}

	release, err := cmd.Flags().GetString("release")
	if err != nil {
		return err
	}
	releaseProperty := release != ""

	// Run info
	runInfo, err := GetRunInfo(cmd)
	if err != nil {
		return err
	}

	config, err := config.GetSelectedConfiguration()
	if err != nil {
		return err
	}

	// Creates or get the build
	var data struct {
		CreateBuildOrGet struct {
			Errors []struct {
				Message string
			}
		}
	}
	if err := client.GraphQLCall(config, `
		mutation BuildSetup(
			$project: String!,
			$branch: String!, 
			$build: String!, 
			$description: String, 
			$runInfo: RunInfoInput,
			$releaseProperty: Boolean!,
			$release: String!
		) {
			createBuildOrGet(input: {
				projectName: $project, 
				branchName: $branch, 
				name: $build, 
				description: $description, 
				runInfo: $runInfo
			}) {
				errors {
				  message
				}
			}
			setBuildReleaseProperty(input: {
				project: $project,
				branch: $branch,
				build: $build,
				release: $release
			}) @include(if: $releaseProperty) {
				errors {
					message
				}
			}
		}
	`, map[string]interface{}{
		"project":         project,
		"branch":          branch,
		"build":           build,
		"description":     description,
		"runInfo":         runInfo,
		"releaseProperty": releaseProperty,
		"release":         release,
	}, &data); err != nil {
		return err
	}

	// Checks errors for the build
	if err := client.CheckDataErrors(data.CreateBuildOrGet.Errors); err != nil {
		return err
	}

	// OK
	return nil
}

func init() {
	buildCmd.AddCommand(buildSetupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildSetupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	buildSetupCmd.Flags().StringP("project", "p", "", "Project name (required)")
	buildSetupCmd.Flags().StringP("branch", "b", "", "Branch name (required)")
	buildSetupCmd.Flags().StringP("build", "n", "", "Build name (required)")
	buildSetupCmd.Flags().StringP("description", "d", "", "Build description")

	// Run info parameters
	InitRunInfoCommandFlags(buildSetupCmd)

	// Release property
	buildSetupCmd.Flags().StringP("release", "r", "", "Build release property")

	buildSetupCmd.MarkFlagRequired("project")
	buildSetupCmd.MarkFlagRequired("branch")
	buildSetupCmd.MarkFlagRequired("build")
}
