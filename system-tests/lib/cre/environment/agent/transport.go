package agent

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"
)

// EncodeForTransport sanitizes arbitrary structs for JSON transport by round-tripping through TOML.
// This drops fields intentionally excluded from TOML (for example runtime handles with toml:"-").
func EncodeForTransport(v any) (map[string]any, error) {
	b, err := toml.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transport payload to TOML: %w", err)
	}

	var payload map[string]any
	if err := toml.Unmarshal(b, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transport payload from TOML: %w", err)
	}

	return payload, nil
}

// DecodeFromTransport decodes sanitized transport payload into a target type using TOML round-trip.
func DecodeFromTransport[T any](payload map[string]any) (*T, error) {
	b, err := toml.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transport payload to TOML: %w", err)
	}

	var out T
	if err := toml.Unmarshal(b, &out); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transport payload into target: %w", err)
	}

	return &out, nil
}
