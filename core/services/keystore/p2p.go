package keystore

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type P2P interface {
	Get(id p2pkey.PeerID) (p2pkey.KeyV2, error)
	GetAll() ([]p2pkey.KeyV2, error)
	Create(ctx context.Context) (p2pkey.KeyV2, error)
	Add(ctx context.Context, key p2pkey.KeyV2) error
	Delete(ctx context.Context, id p2pkey.PeerID) (p2pkey.KeyV2, error)
	Import(ctx context.Context, keyJSON []byte, password string) (p2pkey.KeyV2, error)
	Export(id p2pkey.PeerID, password string) ([]byte, error)
	EnsureKey(ctx context.Context) error

	GetOrFirst(id p2pkey.PeerID) (p2pkey.KeyV2, error)
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

func (ks *p2p) Get(id p2pkey.PeerID) (p2pkey.KeyV2, error) {
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

func (ks *p2p) Create(ctx context.Context) (p2pkey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return p2pkey.KeyV2{}, ErrLocked
	}
	key, err := p2pkey.NewV2()
	if err != nil {
		return p2pkey.KeyV2{}, err
	}
	return key, ks.safeAddKey(ctx, key)
}

func (ks *p2p) Add(ctx context.Context, key p2pkey.KeyV2) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if _, found := ks.keyRing.P2P[key.ID()]; found {
		return fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return ks.safeAddKey(ctx, key)
}

func (ks *p2p) Delete(ctx context.Context, id p2pkey.PeerID) (p2pkey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return p2pkey.KeyV2{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return p2pkey.KeyV2{}, err
	}
	err = ks.safeRemoveKey(ctx, key, func(ds sqlutil.DataSource) error {
		_, err2 := ds.ExecContext(ctx, `DELETE FROM p2p_peers WHERE peer_id = $1`, key.ID())
		return err2
	})
	return key, err
}

func (ks *p2p) Import(ctx context.Context, keyJSON []byte, password string) (p2pkey.KeyV2, error) {
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
	return key, ks.keyManager.safeAddKey(ctx, key)
}

func (ks *p2p) Export(id p2pkey.PeerID, password string) ([]byte, error) {
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

func (ks *p2p) EnsureKey(ctx context.Context) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}

	if len(ks.keyRing.P2P) > 0 {
		return nil
	}

	key, err := p2pkey.NewV2()
	if err != nil {
		return err
	}

	ks.logger.Infof("Created P2P key with ID %s", key.ID())

	return ks.safeAddKey(ctx, key)
}

var (
	ErrNoP2PKey = errors.New("no p2p keys exist")
)

func (ks *p2p) GetOrFirst(id p2pkey.PeerID) (p2pkey.KeyV2, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return p2pkey.KeyV2{}, ErrLocked
	}
	if id != (p2pkey.PeerID{}) {
		return ks.getByID(id)
	} else if len(ks.keyRing.P2P) == 1 {
		ks.logger.Warn("No P2P.PeerID set, defaulting to first key in database")
		for _, key := range ks.keyRing.P2P {
			return key, nil
		}
	} else if len(ks.keyRing.P2P) == 0 {
		return p2pkey.KeyV2{}, ErrNoP2PKey
	}
	possibleKeys := make([]string, 0, len(ks.keyRing.P2P))
	for _, key := range ks.keyRing.P2P {
		possibleKeys = append(possibleKeys, key.ID())
	}
	// To avoid ambiguity, we require the user to specify a peer ID if there are multiple keys
	return p2pkey.KeyV2{}, errors.New(
		"multiple p2p keys found but peer ID was not set - you must specify a P2P.PeerID " +
			"config var if you have more than one key, or delete the keys you aren't using" +
			" (possible keys: " + strings.Join(possibleKeys, ", ") + ")",
	)
}

func (ks *p2p) getByID(id p2pkey.PeerID) (p2pkey.KeyV2, error) {
	key, found := ks.keyRing.P2P[id.Raw()]
	if !found {
		return p2pkey.KeyV2{}, KeyNotFoundError{ID: id.String(), KeyType: "P2P"}
	}
	return key, nil
}
