package keystore

import (
	"context"
	"fmt"
	"math/big"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
)

type VRF interface {
	Get(id string) (vrfkey.KeyV2, error)
	GetAll() ([]vrfkey.KeyV2, error)
	Create(ctx context.Context) (vrfkey.KeyV2, error)
	Add(ctx context.Context, key vrfkey.KeyV2) error
	Delete(ctx context.Context, id string) (vrfkey.KeyV2, error)
	Import(ctx context.Context, keyJSON []byte, password string) (vrfkey.KeyV2, error)
	Export(id string, password string) ([]byte, error)

	GenerateProof(id string, seed *big.Int) (vrfkey.Proof, error)
}

var (
	ErrMissingVRFKey = errors.New("unable to find VRF key")
)

type vrf struct {
	*keyManager
}

var _ VRF = &vrf{}

func newVRFKeyStore(km *keyManager) *vrf {
	return &vrf{
		km,
	}
}

func (ks *vrf) Get(id string) (vrfkey.KeyV2, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return vrfkey.KeyV2{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *vrf) GetAll() (keys []vrfkey.KeyV2, _ error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, key := range ks.keyRing.VRF {
		keys = append(keys, key)
	}
	return keys, nil
}

func (ks *vrf) Create(ctx context.Context) (vrfkey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return vrfkey.KeyV2{}, ErrLocked
	}
	key, err := vrfkey.NewV2()
	if err != nil {
		return vrfkey.KeyV2{}, err
	}
	return key, ks.safeAddKey(ctx, key)
}

func (ks *vrf) Add(ctx context.Context, key vrfkey.KeyV2) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if _, found := ks.keyRing.VRF[key.ID()]; found {
		return fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return ks.safeAddKey(ctx, key)
}

func (ks *vrf) Delete(ctx context.Context, id string) (vrfkey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return vrfkey.KeyV2{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return vrfkey.KeyV2{}, err
	}
	err = ks.safeRemoveKey(ctx, key)
	return key, err
}

func (ks *vrf) Import(ctx context.Context, keyJSON []byte, password string) (vrfkey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return vrfkey.KeyV2{}, ErrLocked
	}
	key, err := vrfkey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return vrfkey.KeyV2{}, errors.Wrap(err, "VRFKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.VRF[key.ID()]; found {
		return vrfkey.KeyV2{}, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return key, ks.keyManager.safeAddKey(ctx, key)
}

func (ks *vrf) Export(id string, password string) ([]byte, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return nil, err
	}
	return key.ToEncryptedJSON(password, ks.scryptParams)
}

func (ks *vrf) GenerateProof(id string, seed *big.Int) (vrfkey.Proof, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return vrfkey.Proof{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return vrfkey.Proof{}, err
	}
	return key.GenerateProof(seed)
}

func (ks *vrf) getByID(id string) (vrfkey.KeyV2, error) {
	key, found := ks.keyRing.VRF[id]
	if !found {
		return vrfkey.KeyV2{}, KeyNotFoundError{ID: id, KeyType: "VRF"}
	}
	return key, nil
}
