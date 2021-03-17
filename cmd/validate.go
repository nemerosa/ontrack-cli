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
	"bytes"
	"errors"
	"text/template"

	"github.com/spf13/cobra"

	client "ontrack-cli/client"
	config "ontrack-cli/config"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validates a build for a validation stamp",
	Long: `Validates a build for a validation stamp.

The simplest form is:

    ontrack-cli validate -p PROJECT -b BRANCH -n BUILD -v VALIDATION -s STATUS

where 'STATUS' is a valid Ontrack status, like 'PASSED', 'WARNING' or 'FAILED'.

In case there is some data to be passed to the validation:

	ontrack-cli validate -p PROJECT -b BRANCH -n BUILD -v VALIDATION \
		--data-type net.nemerosa.ontrack.extension.general.validation.TestSummaryValidationDataType \
		--data {passed: 1, skipped: 2, failed: 3}

In this case, there is no need to pass the status but it could still be forced using the '-s STATUS' flag.

Note that subcommands, dedicated to the most common types are also available. For example:

    ontrack-cli validate -p PROJECT -b BRANCH -n BUILD -v VALIDATION tests --passed 1 --skipped 2 --failed 3

Type 'ontrack-cli validate --help' to get a list of all options.
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

		build, err := cmd.Flags().GetString("build")
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

		status, err := cmd.Flags().GetString("status")
		if err != nil {
			return err
		}

		dataType, err := cmd.Flags().GetString("data-type")
		if err != nil {
			return err
		}

		data, err := cmd.Flags().GetString("data")
		if err != nil {
			return err
		}

		// Status is required if no data
		if status == "" {
			if dataType == "" {
				return errors.New("Status is required if no data is provided.")
			}
		}

		// Data type is required if some data is provided
		if data != "" {
			if dataType == "" {
				return errors.New("Data type is required if some data is provided")
			}
		}

		// Null conversions
		var inputs struct {
			Status      string
			Description string
			DataType    string
			Data        string
		}

		if status == "" {
			inputs.Status = "null"
		} else {
			inputs.Status = `"` + status + `"`
		}

		if description == "" {
			inputs.Description = "null"
		} else {
			inputs.Description = `"` + description + `"`
		}

		if dataType == "" {
			inputs.DataType = "null"
		} else {
			inputs.DataType = `"` + dataType + `"`
		}

		if data == "" {
			inputs.Data = "null"
		} else {
			inputs.Data = data // Pure JSON
		}

		// Run info
		runInfo, err := GetRunInfo(cmd)
		if err != nil {
			return err
		}

		// Query template
		tmpl, err := template.New("mutation").Parse(`
			mutation CreateValidationRun(
				$project: String!,
				$branch: String!,
				$build: String!,
				$validationStamp: String!,
				$runInfo: RunInfoInput
			) {
				createValidationRun(input: {
					project: $project,
					branch: $branch,
					build: $build,
					validationStamp: $validationStamp,
					validationRunStatus: {{ .Status }},
					description: {{ .Description }}
					dataTypeId: {{ .DataType }},
					data: {{ .Data }},
					runInfo: $runInfo
				}) {
					errors {
						message
					}
				}
			}
		`)
		if err != nil {
			return err
		}

		// Query rendering
		var query bytes.Buffer
		if err := tmpl.Execute(&query, inputs); err != nil {
			return err
		}

		// fmt.Println(query.String())

		// Get the configuration
		cfg, err := config.GetSelectedConfiguration()
		if err != nil {
			return err
		}

		// Mutation payload
		var payload struct {
			CreateValidationRun struct {
				Errors []struct {
					Message string
				}
			}
		}

		// Runs the mutation
		if err := client.GraphQLCall(cfg, query.String(), map[string]interface{}{
			"project":             project,
			"branch":              branch,
			"build":               build,
			"validationStamp":     validation,
			"validationRunStatus": status,
			"description":         description,
			"dataTypeId":          dataType,
			"data":                data,
			"runInfo":             runInfo,
		}, &payload); err != nil {
			return err
		}

		// Checks for errors
		if err := client.CheckDataErrors(payload.CreateValidationRun.Errors); err != nil {
			return err
		}

		// OK
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	validateCmd.PersistentFlags().StringP("project", "p", "", "Name of the project")
	validateCmd.PersistentFlags().StringP("branch", "b", "", "Name of the branch")
	validateCmd.PersistentFlags().StringP("build", "n", "", "Name of the build")
	validateCmd.PersistentFlags().StringP("validation", "v", "", "Name of the validation stamp")
	validateCmd.PersistentFlags().StringP("description", "d", "", "Description for the validation")

	validateCmd.MarkPersistentFlagRequired("project")
	validateCmd.MarkPersistentFlagRequired("branch")
	validateCmd.MarkPersistentFlagRequired("build")
	validateCmd.MarkPersistentFlagRequired("validation")

	// Run info arguments
	InitRunInfoCommandFlags(validateCmd)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	validateCmd.Flags().StringP("status", "s", "", "ID of the status (required if no data)")
	validateCmd.Flags().StringP("data-type", "t", "", "FQCN of the validation data type")
	validateCmd.Flags().StringP("data", "o", "", "JSON representation of the validation data")
}
