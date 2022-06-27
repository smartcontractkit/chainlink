package p2pkey

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestP2PKeys_PeerID(t *testing.T) {
	t.Run("make peer ID", func(t *testing.T) {
		id, err := MakePeerID("11")
		require.NoError(t, err)
		_, err = MakePeerID("invalid")
		assert.Error(t, err)

		assert.Equal(t, "p2p_11", id.String())
	})

	t.Run("unmarshals new ID", func(t *testing.T) {
		id, err := MakePeerID("11")
		require.NoError(t, err)
		fakeKey := MustNewV2XXXTestingOnly(big.NewInt(1))

		err = id.UnmarshalString(fakeKey.ID())
		require.NoError(t, err)

		assert.Equal(t, "p2p_"+fakeKey.ID(), id.String())
	})

	t.Run("scans new ID", func(t *testing.T) {
		id, err := MakePeerID("11")
		require.NoError(t, err)
		fakeKey := MustNewV2XXXTestingOnly(big.NewInt(1))

		err = id.Scan(fakeKey.ID())
		require.NoError(t, err)

		assert.Equal(t, "p2p_"+fakeKey.ID(), id.String())

		err = id.Scan(12)
		assert.Error(t, err)
		assert.Equal(t, "", id.String())
	})
}

func TestPeerID_marshal(t *testing.T) {
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
