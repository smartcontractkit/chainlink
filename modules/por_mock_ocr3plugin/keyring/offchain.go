package keyring

import (
	"crypto/ed25519"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"golang.org/x/crypto/curve25519"
)

// DummyOffchainKeyring is not intended to be used in production.
type DummyOffchainKeyring struct {
	OffchainPrivateKey         ed25519.PrivateKey
	ConfigEncryptionPrivateKey [curve25519.ScalarSize]byte
}

var _ types.OffchainKeyring = &DummyOffchainKeyring{}

func (ring *DummyOffchainKeyring) OffchainSign(msg []byte) (signature []byte, err error) {
	sig := ed25519.Sign(ring.OffchainPrivateKey, msg)
	return sig, nil
}

func (ring *DummyOffchainKeyring) ConfigDiffieHellman(
	point [curve25519.PointSize]byte,
) (
	sharedPoint [curve25519.PointSize]byte,
	err error,
) {
	p, err := curve25519.X25519(ring.ConfigEncryptionPrivateKey[:], point[:])
	if err != nil {
		return [curve25519.PointSize]byte{}, err
	}
	copy(sharedPoint[:], p)
	return sharedPoint, nil
}

func (ring *DummyOffchainKeyring) OffchainPublicKey() types.OffchainPublicKey {
	var ocpk types.OffchainPublicKey
	pubKey := ring.OffchainPrivateKey.Public().(ed25519.PublicKey)
	if len(ocpk) != len(pubKey) {
		// assertion
		panic("OffchainPublicKey length mismatch")
	}
	copy(ocpk[:], pubKey)
	return ocpk
}

func (ring *DummyOffchainKeyring) ConfigEncryptionPublicKey() types.ConfigEncryptionPublicKey {
	rv, err := curve25519.X25519(ring.ConfigEncryptionPrivateKey[:], curve25519.Basepoint)
	if err != nil {
		panic("failure while computing public key: " + err.Error())
	}
	var rvFixed [curve25519.PointSize]byte
	copy(rvFixed[:], rv)
	return rvFixed
}
