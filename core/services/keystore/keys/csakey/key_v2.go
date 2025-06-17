package csakey

import (
	"crypto"
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"

	"github.com/smartcontractkit/wsrpc/credentials"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

func KeyFor(raw internal.Raw) KeyV2 {
	privKey := ed25519.PrivateKey(internal.Bytes(raw))
	return KeyV2{
		raw:       raw,
		signer:    &privKey,
		PublicKey: privKey.Public().(ed25519.PublicKey),
	}
}

type KeyV2 struct {
	raw    internal.Raw
	signer crypto.Signer

	PublicKey ed25519.PublicKey
	Version   int
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
		raw:       internal.NewRaw(privKey),
		signer:    &privKey,
		PublicKey: pubKey,
		Version:   2,
	}, nil
}

func MustNewV2XXXTestingOnly(k *big.Int) KeyV2 {
	seed := make([]byte, ed25519.SeedSize)
	copy(seed, k.Bytes())
	privKey := ed25519.NewKeyFromSeed(seed)
	return KeyV2{
		raw:       internal.NewRaw(privKey),
		signer:    privKey,
		PublicKey: privKey.Public().(ed25519.PublicKey),
		Version:   2,
	}
}

func (k KeyV2) ID() string {
	return k.PublicKeyString()
}

func (k KeyV2) PublicKeyString() string {
	return hex.EncodeToString(k.PublicKey)
}

func (k KeyV2) Raw() internal.Raw {
	return k.raw
}

func (k KeyV2) Public() crypto.PublicKey { return k.PublicKey }

func (k KeyV2) Sign(rand io.Reader, message []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	return k.signer.Sign(rand, message, opts)
}
