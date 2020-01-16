package vrfkey

import (
	"crypto/ecdsa"

	"chainlink/core/services/signatures/secp256k1"
	"chainlink/core/services/vrf"

	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
)

var suite = secp256k1.NewBlakeKeccackSecp256k1()

// newPrivateKey(k) is k wrapped in a PrivateKey along with corresponding
// PublicKey, or an error. Internal use only. Use cltest.StoredVRFKey for stable
// testing key, or CreateKey if you don't need determinism.
func newPrivateKey(rawKey *big.Int) (*PrivateKey, error) {
	if rawKey.Cmp(secp256k1.GroupOrder) != -1 || rawKey.Cmp(big.NewInt(0)) == -1 {
		return nil, fmt.Errorf("secret key must be in {0, ..., #secp256k1 - 1}")
	}
	sk := &PrivateKey{}
	sk.k = secp256k1.IntToScalar(rawKey)
	pk := secp256k1.LongMarshal(suite.Point().Mul(sk.k, nil))
	if len(pk) != UncompressedPublicKeyLength {
		panic(fmt.Errorf("public key %x has wrong length", pk))
	}
	if l := copy(sk.PublicKey[:], pk[:]); l != UncompressedPublicKeyLength {
		panic(fmt.Errorf("failed to copy correct length in serialized public key"))
	}
	return sk, nil
}

// k.MarshaledProof(seed) is a VRF proof of randomness using k and seed, in the
// form required by VRF.sol's randomValueFromVRFProof
func (k *PrivateKey) MarshaledProof(seed *big.Int) (vrf.MarshaledProof, error) {
	proof, err := vrf.GenerateProof(secp256k1.ToInt(k.k), seed)
	if err != nil {
		return vrf.MarshaledProof{}, err
	}
	rv, err := proof.MarshalForSolidityVerifier()
	if err != nil {
		return vrf.MarshaledProof{}, err
	}
	return rv, nil
}

// gethKey returns the geth keystore representation of k. Do not abuse this to
// convert a VRF key to an ethereum key!
func (k *PrivateKey) gethKey() *keystore.Key {
	return &keystore.Key{
		Address:    k.PublicKey.Address(),
		PrivateKey: &ecdsa.PrivateKey{D: secp256k1.ToInt(k.k)},
	}
}

// fromGethKey returns the vrfkey representation of gethKey. Do not abuse this
// to convert an ethereum key into a VRF key!
func fromGethKey(gethKey *keystore.Key) *PrivateKey {
	k := secp256k1.IntToScalar(gethKey.PrivateKey.D)
	var publicKey PublicKey
	copy(publicKey[:], secp256k1.LongMarshal(secp256k1.ScalarToPublicPoint(k)))
	return &PrivateKey{k, publicKey}
}

// CreateKey makes a new VRF proving key from cryptographically secure entropy
func CreateKey() (key *PrivateKey) {
	sk := suite.Scalar().Pick(suite.RandomStream())
	k, err := newPrivateKey(secp256k1.ToInt(sk))
	if err != nil {
		panic(errors.Wrapf(err, "should not be possible to error, here"))
	}
	return k
}

// NewPrivateKeyXXXTestingOnly is for testing purposes only!
func NewPrivateKeyXXXTestingOnly(k *big.Int) (*PrivateKey, error) {
	return newPrivateKey(k)
}

// String reduces the risk of accidentally logging the private key
func (k *PrivateKey) String() string {
	return fmt.Sprintf("PrivateKey{k: <redacted>, PublicKey: 0x%x}", k.PublicKey)
}

// GoStringer reduces the risk of accidentally logging the private key
func (k *PrivateKey) GoStringer() string {
	return k.String()
}
