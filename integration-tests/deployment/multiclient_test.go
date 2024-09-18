package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiClient(t *testing.T) {
	// Expect an error if no RPCs supplied.
	_, err := NewMultiClient([]RPC{})
	require.Error(t, err)

	// Expect defaults to be set if not provided.
	mc, err := NewMultiClient([]RPC{{HTTPURL: "http://localhost:8545"}})
	require.NoError(t, err)
	assert.Equal(t, mc.RetryConfig.Attempts, uint(RPC_DEFAULT_RETRY_ATTEMPTS))
	assert.Equal(t, mc.RetryConfig.Delay, RPC_DEFAULT_RETRY_DELAY)

	// Expect second client to be set as backup.
	mc, err = NewMultiClient([]RPC{
		{HTTPURL: "http://localhost:8545"},
		{HTTPURL: "http://localhost:8546"},
	})
	require.NoError(t, err)
	require.Equal(t, len(mc.Backups), 1)
	assert.Equal(t, mc.Backups[0], "http://localhost:8546")
}
