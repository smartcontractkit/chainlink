package testutils

import (
	"database/sql"
	_ "embed"
	"fmt"
	"net/url"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/scylladb/go-reflectx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
	"github.com/smartcontractkit/chainlink/v2/core/store/migrate/plugins/relayer/evm"
	"github.com/smartcontractkit/chainlink/v2/internal/testdb"
)

var evmMu sync.RWMutex
var migratedDBs = map[string]*sqlx.DB{}
var evmHeavyDB *sqlx.DB
var testDBURL string

var initTestDB = sync.OnceFunc(newTestDB)

// func newTestDB() {
func newTestDB() {
	// hack to get the migrations to run
	eurl := string(env.DatabaseURL.Get())
	if eurl == "" {
		panic("you must provide a CL_DATABASE_URL environment variable")
	}

	uuid := uuid.NewString()
	durl, err := url.Parse(eurl)
	if err != nil {
		panic(fmt.Sprintf("failed to parse database URL '%s': %v", eurl, err))
	}
	testDBURL, err = testdb.CreateOrReplace(*durl, "evm_hacking_test_db_"+uuid[:8], false)
	if err != nil {
		panic(fmt.Sprintf("failed to create or replace database for EVM tests: %v", err))
	}
	evmHeavyDB = sqlx.MustOpen(string(dialects.Postgres), testDBURL)
	if evmHeavyDB == nil {
		panic("failed to open EVM heavy database")
	}
	_, err = evmHeavyDB.DB.Exec(evmInitialState)
	if err != nil {
		panic(fmt.Sprintf("failed to exec SQL for the initial state of the EVM database: %v", err))
	}

}

//go:embed evm_initial_state.sql
var evmInitialState string

// NewEVMDB creates a new EVM database with the given schema and chainId
func NewDB(t testing.TB, cfg evm.Cfg) *sqlx.DB {
	testutils.SkipShortDB(t)
	initTestDB()

	// fetch the db for the schema
	id := fmt.Sprintf("%s_%s", cfg.Schema, cfg.ChainID.String())
	evmMu.RLock()
	_, ok := migratedDBs[id]
	evmMu.RUnlock()
	if !ok {
		evmMu.Lock()
		defer evmMu.Unlock()
		// need to check again in case another goroutine has already migrated the db
		// while we were waiting for the write lock, which is more expensive than the optimistic read lock
		_, exists := migratedDBs[id]
		//var
		if !exists {
			err := evm.Migrate(testutils.Context(t), evmHeavyDB, cfg)
			require.NoError(t, err, "failed to migrate EVM database for cfg %v", cfg)
			migratedDBs[id] = evmHeavyDB

			sql.Register(id, pgtest.NewTxDriver(testDBURL))
			sqlx.BindDriver(id, sqlx.DOLLAR)
		}
	}
	db, err := sqlx.Open(id, uuid.NewString())
	require.NoError(t, err, "failed to open EVM database for cfg %v with driver id %s", cfg, id)
	t.Cleanup(func() {
		assert.NoError(t, db.Close())
		assert.NoError(t, evmHeavyDB.Close())
	})
	db.MapperFunc(reflectx.CamelToSnakeASCII)

	return db
}
