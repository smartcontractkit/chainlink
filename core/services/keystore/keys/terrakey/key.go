package terrakey

import (
	"fmt"
	"io"

	"github.com/smartcontractkit/terra.go/msg"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type Raw []byte

func (raw Raw) Key() Key {
	privKey2 := secp256k1.PrivKey(raw)
	return Key{
		k2: privKey2,
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
	k2 secp256k1.PrivKey
}

func New() Key {
	privKey2 := secp256k1.GenPrivKey()
	return Key{
		k2: privKey2,
	}
}

func MustNewInsecure(reader io.Reader) Key {
	seed := make([]byte, 32)
	_, err := reader.Read(seed)
	if err != nil {
		panic(err)
	}
	privKey := secp256k1.GenPrivKeySecp256k1(seed)
	return Key{
		k2: privKey,
	}
}

func (key Key) ID() string {
	return key.PublicKeyStr()
}

func (key Key) PublicKey() (pubKey crypto.PubKey) {
	return key.k2.PubKey()
}

// PublicKeyStr returns the terra address of the public key
func (key Key) PublicKeyStr() string {
	addr := msg.AccAddress(key.k2.PubKey().Address())
	return addr.String()
}

func (key Key) Raw() Raw {
	return key.k2.Bytes()
}

// ToPrivKey returns the key usable for signing.
func (key Key) ToPrivKey() secp256k1.PrivKey {
	return key.k2
}

func (key Key) String() string {
	return fmt.Sprintf("TerraKey{PrivateKey: <redacted>, Public Key: %s}", key.PublicKeyStr())
}

func (key Key) GoString() string {
	return key.String()
}
