package terrakey

import (
	cryptorand "crypto/rand"
	"fmt"
	"io"

	cosmosed25519 "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/terra-project/terra.go/msg"
)

type Raw []byte

func (raw Raw) Key() Key {
	return Key{
		PrivKey: cosmosed25519.GenPrivKeyFromSecret(raw),
		secret:  raw,
	}
}

func (raw Raw) String() string {
	return "<Terra Raw Private Key>"
}

func (raw Raw) GoString() string {
	return raw.String()
}

var _ fmt.GoStringer = &Key{}

type Key struct {
	*cosmosed25519.PrivKey
	secret []byte
}

func New() Key {
	return newFrom(cryptorand.Reader)
}

func MustNewInsecure(reader io.Reader) Key {
	return newFrom(reader)
}

func newFrom(reader io.Reader) Key {
	secret := make([]byte, 32)
	_, err := io.ReadFull(reader, secret)
	if err != nil {
		panic(err)
	}
	return Key{
		PrivKey: cosmosed25519.GenPrivKeyFromSecret(secret),
		secret:  secret,
	}
}

func (key Key) ID() string {
	return key.PublicKeyStr()
}

func (key Key) PublicKey() (pubKey cryptotypes.PubKey) {
	return key.PubKey()
}

// PublicKeyStr returns the terra address of the public key
func (key Key) PublicKeyStr() string {
	addr := msg.AccAddress(key.PubKey().Address())
	return addr.String()
}

func (key Key) Raw() Raw {
	return key.secret
}

func (key Key) String() string {
	return fmt.Sprintf("TerraKey{PrivateKey: <redacted>, Public Key: %s}", key.PublicKeyStr())
}

func (key Key) GoString() string {
	return key.String()
}
