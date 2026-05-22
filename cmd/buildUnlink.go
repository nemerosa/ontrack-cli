package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"yontrack/client"
	"yontrack/config"
)

var buildUnlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "Removes links from a build",
	Long: `Removes dependency links from a source build.

The source build can be identified by one of:
  - project name + build name:    --from-project PROJECT --from-build BUILD
  - project name + version label: --from-project PROJECT --from-version VERSION
  - build ID:                     --from-id ID

--from-project defaults to YONTRACK_PROJECT_NAME and --from-build defaults to YONTRACK_BUILD_NAME.

Deletion scope:
  - Omit --to-project to remove all links from the source build.
  - Provide --to-project to remove only links targeting that project.
  - Add --qualifier (-q) to further restrict to a specific qualifier (requires --to-project).

Examples:

  yontrack build unlink --from-project proj-a --from-build 42
  yontrack build unlink --from-project proj-a --from-build 42 --to-project proj-b
  yontrack build unlink --from-project proj-a --from-build 42 --to-project proj-b --qualifier integration
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return buildUnlink(cmd)
	},
}

func buildUnlink(cmd *cobra.Command) error {
	// Read from-side flags
	fromProject, _ := cmd.Flags().GetString("from-project")
	if fromProject == "" {
		fromProject = os.Getenv("YONTRACK_PROJECT_NAME")
	}
	fromBuild, _ := cmd.Flags().GetString("from-build")
	if fromBuild == "" {
		fromBuild = os.Getenv("YONTRACK_BUILD_NAME")
	}
	fromVersion, err := cmd.Flags().GetString("from-version")
	if err != nil {
		return err
	}
	fromID, err := cmd.Flags().GetInt("from-id")
	if err != nil {
		return err
	}

	// Read scope flags
	toProject, err := cmd.Flags().GetString("to-project")
	if err != nil {
		return err
	}
	qualifier, err := cmd.Flags().GetString("qualifier")
	if err != nil {
		return err
	}

	// Validate from-side: exactly one identification method
	fromMethods := countSet(fromID > 0, fromBuild != "", fromVersion != "")
	if fromMethods == 0 {
		return fmt.Errorf("source build must be identified by one of: --from-id, --from-project + --from-build, or --from-project + --from-version")
	}
	if fromMethods > 1 {
		return fmt.Errorf("source build: only one of --from-id, --from-build, or --from-version may be specified at a time")
	}
	if (fromBuild != "" || fromVersion != "") && fromProject == "" {
		return fmt.Errorf("--from-project (or YONTRACK_PROJECT_NAME) is required when using --from-build or --from-version")
	}

	// Validate scope: qualifier requires to-project
	if qualifier != "" && toProject == "" {
		return fmt.Errorf("--qualifier requires --to-project to be specified")
	}

	cfg, err := config.GetSelectedConfiguration()
	if err != nil {
		return err
	}

	// Resolve from-side to ID
	resolvedFromID := fromID
	if resolvedFromID == 0 {
		resolvedFromID, err = resolveBuildID(cfg, fromProject, fromBuild, fromVersion)
		if err != nil {
			return fmt.Errorf("source build: %w", err)
		}
	}

	// Build optional variables (null means "match any")
	var projectVar interface{}
	if toProject != "" {
		projectVar = toProject
	}
	var qualifierVar interface{}
	if qualifier != "" {
		qualifierVar = qualifier
	}

	var data struct {
		DeleteBuildLinks struct {
			Errors []struct {
				Message string
			}
		}
	}
	if err := client.GraphQLCall(cfg, `
		mutation DeleteBuildLinks(
			$fromBuild: Int!,
			$project: String,
			$qualifier: String
		) {
			deleteBuildLinks(input: {
				fromBuild: $fromBuild,
				project: $project,
				qualifier: $qualifier
			}) {
				errors {
					message
				}
			}
		}
	`, map[string]interface{}{
		"fromBuild": resolvedFromID,
		"project":   projectVar,
		"qualifier": qualifierVar,
	}, &data); err != nil {
		return err
	}

	return client.CheckDataErrors(data.DeleteBuildLinks.Errors)
}

func init() {
	buildCmd.AddCommand(buildUnlinkCmd)

	// From-side flags
	buildUnlinkCmd.Flags().String("from-project", "", "Project of the source build (defaults to YONTRACK_PROJECT_NAME)")
	buildUnlinkCmd.Flags().String("from-build", "", "Name of the source build (defaults to YONTRACK_BUILD_NAME)")
	buildUnlinkCmd.Flags().String("from-version", "", "Version (release label) of the source build")
	buildUnlinkCmd.Flags().Int("from-id", 0, "ID of the source build")

	// Scope flags
	buildUnlinkCmd.Flags().String("to-project", "", "Target project name (omit to remove links to all projects)")
	buildUnlinkCmd.Flags().StringP("qualifier", "q", "", "Qualifier to match (requires --to-project)")
}
