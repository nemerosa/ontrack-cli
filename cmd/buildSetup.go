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
	"regexp"

	"github.com/spf13/cobra"

	client "ontrack-cli/client"
	config "ontrack-cli/config"
)

var buildSetupProject string
var buildSetupBranch string
var buildSetupBuild string
var buildSetupDescription string

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
		return buildSetup()
	},
}

func buildSetup() error {
	config, err := config.GetSelectedConfiguration()
	if err != nil {
		return err
	}
	// Normalizing the name of the branch
	re := regexp.MustCompile("[^A-Za-z0-9\\._-]")
	normalizedBranchName := re.ReplaceAllString(buildSetupBranch, "-")
	// Creates or get the build
	var data struct {
		CreateBuildOrGet struct {
			Build struct {
				ID int
			}
			Errors []struct {
				Message string
			}
		}
	}
	if err := client.GraphQLCall(config, `
		mutation BuildSetup($project: String!, $branch: String!, $build: String!, $description: String) {
			createBuildOrGet(input: {projectName: $project, branchName: $branch, name: $build, description: $description}) {
				build {
				  id
				}
				errors {
				  message
				}
			}
		}
	`, map[string]interface{}{
		"project":     buildSetupProject,
		"branch":      normalizedBranchName,
		"build":       buildSetupBuild,
		"description": buildSetupDescription,
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
	buildSetupCmd.Flags().StringVarP(&buildSetupProject, "project", "p", "", "Project name (required)")
	buildSetupCmd.Flags().StringVarP(&buildSetupBranch, "branch", "b", "", "Branch name (required)")
	buildSetupCmd.Flags().StringVarP(&buildSetupBuild, "build", "n", "", "Build name (required)")
	buildSetupCmd.Flags().StringVarP(&buildSetupDescription, "description", "d", "", "Build description")

	buildSetupCmd.MarkFlagRequired("project")
	buildSetupCmd.MarkFlagRequired("branch")
	buildSetupCmd.MarkFlagRequired("build")
}
