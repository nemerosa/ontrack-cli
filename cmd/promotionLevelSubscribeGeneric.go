package cmd

import (
	"github.com/spf13/cobra"
	"yontrack/client"
	"yontrack/config"
)

var promotionLevelSubscribeGenericCmd = &cobra.Command{
	Use:   "generic",
	Short: "Subscribe to events on a promotion level using a generic notification",
	Long: `Subscribe to events on a promotion level using a generic notification.
	
You can use subcommands to subscribe to events on a promotion level. For example:

	yontrack promotion-level subscribe --project PROJECT --branch BRANCH --promotion LEVEL \ 
        --name "My subscription" \
        generic --channel "slack" --channel-config '{"channel":"#my-channel", "type": "SUCCESS"}'
	
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

		channel_config, err := cmd.Flags().GetString("channel-config")
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
			channel,
			channel_config,
			template,
		)
	},
}

func init() {
	promotionLevelSubscribeCmd.AddCommand(promotionLevelSubscribeGenericCmd)

	promotionLevelSubscribeGenericCmd.Flags().StringP("channel", "c", "", "Name of notification channel: mail, slack, etc.")
	promotionLevelSubscribeGenericCmd.Flags().StringP("channel-config", "v", "", "JSON configuration for the notification channel. For example: {\"channel\":\"#my-channel\", \"type\": \"SUCCESS\"} for Slack.")

	err := promotionLevelSubscribeGenericCmd.MarkFlagRequired("channel")
	if err != nil {
		return
	}
	err = promotionLevelSubscribeGenericCmd.MarkFlagRequired("channel-config")
	if err != nil {
		return
	}
}
