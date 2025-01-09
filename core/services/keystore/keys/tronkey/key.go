package tronkey

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

// Tron uses the same elliptic curve cryptography as Ethereum (ECDSA with secp256k1)
var curve = crypto.S256()

// Raw represents the Tron private key
type Raw []byte

// Key generates a public-private key pair from the raw private key
func (raw Raw) Key() Key {
	var privKey ecdsa.PrivateKey
	d := big.NewInt(0).SetBytes(raw)
	privKey.PublicKey.Curve = curve
	privKey.D = d
	privKey.PublicKey.X, privKey.PublicKey.Y = curve.ScalarBaseMult(d.Bytes())
	return Key{
		pubKey:  &privKey.PublicKey,
		privKey: &privKey,
	}
}

func (raw Raw) String() string {
	return "<Tron Raw Private Key>"
}

func (raw Raw) GoString() string {
	return raw.String()
}

var _ fmt.GoStringer = &Key{}

type Key struct {
	privKey *ecdsa.PrivateKey
	pubKey  *ecdsa.PublicKey
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
		privKey: privKeyECDSA,
		pubKey:  &privKeyECDSA.PublicKey,
	}, nil
}

func (key Key) ID() string {
	return key.Base58Address()
}

func (key Key) Raw() Raw {
	return key.privKey.D.Bytes()
}

func (key Key) ToEcdsaPrivKey() *ecdsa.PrivateKey {
	return key.privKey
}

func (key Key) String() string {
	return fmt.Sprintf("TronKey{PrivateKey: <redacted>, Address: %s}", key.Base58Address())
}

// GoString wraps String()
func (key Key) GoString() string {
	return key.String()
}

// Sign is used to sign a message
func (key Key) Sign(msg []byte) ([]byte, error) {
	return crypto.Sign(msg, key.privKey)
}

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
