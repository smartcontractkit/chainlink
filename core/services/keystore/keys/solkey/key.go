package solkey

import (
	"crypto"
	"crypto/ed25519"
	crypto_rand "crypto/rand"
	"io"

	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

func KeyFor(raw internal.Raw) Key {
	privKey := ed25519.NewKeyFromSeed(internal.Bytes(raw))
	return Key{
		raw:    raw,
		signFn: privKey.Sign,
		pubKey: privKey.Public().(ed25519.PublicKey),
	}
}

// Key represents Solana key
type Key struct {
	raw    internal.Raw
	signFn func(io.Reader, []byte, crypto.SignerOpts) ([]byte, error)
	pubKey ed25519.PublicKey
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

// GetPublic get Key's public key
func (key Key) GetPublic() ed25519.PublicKey {
	return key.pubKey
}

// PublicKeyStr return base58 encoded public key
func (key Key) PublicKeyStr() string {
	return base58.Encode(key.pubKey)
}

// Raw from private key
func (key Key) Raw() internal.Raw { return key.raw }

// Sign is used to sign a message
func (key Key) Sign(msg []byte) ([]byte, error) {
	return key.signFn(crypto_rand.Reader, msg, crypto.Hash(0))
}

// PublicKey copies public key slice
func (key Key) PublicKey() (pubKey solana.PublicKey) {
	copy(pubKey[:], key.pubKey)
	return
}
