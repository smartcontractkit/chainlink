package store_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/migrations"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_SquashMigrationUpgrade(t *testing.T) {
	_, orm, cleanup := heavyweight.FullTestORM(t, "migrationssquash", false)
	defer cleanup()
	db := orm.DB

	// Latest migrations should work fine.
	static.Version = "0.9.11"
	err := migrations.MigrateUp(db, "")
	require.NoError(t, err)
	err = store.CheckSquashUpgrade(db)
	require.NoError(t, err)
	static.Version = "unset"
}

func TestStore_Start(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	assert.NoError(t, store.Start())
}

func TestStore_Close(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore(t)
	defer cleanup()

	assert.NoError(t, s.Close())
}
