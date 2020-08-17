package vrfkey

import (
	"crypto/ecdsa"

	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/services/vrf"

	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3"
)

// PrivateKey represents the secret used to construct a VRF proof.
//
// Don't serialize directly, use Encrypt method, with user-supplied passphrase.
// The unencrypted PrivateKey struct should only live in-memory.
//
// Only use it if you absolutely need it (i.e., for a novel crypto protocol.)
// Implement whatever cryptography you need on this struct, so your callers
// don't need to know the secret key explicitly. (See, e.g., MarshaledProof.)
type PrivateKey struct {
	k         kyber.Scalar
	PublicKey PublicKey
}

var suite = secp256k1.NewBlakeKeccackSecp256k1()

// newPrivateKey(k) is k wrapped in a PrivateKey along with corresponding
// PublicKey, or an error. Internal use only. Use cltest.StoredVRFKey for stable
// testing key, or CreateKey if you don't need determinism.
func newPrivateKey(rawKey *big.Int) (*PrivateKey, error) {
	if rawKey.Cmp(secp256k1.GroupOrder) >= 0 || rawKey.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("secret key must be in {1, ..., #secp256k1 - 1}")
	}
	sk := &PrivateKey{}
	sk.k = secp256k1.IntToScalar(rawKey)
	pk, err := suite.Point().Mul(sk.k, nil).MarshalBinary()
	if err != nil {
		panic(errors.Wrapf(err, "could not marshal public key"))
	}
	if len(pk) != CompressedPublicKeyLength {
		panic(fmt.Errorf("public key %x has wrong length", pk))
	}
	if l := copy(sk.PublicKey[:], pk[:]); l != CompressedPublicKeyLength {
		panic(fmt.Errorf("failed to copy correct length in serialized public key"))
	}
	return sk, nil
}

// MarshaledProof is a VRF proof of randomness using i.Key and seed, in the form
// required by VRFCoordinator.sol's fulfillRandomnessRequest
func (k *PrivateKey) MarshaledProof(i vrf.PreSeedData) (
	vrf.MarshaledOnChainResponse, error) {
	return vrf.GenerateProofResponse(secp256k1.ScalarToHash(k.k), i)
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
	secretKey := secp256k1.IntToScalar(gethKey.PrivateKey.D)
	rawPublicKey, err := secp256k1.ScalarToPublicPoint(secretKey).MarshalBinary()
	if err != nil {
		panic(err) // Only way this can happen is out-of-memory failure
	}
	var publicKey PublicKey
	copy(publicKey[:], rawPublicKey)
	return &PrivateKey{secretKey, publicKey}
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
func NewPrivateKeyXXXTestingOnly(k *big.Int) *PrivateKey {
	rv, err := newPrivateKey(k)
	if err != nil {
		panic(err)
	}
	return rv
}

// String reduces the risk of accidentally logging the private key
func (k *PrivateKey) String() string {
	return fmt.Sprintf("PrivateKey{k: <redacted>, PublicKey: %s}", k.PublicKey)
}

// GoStringer reduces the risk of accidentally logging the private key
func (k *PrivateKey) GoStringer() string {
	return k.String()
}
