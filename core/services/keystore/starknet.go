package keystore

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo"
	"github.com/pkg/errors"

	stark "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/keys"
	pkgstarknet "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/starknet"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/starkkey"
)

//go:generate mockery --quiet --name StarkNet --output ./mocks/ --case=underscore --filename starknet.go

type StarkNet interface {
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

var _ StarkNet = &starknet{}

func newStarkNetKeyStore(km *keyManager) *starknet {
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
	for _, key := range ks.keyRing.StarkNet {
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
	if _, found := ks.keyRing.StarkNet[key.ID()]; found {
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
		return stark.Key{}, errors.Wrap(err, "StarkNetKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.StarkNet[key.ID()]; found {
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

func (ks *starknet) getByID(id string) (stark.Key, error) {
	key, found := ks.keyRing.StarkNet[id]
	if !found {
		return stark.Key{}, KeyNotFoundError{ID: id, KeyType: "StarkNet"}
	}
	return key, nil
}

// StarkSigner adapts Starknet to [loop.Keystore].
type StarkNetSigner struct {
	StarkNet
}

func (s *StarkNetSigner) Accounts(ctx context.Context) ([]string, error) {
	var accounts []string
	ks, err := s.GetAll()
	if err != nil {
		return nil, err
	}
	for _, k := range ks {
		accounts = append(accounts, k.StarkKeyStr())
	}
	return accounts, nil
}

const byteLen = 32

func (s *StarkNetSigner) Sign(_ context.Context, id string, msg []byte) ([]byte, error) {
	k, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	// big.Int requires big endian
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, msg)
	if err != nil {
		return nil, fmt.Errorf("error transcoding input message to big endian: %w", err)
	}

	hsh, err := caigo.Curve.ComputeHashOnElements([]*big.Int{
		new(big.Int).SetBytes(buf.Bytes()),
	})
	if err != nil {
		return nil, fmt.Errorf("error computing hash: %w", err)
	}
	x, y, err := caigo.Curve.Sign(hsh, k.ToPrivKey())
	if err != nil {
		return nil, fmt.Errorf("error signing hash: %w", err)
	}
	result := new(bytes.Buffer)
	n, err := result.Write(pkgstarknet.PadBytes(x.Bytes(), byteLen))
	if err != nil {
		return nil, fmt.Errorf("error encoding 'x' component of signature: %w", err)
	}
	if n != byteLen {
		return nil, fmt.Errorf("truncated write of 'x' component of signature. wrote %d expected %d", n, byteLen)
	}
	n, err = result.Write(pkgstarknet.PadBytes(y.Bytes(), byteLen))
	if err != nil {
		return nil, fmt.Errorf("error encoding 'x' component of signature: %w", err)
	}
	if n != byteLen {
		return nil, fmt.Errorf("truncated write of 'x' component of signature. wrote %d expected %d", n, byteLen)
	}
	return result.Bytes(), nil
}
