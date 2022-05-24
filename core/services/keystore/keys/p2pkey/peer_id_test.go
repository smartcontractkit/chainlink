package p2pkey

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPeerID(t *testing.T) {
	id, err := MakePeerID("p2p_12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw")
	require.NoError(t, err)
	t.Run("json", func(t *testing.T) {
		b, err := id.MarshalJSON()
		require.NoError(t, err)
		var got PeerID
		require.NoError(t, got.UnmarshalJSON(b))
		require.Equal(t, id, got)
	})
	t.Run("db", func(t *testing.T) {
		v, err := id.Value()
		require.NoError(t, err)
		var got PeerID
		require.NoError(t, got.Scan(v))
		require.Equal(t, id, got)
	})
	t.Run("text", func(t *testing.T) {
		s, err := id.MarshalText()
		require.NoError(t, err)
		var got PeerID
		require.NoError(t, got.UnmarshalText(s))
		require.Equal(t, id, got)
	})
}
