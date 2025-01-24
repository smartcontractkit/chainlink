package helpers_test

import (
	"reflect"
	"testing"

	"github.com/smartcontractkit/chainlink/deployment/helpers"
)

func TestAddValueToNestedMap(t *testing.T) {
	tests := []struct {
		name     string
		mapping  map[string]map[string]int
		key1     string
		key2     string
		value    int
		expected map[string]map[string]int
	}{
		{
			name:     "Add to empty map",
			mapping:  nil,
			key1:     "group1",
			key2:     "item1",
			value:    42,
			expected: map[string]map[string]int{"group1": {"item1": 42}},
		},
		{
			name: "Add to existing nested map",
			mapping: map[string]map[string]int{
				"group1": {"item1": 10},
			},
			key1:     "group1",
			key2:     "item2",
			value:    20,
			expected: map[string]map[string]int{"group1": {"item1": 10, "item2": 20}},
		},
		{
			name: "Add to a new key in top-level map",
			mapping: map[string]map[string]int{
				"group1": {"item1": 10},
			},
			key1:     "group2",
			key2:     "item1",
			value:    30,
			expected: map[string]map[string]int{"group1": {"item1": 10}, "group2": {"item1": 30}},
		},
		{
			name: "Overwrite existing value in nested map",
			mapping: map[string]map[string]int{
				"group1": {"item1": 10},
			},
			key1:     "group1",
			key2:     "item1",
			value:    50,
			expected: map[string]map[string]int{"group1": {"item1": 50}},
		},
		{
			name:     "Add to nil nested map",
			mapping:  map[string]map[string]int{"group1": nil},
			key1:     "group1",
			key2:     "item1",
			value:    60,
			expected: map[string]map[string]int{"group1": {"item1": 60}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helpers.AddValueToNestedMap(tt.mapping, tt.key1, tt.key2, tt.value)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}
