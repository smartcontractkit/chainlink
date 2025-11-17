package keystore

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/suikey"
)

// Sui is the interface for the Sui keystore
type Sui interface {
	Get(id string) (suikey.Key, error)
	GetAll() ([]suikey.Key, error)
	Create(ctx context.Context) (suikey.Key, error)
	Add(ctx context.Context, key suikey.Key) error
	Delete(ctx context.Context, id string) (suikey.Key, error)
	Import(ctx context.Context, keyJSON []byte, password string) (suikey.Key, error)
	Export(id string, password string) ([]byte, error)
	EnsureKey(ctx context.Context) error
	Sign(ctx context.Context, id string, msg []byte) ([]byte, error)
}

type sui struct {
	*keyManager
}

var _ Sui = &sui{}

func newSuiKeyStore(km *keyManager) *sui {
	return &sui{
		keyManager: km,
	}
}

func (ks *sui) Get(id string) (suikey.Key, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return suikey.Key{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *sui) GetAll() ([]suikey.Key, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	accounts := []suikey.Key{}
	for _, key := range ks.keyRing.Sui {
		accounts = append(accounts, key)
	}
	return accounts, nil
}

func (ks *sui) Create(ctx context.Context) (suikey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return suikey.Key{}, ErrLocked
	}
	key, err := suikey.New()
	if err != nil {
		return suikey.Key{}, err
	}
	return key, ks.safeAddKey(ctx, key)
}

func (ks *sui) Add(ctx context.Context, key suikey.Key) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if _, found := ks.keyRing.Sui[key.ID()]; found {
		return fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return ks.safeAddKey(ctx, key)
}

func (ks *sui) Delete(ctx context.Context, id string) (suikey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return suikey.Key{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return suikey.Key{}, err
	}
	err = ks.safeRemoveKey(ctx, key)
	return key, err
}

func (ks *sui) Import(ctx context.Context, keyJSON []byte, password string) (suikey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return suikey.Key{}, ErrLocked
	}
	key, err := suikey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return suikey.Key{}, err
	}
	if _, found := ks.keyRing.Sui[key.ID()]; found {
		return suikey.Key{}, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	err = ks.safeAddKey(ctx, key)
	return key, err
}

func (ks *sui) Export(id string, password string) ([]byte, error) {
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

func (ks *sui) EnsureKey(ctx context.Context) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if len(ks.keyRing.Sui) > 0 {
		return nil
	}

	key, err := suikey.New()
	if err != nil {
		return err
	}

	// ks.logger.Infof("Created Sui key with ID %s", key.ID())

	return ks.safeAddKey(ctx, key)
}

func (ks *sui) Sign(_ context.Context, id string, msg []byte) ([]byte, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return nil, err
	}
	return key.Sign(msg)
}

func (ks *sui) getByID(id string) (suikey.Key, error) {
	key, found := ks.keyRing.Sui[id]
	if !found {
		return suikey.Key{}, KeyNotFoundError{ID: id, KeyType: "Sui"}
	}
	return key, nil
}

// TODO: the approach below is deprecated, replace it
type SuiLoopSinger struct {
	Sui
	core.UnimplementedKeystore
}

var _ loop.Keystore = &SuiLoopSinger{}

// Returns a list of Sui Public Keys
func (s *SuiLoopSinger) Accounts(ctx context.Context) (accounts []string, err error) {
	ks, err := s.GetAll()
	if err != nil {
		return nil, err
	}
	for _, k := range ks {
		accounts = append(accounts, k.ID())
	}
	return
}

func (s *SuiLoopSinger) Sign(ctx context.Context, account string, data []byte) (signed []byte, err error) {
	return s.Sui.Sign(ctx, account, data)
}
