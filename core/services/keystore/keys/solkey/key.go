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

type Raw []byte

func (raw Raw) Key() Key {
	privKey := ed25519.NewKeyFromSeed(raw)
	pubKey := make([]byte, ed25519.PublicKeySize)
	copy(pubKey, privKey[ed25519.PublicKeySize:])
	return Key{
		privkey: privKey,
		pubKey:  pubKey,
	}
}

func (raw Raw) String() string {
	return "<Solana Raw Private Key>"
}

func (raw Raw) GoString() string {
	return raw.String()
}

var _ fmt.GoStringer = &Key{}

type Key struct {
	privkey ed25519.PrivateKey
	pubKey  ed25519.PublicKey
}

func New() (Key, error) {
	return newFrom(crypto_rand.Reader)
}

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

func (key Key) ID() string {
	return key.PublicKeyStr()
}

func (key Key) GetPublic() ed25519.PublicKey {
	return key.pubKey
}

func (key Key) PublicKeyStr() string {
	return base58.Encode(key.pubKey)
}

func (key Key) Raw() Raw {
	return key.privkey.Seed()
}

func (key Key) String() string {
	return fmt.Sprintf("SolanaKey{PrivateKey: <redacted>, Public Key: %s}", key.PublicKeyStr())
}

func (key Key) GoString() string {
	return key.String()
}

func (key Key) Sign(msg []byte) ([]byte, error) {
	return key.privkey.Sign(crypto_rand.Reader, msg, crypto.Hash(0))
}

func (key Key) PublicKey() (pubKey solana.PublicKey) {
	copy(pubKey[:], key.pubKey)
	return
}
