package streams

type LLOStreamID uint32

type LLOTriggerConfig struct {
	// The IDs of the data feeds (LLO streams) that will be included in the trigger event.
	// The DON filters reports to these streams when report identity is available.
	StreamIDs []LLOStreamID `json:"streamIds" yaml:"streamIds" mapstructure:"streamIds"`

	// The interval in seconds after which a new trigger event is generated.
	MaxFrequencyMs uint64 `json:"maxFrequencyMs" yaml:"maxFrequencyMs" mapstructure:"maxFrequencyMs"`

	// AcceptedReportFormats are report encoding formats the workflow accepts (e.g. 5 = CapabilityTrigger, 7 = EVMABIEncodeUnpackedExpr).
	// The DON only sends reports whose report_format is in this list. Empty means accept all (backward compatible).
	AcceptedReportFormats []uint32 `json:"acceptedReportFormats" yaml:"acceptedReportFormats" mapstructure:"acceptedReportFormats"`

	// TransmissionWindowMs is the transmission window in ms. When > 0, DON delays pushing to this workflow until the next wall-clock boundary.
	// 0 = use capability-level default or send immediately. Defined by the workflow.
	TransmissionWindowMs int `json:"transmissionWindowMs" yaml:"transmissionWindowMs" mapstructure:"transmissionWindowMs"`
}
