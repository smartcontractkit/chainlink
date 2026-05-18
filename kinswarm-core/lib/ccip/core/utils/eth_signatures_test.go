package utils

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSignersEthAddress_Success(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	msg := []byte("test message")
	sig, err := GenerateEthSignature(privateKey, msg)
	assert.NoError(t, err)

	recoveredAddress, err := GetSignersEthAddress(msg, sig)
	assert.NoError(t, err)
	assert.Equal(t, address, recoveredAddress)
}

func TestGetSignersEthAddress_InvalidSignatureLength(t *testing.T) {
	msg := []byte("test message")
	sig := []byte("invalid signature length")
	_, err := GetSignersEthAddress(msg, sig)
	assert.EqualError(t, err, "invalid signature: signature length must be 65 bytes")
}

func TestGenerateEthPrefixedMsgHash(t *testing.T) {
	msg := []byte("test message")
	expectedPrefix := "\x19Ethereum Signed Message:\n"
	expectedHash := crypto.Keccak256Hash([]byte(expectedPrefix + "12" + string(msg)))

	hash := GenerateEthPrefixedMsgHash(msg)
	assert.Equal(t, expectedHash, hash)
}

func TestGenerateEthSignature(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)

	msg := []byte("test message")
	signature, err := GenerateEthSignature(privateKey, msg)
	assert.NoError(t, err)
	assert.Len(t, signature, 65)

	recoveredPub, err := crypto.SigToPub(GenerateEthPrefixedMsgHash(msg).Bytes(), signature)
	assert.NoError(t, err)
	assert.Equal(t, privateKey.PublicKey, *recoveredPub)
}
