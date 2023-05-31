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

func TestOffchainKeyring_NaclBoxSealAnonymous(t *testing.T) {
	kr, err := newOffchainKeyring(cryptorand.Reader, cryptorand.Reader)
	require.NoError(t, err)

	originalMessage := []byte("test")

	encryptedMessage := naclBoxSealAnonymous(t, kr.ConfigEncryptionPublicKey(), originalMessage)

	decryptedMessage, err := kr.NaclBoxOpenAnonymous(encryptedMessage)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, originalMessage, decryptedMessage)
}

func TestOffchainKeyring_NaclBoxSealAnonymous_ShortCiphertext(t *testing.T) {
	kr, err := newOffchainKeyring(cryptorand.Reader, cryptorand.Reader)
	require.NoError(t, err)

	shortMessage := []byte("short")

	_, err = kr.NaclBoxOpenAnonymous(shortMessage)
	assert.Equal(t, err.Error(), "ciphertext too short")
}

func TestOffchainKeyring_NaclBoxSealAnonymous_FailedDecryption(t *testing.T) {
	kr, err := newOffchainKeyring(cryptorand.Reader, cryptorand.Reader)
	require.NoError(t, err)

	invalid := []byte("invalidEncryptedMessage")

	_, err = kr.NaclBoxOpenAnonymous(invalid)
	assert.Equal(t, err.Error(), "decryption failed")
}

func naclBoxSealAnonymous(t *testing.T, peerPublicKey [curve25519.PointSize]byte, plaintext []byte) []byte {
	t.Helper()

	ciphertext, err := box.SealAnonymous(nil, plaintext, &peerPublicKey, cryptorand.Reader)
	if err != nil {
		t.Fatalf("encryption failed")
		return nil
	}

	return ciphertext
}
