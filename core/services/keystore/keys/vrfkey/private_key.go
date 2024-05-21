package vrfkey

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3"

	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
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
	PublicKey secp256k1.PublicKey
}

func (k PrivateKey) ToV2() KeyV2 {
	return KeyV2{
		k:         &k.k,
		PublicKey: k.PublicKey,
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
	var publicKey secp256k1.PublicKey
	copy(publicKey[:], rawPublicKey)
	return &PrivateKey{secretKey, publicKey}
}

func (k *PrivateKey) String() string {
	return fmt.Sprintf("PrivateKey{k: <redacted>, PublicKey: %s}", k.PublicKey)
}

// GoString reduces the risk of accidentally logging the private key
func (k *PrivateKey) GoString() string {
	return k.String()
}

// Decrypt returns the PrivateKey in e, decrypted via auth, or an error
func Decrypt(e EncryptedVRFKey, auth string) (*PrivateKey, error) {
	// NOTE: We do this shuffle to an anonymous struct
	// solely to add a throwaway UUID, so we can leverage
	// the keystore.DecryptKey from the geth which requires it
	// as of 1.10.0.
	keyJSON, err := json.Marshal(struct {
		Address string              `json:"address"`
		Crypto  keystore.CryptoJSON `json:"crypto"`
		Version int                 `json:"version"`
		Id      string              `json:"id"`
	}{
		Address: e.VRFKey.Address,
		Crypto:  e.VRFKey.Crypto,
		Version: e.VRFKey.Version,
		Id:      uuid.New().String(),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "while marshaling key for decryption")
	}
	gethKey, err := keystore.DecryptKey(keyJSON, adulteratedPassword(auth))
	if err != nil {
		return nil, errors.Wrapf(err, "could not decrypt VRF key %s",
			e.PublicKey.String())
	}
	return fromGethKey(gethKey), nil
}
