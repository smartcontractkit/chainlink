package ethkey

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

func TestEthKeyV2_ToKey(t *testing.T) {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	require.NoError(t, err)

	k := KeyFor(internal.NewRaw(privateKeyECDSA.D.Bytes()))

	assert.Equal(t, k.getPK(), privateKeyECDSA)
	assert.Equal(t, k.getPK().PublicKey.X, privateKeyECDSA.PublicKey.X)
	assert.Equal(t, k.getPK().PublicKey.Y, privateKeyECDSA.PublicKey.Y)
	assert.Equal(t, types.EIP55AddressFromAddress(crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)).Hex(), k.ID())
}

func TestEthKeyV2_NewV2(t *testing.T) {
	keyV2, err := NewV2()
	require.NoError(t, err)

	assert.NotZero(t, keyV2.Address)
	assert.NotNil(t, keyV2.getPK())
	assert.Equal(t, keyV2.Address.Hex(), keyV2.ID())
}
