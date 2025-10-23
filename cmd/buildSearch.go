package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"yontrack/client"
	config "yontrack/config"

	"github.com/spf13/cobra"
)

// buildSearchCmd represents the buildSearch command
var buildSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Searches for builds",
	Long: `Searches for builds.

Builds can be searched on a project using the '--project' flag:

    yontrack build search --project PROJECT

or on a branch using the '--branch' flag:

    yontrack build search --project PROJECT --branch BRANCH

In both cases, several criteria are available - see 'yontrack build search --help' to get their list. For example,
to look for a build using its commit:

    yontrack build search --project PROJECT --branch BRANCH --commit commit

By default, only the build names are printed, one per line.

You can change the display options using additional flags - see 'yontrack build search --help' to get their list.`,
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

		// Project vs. branch search
		if branch == "" {
			return projectSearch(cmd, project)
		} else {
			return branchSearch(cmd, project, branch)
		}
	},
}

type build struct {
	Id          string
	Name        string
	DisplayName string
	Branch      struct {
		Id          string
		Name        string
		DisplayName string
		Project     struct {
			Id   string
			Name string
		}
	}
}

type buildList struct {
	Builds []build
}

func projectSearch(cmd *cobra.Command, project string) error {
	// Query
	query := `
		query BuildProjectSearch(
			$project: String!,
			$buildProjectFilter: BuildSearchForm!
		) {
			builds(
				project: $project,
				buildProjectFilter: $buildProjectFilter
			) {
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
	`

	// Search form
	form := make(map[string]interface{})
	if err := fillFormWithProperty(cmd, &form, "property", "propertyValue"); err != nil {
		return err
	}
	if err := fillFormWithCount(cmd, &form, "maximumCount"); err != nil {
		return err
	}

	if err := fillFormWithWithPromotion(cmd, &form, "promotionName"); err != nil {
		return err
	}

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	} else if name != "" {
		form["buildName"] = name
		nameExact, err := cmd.Flags().GetBool("name-exact")
		if err != nil {
			return err
		} else if nameExact {
			form["buildExactMatch"] = true
		}
	}

	// Gets the configuration
	cfg, err := config.GetSelectedConfiguration()
	if err != nil {
		return err
	}

	// Result data
	var data buildList

	// Call
	if err := client.GraphQLCall(cfg, query, map[string]interface{}{
		"project":            project,
		"buildProjectFilter": form,
	}, &data); err != nil {
		return err
	}

	// Displaying the data
	return displayBuilds(cmd, &data)
}

func branchSearch(cmd *cobra.Command, project string, branch string) error {
	// Query
	query := `
		query BuildBranchSearch(
			$project: String!,
			$branch: String!,
			$buildBranchFilter: StandardBuildFilter!
		) {
			builds(
				project: $project,
				branch: $branch,
				buildBranchFilter: $buildBranchFilter
			) {
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
	`

	// Search form
	form := make(map[string]interface{})
	if err := fillFormWithProperty(cmd, &form, "withProperty", "withPropertyValue"); err != nil {
		return err
	}
	if err := fillFormWithCount(cmd, &form, "count"); err != nil {
		return err
	}

	if err := fillFormWithWithPromotion(cmd, &form, "withPromotionLevel"); err != nil {
		return err
	}

	// Gets the configuration
	cfg, err := config.GetSelectedConfiguration()
	if err != nil {
		return err
	}

	// Result data
	var data buildList

	// Call
	if err := client.GraphQLCall(cfg, query, map[string]interface{}{
		"project":           project,
		"branch":            branch,
		"buildBranchFilter": form,
	}, &data); err != nil {
		return err
	}

	// Displaying the data
	return displayBuilds(cmd, &data)
}

func displayBuilds(cmd *cobra.Command, data *buildList) error {
	displayBranch, err := cmd.Flags().GetBool("display-branch")
	if err != nil {
		return err
	}
	displayId, err := cmd.Flags().GetBool("display-id")
	if err != nil {
		return err
	}
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	if output == "" {
		for _, build := range data.Builds {
			if displayBranch {
				fmt.Printf(build.Branch.Name)
			} else if displayId {
				fmt.Println(build.Id)
			} else {
				fmt.Println(build.Name)
			}
		}
	} else {
		count, err := cmd.Flags().GetInt("count")
		if err != nil {
			return err
		}
		if count == 1 {
			if len(data.Builds) == 0 {
				return fmt.Errorf("no build found")
			} else {
				build := data.Builds[0]
				return displayBuildOutput(build, output)
			}
		} else {
			return fmt.Errorf("output not supported with multiple builds. Set count to 1 or remove the output flag")
		}
	}

	return nil
}

func displayBuildOutput(build build, output string) error {
	if output == "env" {
		return displayBuildEnv(build)
	} else if output == "json" {
		return displayBuildJson(build)
	} else {
		return fmt.Errorf("unknown output format: %s", output)
	}
}

func displayBuildEnv(build build) error {
	_, _ = fmt.Fprintf(os.Stdout, "export YONTRACK_PROJECT_ID=%s\n", build.Branch.Project.Id)
	_, _ = fmt.Fprintf(os.Stdout, "export YONTRACK_PROJECT_NAME=%s\n", build.Branch.Project.Name)
	_, _ = fmt.Fprintf(os.Stdout, "export YONTRACK_BRANCH_ID=%s\n", build.Branch.Id)
	_, _ = fmt.Fprintf(os.Stdout, "export YONTRACK_BRANCH_NAME=%s\n", build.Branch.Name)
	_, _ = fmt.Fprintf(os.Stdout, "export YONTRACK_BUILD_ID=%s\n", build.Id)
	_, _ = fmt.Fprintf(os.Stdout, "export YONTRACK_BUILD_NAME=%s\n", build.Name)
	return nil
}

func displayBuildJson(build build) error {
	jsonBytes, err := json.MarshalIndent(build, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal build to JSON: %w", err)
	}
	_, err = fmt.Fprintf(os.Stdout, "%s\n", jsonBytes)
	return err
}

func fillFormWithCount(cmd *cobra.Command, form *map[string]interface{}, countField string) error {
	count, err := cmd.Flags().GetInt("count")
	if err != nil {
		return err
	}

	if count > 0 {
		(*form)[countField] = count
	}

	return nil
}

func fillFormWithProperty(cmd *cobra.Command, form *map[string]interface{}, propertyTypeField string, propertyValueField string) error {
	commit, err := cmd.Flags().GetString("commit")
	if err != nil {
		return err
	}

	if commit != "" {
		(*form)[propertyTypeField] = "net.nemerosa.ontrack.extension.git.property.GitCommitPropertyType"
		(*form)[propertyValueField] = commit
	}

	return nil
}

func fillFormWithWithPromotion(cmd *cobra.Command, form *map[string]interface{}, fieldName string) error {
	return fillForm(cmd, form, "with-promotion", fieldName)
}

func fillForm(cmd *cobra.Command, form *map[string]interface{}, argName string, fieldName string) error {
	value, err := cmd.Flags().GetString(argName)
	if err != nil {
		return err
	}

	if value != "" {
		(*form)[fieldName] = value
	}

	return nil
}

func init() {
	buildCmd.AddCommand(buildSearchCmd)

	buildSearchCmd.Flags().StringP("project", "p", "", "Name of the project")
	_ = buildSearchCmd.MarkFlagRequired("project")

	buildSearchCmd.Flags().StringP("branch", "b", "", "Name of the branch")

	// Criteria
	buildSearchCmd.Flags().Int("count", 10, "Number of builds to return")
	buildSearchCmd.Flags().String("with-promotion", "", "Builds must have this promotion")
	buildSearchCmd.Flags().String("name", "", "Builds must have this name or match this regular expression")
	buildSearchCmd.Flags().Bool("name-exact", true, "If present together with the `name` flag, requires an exact match.")

	// Property criteria
	buildSearchCmd.Flags().String("commit", "", "Commit for the build")

	// Display options
	buildSearchCmd.Flags().StringP("output", "o", "", "How to output the search results (env, json). Incompatible with the `display` options.")
	buildSearchCmd.Flags().Bool("display-id", false, "Displays the build ID instead of its name.")
	buildSearchCmd.Flags().Bool("display-branch", false, "Displays the build branch name instead of its name.")
}
