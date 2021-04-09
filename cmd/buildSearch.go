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
	"fmt"
	"ontrack-cli/client"
	config "ontrack-cli/config"

	"github.com/spf13/cobra"
)

// buildSearchCmd represents the buildSearch command
var buildSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Searches for builds",
	Long: `Searches for builds.

Builds can be searched on a project using the '--project' flag:

    ontrack-cli build search --project PROJECT

or on a branch using the '--branch' flag:

    ontrack-cli build search --project PROJECT --branch BRANCH

In both cases, several criteria are available - see 'ontrack-cli build search --help' to get their list. For example,
to look for a build using its commit:

    ontrack-cli build search --project PROJECT --branch BRANCH --commit commit

By default, only the build names are printed, one per line.

You can change the display options using additional flags - see 'ontrack-cli build search --help' to get their list.`,
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

type buildList struct {
	Builds []struct {
		Name   string
		Branch struct {
			Name string
		}
	}
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
				name
				branch {
					name
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
				name
				branch {
					name
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

	for _, build := range data.Builds {
		if displayBranch {
			fmt.Printf("%s/", build.Branch.Name)
		}
		fmt.Println(build.Name)
	}

	return nil
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildSearchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildSearchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	buildSearchCmd.Flags().StringP("project", "p", "", "Name of the project")
	buildSearchCmd.MarkFlagRequired("project")
	buildSearchCmd.Flags().StringP("branch", "b", "", "Name of the branch")

	// Criteria
	buildSearchCmd.Flags().Int("count", 10, "Number of builds to return")
	buildSearchCmd.Flags().String("with-promotion", "", "Builds must have this promotion")

	// Property criteria
	buildSearchCmd.Flags().String("commit", "", "Commit for the build")

	// Display criteria
	buildSearchCmd.Flags().Bool("display-branch", false, "Displays branch information: <branch name>/<build name>. Used only for project-based searches.")
}
