package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"yontrack/client"
	"yontrack/config"
	"yontrack/utils"

	"github.com/spf13/cobra"
)

var ciConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Injection of CI configuration",
	Long: `Injection of CI configuration.

You can create the project, branch and build based on a unique "CI configuration" file.

    yontrack ci config

The YAML configuration files is located by default at .yontrack/ci.yaml and this can 
be changed using the --file option:

	yontrack ci config --file .yontrack/ci.yaml

By default, no environment variable is passed to the configuration, for obvious security reasons.
They need to be passed explicitly using the --env options:

	yontrack ci config \
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

		_, _ = fmt.Fprintf(os.Stderr, "File: %s\n", file)

		ci, err := cmd.Flags().GetString("ci")
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintf(os.Stderr, "CI: %s\n", ci)

		scm, err := cmd.Flags().GetString("scm")
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintf(os.Stderr, "SCM: %s\n", scm)

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintf(os.Stderr, "Output: %s\n", output)

		// Env vars from a file
		envFile, err := cmd.Flags().GetString("env-file")
		if err != nil {
			return err
		}

		// Start with env vars from file (if provided)
		envVars := make(map[string]string)
		if envFile != "" {
			fileEnvVars, err := utils.ReadEnvFile(envFile)
			if err != nil {
				return fmt.Errorf("failed to read env file: %w", err)
			}
			// Add all vars from file
			for key, value := range fileEnvVars {
				envVars[key] = value
			}
		}

		// Get the env values from --env-all flags and merge
		envAll, err := cmd.Flags().GetStringSlice("env-all")
		if err != nil {
			return err
		}

		for _, prefix := range envAll {
			// Iterate through all environment variables
			for _, envVar := range os.Environ() {
				// Split on the first '=' to get key and value
				parts := splitOnce(envVar, '=')
				if len(parts) == 2 {
					key := parts[0]
					value := parts[1]
					// Check if the key starts with the prefix
					if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
						envVars[key] = value
					}
				}
			}
		}

		// Get the env values from --env flags and merge (these take priority)
		cmdEnvVars, err := getEnvMap(cmd)
		if err != nil {
			return err
		}
		// Override/add vars from --env flags (they have highest priority)
		for key, value := range cmdEnvVars {
			envVars[key] = value
		}

		// Print the env variables
		for key, value := range envVars {
			_, _ = fmt.Fprintf(os.Stderr, "Env: %s=%s\n", key, value)
		}

		// Convert envVars map to a list of CIEnv objects
		envList := make([]map[string]interface{}, 0, len(envVars))
		for key, value := range envVars {
			envList = append(envList, map[string]interface{}{
				"name":  key,
				"value": value,
			})
		}

		// Reading the configuration file
		contentBytes, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read configuration file: %w", err)
		}
		initialConfigContent := string(contentBytes)
		configContent, err := utils.ExpandConfig(initialConfigContent)
		if err != nil {
			return fmt.Errorf("failed to expand configuration: %w", err)
		}

		_, _ = fmt.Fprintf(os.Stderr, "Configuration content:\n%s\n", configContent)

		// Configuration
		config, err := config.GetSelectedConfiguration()
		if err != nil {
			return err
		}

		// Returned data
		var data struct {
			ConfigureBuild struct {
				Build struct {
					ID          string
					Name        string
					DisplayName string
					Branch      struct {
						ID          string
						Name        string
						DisplayName string
						Project     struct {
							ID   string
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

		// Checks errors
		if err := client.CheckDataErrors(data.ConfigureBuild.Errors); err != nil {
			return err
		}

		// Output
		if output != "" {
			build := data.ConfigureBuild.Build
			if output == "env" {
				_, _ = fmt.Fprintf(os.Stdout, "export YONTRACK_PROJECT_ID=%s\n", build.Branch.Project.ID)
				_, _ = fmt.Fprintf(os.Stdout, "export YONTRACK_PROJECT_NAME=%s\n", build.Branch.Project.Name)
				_, _ = fmt.Fprintf(os.Stdout, "export YONTRACK_BRANCH_ID=%s\n", build.Branch.ID)
				_, _ = fmt.Fprintf(os.Stdout, "export YONTRACK_BRANCH_NAME=%s\n", build.Branch.Name)
				_, _ = fmt.Fprintf(os.Stdout, "export YONTRACK_BUILD_ID=%s\n", build.ID)
				_, _ = fmt.Fprintf(os.Stdout, "export YONTRACK_BUILD_NAME=%s\n", build.Name)
			} else if output == "json" {
				jsonBytes, err := json.MarshalIndent(build, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal build to JSON: %w", err)
				}
				_, _ = fmt.Fprintf(os.Stdout, "%s\n", jsonBytes)
			} else {
				return fmt.Errorf("unsupported output type %s", output)
			}
		}

		// OK
		return nil
	},
}

func init() {
	ciCmd.AddCommand(ciConfigCmd)

	ciConfigCmd.Flags().StringP("file", "f", ".yontrack/ci.yaml", "Configuration file")
	ciConfigCmd.Flags().StringSliceP("env", "e", []string{}, "Environment variables in KEY=VALUE format (can be used multiple times)")
	ciConfigCmd.Flags().StringSlice("env-all", []string{}, "Uses the specified prefix to select environment variables to inject.")
	ciConfigCmd.Flags().String("env-file", "", "Path to an env file containing key/values (one per line, using the KEY=VALUE format)")
	ciConfigCmd.Flags().String("ci", "", "ID of the CI engine to use. If not specified, Yontrack will try to guess it based on the provided environment variables.")
	ciConfigCmd.Flags().String("scm", "", "ID of the SCM engine to use. If not specified, Yontrack will try to guess it based on the provided environment variables.")
	ciConfigCmd.Flags().StringP("output", "o", "", "Output of the command: env, json.")

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
