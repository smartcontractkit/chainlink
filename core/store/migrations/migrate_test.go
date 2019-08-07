package migrations_test

import (
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560924400"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1565210496"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/smartcontractkit/chainlink/core/store/assets"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/migrations"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1559081901"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1559767166"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560433987"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560791143"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560881846"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560881855"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560886530"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1565139192"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormigrate "gopkg.in/gormigrate.v1"
)

func bootstrapORM(t *testing.T) (*orm.ORM, func()) {
	tc, cleanup := cltest.NewConfig(t)
	config := tc.Config

	require.NoError(t, os.MkdirAll(config.RootDir(), 0700))
	cltest.WipePostgresDatabase(t, tc.Config)

	orm, err := orm.NewORM(orm.NormalizedDatabaseURL(config), config.DatabaseTimeout())
	require.NoError(t, err)

	return orm, func() {
		assert.NoError(t, orm.Close())
		cleanup()
		os.RemoveAll(config.RootDir())
	}
}

func TestMigrate_Migrations(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	db := orm.DB

	require.NoError(t, migration0.Migrate(db))
	require.NoError(t, migration1559081901.Migrate(db))

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

func TestMigrate_Migration1560791143(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	db := orm.DB

	require.NoError(t, migration0.Migrate(db))

	tx := migration0.Tx{
		ID:       1337,
		Data:     make([]byte, 10),
		Value:    models.NewBig(big.NewInt(1)),
		GasPrice: models.NewBig(big.NewInt(127)),
	}
	require.NoError(t, db.Create(&tx).Error)

	require.NoError(t, migration1559081901.Migrate(db))

	txFound := models.Tx{}
	require.NoError(t, db.Where("id = ?", tx.ID).Find(&txFound).Error)

	require.NoError(t, migration1560791143.Migrate(db))

	txNoID := models.Tx{
		Data:     make([]byte, 10),
		Value:    models.NewBig(big.NewInt(2)),
		GasPrice: models.NewBig(big.NewInt(119)),
	}
	require.NoError(t, db.Create(&txNoID).Error)
	assert.Equal(t, uint64(1338), txNoID.ID)

	noIDTxFound := models.Tx{}
	require.NoError(t, db.Where("id = ?", tx.ID).Find(&noIDTxFound).Error)
}

func TestMigrate_Migration1560881855(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	db := orm.DB

	require.NoError(t, migration0.Migrate(db))
	require.NoError(t, migration1559081901.Migrate(db))
	require.NoError(t, migration1559767166.Migrate(db))
	require.NoError(t, migration1560433987.Migrate(db))
	require.NoError(t, migration1560791143.Migrate(db))
	require.NoError(t, migration1560881846.Migrate(db))
	require.NoError(t, migration1560886530.Migrate(db))
	require.NoError(t, migration1560924400.Migrate(db))
	require.NoError(t, migration1565139192.Migrate(db))

	j := models.NewJob()
	i := models.Initiator{Type: models.InitiatorWeb}
	j.Initiators = []models.Initiator{i}
	j.Tasks = []models.TaskSpec{
		cltest.NewTask(t, "noop"),
	}
	assert.NoError(t, db.Create(&j).Error)
	initr := j.Initiators[0]
	jr := j.NewRun(initr)
	data := `{"result":"921.02"}`
	jr.Result = cltest.RunResultWithData(data)
	jr.Overrides.Amount = assets.NewLink(2)
	befCreation := time.Now()
	require.NoError(t, db.Create(&jr).Error)
	aftCreation := time.Now()

	// placement of this migration is important, as it makes sure backfilling
	//  is done if there's already a RunResult with nonzero link reward
	require.NoError(t, migration1560881855.Migrate(db))

	rewFound := models.LinkEarned{}
	require.NoError(t, db.Find(&rewFound).Error)
	assert.Equal(t, j.ID, rewFound.JobSpecID)
	assert.Equal(t, jr.ID, rewFound.JobRunID)
	assert.Equal(t, assets.NewLink(2), rewFound.Earned)
	assert.True(t, true, rewFound.EarnedAt.After(aftCreation), rewFound.EarnedAt.Before(befCreation))
}

func TestMigrate_Migration1560881846(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	db := orm.DB

	require.NoError(t, migration0.Migrate(db))
	require.NoError(t, migration1559081901.Migrate(db))
	require.NoError(t, migration1559767166.Migrate(db))
	require.NoError(t, migration1560433987.Migrate(db))
	require.NoError(t, migration1560791143.Migrate(db))
	require.NoError(t, migration1560881846.Migrate(db))

	head := migration0.Head{
		HashRaw: "dad0000000000000000000000000000000000000000000000000000000000b0d",
		Number:  8616460799,
	}
	require.NoError(t, db.Create(&head).Error)

	require.NoError(t, migration1560886530.Migrate(db))

	headFound := models.Head{}
	require.NoError(t, db.Where("id = (SELECT MAX(id) FROM heads)").Find(&headFound).Error)
	assert.Equal(t, "0xdad0000000000000000000000000000000000000000000000000000000000b0d", headFound.Hash.Hex())
	assert.Equal(t, int64(8616460799), headFound.Number)
}

func TestMigrate_Migration1565139192(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()
	db := orm.DB

	require.NoError(t, migration0.Migrate(db))
	require.NoError(t, migration1565139192.Migrate(db))
	specNoPayment := models.NewJobFromRequest(models.JobSpecRequest{})
	specWithPayment := models.NewJobFromRequest(models.JobSpecRequest{
		MinPayment: assets.NewLink(5),
	})
	specOneFound := models.JobSpec{}
	specTwoFound := models.JobSpec{}

	require.NoError(t, db.Create(&specWithPayment).Error)
	require.NoError(t, db.Create(&specNoPayment).Error)
	require.NoError(t, db.Where("id = ?", specNoPayment.ID).Find(&specOneFound).Error)
	require.Equal(t, assets.NewLink(0), specNoPayment.MinPayment)
	require.NoError(t, db.Where("id = ?", specWithPayment.ID).Find(&specTwoFound).Error)
	require.Equal(t, assets.NewLink(5), specWithPayment.MinPayment)
}

func TestMigrate_Migration1565210496(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	db := orm.DB
	db.LogMode(true)

	require.NoError(t, migration0.Migrate(db))

	jobSpec := models.JobSpec{
		ID:        utils.NewBytes32ID(),
		CreatedAt: time.Now(),
	}
	require.NoError(t, db.Create(&jobSpec).Error)
	jobRun := migration0.JobRun{
		ID:             utils.NewBytes32ID(),
		JobSpecID:      jobSpec.ID,
		ObservedHeight: "115792089237316195423570985008687907853269984665640564039457584007913129639936",
		CreationHeight: "0",
	}
	require.NoError(t, db.Create(&jobRun).Error)

	require.NoError(t, migration1565210496.Migrate(db))

	jobRunFound := models.JobRun{}
	require.NoError(t, db.Where("id = ?", jobRun.ID).Find(&jobRunFound).Error)
	assert.Equal(t, "115792089237316195423570985008687907853269984665640564039457584007913129639936", jobRunFound.ObservedHeight.String())
	assert.Equal(t, "0", jobRunFound.CreationHeight.String())
}

func TestMigrate_NewerVersionGuard(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	db := orm.DB

	// Do full migrations
	require.NoError(t, migrations.Migrate(db))

	// Add a fictional future migration
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID:      "9223372036854775807",
			Migrate: migration0.Migrate,
		},
	})
	require.NoError(t, m.Migrate())

	// Run migrations again, should error
	require.Error(t, migrations.Migrate(db))
}
