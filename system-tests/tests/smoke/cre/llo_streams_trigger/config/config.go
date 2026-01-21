package config

// LLOConsumerConfig defines the configuration for the LLO consumer workflow
// This config is used by the test to pass configuration to the workflow
type LLOConsumerConfig struct {
	// StreamIDs is the list of stream IDs to subscribe to
	// Stream 1: TEST/USD (Format 5)
	// Stream 4: DATA/USD (Format 7)
	StreamIDs []uint32 `yaml:"stream_ids"`

	// MaxFrequencyMs is the maximum frequency for receiving events in milliseconds
	MaxFrequencyMs uint64 `yaml:"max_frequency_ms"`
}

// DefaultLLOConsumerConfig returns the default configuration for E2E testing
func DefaultLLOConsumerConfig() LLOConsumerConfig {
	return LLOConsumerConfig{
		StreamIDs:      []uint32{1, 4}, // Both Format 5 and Format 7 streams
		MaxFrequencyMs: 1000,           // 1 second
	}
}
