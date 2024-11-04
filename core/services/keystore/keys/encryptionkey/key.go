package encryptionkey

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/nacl/box"
)

type Raw []byte

func (raw Raw) Key() Key {
	privKey := [32]byte(raw)
	return Key{
		privateKey: &privKey,
		publicKey:  curve25519PubKeyFromPrivKey(privKey),
	}
}

func (raw Raw) String() string {
	return "<Encryption Raw Private Key>"
}

func (raw Raw) GoString() string {
	return raw.String()
}

func (raw Raw) Bytes() []byte {
	return ([]byte)(raw)
}

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

func (k Key) PublicKey() [curve25519.PointSize]byte {
	return *k.publicKey
}

func (k Key) PublicKeyString() string {
	return hex.EncodeToString(k.publicKey[:])
}

func (k Key) ID() string {
	return k.PublicKeyString()
}

func (k Key) Raw() Raw {
	raw := make([]byte, curve25519.PointSize)
	copy(raw, k.privateKey[:])
	return Raw(raw)
}

func (k Key) String() string {
	return fmt.Sprintf("EncryptionKey{PrivateKey: <redacted>, PublicKey: %s}", *k.publicKey)
}

func (k Key) GoString() string {
	return k.String()
}

// Encrypt encrypts a message using the public key
func (k Key) Encrypt(plaintext []byte) ([]byte, error) {
	publicKey := k.PublicKey()
	encrypted, err := box.SealAnonymous(nil, plaintext, &publicKey, cryptorand.Reader)
	if err != nil {
		return nil, err
	}

	return encrypted, nil
}

// Decrypt decrypts a message that was encrypted using the private key
func (k Key) Decrypt(ciphertext []byte) (plaintext []byte, err error) {
	publicKey := k.PublicKey()
	decrypted, success := box.OpenAnonymous(nil, ciphertext, &publicKey, k.privateKey)
	if !success {
		return nil, errors.New("decryption failed")
	}

	return decrypted, nil
}

func curve25519PubKeyFromPrivKey(privKey [curve25519.PointSize]byte) *[curve25519.PointSize]byte {
	var publicKey [curve25519.PointSize]byte

	// Derive the public key
	curve25519.ScalarBaseMult(&publicKey, &privKey)

	return &publicKey
}
