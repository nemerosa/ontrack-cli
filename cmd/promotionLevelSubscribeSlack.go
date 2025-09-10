package cmd

import (
	"github.com/spf13/cobra"
	"ontrack-cli/client"
	"ontrack-cli/config"
)

var promotionLevelSubscribeSlackCmd = &cobra.Command{
	Use:   "slack",
	Short: "Subscribe to events on a promotion level using Slack",
	Long: `Subscribe to events on a promotion level using Slack.
	
You can use subcommands to subscribe to events on a promotion level. For example:

	ontrack-cli promotion-level subscribe --project PROJECT --branch BRANCH --promotion LEVEL \ 
        --name "My subscription" \
        slack --channel "#my-channel" --type SUCCESS
	
By default, the subscription listens to "new_promotion_run" events.
`,
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

		promotion, err := cmd.Flags().GetString("promotion")
		if err != nil {
			return err
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		template, err := cmd.Flags().GetString("template")
		if err != nil {
			return err
		}

		channel, err := cmd.Flags().GetString("channel")
		if err != nil {
			return err
		}

		type_name, err := cmd.Flags().GetString("type")
		if err != nil {
			return err
		}

		// Configuration
		cfg, err := config.GetSelectedConfiguration()
		if err != nil {
			return err
		}

		// Subscription configuration
		return client.SubscribePromotionLevel(
			cfg,
			project,
			branch,
			promotion,
			name,
			[]string{"new_promotion_run"},
			"slack",
			map[string]interface{}{
				"channel": channel,
				"type":    type_name,
			},
			template,
		)
	},
}

func init() {
	promotionLevelSubscribeCmd.AddCommand(promotionLevelSubscribeSlackCmd)

	promotionLevelSubscribeSlackCmd.Flags().StringP("channel", "c", "", "Name of Slack channel")
	promotionLevelSubscribeSlackCmd.Flags().StringP("type", "t", "INFO", "Slack message type: INFO, SUCCESS, WARNING or ERROR")

	err := promotionLevelSubscribeSlackCmd.MarkPersistentFlagRequired("channel")
	if err != nil {
		return
	}
}
