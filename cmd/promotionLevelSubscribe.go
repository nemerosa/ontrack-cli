package cmd

import (
	"github.com/spf13/cobra"
)

var promotionLevelSubscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe to events on a promotion level",
	Long: `Subscribe to events on a promotion level.
	
You can use subcommands to subscribe to events on a promotion level. For example:

	yontrack promotion-level subscribe --project PROJECT --branch BRANCH --promotion LEVEL \ 
        --name "My subscription" \
        slack --channel "#my-channel" --type SUCCESS

Generic notifications can also be used:

	yontrack promotion-level subscribe --project PROJECT --branch BRANCH --promotion LEVEL \ 
        --name "My subscription" \
        generic --channel "slack" --channel-config '{"channel":"#my-channel", "type": "SUCCESS"}'
	
By default, the subscription listens to "new_promotion_run" events.
`,
}

func init() {
	promotionLevelCmd.AddCommand(promotionLevelSubscribeCmd)

	promotionLevelSubscribeCmd.PersistentFlags().StringP("promotion", "l", "", "Name of the promotion level")
	promotionLevelSubscribeCmd.PersistentFlags().StringP("name", "n", "", "Name of the subscription")
	promotionLevelSubscribeCmd.PersistentFlags().StringP("template", "", "", "Custom template for the notification")

	err := promotionLevelSubscribeCmd.MarkPersistentFlagRequired("promotion")
	if err != nil {
		return
	}
	err = promotionLevelSubscribeCmd.MarkPersistentFlagRequired("name")
	if err != nil {
		return
	}
}
