package csakey

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestExport_ToEncryptedJSON(t *testing.T) {
	key, err := NewV2()
	require.NoError(t, err)

	json, err := key.ToEncryptedJSON("password", utils.DefaultScryptParams)
	require.NoError(t, err)

	assert.NotEmpty(t, json)
}

func TestExport_FromEncryptedJSON(t *testing.T) {
	key, err := NewV2()
	require.NoError(t, err)

	json, err := key.ToEncryptedJSON("password", utils.DefaultScryptParams)
	require.NoError(t, err)
	require.NotEmpty(t, json)

	kV2, err := FromEncryptedJSON(json, "password")
	require.NoError(t, err)

	require.Equal(t, key.String(), kV2.String())
}
