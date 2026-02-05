package config

// LLOConsumerConfig defines the configuration for the LLO consumer workflow
type LLOConsumerConfig struct {
	// StreamIDs is the list of stream IDs to subscribe to
	// Stream 1: TEST/USD (Format 5)
	// Stream 4: DATA/USD (Format 7)
	StreamIDs []uint32 `yaml:"stream_ids"`

	// MaxFrequencyMs is the maximum frequency for receiving events in milliseconds
	MaxFrequencyMs uint64 `yaml:"max_frequency_ms"`

	// TriggerCapabilityID is the capability ID to subscribe to (default "streams-trigger@2.0.0").
	// Use "mock@1.0.0" when running with mock-only capabilities DON (LLO mock test).
	TriggerCapabilityID string `yaml:"trigger_capability_id"`

	// AcceptedReportFormats are report formats the workflow accepts; the Streams DON filters so only these are sent. e.g. [5, 7]. Empty = accept all.
	AcceptedReportFormats []uint32 `yaml:"accepted_report_formats"`

	// TransmissionWindowMs: when > 0, the DON delays pushing to this workflow until the next wall-clock boundary. 0 = use runner default or immediate. Defined by the workflow.
	TransmissionWindowMs uint64 `yaml:"transmission_window_ms"`

	// ExpectedReportFormat is the report encoding fallback when report does not carry report_format (e.g. mock): 5 = protobuf, 7 = ABI.
	ExpectedReportFormat int `yaml:"expected_report_format"`
}

// DefaultLLOConsumerConfig returns the default configuration for E2E testing
func DefaultLLOConsumerConfig() LLOConsumerConfig {
	return LLOConsumerConfig{
		StreamIDs:      []uint32{1, 4}, // Both Format 5 and Format 7 streams
		MaxFrequencyMs: 1000,           // 1 second
	}
}
