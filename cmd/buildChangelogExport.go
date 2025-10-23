package cmd

import (
	"fmt"
	"strings"
	"yontrack/client"
	"yontrack/config"

	"github.com/spf13/cobra"
)

// buildChangelogExportCmd represents the buildChangelog command
var buildChangelogExportCmd = &cobra.Command{
	Use:     "export",
	Aliases: []string{"format"},
	Short:   "Formats the change log between two builds",
	Long: `Formats the change log between two builds.

For example:

    yontrack build changelog --from 1 --to 2

Additional options are available for the formatting:

	yontrack build changelog --from 1 --to 2 \
		--format markdown \
		--grouping "Bugs=bug|Features=features" \
		--alt-group "Misc" \
		--exclude delivery

The change log is available directly in the standard output.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Boundaries

		from, err := cmd.Flags().GetInt("from")
		if err != nil {
			return err
		}
		to, err := cmd.Flags().GetInt("to")
		if err != nil {
			return err
		}

		// Format options

		format, err := cmd.Flags().GetString("format")
		if err != nil {
			return err
		}

		grouping, err := cmd.Flags().GetString("grouping")
		if err != nil {
			return err
		}

		altGroup, err := cmd.Flags().GetString("alt-group")
		if err != nil {
			return err
		}

		exclude, err := cmd.Flags().GetString("exclude")
		if err != nil {
			return err
		}

		// Query

		query := `
			query ExportChangeLog(
				$from: Int!,
				$to: Int!,
				$request: IssueChangeLogExportRequest!
			) {
				gitChangeLog(from: $from, to: $to) {
					export(request: $request)
				}
			}
		`

		// Request

		request := make(map[string]interface{})

		if format != "" {
			request["format"] = format
		}
		if grouping != "" {
			request["grouping"] = grouping
		}
		if altGroup != "" {
			request["altGroup"] = altGroup
		}
		if exclude != "" {
			request["exclude"] = exclude
		}

		// Data

		var data struct {
			GitChangeLog struct {
				Export string
			}
		}

		// Getting the configuration

		cfg, err := config.GetSelectedConfiguration()
		if err != nil {
			return err
		}

		// Call

		if err := client.GraphQLCall(cfg, query, map[string]interface{}{
			"from":    from,
			"to":      to,
			"request": request,
		}, &data); err != nil {
			return err
		}

		// Print the exported change log
		fmt.Println(strings.TrimSpace(data.GitChangeLog.Export))

		// OK
		return nil
	},
}

func init() {
	buildChangelogCmd.AddCommand(buildChangelogExportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildChangelogExportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildChangelogExportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	buildChangelogExportCmd.Flags().String("format", "", "Format of the changelog: text (default), markdown or html")
	buildChangelogExportCmd.Flags().String("grouping", "", "Grouping specification (see Ontrack doc)")
	buildChangelogExportCmd.Flags().String("alt-group", "", "Name of the group for unclassified issues")
	buildChangelogExportCmd.Flags().String("exclude", "", "Comma separated list of issue types to ignore")
}
