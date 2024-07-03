package solkey

import (
	"crypto"
	"crypto/ed25519"
	crypto_rand "crypto/rand"
	"fmt"
	"io"

	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)

// Raw represents the Solana private key
type Raw []byte

// Key gets the Key
func (raw Raw) Key() Key {
	privKey := ed25519.NewKeyFromSeed(raw)
	pubKey := make([]byte, ed25519.PublicKeySize)
	copy(pubKey, privKey[ed25519.PublicKeySize:])
	return Key{
		privkey: privKey,
		pubKey:  pubKey,
	}
}

// String returns description
func (raw Raw) String() string {
	return "<Solana Raw Private Key>"
}

// GoString wraps String()
func (raw Raw) GoString() string {
	return raw.String()
}

var _ fmt.GoStringer = &Key{}

// Key represents Solana key
type Key struct {
	privkey ed25519.PrivateKey
	pubKey  ed25519.PublicKey
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
		privkey: priv,
		pubKey:  pub,
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
func (key Key) Raw() Raw {
	return key.privkey.Seed()
}

// String is the print-friendly format of the Key
func (key Key) String() string {
	return fmt.Sprintf("SolanaKey{PrivateKey: <redacted>, Public Key: %s}", key.PublicKeyStr())
}

// GoString wraps String()
func (key Key) GoString() string {
	return key.String()
}

// Sign is used to sign a message
func (key Key) Sign(msg []byte) ([]byte, error) {
	return key.privkey.Sign(crypto_rand.Reader, msg, crypto.Hash(0))
}

// PublicKey copies public key slice
func (key Key) PublicKey() (pubKey solana.PublicKey) {
	copy(pubKey[:], key.pubKey)
	return
}
