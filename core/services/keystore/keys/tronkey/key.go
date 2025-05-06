package tronkey

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

// Tron uses the same elliptic curve cryptography as Ethereum (ECDSA with secp256k1)
var curve = crypto.S256()

// Key generates a public-private key pair from the raw private key
func KeyFor(raw internal.Raw) Key {
	var privKey ecdsa.PrivateKey
	d := big.NewInt(0).SetBytes(internal.Bytes(raw))
	privKey.PublicKey.Curve = curve
	privKey.D = d
	privKey.PublicKey.X, privKey.PublicKey.Y = curve.ScalarBaseMult(d.Bytes())
	return Key{
		raw:    raw,
		signFn: func(bytes []byte) ([]byte, error) { return crypto.Sign(bytes, &privKey) },
		pubKey: &privKey.PublicKey,
	}
}

type Key struct {
	raw    internal.Raw
	signFn func([]byte) ([]byte, error)

	pubKey *ecdsa.PublicKey
}

func New() (Key, error) {
	return newFrom(rand.Reader)
}

// MustNewInsecure return Key if no error
// This insecure function is used for testing purposes only
func MustNewInsecure(reader io.Reader) Key {
	key, err := newFrom(reader)
	if err != nil {
		panic(err)
	}
	return key
}

func newFrom(reader io.Reader) (Key, error) {
	privKeyECDSA, err := ecdsa.GenerateKey(curve, reader)
	if err != nil {
		return Key{}, err
	}
	return Key{
		raw:    internal.NewRaw(privKeyECDSA.D.Bytes()),
		signFn: func(bytes []byte) ([]byte, error) { return crypto.Sign(bytes, privKeyECDSA) },
		pubKey: &privKeyECDSA.PublicKey,
	}, nil
}

func (key Key) ID() string {
	return key.Base58Address()
}

func (key Key) Raw() internal.Raw { return key.raw }

// Sign is used to sign a message
func (key Key) Sign(msg []byte) ([]byte, error) { return key.signFn(msg) }

// PublicKeyStr returns the public key as a hexadecimal string
func (key Key) PublicKeyStr() string {
	pubKeyBytes := crypto.FromECDSAPub(key.pubKey)
	return hex.EncodeToString(pubKeyBytes)
}

// Base58Address returns the Tron address in Base58 format with checksum
func (key Key) Base58Address() string {
	address := PubkeyToAddress(*key.pubKey)
	return address.String()
}
