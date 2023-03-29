package keystore

import (
	"fmt"

	"github.com/pkg/errors"

	stark "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/keys"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/starkkey"
)

//go:generate mockery --quiet --name Starknet --output ./mocks/ --case=underscore --filename starknet.go

type Starknet interface {
	Get(id string) (stark.Key, error)
	GetAll() ([]stark.Key, error)
	Create() (stark.Key, error)
	Add(key stark.Key) error
	Delete(id string) (stark.Key, error)
	Import(keyJSON []byte, password string) (stark.Key, error)
	Export(id string, password string) ([]byte, error)
	EnsureKey() error
}

type starknet struct {
	*keyManager
}

var _ Starknet = &starknet{}

func newStarknetKeyStore(km *keyManager) *starknet {
	return &starknet{
		km,
	}
}

func (ks *starknet) Get(id string) (stark.Key, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return stark.Key{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *starknet) GetAll() (keys []stark.Key, _ error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, key := range ks.keyRing.Starknet {
		keys = append(keys, key)
	}
	return keys, nil
}

func (ks *starknet) Create() (stark.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return stark.Key{}, ErrLocked
	}
	key, err := stark.New()
	if err != nil {
		return stark.Key{}, err
	}
	return key, ks.safeAddKey(key)
}

func (ks *starknet) Add(key stark.Key) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if _, found := ks.keyRing.Starknet[key.ID()]; found {
		return fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return ks.safeAddKey(key)
}

func (ks *starknet) Delete(id string) (stark.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return stark.Key{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return stark.Key{}, err
	}
	err = ks.safeRemoveKey(key)
	return key, err
}

func (ks *starknet) Import(keyJSON []byte, password string) (stark.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return stark.Key{}, ErrLocked
	}
	key, err := starkkey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return stark.Key{}, errors.Wrap(err, "StarknetKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.Starknet[key.ID()]; found {
		return stark.Key{}, fmt.Errorf("key with ID %s already exists", key.ID())
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
	if len(ks.keyRing.Starknet) > 0 {
		return nil
	}

	key, err := stark.New()
	if err != nil {
		return err
	}

	ks.logger.Infof("Created Starknet key with ID %s", key.ID())

	return ks.safeAddKey(key)
}

var (
	ErrNoStarknetKey = errors.New("no starknet keys exist")
)

func (ks *starknet) getByID(id string) (stark.Key, error) {
	key, found := ks.keyRing.Starknet[id]
	if !found {
		return stark.Key{}, KeyNotFoundError{ID: id, KeyType: "Starknet"}
	}
	return key, nil
}
