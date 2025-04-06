package datastore

import (
	"testing"

	require "github.com/stretchr/testify/require"
)

func TestDomainMetadata_Clone(t *testing.T) {
	original := DomainMetadata[DefaultMetadata]{
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

func TestDomainMetadata_Key(t *testing.T) {
	domainMetadata := DomainMetadata[DefaultMetadata]{
		Domain:      "example.com",
		Environment: "production",
		Metadata:    DefaultMetadata{Data: "test data"},
	}

	key := domainMetadata.Key()
	expectedKey := NewDomainMetadataKey("example.com", "production")

	require.Equal(t, expectedKey, key)
}
