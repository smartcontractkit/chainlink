package encryptionkey

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/nacl/box"
)

type Key struct {
	privateKey *[curve25519.PointSize]byte
	publicKey  *[curve25519.PointSize]byte
}

func New() (Key, error) {
	publicKey, privKey, err := box.GenerateKey(cryptorand.Reader)
	if err != nil {
		return Key{}, err
	}

	return Key{
		privateKey: privKey,
		publicKey:  publicKey,
	}, nil
}

func (k *Key) PublicKey() [curve25519.PointSize]byte {
	return *k.publicKey
}

func (k *Key) PublicKeyString() string {
	return hex.EncodeToString(k.publicKey[:])
}

// Encrypt encrypts a message using the public key
func (k *Key) Encrypt(plaintext []byte) ([]byte, error) {
	publicKey := k.PublicKey()
	encrypted, err := box.SealAnonymous(nil, plaintext, &publicKey, cryptorand.Reader)
	if err != nil {
		return nil, err
	}

	return encrypted, nil
}

// Decrypt decrypts a message that was encrypted using the private key
func (k *Key) Decrypt(ciphertext []byte) (plaintext []byte, err error) {
	publicKey := k.PublicKey()
	decrypted, success := box.OpenAnonymous(nil, ciphertext, &publicKey, k.privateKey)
	if !success {
		return nil, errors.New("decryption failed")
	}

	return decrypted, nil
}
