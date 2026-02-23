package keystore

import (
	"context"
	"errors"
	"sync"

	"github.com/smartcontractkit/chainlink-common/keystore"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/models"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

// memoryORM is an in-memory version of the keystore. This is
// only intended to be used in tests to avoid DB lock contention on
// the single DB row that stores the key material.
//
// Note: we store `q` on the struct since `saveEncryptedKeyRing` needs
// to support DB callbacks.
type memoryORM struct {
	keyRing *models.EncryptedKeyRing
	ds      sqlutil.DataSource
	mu      sync.RWMutex
}

func (o *memoryORM) isEmpty(ctx context.Context) (bool, error) {
	return false, nil
}

func (o *memoryORM) saveEncryptedKeyRing(ctx context.Context, kr *models.EncryptedKeyRing, callbacks ...func(sqlutil.DataSource) error) (err error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.keyRing = kr
	for _, c := range callbacks {
		err = errors.Join(err, c(o.ds))
	}
	return
}

func (o *memoryORM) getEncryptedKeyRing(ctx context.Context) (models.EncryptedKeyRing, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if o.keyRing == nil {
		return models.EncryptedKeyRing{}, nil
	}
	return *o.keyRing, nil
}

func newInMemoryORM(ds sqlutil.DataSource) *memoryORM {
	return &memoryORM{ds: ds}
}

// NewInMemory sets up a keystore which NOOPs attempts to access the `encrypted_key_rings` table. Accessing `evm.key_states`
// will still hit the DB.
func NewInMemory(ds sqlutil.DataSource, scryptParams keystore.ScryptParams, logf Logf) *master {
	dbORM := NewORM(ds)
	memoryORM := newInMemoryORM(ds)

	km := &keyManager{
		orm:          memoryORM,
		keystateORM:  dbORM,
		scryptParams: scryptParams,
		lock:         &sync.RWMutex{},
		announce:     announcer(logf),
	}

	return &master{
		keyManager:   km,
		cosmos:       newCosmosKeyStore(km),
		csa:          newCSAKeyStore(km),
		eth:          newEthKeyStore(km, dbORM, ds),
		ocr:          newOCRKeyStore(km),
		ocr2:         newOCR2KeyStore(km),
		p2p:          newP2PKeyStore(km),
		solana:       newSolanaKeyStore(km),
		starknet:     newStarkNetKeyStore(km),
		sui:          newSuiKeyStore(km),
		aptos:        newAptosKeyStore(km),
		tron:         newTronKeyStore(km),
		ton:          newTONKeyStore(km),
		vrf:          newVRFKeyStore(km),
		workflow:     newWorkflowKeyStore(km),
		dkgRecipient: newDKGRecipientKeyStore(km),
	}
}
