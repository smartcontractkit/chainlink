package keystore

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/starkkey"
)

//go:generate mockery --name StarkNet --output ./mocks/ --case=underscore --filename starknet.go

type StarkNet interface {
	Get(id string) (starkkey.Key, error)
	GetAll() ([]starkkey.Key, error)
	Create() (starkkey.Key, error)
	Add(key starkkey.Key) error
	Delete(id string) (starkkey.Key, error)
	Import(keyJSON []byte, password string) (starkkey.Key, error)
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

func (ks *starknet) Get(id string) (starkkey.Key, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return starkkey.Key{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *starknet) GetAll() (keys []starkkey.Key, _ error) {
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

func (ks *starknet) Create() (starkkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return starkkey.Key{}, ErrLocked
	}
	key, err := starkkey.New()
	if err != nil {
		return starkkey.Key{}, err
	}
	return key, ks.safeAddKey(key)
}

func (ks *starknet) Add(key starkkey.Key) error {
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

func (ks *starknet) Delete(id string) (starkkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return starkkey.Key{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return starkkey.Key{}, err
	}
	err = ks.safeRemoveKey(key)
	return key, err
}

func (ks *starknet) Import(keyJSON []byte, password string) (starkkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return starkkey.Key{}, ErrLocked
	}
	key, err := starkkey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return starkkey.Key{}, errors.Wrap(err, "StarkNetKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.StarkNet[key.ID()]; found {
		return starkkey.Key{}, fmt.Errorf("key with ID %s already exists", key.ID())
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
	return key.ToEncryptedJSON(password, ks.scryptParams)
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

	key, err := starkkey.New()
	if err != nil {
		return err
	}

	ks.logger.Infof("Created StarkNet key with ID %s", key.ID())

	return ks.safeAddKey(key)
}

var (
	ErrNoStarkNetKey = errors.New("no starknet keys exist")
)

func (ks *starknet) getByID(id string) (starkkey.Key, error) {
	key, found := ks.keyRing.StarkNet[id]
	if !found {
		return starkkey.Key{}, KeyNotFoundError{ID: id, KeyType: "StarkNet"}
	}
	return key, nil
}
