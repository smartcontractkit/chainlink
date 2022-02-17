package ocr2key

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExport(t *testing.T) {
	var tt = []struct {
		chain chaintype.ChainType
	}{
		{chain: chaintype.EVM},
		{chain: chaintype.Terra},
		{chain: chaintype.Solana},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(string(tc.chain), func(t *testing.T) {
			kb, err := New(tc.chain)
			require.NoError(t, err)
			ej, err := ToEncryptedJSON(kb, "blah", utils.FastScryptParams)
			require.NoError(t, err)
			kbAfter, err := FromEncryptedJSON(ej, "blah")
			assert.Equal(t, kbAfter.ID(), kb.ID())
			assert.Equal(t, kbAfter.PublicKey(), kb.PublicKey())
			assert.Equal(t, kbAfter.OffchainPublicKey(), kb.OffchainPublicKey())
			assert.Equal(t, kbAfter.MaxSignatureLength(), kb.MaxSignatureLength())
			assert.Equal(t, kbAfter.Raw(), kb.Raw())
			assert.Equal(t, kbAfter.ConfigEncryptionPublicKey(), kb.ConfigEncryptionPublicKey())
			assert.Equal(t, kbAfter.ChainType(), kb.ChainType())
		})
	}
}
