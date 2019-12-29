package store

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"chainlink/core/services/vrf"
	"chainlink/core/store/models/vrf_key"
)

// VRFKeyStore tracks auxillary VRF secret keys, and generates their VRF proofs
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

type InMemoryKeyStore = map[vrf_key.PublicKey]vrf_key.PrivateKey

// NewVRFKeyStore returns an empty VRFKeyStore
func NewVRFKeyStore(store *Store) *VRFKeyStore {
	return &VRFKeyStore{
		lock:  sync.RWMutex{},
		keys:  make(InMemoryKeyStore),
		store: store,
	}
}

// GenerateProof(k, seed) is marshaled randomness proof given k and seed.
//
// k must have already been unlocked in ks.
func (ks *VRFKeyStore) GenerateProof(k *vrf_key.PublicKey, seed *big.Int) (
	vrf.MarshaledProof, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	privateKey, found := ks.keys[*k]
	if !found {
		return vrf.MarshaledProof{}, fmt.Errorf("key %s has not been unlocked", k)
	}
	return privateKey.MarshaledProof(seed)
}

// Unlock tries to unlock each vrf key in the db, using the given pass phrase
func (ks *VRFKeyStore) Unlock(phrase string) (keysUnlocked []vrf_key.PublicKey, merr error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	keys, err := ks.get(nil)
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
func (ks *VRFKeyStore) forget(k *vrf_key.PublicKey) error {
	if _, found := ks.keys[*k]; !found {
		return fmt.Errorf("public key %s is not unlocked; can't forget it", k)
	} else {
		delete(ks.keys, *k)
		return nil
	}
}

func (ks *VRFKeyStore) Forget(k *vrf_key.PublicKey) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	return ks.forget(k)
}
