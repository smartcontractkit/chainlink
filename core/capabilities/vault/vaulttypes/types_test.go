package vaulttypes

import (
	"crypto/ecdsa"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestIsGatewaySecretsMethod(t *testing.T) {
	t.Parallel()

	for _, method := range GatewaySecretsMethods {
		assert.True(t, IsGatewaySecretsMethod(method), method)
	}
	assert.False(t, IsGatewaySecretsMethod(MethodPublicKeyGet))
	assert.False(t, IsGatewaySecretsMethod(MethodSecretsGet))
	assert.False(t, IsGatewaySecretsMethod("vault.unsupported"))
}

func TestValidateNodeSignatures(t *testing.T) {
	t.Parallel()

	digest := "6acc256f1ba5bf28f834040ecd80ff28d316aa1806ef21ace295f45226b0a66d"

	newKeys := func(n int) ([]*ecdsa.PrivateKey, []common.Address) {
		keys := make([]*ecdsa.PrivateKey, n)
		addrs := make([]common.Address, n)
		for i := range keys {
			key, err := crypto.GenerateKey()
			require.NoError(t, err)
			keys[i] = key
			addrs[i] = crypto.PubkeyToAddress(key.PublicKey)
		}
		return keys, addrs
	}

	signDigest := func(t *testing.T, key *ecdsa.PrivateKey) []byte {
		sig, err := utils.GenerateEthSignature(key, []byte(digest))
		require.NoError(t, err)
		return sig
	}

	t.Run("valid F+1 signatures", func(t *testing.T) {
		t.Parallel()
		keys, addrs := newKeys(3)
		sigs := [][]byte{signDigest(t, keys[0]), signDigest(t, keys[1]), signDigest(t, keys[2])}
		assert.NoError(t, ValidateNodeSignatures(digest, sigs, addrs, 2))
		assert.NoError(t, ValidateNodeSignatures(digest, sigs, addrs, 3))
	})

	t.Run("insufficient signatures", func(t *testing.T) {
		t.Parallel()
		keys, addrs := newKeys(3)
		sigs := [][]byte{signDigest(t, keys[0])}
		require.Error(t, ValidateNodeSignatures(digest, sigs, addrs, 2))
	})

	t.Run("signature from non-member ignored", func(t *testing.T) {
		t.Parallel()
		keys, addrs := newKeys(2)
		outsider, _ := newKeys(1)
		sigs := [][]byte{signDigest(t, keys[0]), signDigest(t, outsider[0])}
		require.NoError(t, ValidateNodeSignatures(digest, sigs, addrs, 1))
		require.Error(t, ValidateNodeSignatures(digest, sigs, addrs, 2))
	})

	t.Run("malformed signatures ignored", func(t *testing.T) {
		t.Parallel()
		keys, addrs := newKeys(2)
		sigs := [][]byte{signDigest(t, keys[0]), {0x01, 0x02}}
		require.NoError(t, ValidateNodeSignatures(digest, sigs, addrs, 1))
		require.Error(t, ValidateNodeSignatures(digest, sigs, addrs, 2))
	})

	t.Run("duplicate signatures count once", func(t *testing.T) {
		t.Parallel()
		keys, addrs := newKeys(2)
		sig := signDigest(t, keys[0])
		sigs := [][]byte{sig, sig}
		require.NoError(t, ValidateNodeSignatures(digest, sigs, addrs, 1))
		require.Error(t, ValidateNodeSignatures(digest, sigs, addrs, 2))
	})

	t.Run("signature over wrong digest rejected", func(t *testing.T) {
		t.Parallel()
		keys, addrs := newKeys(2)
		sig, err := utils.GenerateEthSignature(keys[0], []byte("a8e6e0e131f05fc98fedc4826550a13bbdeb93b01bd9b97aa462a2411860e7d9"))
		require.NoError(t, err)
		require.Error(t, ValidateNodeSignatures(digest, [][]byte{sig}, addrs, 1))
	})
}
