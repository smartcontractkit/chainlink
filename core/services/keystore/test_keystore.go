package keystore

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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
	ds      sqlutil.DataSource
	mu      sync.RWMutex
}

func (o *memoryORM) isEmpty(ctx context.Context) (bool, error) {
	return false, nil
}

func (o *memoryORM) saveEncryptedKeyRing(ctx context.Context, kr *encryptedKeyRing, callbacks ...func(sqlutil.DataSource) error) (err error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.keyRing = kr
	for _, c := range callbacks {
		err = errors.Join(err, c(o.ds))
	}
	return
}

func (o *memoryORM) getEncryptedKeyRing(ctx context.Context) (encryptedKeyRing, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if o.keyRing == nil {
		return encryptedKeyRing{}, nil
	}
	return *o.keyRing, nil
}

func newInMemoryORM(ds sqlutil.DataSource) *memoryORM {
	return &memoryORM{ds: ds}
}

// NewInMemory sets up a keystore which NOOPs attempts to access the `encrypted_key_rings` table. Accessing `evm.key_states`
// will still hit the DB.
func NewInMemory(ds sqlutil.DataSource, scryptParams utils.ScryptParams, lggr logger.Logger) *master {
	dbORM := NewORM(ds, lggr)
	memoryORM := newInMemoryORM(ds)

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
		eth:        newEthKeyStore(km, dbORM, ds),
		ocr:        newOCRKeyStore(km),
		ocr2:       newOCR2KeyStore(km),
		p2p:        newP2PKeyStore(km),
		solana:     newSolanaKeyStore(km),
		starknet:   newStarkNetKeyStore(km),
		aptos:      newAptosKeyStore(km),
		tron:       newTronKeyStore(km),
		vrf:        newVRFKeyStore(km),
		workflow:   newWorkflowKeyStore(km),
	}
}

// TestKeystore is a test keystore that wraps the master keystore
// and provides a helper function to generate keys for testing.
type TestKeystore struct {
	t *testing.T
	Master
}

func NewTestKeyStore(t *testing.T) *TestKeystore {
	t.Helper()
	db := pgtest.NewSqlxDB(t)
	ks := New(db, utils.FastScryptParams, logger.TestLogger(t))

	err := ks.Unlock(tests.Context(t), "password")
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Close()
	})
	return &TestKeystore{
		t:      t,
		Master: ks,
	}

}

// ImportableKey is a struct that holds the JSON representation of a key. It echos the core TOML secret [toml.ImportableKey].
type ImportableKey struct {
	JSON     string // JSON representation of the key; the format depends on the key type
	Password string // Password used to encrypt the key
}

// ImportableEthKey is a struct that holds the JSON representation of an Ethereum key. It echos the core TOML secret [toml.ImportableEthKey].
type ImportableEthKey struct {
	EVMChainID uint64 // Chain ID for the Ethereum key. NOT the chain selector
	ImportableKey
}

func (ks *TestKeystore) GenerateEthKeys(chainIDs ...*big.Int) []ImportableEthKey {
	t := ks.t
	ctx := tests.Context(t)
	out := make([]ImportableEthKey, len(chainIDs))
	for i, chainID := range chainIDs {
		k, err := ks.Eth().Create(ctx, chainID)
		require.NoError(t, err)
		json, err := ks.Eth().Export(ctx, k.ID(), "password")
		require.NoError(t, err)
		out[i] = ImportableEthKey{
			EVMChainID: chainID.Uint64(),
			ImportableKey: ImportableKey{
				JSON:     string(json),
				Password: "password",
			},
		}
	}

	return out
}

func (ks *TestKeystore) GenerateP2PKey() ImportableKey {
	t := ks.t
	ctx := tests.Context(t)
	k, err := ks.P2P().Create(ctx)
	require.NoError(t, err)
	json, err := ks.P2P().Export(k.PeerID(), "password")
	require.NoError(t, err)
	return ImportableKey{
		JSON:     string(json),
		Password: "password",
	}
}
