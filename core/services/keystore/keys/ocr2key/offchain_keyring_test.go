package ocr2key

import (
	"bytes"
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/nacl/box"
)

func TestOffchainKeyring(t *testing.T) {
	kr, err := newOffchainKeyring(cryptorand.Reader, cryptorand.Reader)
	require.NoError(t, err)
	pubKey := kr.OffchainPublicKey()
	assert.True(t, bytes.Equal(kr.signingKey.Public().(ed25519.PublicKey), pubKey[:]))
}

func TestOffchainKeyring_Decrypt(t *testing.T) {
	kr, err := newOffchainKeyring(cryptorand.Reader, cryptorand.Reader)
	require.NoError(t, err)

	originalMessage := []byte("test")

	encryptedMessage, err := Encrypt(kr.ConfigEncryptionPublicKey(), originalMessage)
	if err != nil {
		t.Fatal(err)
	}

	decryptedMessage, err := kr.Decrypt(encryptedMessage)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, originalMessage, decryptedMessage)
}

func Encrypt(peerPublicKey [curve25519.PointSize]byte, plaintext []byte) (ciphertext []byte, err error) {
	ciphertext, err = box.SealAnonymous(nil, plaintext, &peerPublicKey, cryptorand.Reader)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}
