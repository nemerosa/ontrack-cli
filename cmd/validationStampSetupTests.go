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
	client "ontrack-cli/client"
	config "ontrack-cli/config"

	"github.com/spf13/cobra"
)

// validationStampSetupTestsCmd represents the validationStampSetupTests command
var validationStampSetupTestsCmd = &cobra.Command{
	Use:   "tests",
	Short: "Setup of a tests validation stamp",
	Long: `Setup of a tests validation stamp.

For example:

	ontrack-cli vs setup tests --project PROJECT --branch BRANCH --validation STAMP \
		--warning-if-skipped
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
		branch = NormalizeBranchName(branch)

		validation, err := cmd.Flags().GetString("validation")
		if err != nil {
			return err
		}

		description, err := cmd.Flags().GetString("description")
		if err != nil {
			return err
		}

		warningIfSkipped, err := cmd.Flags().GetBool("warning-if-skipped")
		if err != nil {
			return err
		}

		cfg, err := config.GetSelectedConfiguration()
		if err != nil {
			return err
		}

		var data struct {
			SetupTestSummaryValidationStamp struct {
				Errors []struct {
					Message string
				}
			}
		}
		if err := client.GraphQLCall(cfg, `
			mutation SetupTestSummaryValidationStamp(
				$project: String!,
				$branch: String!,
				$validation: String!,
				$description: String,
				$warningIfSkipped: Boolean!
			) {
				setupTestSummaryValidationStamp(input: {
					project: $project,
					branch: $branch,
					validation: $validation,
					description: $description,
					warningIfSkipped: $warningIfSkipped
				}) {
					errors {
						message
					}
				}
			}
		`, map[string]interface{}{
			"project":          project,
			"branch":           branch,
			"validation":       validation,
			"description":      description,
			"warningIfSkipped": warningIfSkipped,
		}, &data); err != nil {
			return err
		}

		if err := client.CheckDataErrors(data.SetupTestSummaryValidationStamp.Errors); err != nil {
			return err
		}

		// OK
		return nil
	},
}

func init() {
	validationStampSetupCmd.AddCommand(validationStampSetupTestsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validationStampSetupTestsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validationStampSetupTestsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	validationStampSetupTestsCmd.Flags().BoolP("warning-if-skipped", "w", false, "Configures the validation stamp to fail if there are some skipped tests.")
}
