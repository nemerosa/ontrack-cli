package client

// RunInfo defines the input for the run info for a build or validation run
type RunInfo struct {
	SourceType  string `json:"sourceType"`
	SourceURI   string `json:"sourceUri"`
	TriggerType string `json:"triggerType"`
	TriggerData string `json:"triggerData"`
	RunTime     int    `json:"runTime"`
}
