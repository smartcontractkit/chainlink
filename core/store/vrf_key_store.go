package store

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/store/models/vrfkey"
)

// VRFKeyStore tracks auxiliary VRF secret keys, and generates their VRF proofs
//
// VRF proofs need access to the actual secret key, which geth does not expose.
// Similar to the way geth's KeyStore exposes signing capability, VRFKeyStore
// exposes VRF proof generation without the caller needing explicit knowledge of
// the secret key.
type VRFKeyStore struct {
	lock  sync.RWMutex
	keys  InMemoryKeyStore
	store *Store
}

type InMemoryKeyStore = map[vrfkey.PublicKey]vrfkey.PrivateKey

// NewVRFKeyStore returns an empty VRFKeyStore
func NewVRFKeyStore(store *Store) *VRFKeyStore {
	return &VRFKeyStore{
		lock:  sync.RWMutex{},
		keys:  make(InMemoryKeyStore),
		store: store,
	}
}

// GenerateProof is marshaled randomness proof given k and VRF input seed
// computed from the SeedData
//
// Key must have already been unlocked in ks, as constructing the VRF proof
// requires the secret key.
func (ks *VRFKeyStore) GenerateProof(k vrfkey.PublicKey, i vrf.PreSeedData) (
	vrf.MarshaledOnChainResponse, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	privateKey, found := ks.keys[k]
	if !found {
		return vrf.MarshaledOnChainResponse{}, fmt.Errorf(
			"key %s has not been unlocked", k)
	}
	return privateKey.MarshaledProof(i)
}

// Unlock tries to unlock each vrf key in the db, using the given pass phrase,
// and returns any keys it manages to unlock, and any errors which result.
func (ks *VRFKeyStore) Unlock(phrase string) (keysUnlocked []vrfkey.PublicKey,
	merr error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	keys, err := ks.get()
	if err != nil {
		return nil, errors.Wrap(err, "while retrieving vrf keys from db")
	}
	for _, k := range keys {
		key, err := k.Decrypt(phrase)
		if err != nil {
			merr = multierr.Append(merr, err)
			continue
		}
		ks.keys[key.PublicKey] = *key
		keysUnlocked = append(keysUnlocked, key.PublicKey)
	}
	return keysUnlocked, merr
}

// Forget removes the in-memory copy of the secret key of k, or errors if not
// present. Caller is responsible for taking ks.lock.
func (ks *VRFKeyStore) forget(k vrfkey.PublicKey) error {
	if _, found := ks.keys[k]; !found {
		return fmt.Errorf("public key %s is not unlocked; can't forget it", k)
	}

	delete(ks.keys, k)
	return nil
}

func (ks *VRFKeyStore) Forget(k vrfkey.PublicKey) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	return ks.forget(k)
}
