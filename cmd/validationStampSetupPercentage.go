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
	client "yontrack/client"
	config "yontrack/config"
	"yontrack/utils"

	"github.com/spf13/cobra"
)

// validationStampSetupPercentageCmd represents the validationStampSetupPercentage command
var validationStampSetupPercentageCmd = &cobra.Command{
	Use:     "percentage",
	Aliases: []string{"percent"},
	Short:   "Setup of a percentage validation stamp",
	Long: `Setup of a percentage validation stamp.

For example:

	yontrack vs setup percentage --project PROJECT --branch BRANCH --validation STAMP \
		--warning-if-skipped
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, branch, err := utils.GetProjectBranchFlags(cmd, false, true)
		if err != nil {
			return err
		}

		validation, err := cmd.Flags().GetString("validation")
		if err != nil {
			return err
		}

		description, err := cmd.Flags().GetString("description")
		if err != nil {
			return err
		}

		warning, err := cmd.Flags().GetInt("warning")
		if err != nil {
			return err
		}
		var warningValue *int
		if warning > 0 {
			warningValue = &warning
		} else {
			warningValue = nil
		}

		failure, err := cmd.Flags().GetInt("failure")
		if err != nil {
			return err
		}
		var failureValue *int
		if failure > 0 {
			failureValue = &failure
		} else {
			failureValue = nil
		}

		okIfGreater, err := cmd.Flags().GetBool("ok-if-greater")
		if err != nil {
			return err
		}

		cfg, err := config.GetSelectedConfiguration()
		if err != nil {
			return err
		}

		var data struct {
			SetupPercentageValidationStamp struct {
				Errors []struct {
					Message string
				}
			}
		}
		if err := client.GraphQLCall(cfg, `
			mutation SetupPercentageValidationStamp(
				$project: String!,
				$branch: String!,
				$validation: String!,
				$description: String,
				$warning: Int,
				$failure: Int,
				$okIfGreater: Boolean!
			) {
				setupPercentageValidationStamp(input: {
					project: $project,
					branch: $branch,
					validation: $validation,
					description: $description,
					warningThreshold: $warning,
					failureThreshold: $failure,
					okIfGreater: $okIfGreater
				}) {
					errors {
						message
					}
				}
			}
		`, map[string]interface{}{
			"project":     project,
			"branch":      branch,
			"validation":  validation,
			"description": description,
			"warning":     warningValue,
			"failure":     failureValue,
			"okIfGreater": okIfGreater,
		}, &data); err != nil {
			return err
		}

		if err := client.CheckDataErrors(data.SetupPercentageValidationStamp.Errors); err != nil {
			return err
		}

		// OK
		return nil
	},
}

func init() {
	validationStampSetupCmd.AddCommand(validationStampSetupPercentageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validationStampSetupPercentageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validationStampSetupPercentageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	validationStampSetupPercentageCmd.Flags().IntP("warning", "w", 0, "Threshold value for a warning")
	validationStampSetupPercentageCmd.Flags().IntP("failure", "f", 0, "Threshold value for a failure")
	validationStampSetupPercentageCmd.Flags().BoolP("ok-if-greater", "o", false, "Direction of the value scale")

}
