package keystore

import (
	"fmt"

	"github.com/pkg/errors"

	stark "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/keys"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/starkkey"
)

//go:generate mockery --name StarkNet --output ./mocks/ --case=underscore --filename starknet.go

type StarkNet interface {
	Get(id string) (stark.StarkKey, error)
	GetAll() ([]stark.StarkKey, error)
	Create() (stark.StarkKey, error)
	Add(key stark.StarkKey) error
	Delete(id string) (stark.StarkKey, error)
	Import(keyJSON []byte, password string) (stark.StarkKey, error)
	Export(id string, password string) ([]byte, error)
	EnsureKey() error
}

type starknet struct {
	*keyManager
}

var _ StarkNet = &starknet{}

func newStarkNetKeyStore(km *keyManager) *starknet {
	return &starknet{
		km,
	}
}

func (ks *starknet) Get(id string) (stark.StarkKey, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return stark.StarkKey{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *starknet) GetAll() (keys []stark.StarkKey, _ error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, key := range ks.keyRing.StarkNet {
		keys = append(keys, key)
	}
	return keys, nil
}

func (ks *starknet) Create() (stark.StarkKey, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return stark.StarkKey{}, ErrLocked
	}
	key, err := stark.New()
	if err != nil {
		return stark.StarkKey{}, err
	}
	return key, ks.safeAddKey(key)
}

func (ks *starknet) Add(key stark.StarkKey) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if _, found := ks.keyRing.StarkNet[key.ID()]; found {
		return fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return ks.safeAddKey(key)
}

func (ks *starknet) Delete(id string) (stark.StarkKey, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return stark.StarkKey{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return stark.StarkKey{}, err
	}
	err = ks.safeRemoveKey(key)
	return key, err
}

func (ks *starknet) Import(keyJSON []byte, password string) (stark.StarkKey, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return stark.StarkKey{}, ErrLocked
	}
	key, err := starkkey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return stark.StarkKey{}, errors.Wrap(err, "StarkNetKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.StarkNet[key.ID()]; found {
		return stark.StarkKey{}, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return key, ks.keyManager.safeAddKey(key)
}

func (ks *starknet) Export(id string, password string) ([]byte, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return nil, err
	}
	return starkkey.ToEncryptedJSON(key, password, ks.scryptParams)
}

func (ks *starknet) EnsureKey() error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if len(ks.keyRing.StarkNet) > 0 {
		return nil
	}

	key, err := stark.New()
	if err != nil {
		return err
	}

	ks.logger.Infof("Created StarkNet key with ID %s", key.ID())

	return ks.safeAddKey(key)
}

var (
	ErrNoStarkNetKey = errors.New("no starknet keys exist")
)

func (ks *starknet) getByID(id string) (stark.StarkKey, error) {
	key, found := ks.keyRing.StarkNet[id]
	if !found {
		return stark.StarkKey{}, KeyNotFoundError{ID: id, KeyType: "StarkNet"}
	}
	return key, nil
}
