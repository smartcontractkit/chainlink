package csakey

import (
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/wsrpc/credentials"
)

type Raw []byte

func (raw Raw) Key() KeyV2 {
	privKey := ed25519.PrivateKey(raw)
	return KeyV2{
		privateKey: &privKey,
		PublicKey:  ed25519PubKeyFromPrivKey(privKey),
	}
}

func (raw Raw) String() string {
	return "<CSA Raw Private Key>"
}

func (raw Raw) GoString() string {
	return raw.String()
}

func (raw Raw) Bytes() []byte {
	return ([]byte)(raw)
}

var _ fmt.GoStringer = &KeyV2{}

type KeyV2 struct {
	privateKey *ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
	Version    int
}

func (k KeyV2) StaticSizedPublicKey() (sspk credentials.StaticSizedPublicKey) {
	if len(k.PublicKey) != ed25519.PublicKeySize {
		panic(fmt.Sprintf("expected ed25519.PublicKey to have len %d but got len %d", ed25519.PublicKeySize, len(k.PublicKey)))
	}
	copy(sspk[:], k.PublicKey)
	return sspk
}

func NewV2() (KeyV2, error) {
	pubKey, privKey, err := ed25519.GenerateKey(cryptorand.Reader)
	if err != nil {
		return KeyV2{}, err
	}
	return KeyV2{
		privateKey: &privKey,
		PublicKey:  pubKey,
		Version:    2,
	}, nil
}

func MustNewV2XXXTestingOnly(k *big.Int) KeyV2 {
	seed := make([]byte, ed25519.SeedSize)
	copy(seed, k.Bytes())
	privKey := ed25519.NewKeyFromSeed(seed)
	return KeyV2{
		privateKey: &privKey,
		PublicKey:  ed25519PubKeyFromPrivKey(privKey),
		Version:    2,
	}
}

func (k KeyV2) ID() string {
	return k.PublicKeyString()
}

func (k KeyV2) PublicKeyString() string {
	return hex.EncodeToString(k.PublicKey)
}

func (k KeyV2) Raw() Raw {
	return Raw(*k.privateKey)
}

func (k KeyV2) String() string {
	return fmt.Sprintf("CSAKeyV2{PrivateKey: <redacted>, PublicKey: %s}", k.PublicKey)
}

func (k KeyV2) GoString() string {
	return k.String()
}

func ed25519PubKeyFromPrivKey(privKey ed25519.PrivateKey) ed25519.PublicKey {
	publicKey := make([]byte, ed25519.PublicKeySize)
	copy(publicKey, privKey[32:])
	return publicKey
}
