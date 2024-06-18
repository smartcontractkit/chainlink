package testutils

import (
	"database/sql"
	_ "embed"
	"fmt"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/scylladb/go-reflectx"
	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/store/migrate/plugins/relayer/evm"
)

var evmMu sync.RWMutex
var migratedDBs = map[string]*sqlx.DB{}

//go:embed evm_initial_state.sql
var evmInitialState string

// NewEVMDB creates a new EVM database with the given schema and chainId
func NewDB(t testing.TB, cfg evm.Cfg) *sqlx.DB {
	require.NotEmpty(t, evmInitialState, "evm initial state must not be empty")
	testutils.SkipShortDB(t)
	// fetch the db for the schema
	id := fmt.Sprintf("%s_%d", cfg.Schema, cfg.ChainID)
	evmMu.RLock()
	_, ok := migratedDBs[id]
	evmMu.RUnlock()
	if !ok {
		evmMu.Lock()
		defer evmMu.Unlock()
		// need to check again in case another goroutine has already migrated the db
		// while we were waiting for the write lock, which is more expensive than the optimistic read lock
		_, exists := migratedDBs[id]
		if !exists {
			c, evmHeavyDB := heavyweight.FullTestDBEmptyV2(t, nil)
			// run migrations to mutate the db
			// we have to setup the minimal tables for the migrations to work
			// must load the initial state, derivied from the core migrations at v244
			// because the evm migrations try to move try from the core schema to the evm schema
			//			b, err := os.ReadFile("../testdata/evm_initial_state.sql")
			//			require.NoError(t, err, "failed to read initial state for the evm migrations")
			_, err := evmHeavyDB.DB.Exec(evmInitialState)
			require.NoError(t, err, "failed to exec SQL for the initial state of the EVM database")
			// now we can run the migrations
			err = evm.Migrate(testutils.Context(t), evmHeavyDB.DB, cfg)
			require.NoError(t, err, "failed to migrate EVM database for cfg %v", cfg)
			migratedDBs[id] = evmHeavyDB
			url := c.Database().URL()
			sql.Register(id, pgtest.NewTxDriver(url.String()))
			sqlx.BindDriver(id, sqlx.DOLLAR)
		}
	}
	db, err := sqlx.Open(id, uuid.NewString())
	require.NoError(t, err, "failed to open EVM database for cfg %v with driver id %s", cfg, id)
	t.Cleanup(func() { assert.NoError(t, db.Close()) })
	db.MapperFunc(reflectx.CamelToSnakeASCII)

	return db
}
