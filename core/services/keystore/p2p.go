package keystore

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"gorm.io/gorm"
)

type P2P interface {
	Get(id string) (p2pkey.KeyV2, error)
	GetAll() ([]p2pkey.KeyV2, error)
	Create() (p2pkey.KeyV2, error)
	Add(key p2pkey.KeyV2) error
	Delete(id string) (p2pkey.KeyV2, error)
	Import(keyJSON []byte, password string) (p2pkey.KeyV2, error)
	Export(id string, password string) ([]byte, error)
	EnsureKey() (p2pkey.KeyV2, bool, error)

	GetV1KeysAsV2() ([]p2pkey.KeyV2, error)

	GetOrFirst(id string) (p2pkey.KeyV2, error)
}

type p2p struct {
	*keyManager
}

var _ P2P = &p2p{}

func newP2PKeyStore(km *keyManager) *p2p {
	return &p2p{
		km,
	}
}

func (ks *p2p) Get(id string) (p2pkey.KeyV2, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return p2pkey.KeyV2{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *p2p) GetAll() (keys []p2pkey.KeyV2, _ error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, key := range ks.keyRing.P2P {
		keys = append(keys, key)
	}
	return keys, nil
}

func (ks *p2p) Create() (p2pkey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return p2pkey.KeyV2{}, ErrLocked
	}
	key, err := p2pkey.NewV2()
	if err != nil {
		return p2pkey.KeyV2{}, err
	}
	return key, ks.safeAddKey(key)
}

func (ks *p2p) Add(key p2pkey.KeyV2) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if _, found := ks.keyRing.P2P[key.ID()]; found {
		return fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return ks.safeAddKey(key)
}

func (ks *p2p) Delete(id string) (p2pkey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return p2pkey.KeyV2{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return p2pkey.KeyV2{}, err
	}
	err = ks.safeRemoveKey(key, func(db *gorm.DB) error {
		return db.Exec(`DELETE FROM p2p_peers WHERE peer_id = ?`, key.ID()).Error
	})
	return key, err
}

func (ks *p2p) Import(keyJSON []byte, password string) (p2pkey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return p2pkey.KeyV2{}, ErrLocked
	}
	key, err := p2pkey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return p2pkey.KeyV2{}, errors.Wrap(err, "P2PKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.P2P[key.ID()]; found {
		return p2pkey.KeyV2{}, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return key, ks.keyManager.safeAddKey(key)
}

func (ks *p2p) Export(id string, password string) ([]byte, error) {
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

func (ks *p2p) EnsureKey() (p2pkey.KeyV2, bool, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return p2pkey.KeyV2{}, false, ErrLocked
	}
	if len(ks.keyRing.P2P) > 0 {
		return p2pkey.KeyV2{}, true, nil
	}
	key, err := p2pkey.NewV2()
	if err != nil {
		return p2pkey.KeyV2{}, false, err
	}
	return key, false, ks.safeAddKey(key)
}

func (ks *p2p) GetV1KeysAsV2() (keys []p2pkey.KeyV2, _ error) {
	v1Keys, err := ks.orm.GetEncryptedV1P2PKeys()
	if err != nil {
		return keys, err
	}
	for _, keyV1 := range v1Keys {
		pk, err := keyV1.Decrypt(ks.password)
		if err != nil {
			return keys, err
		}
		keys = append(keys, pk.ToV2())
	}
	return keys, nil
}

func (ks *p2p) GetOrFirst(id string) (p2pkey.KeyV2, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return p2pkey.KeyV2{}, ErrLocked
	}
	if id != "" {
		return ks.getByID(id)
	} else if len(ks.keyRing.P2P) == 1 {
		logger.Warn("No P2P_PEER_ID set, defaulting to first key in database")
		for _, key := range ks.keyRing.P2P {
			return key, nil
		}
	} else if len(ks.keyRing.P2P) == 0 {
		return p2pkey.KeyV2{}, errors.New("no p2p keys exist")
	}
	return p2pkey.KeyV2{}, errors.New(
		"multiple p2p keys found but peer ID was not set - you must specify a P2P_PEER_ID" +
			"env var or set the peer id in the job spec if you have more than one key",
	)
}

func (ks *p2p) getByID(id string) (p2pkey.KeyV2, error) {
	key, found := ks.keyRing.P2P[id]
	if !found {
		return p2pkey.KeyV2{}, fmt.Errorf("unable to find P2P key with id %s", id)
	}
	return key, nil
}
