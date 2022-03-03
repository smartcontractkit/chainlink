package keystore

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/terrakey"
)

//go:generate mockery --name Terra --output ./mocks/ --case=underscore --filename terra.go

type Terra interface {
	Get(id string) (terrakey.Key, error)
	GetAll() ([]terrakey.Key, error)
	Create() (terrakey.Key, error)
	Add(key terrakey.Key) error
	Delete(id string) (terrakey.Key, error)
	Import(keyJSON []byte, password string) (terrakey.Key, error)
	Export(id string, password string) ([]byte, error)
	EnsureKey() error
}

type terra struct {
	*keyManager
}

var _ Terra = &terra{}

func newTerraKeyStore(km *keyManager) *terra {
	return &terra{
		km,
	}
}

func (ks *terra) Get(id string) (terrakey.Key, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return terrakey.Key{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *terra) GetAll() (keys []terrakey.Key, _ error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, key := range ks.keyRing.Terra {
		keys = append(keys, key)
	}
	return keys, nil
}

func (ks *terra) Create() (terrakey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return terrakey.Key{}, ErrLocked
	}
	key := terrakey.New()
	return key, ks.safeAddKey(key)
}

func (ks *terra) Add(key terrakey.Key) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if _, found := ks.keyRing.Terra[key.ID()]; found {
		return fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return ks.safeAddKey(key)
}

func (ks *terra) Delete(id string) (terrakey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return terrakey.Key{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return terrakey.Key{}, err
	}
	err = ks.safeRemoveKey(key)
	return key, err
}

func (ks *terra) Import(keyJSON []byte, password string) (terrakey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return terrakey.Key{}, ErrLocked
	}
	key, err := terrakey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return terrakey.Key{}, errors.Wrap(err, "TerraKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.Terra[key.ID()]; found {
		return terrakey.Key{}, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return key, ks.keyManager.safeAddKey(key)
}

func (ks *terra) Export(id string, password string) ([]byte, error) {
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

func (ks *terra) EnsureKey() error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}

	if len(ks.keyRing.Terra) > 0 {
		return nil
	}

	key := terrakey.New()

	ks.logger.Infof("Created Terra key with ID %s", key.ID())

	return ks.safeAddKey(key)
}

var (
	ErrNoTerraKey = errors.New("no terra keys exist")
)

func (ks *terra) getByID(id string) (terrakey.Key, error) {
	key, found := ks.keyRing.Terra[id]
	if !found {
		return terrakey.Key{}, KeyNotFoundError{ID: id, KeyType: "Terra"}
	}
	return key, nil
}
