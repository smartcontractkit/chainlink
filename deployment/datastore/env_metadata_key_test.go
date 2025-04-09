package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvMetadataKey(t *testing.T) {
	tests := []struct {
		name           string
		domain         string
		environment    string
		otherDomain    string
		otherEnv       string
		expectedEquals bool
	}{
		{
			name:           "Equal keys",
			domain:         "example.com",
			environment:    "production",
			otherDomain:    "example.com",
			otherEnv:       "production",
			expectedEquals: true,
		},
		{
			name:           "Different environments",
			domain:         "example.com",
			environment:    "production",
			otherDomain:    "example.com",
			otherEnv:       "staging",
			expectedEquals: false,
		},
		{
			name:           "Different domains",
			domain:         "example.com",
			environment:    "production",
			otherDomain:    "another.com",
			otherEnv:       "production",
			expectedEquals: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := NewEnvMetadataKey(tt.domain, tt.environment)
			otherKey := NewEnvMetadataKey(tt.otherDomain, tt.otherEnv)

			assert.Equal(t, tt.domain, key.Domain())
			assert.Equal(t, tt.environment, key.Environment())
			assert.Equal(t, tt.expectedEquals, key.Equals(otherKey))
		})
	}
}
