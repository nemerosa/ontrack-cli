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

// branchSetPropertyGitCmd represents the branchSetPropertyGit command
var branchSetPropertyGitCmd = &cobra.Command{
	Use:   "git",
	Short: "Sets the Git configuration for a branch",
	Long: `Sets the Git configuration for a branch.

For example:

    ontrack-cli branch set-property --project PROJECT --branch BRANCH git --git-branch main

As of now, this also sets the "GitCommitPropertyLink" property by default (builds must have a Git commit).
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

		gitBranch, err := cmd.Flags().GetString("git-branch")
		if err != nil {
			return err
		}

		cfg, err := config.GetSelectedConfiguration()
		if err != nil {
			return err
		}

		var data struct {
			SetBranchGitConfigProperty struct {
				Errors []struct {
					Message string
				}
			}
		}
		if err := client.GraphQLCall(cfg, `
			mutation SetBranchGitConfigProperty(
				$project: String!,
				$branch: String!,
				$gitBranch: String!
			) {
				setBranchGitConfigProperty(input: {
					project: $project,
					branch: $branch,
					gitBranch: $gitBranch
				}) {
					errors {
						message
					}
				}
			}
		`, map[string]interface{}{
			"project":   project,
			"branch":    branch,
			"gitBranch": gitBranch,
		}, &data); err != nil {
			return err
		}

		if err := client.CheckDataErrors(data.SetBranchGitConfigProperty.Errors); err != nil {
			return err
		}

		// OK
		return nil
	},
}

func init() {
	branchSetPropertyCmd.AddCommand(branchSetPropertyGitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// branchSetPropertyGitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// branchSetPropertyGitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	branchSetPropertyGitCmd.Flags().StringP("git-branch", "g", "", "Git branch to associate with the branch")

	branchSetPropertyGitCmd.MarkFlagRequired("git-branch")
}
