package cmd

import (
	"github.com/spf13/cobra"

	"yontrack/client"
	"yontrack/config"
)

// buildSetupCmd represents the buildSetup command
var buildSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Creates a build if it does not exist",
	Long: `Creates a build if it does not exist.

For example, the following command will create the build:

    yontrack build setup --project my-project --branch release/1.0 --build 1

and the same command run a second time won't do anything.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return buildSetup(cmd)
	},
}

func buildSetup(cmd *cobra.Command) error {
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

	description, err := cmd.Flags().GetString("description")
	if err != nil {
		return err
	}

	release, err := cmd.Flags().GetString("release")
	if err != nil {
		return err
	}
	releaseProperty := release != ""

	commit, err := cmd.Flags().GetString("commit")
	if err != nil {
		return err
	}
	commitProperty := commit != ""

	// Run info
	runInfo, err := GetRunInfo(cmd)
	if err != nil {
		return err
	}

	cfg, err := config.GetSelectedConfiguration()
	if err != nil {
		return err
	}

	// Creates or get the build
	var data struct {
		CreateBuildOrGet struct {
			Errors []struct {
				Message string
			}
		}
		SetBuildReleaseProperty struct {
			Errors []struct {
				Message string
			}
		}
		SetBuildGitCommitProperty struct {
			Errors []struct {
				Message string
			}
		}
	}
	if err := client.GraphQLCall(cfg, `
		mutation BuildSetup(
			$project: String!,
			$branch: String!, 
			$build: String!, 
			$description: String, 
			$runInfo: RunInfoInput,
			$releaseProperty: Boolean!,
			$release: String!,
			$commitProperty: Boolean!,
			$commit: String!
		) {
			createBuildOrGet(input: {
				projectName: $project, 
				branchName: $branch, 
				name: $build, 
				description: $description, 
				runInfo: $runInfo
			}) {
				errors {
				  message
				}
			}
			setBuildReleaseProperty(input: {
				project: $project,
				branch: $branch,
				build: $build,
				release: $release
			}) @include(if: $releaseProperty) {
				errors {
					message
				}
			}
			setBuildGitCommitProperty(input: {
				project: $project,
				branch: $branch,
				build: $build,
				commit: $commit
			}) @include(if: $commitProperty) {
				errors {
					message
				}
			}
		}
	`, map[string]interface{}{
		"project":         project,
		"branch":          branch,
		"build":           build,
		"description":     description,
		"runInfo":         runInfo,
		"releaseProperty": releaseProperty,
		"release":         release,
		"commitProperty":  commitProperty,
		"commit":          commit,
	}, &data); err != nil {
		return err
	}

	// Checks errors for the build
	if err := client.CheckDataErrors(data.CreateBuildOrGet.Errors); err != nil {
		return err
	}
	if err := client.CheckDataErrors(data.SetBuildReleaseProperty.Errors); err != nil {
		return err
	}
	if err := client.CheckDataErrors(data.SetBuildGitCommitProperty.Errors); err != nil {
		return err
	}

	// OK
	return nil
}

func init() {
	buildCmd.AddCommand(buildSetupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildSetupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	buildSetupCmd.Flags().StringP("project", "p", "", "Project name (required)")
	buildSetupCmd.Flags().StringP("branch", "b", "", "Branch name (required)")
	buildSetupCmd.Flags().StringP("build", "n", "", "Build name (required)")
	buildSetupCmd.Flags().StringP("description", "d", "", "Build description")

	// Run info parameters
	InitRunInfoCommandFlags(buildSetupCmd)

	// Release property
	buildSetupCmd.Flags().StringP("release", "r", "", "Build release property")

	// Commit property
	buildSetupCmd.Flags().StringP("commit", "c", "", "Build commit property")

	err := buildSetupCmd.MarkFlagRequired("project")
	if err != nil {
		return
	}
	err = buildSetupCmd.MarkFlagRequired("branch")
	if err != nil {
		return
	}
	err = buildSetupCmd.MarkFlagRequired("build")
	if err != nil {
		return
	}
}
