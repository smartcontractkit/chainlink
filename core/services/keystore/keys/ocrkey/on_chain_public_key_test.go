package ocrkey

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestOCRKeys_OnChainPublicKey(t *testing.T) {
	t.Parallel()
	tests.BelongsToCISuite(t, "unit")

	pk, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	require.NoError(t, err)

	publicKey := OnChainPublicKey(pk.PublicKey)

	assert.NotEmpty(t, publicKey.Address())
}
