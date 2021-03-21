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

// promotionLevelCmd represents the promotionLevel command
var promotionLevelCmd = &cobra.Command{
	Use:     "promotion-level",
	Aliases: []string{"promotion", "pl"},
	Short:   "Management of promotion levels",
	Long: `Management of promotion levels.

The simplest way to setup a promotion level for a branch:

    ontrack-cli pl setup -p PROJECT -b BRANCH -l PROMOTION

This will create a promotion level for the branch, or updates it if it exists already.
`,
	// Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(promotionLevelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// promotionLevelCmd.PersistentFlags().String("foo", "", "A help for foo")
	promotionLevelCmd.PersistentFlags().StringP("project", "p", "", "Name of the project")
	promotionLevelCmd.PersistentFlags().StringP("branch", "b", "", "Name of the branch")

	promotionLevelCmd.MarkPersistentFlagRequired("project")
	promotionLevelCmd.MarkPersistentFlagRequired("branch")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// promotionLevelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
