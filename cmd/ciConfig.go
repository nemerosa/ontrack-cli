package cmd

import (
	"fmt"
	"ontrack-cli/client"
	"ontrack-cli/config"
	"os"

	"github.com/spf13/cobra"
)

var ciConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Injection of CI configuration",
	Long: `Injection of CI configuration.

You can create the project, branch and build based on a unique "CI configuration" file.

    ontrack-cli ci config

The YAML configuration files is located by default at .yontrack/ci.yaml and this can 
be changed using the --file option:

	ontrack-cli ci config --file .yontrack/ci.yaml

By default, no environment variable is passed to the configuration, for obvious security reasons.
They need to be passed explicitly using the --env options:

	ontrack-cli ci config \
	  --env GIT_URL=git@github.com:nemerosa/ontrack.git \
	  --env GIT_BRANCH=release/5.0

See the Yontrack documentation for the format of the YAML configuration and the list of needed
environment variables.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}

		fmt.Printf("File: %s\n", file)

		ci, err := cmd.Flags().GetString("ci")
		if err != nil {
			return err
		}
		fmt.Printf("CI: %s\n", ci)

		scm, err := cmd.Flags().GetString("scm")
		if err != nil {
			return err
		}
		fmt.Printf("SCM: %s\n", scm)

		// Get the env values as a map
		envVars, err := getEnvMap(cmd)
		if err != nil {
			return err
		}

		// Print the env variables
		for key, value := range envVars {
			fmt.Printf("Env: %s=%s\n", key, value)
		}

		// Reading the configuration file
		contentBytes, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read configuration file: %w", err)
		}
		configContent := string(contentBytes)

		fmt.Printf("Configuration content:\n%s\n", configContent)

		//project, err := cmd.Flags().GetString("project")
		//if err != nil {
		//	return err
		//}
		//branch, err := cmd.Flags().GetString("branch")
		//if err != nil {
		//	return err
		//}
		//branch = NormalizeBranchName(branch)
		//
		//// Project auto validation stamps
		//autoCreateVS, err := cmd.Flags().GetBool("auto-create-vs")
		//if err != nil {
		//	return err
		//}
		//autoCreateVSAlways, err := cmd.Flags().GetBool("auto-create-vs-always")
		//if err != nil {
		//	return err
		//}
		//if autoCreateVSAlways {
		//	autoCreateVS = true
		//}
		//
		//// Project auto promotion levels
		//autoCreatePL, err := cmd.Flags().GetBool("auto-create-pl")
		//if err != nil {
		//	return err
		//}

		// Configuration
		config, err := config.GetSelectedConfiguration()
		if err != nil {
			return err
		}

		// Returned data
		var data struct {
			ConfigureBuild struct {
				Build struct {
					ID          int
					Name        string
					DisplayName string
					Branch      struct {
						ID          int
						Name        string
						DisplayName string
						Project     struct {
							ID   int
							Name string
						}
					}
				}
				Errors []struct {
					Message string
				}
			}
		}

		if err := client.GraphQLCall(config, `
                mutation OntrackCliCIConfig(
                    $config: String!,
                    $ci: String,
                    $scm: String,
                    $env: [CIEnv!]!,
                ) {
                    configureBuild(input: {
                        config: $config,
                        ci: $ci,
                        scm: $scm,
                        env: $env,
                    }) {
                        errors {
                            message
                            exception
                        }
                        build {
                            id
                            name
                            displayName
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
                }
		`, map[string]interface{}{
			"config": configContent,
			"ci":     ci,
			"scm":    scm,
			"env":    envList,
		}, &data); err != nil {
			return err
		}

		//// Checks errors for the project
		//if err := client.CheckDataErrors(data.CreateProjectOrGet.Errors); err != nil {
		//	return err
		//}
		//// Checks errors for the branch
		//if err := client.CheckDataErrors(data.CreateBranchOrGet.Errors); err != nil {
		//	return err
		//}
		//// Checks errors for the project auto validation stamp propetyu
		//if err := client.CheckDataErrors(data.SetProjectAutoValidationStampProperty.Errors); err != nil {
		//	return err
		//}
		//// Checks errors for the project auto promotion level propetyu
		//if err := client.CheckDataErrors(data.SetProjectAutoPromotionLevelProperty.Errors); err != nil {
		//	return err
		//}

		// OK
		return nil
	},
}

func init() {
	ciCmd.AddCommand(ciConfigCmd)

	ciConfigCmd.Flags().StringP("file", "f", ".yontrack/ci.yaml", "Configuration file")
	ciConfigCmd.Flags().StringSliceP("env", "e", []string{}, "Environment variables in KEY=VALUE format (can be used multiple times)")
	ciConfigCmd.Flags().String("ci", "", "ID of the CI engine to use. If not specified, Yontrack will try to guess it based on the provided environment variables.")
	ciConfigCmd.Flags().String("scm", "", "ID of the SCM engine to use. If not specified, Yontrack will try to guess it based on the provided environment variables.")

	// _ = ciConfigCmd.MarkFlagRequired("file")
}

// getEnvMap parses the --env flags and returns a map of key-value pairs
func getEnvMap(cmd *cobra.Command) (map[string]string, error) {
	envSlice, err := cmd.Flags().GetStringSlice("env")
	if err != nil {
		return nil, err
	}

	envMap := make(map[string]string)
	for _, env := range envSlice {
		// Split on the first '=' only
		parts := splitOnce(env, '=')
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid env format: %s (expected KEY=VALUE)", env)
		}
		envMap[parts[0]] = parts[1]
	}

	return envMap, nil
}

// splitOnce splits a string on the first occurrence of sep
func splitOnce(s string, sep rune) []string {
	for i, c := range s {
		if c == sep {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
}
