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
	"errors"

	"github.com/spf13/cobra"

	"regexp"
	"strconv"

	client "yontrack/client"
	config "yontrack/config"
)

// validationStampSetupCHMLCmd represents the validationStampSetupCHML command
var validationStampSetupCHMLCmd = &cobra.Command{
	Use:   "chml",
	Short: "Setup of a CHML validation stamp",
	Long: `Setup of a CHML validation stamp.

For example:

	yontrack vs setup chml --project PROJECT --branch BRANCH --validation STAMP \
		--warning HIGH=1 \
		--failed CRITICAL=1
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

		warningLevel, warningValue, err := parseLevel(cmd, "warning")
		if err != nil {
			return err
		}

		failedLevel, failedValue, err := parseLevel(cmd, "failed")
		if err != nil {
			return err
		}

		cfg, err := config.GetSelectedConfiguration()
		if err != nil {
			return err
		}

		var data struct {
			SetupCHMLValidationStamp struct {
				Errors []struct {
					Message string
				}
			}
		}
		if err := client.GraphQLCall(cfg, `
			mutation SetupCHMLValidationStamp(
				$project: String!,
				$branch: String!,
				$validation: String!,
				$description: String,
				$warningLevel: CHML!,
				$warningValue: Int!,
				$failedLevel: CHML!,
				$failedValue: Int!
			) {
				setupCHMLValidationStamp(input: {
					project: $project,
					branch: $branch,
					validation: $validation,
					description: $description,
					warningLevel: {
						level: $warningLevel,
						value: $warningValue
					},
					failedLevel: {
						level: $failedLevel,
						value: $failedValue
					}
				}) {
					errors {
						message
					}
				}
			}
		`, map[string]interface{}{
			"project":      project,
			"branch":       branch,
			"validation":   validation,
			"description":  description,
			"warningLevel": warningLevel,
			"warningValue": warningValue,
			"failedLevel":  failedLevel,
			"failedValue":  failedValue,
		}, &data); err != nil {
			return err
		}

		if err := client.CheckDataErrors(data.SetupCHMLValidationStamp.Errors); err != nil {
			return err
		}

		// OK
		return nil
	},
}

func parseLevel(cmd *cobra.Command, name string) (string, int, error) {
	arg, err := cmd.Flags().GetString(name)
	if err != nil {
		return "", 0, err
	}

	re := regexp.MustCompile(`^(CRITICAL|HIGH|MEDIUM|LOW)=(\d+)$`)
	match := re.FindStringSubmatch(arg)
	if match == nil {
		return "", 0, errors.New("Argument " + name + " with value \"" + arg + "\" must match (CRITICAL|HIGH|MEDIUM|LOW)=<value>")
	}

	level := match[1]
	value, err := strconv.Atoi(match[2])
	if err != nil {
		return "", 0, err
	}

	return level, value, nil
}

func init() {
	validationStampSetupCmd.AddCommand(validationStampSetupCHMLCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validationStampSetupCHMLCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validationStampSetupCHMLCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	validationStampSetupCHMLCmd.Flags().StringP("warning", "w", "", "Warning threshold, in the form of 'level=value'")
	validationStampSetupCHMLCmd.Flags().StringP("failed", "f", "", "Failure threshold, in the form of 'level=value'")

	validationStampSetupCHMLCmd.MarkFlagRequired("warning")
	validationStampSetupCHMLCmd.MarkFlagRequired("failed")
}
