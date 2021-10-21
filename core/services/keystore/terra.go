package keystore

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/terrakey"
)

//go:generate mockery --name Terra --output mocks/ --case=underscore

// Terra is the external interface for EthKeyStore
type Terra interface {
	Get(id string) (terrakey.KeyV2, error)
	GetAll() ([]terrakey.KeyV2, error)
	Create() (terrakey.KeyV2, error)
	Add(key terrakey.KeyV2) error
	Delete(id string) (terrakey.KeyV2, error)
	Import(keyJSON []byte, password string) (terrakey.KeyV2, error)
	Export(id string, password string) ([]byte, error)

	// SignTx(fromAddress common.Address, tx *types.Transaction) (*types.Transaction, error)

	GetV1KeysAsV2() ([]terrakey.KeyV2, error)
}

type terra struct {
	*keyManager
}

var _ Terra = terra{}

func newTerraKeyStore(km *keyManager) terra {
	return terra{
		km,
	}
}

func (ks terra) Get(id string) (terrakey.KeyV2, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return terrakey.KeyV2{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks terra) GetAll() (keys []terrakey.KeyV2, _ error) {
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

func (ks terra) Create() (terrakey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return terrakey.KeyV2{}, ErrLocked
	}
	key, err := terrakey.NewV2()
	if err != nil {
		return terrakey.KeyV2{}, err
	}
	return key, ks.safeAddKey(key)
}

func (ks terra) Add(key terrakey.KeyV2) error {
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

func (ks terra) Import(keyJSON []byte, password string) (terrakey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return terrakey.KeyV2{}, ErrLocked
	}
	key, err := terrakey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return terrakey.KeyV2{}, errors.Wrap(err, "EthKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.Terra[key.ID()]; found {
		return terrakey.KeyV2{}, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return key, ks.safeAddKey(key)
}

func (ks terra) Export(id string, password string) ([]byte, error) {
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

func (ks terra) Delete(id string) (terrakey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return terrakey.KeyV2{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return terrakey.KeyV2{}, err
	}
	err = ks.safeRemoveKey(key)
	return key, err
}

func (ks terra) SignTx(address common.Address, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	key, err := ks.getByID(address.Hex())
	if err != nil {
		return nil, err
	}
	signer := types.LatestSignerForChainID(chainID)
	//TODO
	fmt.Println(key)
	return types.SignTx(tx, signer, &ecdsa.PrivateKey{})
}

func (ks terra) GetV1KeysAsV2() (keys []terrakey.KeyV2, _ error) {
	v1Keys, err := ks.orm.GetEncryptedV1TerraKeys()
	if err != nil {
		return keys, err
	}
	for _, keyV1 := range v1Keys {
		err := keyV1.Unlock(ks.password)
		if err != nil {
			return keys, err
		}
		keys = append(keys, keyV1.ToV2())
	}
	return keys, nil
}

// caller must hold lock!
func (ks terra) getByID(id string) (terrakey.KeyV2, error) {
	key, found := ks.keyRing.Terra[id]
	if !found {
		return terrakey.KeyV2{}, fmt.Errorf("unable to find terra key with id %s", id)
	}
	return key, nil
}
