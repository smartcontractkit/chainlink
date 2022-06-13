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
