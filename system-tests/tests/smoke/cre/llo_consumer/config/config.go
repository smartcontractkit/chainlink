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
}

// DefaultLLOConsumerConfig returns the default configuration for E2E testing
func DefaultLLOConsumerConfig() LLOConsumerConfig {
	return LLOConsumerConfig{
		StreamIDs:      []uint32{1, 4}, // Both Format 5 and Format 7 streams
		MaxFrequencyMs: 1000,           // 1 second
	}
}
