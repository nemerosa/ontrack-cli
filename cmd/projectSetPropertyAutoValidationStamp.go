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

	"github.com/spf13/cobra"
)

// projectSetPropertyAutoValidationStampCmd represents the projectSetPropertyAutoValidationStamp command
var projectSetPropertyAutoValidationStampCmd = &cobra.Command{
	Use:     "auto-validation-stamp",
	Aliases: []string{"auto-vs", "avs"},
	Short:   "Sets the auto creation of validation stamps property on a project",
	Long: `Sets the auto creation of validation stamps property on a project.
	
For example, to set a project to create validation stamps only when they are predefined:

	yontrack project set-property -p PROJECT avs
	
For example, to set a project to create validation stamps in all cases:

	yontrack project set-property -p PROJECT avs --auto-create-if-not-predefined

The '--auto-create' option can be used to disable this behaviour altogether:

    yontrack project set-property -p PROJECT avs --auto-create=false
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, err := cmd.Flags().GetString("project")
		if err != nil {
			return err
		}

		autoCreate, err := cmd.Flags().GetBool("auto-create")
		if err != nil {
			return err
		}

		autoCreateIfNotPredefined, err := cmd.Flags().GetBool("auto-create-if-not-predefined")
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
			SetProjectAutoValidationStampProperty struct {
				Errors []struct {
					Message string
				}
			}
		}

		// GraphQL call
		if err := client.GraphQLCall(cfg, `
			mutation SetProjectAutoValidationStampProperty(
				$project: String!,
				$autoCreate: Boolean!,
				$autoCreateIfNotPredefined: Boolean!
			) {
				setProjectAutoValidationStampProperty(input: {
					project: $project,
					isAutoCreate: $autoCreate,
					isAutoCreateIfNotPredefined: $autoCreateIfNotPredefined
				}) {
					errors {
						message
					}
				}
			}
		`, map[string]interface{}{
			"project":                   project,
			"autoCreate":                autoCreate,
			"autoCreateIfNotPredefined": autoCreateIfNotPredefined,
		}, &data); err != nil {
			return err
		}

		if err := client.CheckDataErrors(data.SetProjectAutoValidationStampProperty.Errors); err != nil {
			return err
		}

		// OK
		return nil
	},
}

func init() {
	projectSetPropertyCmd.AddCommand(projectSetPropertyAutoValidationStampCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// projectSetPropertyAutoValidationStampCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// projectSetPropertyAutoValidationStampCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	projectSetPropertyAutoValidationStampCmd.Flags().Bool("auto-create", true, "If validation stamps must be created from predefined validation stamps")
	projectSetPropertyAutoValidationStampCmd.Flags().Bool("auto-create-if-not-predefined", false, "If validation stamps must be created even if predefined validation stamp is not available")
}
