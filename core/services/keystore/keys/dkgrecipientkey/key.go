package dkgrecipientkey

import (
	cryptorand "crypto/rand"
	"encoding/hex"

	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"
	"github.com/smartcontractkit/smdkg/p256keyring"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

var _ internal.Key = &Key{}
var _ dkgocrtypes.P256Keyring = &Key{}

type Key struct {
	raw     internal.Raw
	keyRing dkgocrtypes.P256Keyring
}

func New() (Key, error) {
	keyRing, err := p256keyring.New(cryptorand.Reader)
	if err != nil {
		return Key{}, err
	}
	rawBytes, err := keyRing.MarshalBinary()
	if err != nil {
		return Key{}, err
	}

	return Key{raw: internal.NewRaw(rawBytes), keyRing: keyRing}, nil
}

func (k Key) PublicKey() dkgocrtypes.P256ParticipantPublicKey {
	return k.keyRing.PublicKey()
}

func (k Key) PublicKeyString() string {
	return hex.EncodeToString(k.keyRing.PublicKey()[:])
}

func (k Key) ID() string {
	return k.PublicKeyString()
}

func (k Key) ECDH(publicKey dkgocrtypes.P256ParticipantPublicKey) (dkgocrtypes.P256ECDHSharedSecret, error) {
	return k.keyRing.ECDH(publicKey)
}

func KeyFor(raw internal.Raw) Key {
	keyRing := &p256keyring.P256Keyring{}
	err := keyRing.UnmarshalBinary(internal.Bytes(raw))
	if err != nil {
		panic(err)
	}
	return Key{raw: raw, keyRing: keyRing}
}

func (k Key) Raw() internal.Raw {
	return k.raw
}
