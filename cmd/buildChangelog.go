package cmd

import (
	"github.com/spf13/cobra"
)

// buildChangelogCmd represents the buildChangelog command
var buildChangelogCmd = &cobra.Command{
	Use:     "changelog",
	Aliases: []string{"log"},
	Short:   "Computes a change log between two builds",
	Long: `Computes a change log between two builds.

See the subcommands.

For example:

    yontrack build changelog export --from 1 --to 2`,
	// Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	buildCmd.AddCommand(buildChangelogCmd)

	buildChangelogCmd.PersistentFlags().Int("from", 0, "Build from")
	buildChangelogCmd.PersistentFlags().String("from-promotion", "", "Build from a given promotion")
	buildChangelogCmd.PersistentFlags().Int("to", 0, "Build to")

	_ = buildChangelogCmd.MarkPersistentFlagRequired("to")
}
