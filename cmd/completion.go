package cmd

import (
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [SHELL]",
	Short: "Prints shell completion scripts",
	Long: `Prints shell completion scripts. To install:

BASH
yontrack completion bash > ontrack-completion.bash
sudo cp ontrack-completion.bash /etc/bash_completion.d/
source ~/.bashrc

ZSH
yontrack completion zsh > _yontrack
sudo mkdir -p /usr/local/share/zsh/site-functions
sudo cp _yontrack /usr/local/share/zsh/site-functions/
source ~/.zshrc

FISH
yontrack completion fish > ~/.config/fish/completions/yontrack.fish
source ~/.config/fish/config.fish
`,
	ValidArgs: []string{"bash", "zsh", "fish"},
	Args:      cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			_ = cmd.Root().GenBashCompletion(cmd.OutOrStdout())
		case "zsh":
			_ = cmd.Root().GenZshCompletion(cmd.OutOrStdout())
		case "fish":
			_ = cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
