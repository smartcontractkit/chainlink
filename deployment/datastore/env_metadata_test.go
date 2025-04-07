package datastore

import (
	"testing"

	require "github.com/stretchr/testify/require"
)

func TestEnvMetadata_Clone(t *testing.T) {
	original := EnvMetadata[DefaultMetadata]{
		Domain:      "example.com",
		Environment: "production",
		Metadata:    DefaultMetadata{Data: "test-value"},
	}

	cloned := original.Clone()

	require.Equal(t, original.Domain, cloned.Domain)
	require.Equal(t, original.Environment, cloned.Environment)
	require.Equal(t, original.Metadata, cloned.Metadata)
	require.NotSame(t, &original.Metadata, &cloned.Metadata) // Ensure Metadata is a deep copy
}

func TestEnvMetadata_Key(t *testing.T) {
	envMetadata := EnvMetadata[DefaultMetadata]{
		Domain:      "example.com",
		Environment: "production",
		Metadata:    DefaultMetadata{Data: "test data"},
	}

	key := envMetadata.Key()
	expectedKey := NewEnvMetadataKey("example.com", "production")

	require.Equal(t, expectedKey, key)
}
