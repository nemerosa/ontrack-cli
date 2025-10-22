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

	client "yontrack/client"
	config "yontrack/config"
)

// validationStampCmd represents the validationStamp command
var validationStampCmd = &cobra.Command{
	Use:     "validation-stamp",
	Aliases: []string{"validation", "vs"},
	Short:   "Management of validation stamps",
	Long:    `Management of validation stamps.`,
	// Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(validationStampCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validationStampCmd.PersistentFlags().String("foo", "", "A help for foo")
	validationStampCmd.PersistentFlags().StringP("project", "p", "", "Name of the project")
	validationStampCmd.PersistentFlags().StringP("branch", "b", "", "Name of the branch")

	validationStampCmd.MarkPersistentFlagRequired("project")
	validationStampCmd.MarkPersistentFlagRequired("branch")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validationStampCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Utility method to setup a "tests" validation stamp
func SetupTestValidationStamp(
	project string,
	branch string,
	validation string,
	description string,
	warningIfSkipped bool,
	failWhenNoResults bool,
) error {

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

	err = client.GraphQLCall(cfg, `
		mutation SetupTestSummaryValidationStamp(
			$project: String!,
			$branch: String!,
			$validation: String!,
			$description: String,
			$warningIfSkipped: Boolean!,
			$failWhenNoResults: Boolean!,
		) {
			setupTestSummaryValidationStamp(input: {
				project: $project,
				branch: $branch,
				validation: $validation,
				description: $description,
				warningIfSkipped: $warningIfSkipped,
				failWhenNoResults: $failWhenNoResults,
			}) {
				errors {
					message
				}
			}
		}
	`, map[string]interface{}{
		"project":           project,
		"branch":            branch,
		"validation":        validation,
		"description":       description,
		"warningIfSkipped":  warningIfSkipped,
		"failWhenNoResults": failWhenNoResults,
	}, &data)

	if err != nil {
		return err
	}

	return client.CheckDataErrors(data.SetupTestSummaryValidationStamp.Errors)
}
