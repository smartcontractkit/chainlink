package keystore

import (
	"errors"
	"sync"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// memoryORM is an in-memory version of the keystore. This is
// only intended to be used in tests to avoid DB lock contention on
// the single DB row that stores the key material.
//
// Note: we store `q` on the struct since `saveEncryptedKeyRing` needs
// to support DB callbacks.
type memoryORM struct {
	keyRing *encryptedKeyRing
	q       pg.Queryer
	mu      sync.RWMutex
}

func (o *memoryORM) isEmpty() (bool, error) {
	return false, nil
}

func (o *memoryORM) saveEncryptedKeyRing(kr *encryptedKeyRing, callbacks ...func(pg.Queryer) error) (err error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.keyRing = kr
	for _, c := range callbacks {
		err = errors.Join(err, c(o.q))
	}
	return
}

func (o *memoryORM) getEncryptedKeyRing() (encryptedKeyRing, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if o.keyRing == nil {
		return encryptedKeyRing{}, nil
	}
	return *o.keyRing, nil
}

func newInMemoryORM(q pg.Queryer) *memoryORM {
	return &memoryORM{q: q}
}

// NewInMemory sets up a keystore which NOOPs attempts to access the `encrypted_key_rings` table. Accessing `evm.key_states`
// will still hit the DB.
func NewInMemory(db *sqlx.DB, scryptParams utils.ScryptParams, lggr logger.Logger, cfg pg.QConfig) *master {
	dbORM := NewORM(db, lggr, cfg)
	memoryORM := newInMemoryORM(dbORM.q)

	km := &keyManager{
		orm:          memoryORM,
		keystateORM:  dbORM,
		scryptParams: scryptParams,
		lock:         &sync.RWMutex{},
		logger:       lggr.Named("KeyStore"),
	}

	return &master{
		keyManager: km,
		cosmos:     newCosmosKeyStore(km),
		csa:        newCSAKeyStore(km),
		eth:        newEthKeyStore(km, dbORM, dbORM.q),
		ocr:        newOCRKeyStore(km),
		ocr2:       newOCR2KeyStore(km),
		p2p:        newP2PKeyStore(km),
		solana:     newSolanaKeyStore(km),
		starknet:   newStarkNetKeyStore(km),
		vrf:        newVRFKeyStore(km),
		dkgSign:    newDKGSignKeyStore(km),
		dkgEncrypt: newDKGEncryptKeyStore(km),
	}
}
