package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"yontrack/client"
	"yontrack/config"

	"github.com/spf13/cobra"
)

var ciConfigBranchCmd = &cobra.Command{
	Use:   "config-branch",
	Short: "Injection of CI configuration for a branch only",
	Long: `Injection of CI configuration for a branch only.

You can create the project and branch based on a unique "CI configuration" file.

    yontrack ci config-branch

The YAML configuration files is located by default at .yontrack/ci.yaml and this can 
be changed using the --file option:

	yontrack ci config-branch --file .yontrack/ci.yaml

By default, no environment variable is passed to the configuration, for obvious security reasons.
They need to be passed explicitly using the --env options:

	yontrack ci config-branch \
	  --env GIT_URL=git@github.com:nemerosa/ontrack.git \
	  --env GIT_BRANCH=release/5.0

See the Yontrack documentation for the format of the YAML configuration and the list of needed
environment variables.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ciContext, err := getCIConfigContext(cmd)
		if err != nil {
			return err
		}

		// Configuration
		config, err := config.GetSelectedConfiguration()
		if err != nil {
			return err
		}

		// Returned data
		var data struct {
			ConfigureBranch struct {
				Branch struct {
					ID          string
					Name        string
					DisplayName string
					Project     struct {
						ID   string
						Name string
					}
				}
				Errors []struct {
					Message string
				}
			}
		}

		if err := client.GraphQLCall(config, `
                mutation OntrackCliCIConfigBranch(
                    $config: String!,
                    $ci: String,
                    $scm: String,
                    $env: [CIEnv!]!,
                ) {
                    configureBranch(input: {
                        config: $config,
                        ci: $ci,
                        scm: $scm,
                        env: $env,
                    }) {
                        errors {
                            message
                            exception
                        }
                        branch {
							id
							name
							displayName
							project {
								id
								name
							}
                        }
                    }
                }
		`, map[string]interface{}{
			"config": ciContext.ConfigContent,
			"ci":     ciContext.CI,
			"scm":    ciContext.SCM,
			"env":    ciContext.EnvList,
		}, &data); err != nil {
			return err
		}

		// Checks errors
		if err := client.CheckDataErrors(data.ConfigureBranch.Errors); err != nil {
			return err
		}

		// Output
		if ciContext.Output != "" {
			branch := data.ConfigureBranch.Branch
			if ciContext.Output == "env" {
				_, _ = fmt.Fprintf(os.Stdout, "export YONTRACK_PROJECT_ID=%s\n", branch.Project.ID)
				_, _ = fmt.Fprintf(os.Stdout, "export YONTRACK_PROJECT_NAME=%s\n", branch.Project.Name)
				_, _ = fmt.Fprintf(os.Stdout, "export YONTRACK_BRANCH_ID=%s\n", branch.ID)
				_, _ = fmt.Fprintf(os.Stdout, "export YONTRACK_BRANCH_NAME=%s\n", branch.Name)
			} else if ciContext.Output == "json" {
				jsonBytes, err := json.MarshalIndent(branch, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal branch to JSON: %w", err)
				}
				_, _ = fmt.Fprintf(os.Stdout, "%s\n", jsonBytes)
			} else {
				return fmt.Errorf("unsupported output type %s", ciContext.Output)
			}
		}

		// OK
		return nil
	},
}

func init() {
	ciCmd.AddCommand(ciConfigBranchCmd)
	registerCIConfigFlags(ciConfigBranchCmd)
}
