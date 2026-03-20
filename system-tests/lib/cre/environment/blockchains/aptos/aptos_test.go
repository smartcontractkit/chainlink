package aptos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeNodeURL(t *testing.T) {
	t.Run("adds v1 when path is empty", func(t *testing.T) {
		got, err := NormalizeNodeURL("http://127.0.0.1:8080")
		require.NoError(t, err)
		require.Equal(t, "http://127.0.0.1:8080/v1", got)
	})

	t.Run("preserves v1 path", func(t *testing.T) {
		got, err := NormalizeNodeURL("http://127.0.0.1:8080/v1")
		require.NoError(t, err)
		require.Equal(t, "http://127.0.0.1:8080/v1", got)
	})
}

func TestFaucetURLFromNodeURL(t *testing.T) {
	got, err := FaucetURLFromNodeURL("http://127.0.0.1:8080/v1")
	require.NoError(t, err)
	require.Equal(t, "http://127.0.0.1:8081", got)
}

func TestChainIDUint8(t *testing.T) {
	got, err := ChainIDUint8(4)
	require.NoError(t, err)
	require.Equal(t, uint8(4), got)

	_, err = ChainIDUint8(256)
	require.Error(t, err)
}
