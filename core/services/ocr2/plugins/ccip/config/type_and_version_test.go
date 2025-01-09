package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTypeAndVersion(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedType    string
		expectedVersion string
		expectedError   string
	}{
		{
			name:            "Valid input",
			input:           string(EVM2EVMOnRamp) + " 1.2.0",
			expectedType:    string(EVM2EVMOnRamp),
			expectedVersion: "1.2.0",
		},
		{
			name:            "Empty input",
			input:           "",
			expectedType:    string(Unknown),
			expectedVersion: defaultVersion,
		},
		{
			name:          "Invalid input",
			input:         "InvalidInput",
			expectedError: "invalid type and version InvalidInput",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualType, actualVersion, err := ParseTypeAndVersion(tc.input)

			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedType, actualType)
				assert.Equal(t, tc.expectedVersion, actualVersion)
			}
		})
	}
}
