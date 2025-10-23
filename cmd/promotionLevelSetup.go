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
	"yontrack/utils"

	"github.com/spf13/cobra"

	client "yontrack/client"
	config "yontrack/config"
)

// promotionLevelSetupCmd represents the promotionLevelSetup command
var promotionLevelSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Creates or updates a promotion level",
	Long: `Creates or updates a promotion level.

	yontrack pl setup -p PROJECT -b BRANCH -l PROMOTION

The promotion can be set to be in "auto promotion" mode by using addtional options. For example:

    yontrack pl setup -p PROJECT -b BRANCH -l PROMOTION \
	    --validation VALIDATION_1 \
	    --validation VALIDATION_2 \
		--depends-on IRON \
		--depends-on SILVER
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, branch, err := utils.GetProjectBranchFlags(cmd, false, true)
		if err != nil {
			return err
		}
		promotion, err := cmd.Flags().GetString("promotion")
		if err != nil {
			return err
		}
		description, err := cmd.Flags().GetString("description")
		if err != nil {
			return err
		}

		// Auto promotion
		validations, err := cmd.Flags().GetStringSlice("validation")
		if err != nil {
			return err
		}
		include, err := cmd.Flags().GetString("include")
		if err != nil {
			return err
		}
		exclude, err := cmd.Flags().GetString("exclude")
		if err != nil {
			return err
		}
		promotions, err := cmd.Flags().GetStringSlice("depends-on")
		if err != nil {
			return err
		}

		// Configuration
		cfg, err := config.GetSelectedConfiguration()
		if err != nil {
			return err
		}

		autoPromotion := len(validations) > 0 || len(promotions) > 0 || include != "" || exclude != ""

		return client.SetupPromotionLevel(
			cfg,
			project,
			branch,
			promotion,
			description,
			autoPromotion,
			validations,
			promotions,
			include,
			exclude,
		)
	},
}

func init() {
	promotionLevelCmd.AddCommand(promotionLevelSetupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// promotionLevelSetupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// promotionLevelSetupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	promotionLevelSetupCmd.Flags().StringP("promotion", "l", "", "Name of the promotion level")
	promotionLevelSetupCmd.Flags().StringP("description", "d", "", "Description of the promotion level")

	promotionLevelSetupCmd.Flags().StringSliceP("validation", "v", []string{}, "Validations the promotion level needs")
	promotionLevelSetupCmd.Flags().StringSliceP("depends-on", "o", []string{}, "Promotions the promotion level needs")
	promotionLevelSetupCmd.Flags().StringP("include", "i", "", "Including validation stamps using a regular expression")
	promotionLevelSetupCmd.Flags().StringP("exclude", "x", "", "Excluding validation stamps using a regular expression")

	promotionLevelSetupCmd.MarkPersistentFlagRequired("promotion")
}
