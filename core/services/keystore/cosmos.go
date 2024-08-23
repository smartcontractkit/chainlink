package keystore

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/cosmoskey"
)

type Cosmos interface {
	Get(id string) (cosmoskey.Key, error)
	GetAll() ([]cosmoskey.Key, error)
	Create(ctx context.Context) (cosmoskey.Key, error)
	Add(ctx context.Context, key cosmoskey.Key) error
	Delete(ctx context.Context, id string) (cosmoskey.Key, error)
	Import(ctx context.Context, keyJSON []byte, password string) (cosmoskey.Key, error)
	Export(id string, password string) ([]byte, error)
	EnsureKey(ctx context.Context) error
}

type cosmos struct {
	*keyManager
}

var _ Cosmos = &cosmos{}

func newCosmosKeyStore(km *keyManager) *cosmos {
	return &cosmos{
		km,
	}
}

func (ks *cosmos) Get(id string) (cosmoskey.Key, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return cosmoskey.Key{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *cosmos) GetAll() (keys []cosmoskey.Key, _ error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, key := range ks.keyRing.Cosmos {
		keys = append(keys, key)
	}
	return keys, nil
}

func (ks *cosmos) Create(ctx context.Context) (cosmoskey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return cosmoskey.Key{}, ErrLocked
	}
	key := cosmoskey.New()
	return key, ks.safeAddKey(ctx, key)
}

func (ks *cosmos) Add(ctx context.Context, key cosmoskey.Key) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if _, found := ks.keyRing.Cosmos[key.ID()]; found {
		return fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return ks.safeAddKey(ctx, key)
}

func (ks *cosmos) Delete(ctx context.Context, id string) (cosmoskey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return cosmoskey.Key{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return cosmoskey.Key{}, err
	}
	err = ks.safeRemoveKey(ctx, key)
	return key, err
}

func (ks *cosmos) Import(ctx context.Context, keyJSON []byte, password string) (cosmoskey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return cosmoskey.Key{}, ErrLocked
	}
	key, err := cosmoskey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return cosmoskey.Key{}, errors.Wrap(err, "CosmosKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.Cosmos[key.ID()]; found {
		return cosmoskey.Key{}, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return key, ks.keyManager.safeAddKey(ctx, key)
}

func (ks *cosmos) Export(id string, password string) ([]byte, error) {
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

func (ks *cosmos) EnsureKey(ctx context.Context) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}

	if len(ks.keyRing.Cosmos) > 0 {
		return nil
	}

	key := cosmoskey.New()

	ks.logger.Infof("Created Cosmos key with ID %s", key.ID())

	return ks.safeAddKey(ctx, key)
}

func (ks *cosmos) getByID(id string) (cosmoskey.Key, error) {
	key, found := ks.keyRing.Cosmos[id]
	if !found {
		return cosmoskey.Key{}, KeyNotFoundError{ID: id, KeyType: "Cosmos"}
	}
	return key, nil
}

// CosmosLoopKeystore implements the [github.com/smartcontractkit/chainlink-common/pkg/loop.Keystore] interface and
// handles signing for Cosmos messages.
type CosmosLoopKeystore struct {
	Cosmos
}

var _ loop.Keystore = &CosmosLoopKeystore{}

func (lk *CosmosLoopKeystore) Sign(ctx context.Context, id string, hash []byte) ([]byte, error) {
	k, err := lk.Get(id)
	if err != nil {
		return nil, err
	}
	// loopp spec requires passing nil hash to check existence of id
	if hash == nil {
		return nil, nil
	}

	return k.ToPrivKey().Sign(hash)
}

func (lk *CosmosLoopKeystore) Accounts(ctx context.Context) ([]string, error) {
	keys, err := lk.GetAll()
	if err != nil {
		return nil, err
	}

	accounts := []string{}
	for _, k := range keys {
		accounts = append(accounts, k.PublicKeyStr())
	}

	return accounts, nil
}
