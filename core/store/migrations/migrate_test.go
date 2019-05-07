package migrations_test

import (
	"os"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func bootstrapORM(t *testing.T) (*orm.ORM, func()) {
	tc, cleanup := cltest.NewConfig()
	config := tc.Config

	require.NoError(t, os.MkdirAll(config.RootDir(), 0700))
	cltest.WipePostgresDatabase(tc.Config)

	orm, err := orm.NewORM(config.NormalizedDatabaseURL(), config.DatabaseTimeout())
	require.NoError(t, err)

	return orm, func() {
		assert.NoError(t, orm.Close())
		cleanup()
		os.RemoveAll(config.RootDir())
	}
}

func TestMigrate_Migration0(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	db := orm.DB

	require.NoError(t, migration0.Migrate(db))

	assert.True(t, db.HasTable("bridge_types"))
	assert.True(t, db.HasTable("encumbrances"))
	assert.True(t, db.HasTable("external_initiators"))
	assert.True(t, db.HasTable("heads"))
	assert.True(t, db.HasTable("job_specs"))
	assert.True(t, db.HasTable("initiators"))
	assert.True(t, db.HasTable("job_runs"))
	assert.True(t, db.HasTable("keys"))
	assert.True(t, db.HasTable("run_requests"))
	assert.True(t, db.HasTable("run_results"))
	assert.True(t, db.HasTable("service_agreements"))
	assert.True(t, db.HasTable("sessions"))
	assert.True(t, db.HasTable("task_runs"))
	assert.True(t, db.HasTable("task_specs"))
	assert.True(t, db.HasTable("tx_attempts"))
	assert.True(t, db.HasTable("txes"))
	assert.True(t, db.HasTable("users"))
}
