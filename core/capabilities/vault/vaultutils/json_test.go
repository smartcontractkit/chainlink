package vaultutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestToCanonicalJSON(t *testing.T) {
	testData, err := structpb.NewValue(map[string]any{
		"field1": "value1",
		"field2": 42,
	})
	require.NoError(t, err)

	canonicalJSON, err := ToCanonicalJSON(testData)
	require.NoError(t, err)

	expectedJSON := `{"field1":"value1","field2":42}`
	assert.Equal(t, expectedJSON, string(canonicalJSON)) //nolint:testifylint // testifylint requires use of assert.JSONEq which doesn't do a string match of the JSON, but does a structural match.

	// same data, different order
	testData, err = structpb.NewValue(map[string]any{
		"field2": 42,
		"field1": "value1",
	})
	require.NoError(t, err)

	canonicalJSON, err = ToCanonicalJSON(testData)
	require.NoError(t, err)

	assert.Equal(t, expectedJSON, string(canonicalJSON)) //nolint:testifylint // testifylint requires use of assert.JSONEq which doesn't do a string match of the JSON, but does a structural match.
}
