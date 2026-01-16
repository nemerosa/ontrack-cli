package cmd

import (
	"fmt"
	"os"
	"strings"
	"yontrack/client"
	"yontrack/config"

	"github.com/spf13/cobra"
)

var buildChangelogSemanticCmd = &cobra.Command{
	Use:     "semantic",
	Aliases: []string{},
	Short:   "Semantic changelog between two builds",
	Long: `Semantic changelog between two builds.

For example:

    yontrack build changelog semantic --from-promotion RELEASE --issues true
	`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Boundaries

		from, err := cmd.Flags().GetInt("from")
		if err != nil {
			return err
		}

		fromPromotion, err := cmd.Flags().GetString("from-promotion")
		if err != nil {
			return err
		}

		to, err := cmd.Flags().GetInt("to")
		if err != nil {
			return err
		}

		// Format options

		issues, err := cmd.Flags().GetBool("issues")
		if err != nil {
			return err
		}

		emojis, err := cmd.Flags().GetBool("emojis")
		if err != nil {
			return err
		}

		renderer, err := cmd.Flags().GetString("renderer")
		if err != nil {
			return err
		}

		sections, err := cmd.Flags().GetStringSlice("sections")
		if err != nil {
			return err
		}

		exclude, err := cmd.Flags().GetStringSlice("exclude")
		if err != nil {
			return err
		}

		// Computing the boundaries (to)
		var toId int
		if to != 0 {
			toId = to
		} else {
			toId, err = getToBoundary()
			if err != nil {
				return err
			}
		}

		// Computing the boundaries (from)
		fromId, err := getFromBoundary(toId, from, fromPromotion)
		if err != nil {
			return err
		}
		if fromId == 0 {
			return nil
		}

		// Query

		query := `
			query SemanticChangeLog(
				$from: Int!,
				$to: Int!,
				$renderer: String,
				$issues: Boolean,
				$emojis: Boolean,
				$sections: [String!],
				$exclude: [String!],
			) {
				scmChangeLog(from: $from, to: $to) {
					semantic(
						renderer: $renderer,
						config: {
							issues: $issues,
							emojis: $emojis,
							sections: $sections,
							exclude: $exclude,
						}
					)
				}
			}
		`

		// Data

		var data struct {
			ScmChangeLog struct {
				Semantic string
			}
		}

		// Getting the configuration

		cfg, err := config.GetSelectedConfiguration()
		if err != nil {
			return err
		}

		// Call

		if err := client.GraphQLCall(cfg, query, map[string]interface{}{
			"from":     fromId,
			"to":       toId,
			"renderer": renderer,
			"issues":   issues,
			"emojis":   emojis,
			"sections": sections,
			"exclude":  exclude,
		}, &data); err != nil {
			return err
		}

		// Print the exported change log
		_, _ = fmt.Fprintf(os.Stdout, strings.TrimSpace(data.ScmChangeLog.Semantic))

		// OK
		return nil
	},
}

func init() {
	buildChangelogCmd.AddCommand(buildChangelogSemanticCmd)

	buildChangelogSemanticCmd.Flags().StringP("renderer", "r", "", "Format renderer for the changelog: text (default), markdown or html")
	buildChangelogSemanticCmd.Flags().BoolP("issues", "i", false, "True to include a changelog of issues")
	buildChangelogSemanticCmd.Flags().BoolP("emojis", "e", false, "True to prefix section titles with emojis")
	buildChangelogSemanticCmd.Flags().StringSliceP("sections", "s", []string{}, "Sections to include in the changelog (type=title)")
	buildChangelogSemanticCmd.Flags().StringSliceP("exclude", "X", []string{}, "Types of commit to exclude from the changelog")
}
