package cmd

import (
	"github.com/spf13/cobra"
	"ontrack-cli/client"
	"ontrack-cli/config"
)

var buildAutoversioningCheckCmd = &cobra.Command{
	Use:     "auto-versioning-check",
	Aliases: []string{"av-check"},
	Short:   "Launches the auto-versioning check",
	Long: `Launches the auto-versioning check for a given build.

For example:

    ontrack-cli build auto-versioning-check --project my-project --branch release/1.0 --build 1`,
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

		build, err := cmd.Flags().GetString("build")
		if err != nil {
			return err
		}

		// Configuration
		config, err := config.GetSelectedConfiguration()
		if err != nil {
			return err
		}

		var data struct {
			CheckAutoVersioning struct {
				Errors []struct {
					Message string
				}
			}
		}

		// GraphQL call
		if err := client.GraphQLCall(config, `
			mutation CheckAutoVersioning(
				$project: String!,
				$branch: String!,
				$build: String!,
			) {
				checkAutoVersioning(input: {
					project: $project,
					branch: $branch,
					build: $build,
				}) {
					errors {
					  message
					}
				}
			}
		`, map[string]interface{}{
			"project": project,
			"branch":  branch,
			"build":   build,
		}, &data); err != nil {
			return err
		}

		// Checks errors
		if err := client.CheckDataErrors(data.CheckAutoVersioning.Errors); err != nil {
			return err
		}

		// OK
		return nil
	},
}

func init() {
	buildCmd.AddCommand(buildAutoversioningCheckCmd)

	buildAutoversioningCheckCmd.PersistentFlags().StringP("project", "p", "", "Name of the project")
	buildAutoversioningCheckCmd.PersistentFlags().StringP("branch", "b", "", "Name of the branch")
	buildAutoversioningCheckCmd.PersistentFlags().StringP("build", "n", "", "Name of the build")

	buildAutoversioningCheckCmd.MarkPersistentFlagRequired("project")
	buildAutoversioningCheckCmd.MarkPersistentFlagRequired("branch")
	buildAutoversioningCheckCmd.MarkPersistentFlagRequired("build")
}
