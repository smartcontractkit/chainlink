package ethkey

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

var curve = crypto.S256()

type Raw []byte

func (raw Raw) Key() KeyV2 {
	var privateKey ecdsa.PrivateKey
	d := big.NewInt(0).SetBytes(raw)
	privateKey.PublicKey.Curve = curve
	privateKey.D = d
	privateKey.PublicKey.X, privateKey.PublicKey.Y = curve.ScalarBaseMult(d.Bytes())
	address := EIP55AddressFromAddress(crypto.PubkeyToAddress(privateKey.PublicKey))
	return KeyV2{
		Address:    address,
		privateKey: &privateKey,
	}
}

func (raw Raw) String() string {
	return "<Eth Raw Private Key>"
}

func (raw Raw) GoString() string {
	return raw.String()
}

var _ fmt.GoStringer = &KeyV2{}

type KeyV2 struct {
	Address    EIP55Address
	privateKey *ecdsa.PrivateKey
}

func NewV2() (KeyV2, error) {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return KeyV2{}, err
	}
	return FromPrivateKey(privateKeyECDSA), nil
}

func FromPrivateKey(privKey *ecdsa.PrivateKey) (key KeyV2) {
	address := EIP55AddressFromAddress(crypto.PubkeyToAddress(privKey.PublicKey))
	return KeyV2{
		Address:    address,
		privateKey: privKey,
	}
}

func (key KeyV2) ID() string {
	return key.Address.Hex()
}

func (key KeyV2) Raw() Raw {
	return key.privateKey.D.Bytes()
}

func (key KeyV2) ToEcdsaPrivKey() *ecdsa.PrivateKey {
	return key.privateKey
}

func (key KeyV2) String() string {
	return fmt.Sprintf("EthKeyV2{PrivateKey: <redacted>, Address: %s}", key.Address)
}

func (key KeyV2) GoString() string {
	return key.String()
}

// Cmp uses byte-order address comparison to give a stable comparison between two keys
func (key KeyV2) Cmp(key2 KeyV2) int {
	return bytes.Compare(key.Address.Bytes(), key2.Address.Bytes())
}
