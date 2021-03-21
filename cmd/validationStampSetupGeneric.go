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
	"text/template"

	"github.com/spf13/cobra"

	client "ontrack-cli/client"
	config "ontrack-cli/config"
)

// validationStampSetupGenericCmd represents the validationStampSetupGeneric command
var validationStampSetupGenericCmd = &cobra.Command{
	Use:   "generic",
	Short: "Setup a validation stamp using a generic format",
	Long: `Setup a validation stamp using a generic format.

To create a plain validation stamp (without any data type):

	ontrack-cli vs setup generic --project PROJECT --branch BRANCH --validation STAMP

You can also associate a data type with it, using the JSON representation of the configuration:

    ontrack-cli vs setup generic --project PROJECT --branch BRANCH --validation STAMP \
        --data-type "net.nemerosa.ontrack.extension.general.validation.CHMLValidationDataType" \
        --data-config '{warningLevel: {level: "HIGH",value:1}, failedLevel:{level:"CRITICAL",value:1}}'

Note that specific commands per type are also available, see 'ontrack-cli vs setup'.
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

		dataType, err := cmd.Flags().GetString("data-type")
		if err != nil {
			return err
		}

		dataTypeConfig, err := cmd.Flags().GetString("data-config")
		if err != nil {
			return err
		}

		cfg, err := config.GetSelectedConfiguration()
		if err != nil {
			return err
		}

		tmpl, err := template.New("mutation").Parse(`
			mutation SetupValidationStamp(
				$project: String!,
				$branch: String!,
				$validation: String!,
				$description: String,
				$dataType: String
			) {
				setupValidationStamp(input: {
					project: $project,
					branch: $branch,
					validation: $validation,
					description: $description,
					dataType: $dataType,
					dataTypeConfig: {{ .DataTypeConfig }}
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

		var tmplInput struct {
			DataTypeConfig string
		}

		if dataTypeConfig != "" {
			tmplInput.DataTypeConfig = dataTypeConfig
		} else {
			tmplInput.DataTypeConfig = "null"
		}

		var query bytes.Buffer
		if err := tmpl.Execute(&query, tmplInput); err != nil {
			return err
		}

		var data struct {
			SetupValidationStamp struct {
				Errors []struct {
					Message string
				}
			}
		}
		if err := client.GraphQLCall(cfg, query.String(), map[string]interface{}{
			"project":        project,
			"branch":         branch,
			"validation":     validation,
			"description":    description,
			"dataType":       dataType,
			"dataTypeConfig": dataTypeConfig,
		}, &data); err != nil {
			return err
		}

		// Checks for errors
		if err := client.CheckDataErrors(data.SetupValidationStamp.Errors); err != nil {
			return err
		}

		// OK
		return nil
	},
}

func init() {
	validationStampSetupCmd.AddCommand(validationStampSetupGenericCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validationStampSetupGenericCmd.PersistentFlags().String("foo", "", "A help for foo")

	validationStampSetupGenericCmd.PersistentFlags().StringP("data-type", "t", "", "FQCN of the data type")
	validationStampSetupGenericCmd.PersistentFlags().StringP("data-config", "c", "", "JSON for the data type configuration")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validationStampSetupGenericCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
