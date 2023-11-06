package cmd

import (
	client "ontrack-cli/client"

	"github.com/spf13/cobra"
)

func InitRunInfoCommandFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("source-type", "", "Run info source type")
	cmd.PersistentFlags().String("source-uri", "", "Run info source URI")
	cmd.PersistentFlags().String("trigger-type", "", "Run info trigger type")
	cmd.PersistentFlags().String("trigger-data", "", "Run info trigger data")
	cmd.PersistentFlags().Int("run-time", 0, "Run info run time (in seconds)")
}

func GetRunInfo(cmd *cobra.Command) (*client.RunInfo, error) {
	sourceType, err := cmd.Flags().GetString("source-type")
	if err != nil {
		return nil, err
	}
	sourceURI, err := cmd.Flags().GetString("source-uri")
	if err != nil {
		return nil, err
	}
	triggerType, err := cmd.Flags().GetString("trigger-type")
	if err != nil {
		return nil, err
	}
	triggerData, err := cmd.Flags().GetString("trigger-data")
	if err != nil {
		return nil, err
	}
	runTime, err := cmd.Flags().GetInt("run-time")
	if err != nil {
		return nil, err
	}

	if sourceType != "" || sourceURI != "" || triggerType != "" || triggerData != "" || runTime != 0 {
		var info = client.RunInfo{
			SourceType:  sourceType,
			SourceURI:   sourceURI,
			TriggerType: triggerType,
			TriggerData: triggerData,
			RunTime:     runTime,
		}
		return &info, nil
	} else {
		return nil, nil
	}
}
