package ocrkey

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/offchain-reporting-design/prototype/offchainreporting/to_be_internal/signature"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/curve25519"
)

func TestOCRKeys_Encrypt_Decrypt(t *testing.T) {
	esdcaKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	onChainSigning := signature.OnChainPrivateKey(*esdcaKey)

	_, esdcaPrivKey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)
	offChainSigning := signature.OffChainPrivateKey(esdcaPrivKey)

	var offChainEncryption [curve25519.ScalarSize]byte
	randBytes := make([]byte, curve25519.ScalarSize)
	rand.Read(randBytes)
	copy(offChainEncryption[:], randBytes[:curve25519.ScalarSize])

	pk := OCRPrivateKeys{
		onChainSigning:     &onChainSigning,
		offChainSigning:    &offChainSigning,
		offChainEncryption: &offChainEncryption,
	}
	encryptedPKs, err := pk.Encrypt("password")
	require.NoError(t, err)
	pk2, err := encryptedPKs.Decrypt("password")
	require.NoError(t, err)

	assert.Equal(t, pk.onChainSigning.D, pk2.onChainSigning.D)
	assert.Equal(t, pk.onChainSigning.X, pk2.onChainSigning.X)
	assert.Equal(t, pk.onChainSigning.Y, pk2.onChainSigning.Y)
	assert.Equal(t, pk.offChainSigning.PublicKey(), pk2.offChainSigning.PublicKey())
	assert.Equal(t, pk.offChainEncryption, pk2.offChainEncryption)
}
