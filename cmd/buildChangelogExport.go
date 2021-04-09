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
	"github.com/spf13/cobra"
)

// buildChangelogExportCmd represents the buildChangelog command
var buildChangelogExportCmd = &cobra.Command{
	Use:     "export",
	Aliases: []string{"format"},
	Short:   "Formats the change log between two builds",
	Long: `Formats the change log between two builds.

For example:

    ontrack-cli build changelog --from 1 --to 2

Additional options are available for the formatting:

	ontrack-cli build changelog --from 1 --to 2 \
		--format markdown \
		--grouping "Bugs=bug,Features=features" \
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
		`

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
