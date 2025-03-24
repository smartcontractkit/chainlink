package streams

type LLOStreamID uint32

type LLOTriggerConfig struct {
	// The IDs of the data feeds (LLO streams) that will be included in the trigger event.
	StreamIDs []LLOStreamID `json:"streamIds" yaml:"streamIds" mapstructure:"streamIds"`

	// The interval in seconds after which a new trigger event is generated.
	MaxFrequencyMs uint64 `json:"maxFrequencyMs" yaml:"maxFrequencyMs" mapstructure:"maxFrequencyMs"`
}
