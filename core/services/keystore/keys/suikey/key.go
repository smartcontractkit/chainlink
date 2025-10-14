package suikey

import (
	"crypto"
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/blake2b"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

// Ed25519Scheme Ed25519 signature scheme flag
// https://docs.sui.io/concepts/cryptography/transaction-auth/keys-addresses#address-format
const Ed25519Scheme byte = 0x00

// Key represents a Sui account
type Key struct {
	raw    internal.Raw
	signFn func(io.Reader, []byte, crypto.SignerOpts) ([]byte, error)
	pubKey ed25519.PublicKey
}

// KeyFor creates an Account from a raw key
func KeyFor(raw internal.Raw) Key {
	privKey := ed25519.NewKeyFromSeed(internal.Bytes(raw))
	pubKey := privKey.Public().(ed25519.PublicKey)
	return Key{
		raw:    raw,
		signFn: privKey.Sign,
		pubKey: pubKey,
	}
}

// New creates new Key
func New() (Key, error) {
	return newFrom(cryptorand.Reader)
}

// MustNewInsecure returns an Account if no error
func MustNewInsecure(reader io.Reader) Key {
	key, err := newFrom(reader)
	if err != nil {
		panic(err)
	}
	return key
}

// newFrom creates a new Account from a provided random reader
func newFrom(reader io.Reader) (Key, error) {
	pub, priv, err := ed25519.GenerateKey(reader)
	if err != nil {
		return Key{}, err
	}
	return Key{
		raw:    internal.NewRaw(priv.Seed()),
		signFn: priv.Sign,
		pubKey: pub,
	}, nil
}

// ID gets Account ID
func (s Key) ID() string {
	return s.PublicKeyStr()
}

// Address returns the Sui address
// https://docs.sui.io/concepts/cryptography/transaction-auth/keys-addresses#address-format
func (s Key) Address() string {
	// Prepend the Ed25519 scheme flag to the public key
	flaggedPubKey := make([]byte, 1+len(s.pubKey))
	flaggedPubKey[0] = Ed25519Scheme
	copy(flaggedPubKey[1:], s.pubKey)

	// Hash the flagged public key with Blake2b-256
	addressBytes := blake2b.Sum256(flaggedPubKey)

	// Return the full 32-byte address with 0x prefix
	return hex.EncodeToString(addressBytes[:])
}

// Account is an alias for Address
func (s Key) Account() string {
	return s.Address()
}

// GetPublic gets Account's public key
func (s Key) GetPublic() ed25519.PublicKey {
	return s.pubKey
}

// PublicKeyStr returns hex encoded public key
func (s Key) PublicKeyStr() string {
	return fmt.Sprintf("%064x", s.pubKey)
}

// Raw returns the seed from private key
func (s Key) Raw() internal.Raw { return s.raw }

// Sign is used to sign a message
func (s Key) Sign(msg []byte) ([]byte, error) {
	var noHash crypto.Hash
	return s.signFn(cryptorand.Reader, msg, noHash) // no specific hash function used
}
