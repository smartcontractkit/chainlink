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

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

var secpSigningAlgo, _ = keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), []keyring.SignatureAlgo{hd.Secp256k1})

func KeyFor(raw internal.Raw) Key {
	d := big.NewInt(0).SetBytes(internal.Bytes(raw))
	privKey := secpSigningAlgo.Generate()(d.Bytes())
	return Key{
		raw:    raw,
		signFn: privKey.Sign,
		pubKey: privKey.PubKey(),
	}
}

// Key represents Cosmos key
type Key struct {
	raw    internal.Raw
	signFn func([]byte) ([]byte, error)
	pubKey cryptotypes.PubKey
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

	return Key{
		raw:    internal.NewRaw(rawKey.D.Bytes()),
		signFn: privKey.Sign,
		pubKey: privKey.PubKey(),
	}
}

func (key Key) ID() string {
	return key.PublicKeyStr()
}

func (key Key) PublicKey() (pubKey cryptotypes.PubKey) {
	return key.pubKey
}

func (key Key) PublicKeyStr() string {
	return fmt.Sprintf("%X", key.pubKey.Bytes())
}

func (key Key) Raw() internal.Raw { return key.raw }

func (key Key) Sign(data []byte) ([]byte, error) {
	return key.signFn(data)
}
