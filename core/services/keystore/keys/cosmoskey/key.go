package cosmoskey

import (
	"crypto/ecdsa"
	cryptorand "crypto/rand"
	"fmt"
	"io"
	"math/big"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/ethereum/go-ethereum/crypto"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
)

var secpSigningAlgo, _ = keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), []keyring.SignatureAlgo{hd.Secp256k1})

type Raw []byte

func (raw Raw) Key() Key {
	d := big.NewInt(0).SetBytes(raw)
	privKey := secpSigningAlgo.Generate()(d.Bytes())
	return Key{
		d: d,
		k: privKey,
	}
}

func (raw Raw) String() string {
	return "<Cosmos Raw Private Key>"
}

func (raw Raw) GoString() string {
	return raw.String()
}

var _ fmt.GoStringer = &Key{}

// Key represents Cosmos key
type Key struct {
	d *big.Int
	k cryptotypes.PrivKey
}

// New creates new Key
func New() Key {
	return newFrom(cryptorand.Reader)
}

// MustNewInsecure return Key
func MustNewInsecure(reader io.Reader) Key {
	return newFrom(reader)
}

func newFrom(reader io.Reader) Key {
	rawKey, err := ecdsa.GenerateKey(crypto.S256(), reader)
	if err != nil {
		panic(err)
	}
	privKey := secpSigningAlgo.Generate()(rawKey.D.Bytes())
	if err != nil {
		panic(err)
	}

	return Key{
		d: rawKey.D,
		k: privKey,
	}
}

func (key Key) ID() string {
	return key.PublicKeyStr()
}

func (key Key) PublicKey() (pubKey cryptotypes.PubKey) {
	return key.k.PubKey()
}

// PublicKeyStr returns the cosmos address of the public key
func (key Key) PublicKeyStr() string {
	addr := types.AccAddress(key.k.PubKey().Address())
	return addr.String()
}

func (key Key) Raw() Raw {
	return key.d.Bytes()
}

// ToPrivKey returns the key usable for signing.
func (key Key) ToPrivKey() cryptotypes.PrivKey {
	return key.k
}

func (key Key) String() string {
	return fmt.Sprintf("CosmosKey{PrivateKey: <redacted>, Public Key: %s}", key.PublicKeyStr())
}

func (key Key) GoString() string {
	return key.String()
}
