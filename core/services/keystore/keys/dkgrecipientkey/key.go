package dkgrecipientkey

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	"encoding/json"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"
	"github.com/smartcontractkit/smdkg/dummydkg"
)

var _ internal.Key = &Key{}
var _ dkgocrtypes.P256Keyring = &Key{}

type Key struct {
	raw     internal.Raw
	keyRing dkgocrtypes.P256Keyring
}

func New() (Key, error) {
	keyRing, err := dummydkg.NewP256Keyring(cryptorand.Reader)
	if err != nil {
		return Key{}, err
	}
	rawBytes, err := json.Marshal(keyRing)
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
	panic("not implemented")
}

func (k Key) Raw() internal.Raw {
	panic("not implemented")
}
