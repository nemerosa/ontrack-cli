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
	"yontrack/utils"

	"github.com/spf13/cobra"

	client "yontrack/client"
	config "yontrack/config"
)

// promoteCmd represents the promote command
var promoteCmd = &cobra.Command{
	Use:   "promote",
	Short: "Promotes a build",
	Long: `Promotes a build.
	
	yontrack promote -p PROJECT -b BRANCH -n BUILD -l PROMOTION -d DESCRIPTION
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, branch, build, err := utils.GetProjectBranchBuildFlags(cmd, false, true)
		if err != nil {
			return err
		}
		promotion, err := cmd.Flags().GetString("promotion")
		if err != nil {
			return err
		}
		description, err := cmd.Flags().GetString("description")
		if err != nil {
			return err
		}

		// Configuration
		cfg, err := config.GetSelectedConfiguration()
		if err != nil {
			return err
		}

		// Data
		var data struct {
			CreatePromotionRun struct {
				Errors []struct {
					Message string
				}
			}
		}

		// Call
		if err := client.GraphQLCall(cfg, `
			mutation CreatePromotionRun(
				$project: String!,
				$branch: String!,
				$build: String!,
				$promotion: String!,
				$description: String
			) {
				createPromotionRun(input: {
					project: $project,
					branch: $branch,
					build: $build,
					promotion: $promotion,
					description: $description
				}) {
					errors {
						message
					}
				}
			}
		`, map[string]interface{}{
			"project":     project,
			"branch":      branch,
			"build":       build,
			"promotion":   promotion,
			"description": description,
		}, &data); err != nil {
			return err
		}

		// Error check
		if err := client.CheckDataErrors(data.CreatePromotionRun.Errors); err != nil {
			return err
		}

		// OK
		return nil
	},
}

func init() {
	rootCmd.AddCommand(promoteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// promoteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// promoteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	promoteCmd.Flags().StringP("project", "p", "", "Name of the project")
	promoteCmd.Flags().StringP("branch", "b", "", "Name of the branch")
	promoteCmd.Flags().StringP("build", "n", "", "Name of the build")
	promoteCmd.Flags().StringP("promotion", "l", "", "Name of the promotion level")
	promoteCmd.Flags().StringP("description", "d", "", "Description for the promotion")

	promoteCmd.MarkFlagRequired("project")
	promoteCmd.MarkFlagRequired("branch")
	promoteCmd.MarkFlagRequired("build")
	promoteCmd.MarkFlagRequired("promotion")
}
