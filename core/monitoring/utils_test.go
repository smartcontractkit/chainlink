package monitoring

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
)

func TestKvMapToOtelAttributes(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]string
		expected []attribute.KeyValue
	}{
		{
			name:     "empty map",
			input:    map[string]string{},
			expected: []attribute.KeyValue{},
		},
		{
			name: "single key-value pair",
			input: map[string]string{
				"key1": "value1",
			},
			expected: []attribute.KeyValue{
				attribute.String("key1", "value1"),
			},
		},
		{
			name: "multiple key-value pairs",
			input: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			expected: []attribute.KeyValue{
				attribute.String("key1", "value1"),
				attribute.String("key2", "value2"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := KvMapToOtelAttributes(tt.input)
			assert.ElementsMatch(t, tt.expected, result, "unexpected KeyValue slice")
		})
	}
}
