package aptoskey

import (
	"crypto"
	"crypto/ed25519"
	crypto_rand "crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/sha3"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

// Key represents Aptos key
type Key struct {
	// TODO: store initial Account() derivation to support key rotation
	raw    internal.Raw
	signFn func(io.Reader, []byte, crypto.SignerOpts) ([]byte, error)
	pubKey ed25519.PublicKey
}

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
	return newFrom(crypto_rand.Reader)
}

// MustNewInsecure return Key if no error
func MustNewInsecure(reader io.Reader) Key {
	key, err := newFrom(reader)
	if err != nil {
		panic(err)
	}
	return key
}

// newFrom creates new Key from a provided random reader
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

// ID gets Key ID
func (key Key) ID() string {
	return key.PublicKeyStr()
}

// https://github.com/aptos-foundation/AIPs/blob/main/aips/aip-40.md#long
func (key Key) Account() string {
	authKey := sha3.Sum256(append([]byte(key.pubKey), 0x00))
	return fmt.Sprintf("%064x", authKey)
}

// GetPublic get Key's public key
func (key Key) GetPublic() ed25519.PublicKey {
	return key.pubKey
}

// PublicKeyStr returns hex encoded public key
func (key Key) PublicKeyStr() string {
	return fmt.Sprintf("%064x", key.pubKey)
}

// Raw returns the seed from private key
func (key Key) Raw() internal.Raw { return key.raw }

// Sign is used to sign a message
func (key Key) Sign(msg []byte) ([]byte, error) {
	return key.signFn(crypto_rand.Reader, msg, crypto.Hash(0)) // no specific hash function used
}
