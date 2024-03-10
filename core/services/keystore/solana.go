package keystore

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/solkey"
)

//go:generate mockery --quiet --name Solana --output ./mocks/ --case=underscore --filename solana.go

type Solana interface {
	Get(id string) (solkey.Key, error)
	GetAll() ([]solkey.Key, error)
	Create() (solkey.Key, error)
	Add(key solkey.Key) error
	Delete(id string) (solkey.Key, error)
	Import(keyJSON []byte, password string) (solkey.Key, error)
	Export(id string, password string) ([]byte, error)
	EnsureKey() error
	Sign(ctx context.Context, id string, msg []byte) (signature []byte, err error)
}

// SolanaSigner adapts Solana to [loop.Keystore].
type SolanaSigner struct {
	Solana
}

func (s *SolanaSigner) Accounts(ctx context.Context) (accounts []string, err error) {
	ks, err := s.GetAll()
	if err != nil {
		return nil, err
	}
	for _, k := range ks {
		accounts = append(accounts, k.PublicKeyStr())
	}
	return
}

type solana struct {
	*keyManager
}

var _ Solana = &solana{}

func newSolanaKeyStore(km *keyManager) *solana {
	return &solana{
		km,
	}
}

func (ks *solana) Get(id string) (solkey.Key, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return solkey.Key{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *solana) GetAll() (keys []solkey.Key, _ error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, key := range ks.keyRing.Solana {
		keys = append(keys, key)
	}
	return keys, nil
}

func (ks *solana) Create() (solkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return solkey.Key{}, ErrLocked
	}
	key, err := solkey.New()
	if err != nil {
		return solkey.Key{}, err
	}
	return key, ks.safeAddKey(key)
}

func (ks *solana) Add(key solkey.Key) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if _, found := ks.keyRing.Solana[key.ID()]; found {
		return fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return ks.safeAddKey(key)
}

func (ks *solana) Delete(id string) (solkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return solkey.Key{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return solkey.Key{}, err
	}
	err = ks.safeRemoveKey(key)
	return key, err
}

func (ks *solana) Import(keyJSON []byte, password string) (solkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return solkey.Key{}, ErrLocked
	}
	key, err := solkey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return solkey.Key{}, errors.Wrap(err, "SolanaKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.Solana[key.ID()]; found {
		return solkey.Key{}, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return key, ks.keyManager.safeAddKey(key)
}

func (ks *solana) Export(id string, password string) ([]byte, error) {
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

func (ks *solana) EnsureKey() error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if len(ks.keyRing.Solana) > 0 {
		return nil
	}

	key, err := solkey.New()
	if err != nil {
		return err
	}

	ks.logger.Infof("Created Solana key with ID %s", key.ID())

	return ks.safeAddKey(key)
}

func (ks *solana) Sign(_ context.Context, id string, msg []byte) (signature []byte, err error) {
	k, err := ks.Get(id)
	if err != nil {
		return nil, err
	}
	return k.Sign(msg)
}

func (ks *solana) getByID(id string) (solkey.Key, error) {
	key, found := ks.keyRing.Solana[id]
	if !found {
		return solkey.Key{}, KeyNotFoundError{ID: id, KeyType: "Solana"}
	}
	return key, nil
}
