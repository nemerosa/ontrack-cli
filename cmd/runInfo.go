package cmd

import (
	"github.com/spf13/cobra"
)

func InitCommandFlags(cmd *cobra.Command) {
	cmd.Flags().String("source-type", "", "Run info source type")
	cmd.Flags().String("source-uri", "", "Run info source URI")
	cmd.Flags().String("trigger-type", "", "Run info trigger type")
	cmd.Flags().String("trigger-data", "", "Run info trigger data")
	cmd.Flags().Int("run-time", 0, "Run info run time (in seconds)")
}

func GetRunInfo(cmd *cobra.Command) (*RunInfo, error) {
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
		var info = RunInfo{
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

// RunInfo defines the input for the run info for a build or validation run
type RunInfo struct {
	SourceType  string `json:"sourceType"`
	SourceURI   string `json:"sourceUri"`
	TriggerType string `json:"triggerType"`
	TriggerData string `json:"triggerData"`
	RunTime     int    `json:"runTime"`
}
