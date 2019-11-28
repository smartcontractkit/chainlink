package migrations_test

import (
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"chainlink/core/internal/cltest"
	"chainlink/core/store/assets"
	"chainlink/core/store/migrations"
	"chainlink/core/store/migrations/migration0"
	"chainlink/core/store/migrations/migration1560881855"
	"chainlink/core/store/models"
	"chainlink/core/store/orm"
	"chainlink/core/utils"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormigrate "gopkg.in/gormigrate.v1"
)

func bootstrapORM(t *testing.T) (*orm.ORM, func()) {
	tc, cleanup := cltest.NewConfig(t)
	config := tc.Config

	require.NoError(t, os.MkdirAll(config.RootDir(), 0700))
	cleanupDB := cltest.PrepareTestDB(tc)
	orm, err := orm.NewORM(orm.NormalizedDatabaseURL(config), config.DatabaseTimeout())
	require.NoError(t, err)
	orm.SetLogging(true)

	return orm, func() {
		assert.NoError(t, orm.Close())
		cleanup()
		os.RemoveAll(config.RootDir())
		cleanupDB()
	}
}

func TestMigrate_Migrations(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	err := orm.RawDB(func(db *gorm.DB) error {
		require.NoError(t, migrations.Migrate(db))

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
		return nil
	})
	require.NoError(t, err)
}

func TestMigrate_Migration1560791143(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	err := orm.RawDB(func(db *gorm.DB) error {
		require.NoError(t, migrations.MigrateTo(db, "0"))

		tx := migration0.Tx{
			ID:       1337,
			Data:     make([]byte, 10),
			Value:    utils.NewBig(big.NewInt(1)),
			GasPrice: utils.NewBig(big.NewInt(127)),
		}
		require.NoError(t, db.Create(&tx).Error)

		require.NoError(t, migrations.MigrateTo(db, "1559081901"))
		txFound := models.Tx{}
		require.NoError(t, db.Where("id = ?", tx.ID).Find(&txFound).Error)

		require.NoError(t, migrations.MigrateTo(db, "1560791143"))

		txNoID := models.Tx{
			Data:     make([]byte, 10),
			Value:    utils.NewBig(big.NewInt(2)),
			GasPrice: utils.NewBig(big.NewInt(119)),
		}
		require.NoError(t, db.Create(&txNoID).Error)
		assert.Equal(t, uint64(1338), txNoID.ID)

		noIDTxFound := models.Tx{}
		require.NoError(t, db.Where("id = ?", tx.ID).Find(&noIDTxFound).Error)
		return nil
	})
	require.NoError(t, err)
}

func TestMigrate_Migration1560881855(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	err := orm.RawDB(func(db *gorm.DB) error {
		require.NoError(t, migrations.MigrateTo(db, "1560924400"))

		befCreation := time.Now()
		jobSpecID := uuid.Must(uuid.NewV4())
		jobID := uuid.Must(uuid.NewV4())
		query := fmt.Sprintf(`
INSERT INTO run_results (amount) VALUES (2);
INSERT INTO job_specs (id) VALUES ('%s');
INSERT INTO job_runs (id, job_spec_id, overrides_id) VALUES ('%s', '%s', (SELECT id from run_results order by id DESC limit 1));
`, jobSpecID, jobID, jobSpecID)
		require.NoError(t, db.Exec(query).Error)
		aftCreation := time.Now()

		// placement of this migration is important, as it makes sure backfilling
		//  is done if there's already a RunResult with nonzero link reward
		require.NoError(t, migrations.MigrateTo(db, "1560881855"))

		rowFound := migration1560881855.LinkEarned{}
		require.NoError(t, db.Table("link_earned").Find(&rowFound).Error)
		assert.Equal(t, jobSpecID.String(), rowFound.JobSpecID)
		assert.Equal(t, jobID.String(), rowFound.JobRunID)
		assert.Equal(t, assets.NewLink(2), rowFound.Earned)
		assert.True(t, true, rowFound.EarnedAt.After(aftCreation), rowFound.EarnedAt.Before(befCreation))
		return nil
	})
	require.NoError(t, err)
}

func TestMigrate_Migration1560881846(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	err := orm.RawDB(func(db *gorm.DB) error {
		require.NoError(t, migrations.MigrateTo(db, "1560881846"))

		head := migration0.Head{
			HashRaw: "dad0000000000000000000000000000000000000000000000000000000000b0d",
			Number:  8616460799,
		}
		require.NoError(t, db.Create(&head).Error)

		require.NoError(t, migrations.MigrateTo(db, "1560886530"))

		headFound := models.Head{}
		require.NoError(t, db.Where("id = (SELECT MAX(id) FROM heads)").Find(&headFound).Error)
		assert.Equal(t, "0xdad0000000000000000000000000000000000000000000000000000000000b0d", headFound.Hash.Hex())
		assert.Equal(t, int64(8616460799), headFound.Number)
		return nil
	})
	require.NoError(t, err)
}

func TestMigrate_Migration1565139192(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	err := orm.RawDB(func(db *gorm.DB) error {
		require.NoError(t, migrations.MigrateTo(db, "1565139192"))

		specNoPayment := models.NewJobFromRequest(models.JobSpecRequest{})
		specWithPayment := models.NewJobFromRequest(models.JobSpecRequest{
			MinPayment: *assets.NewLink(5),
		})
		specOneFound := models.JobSpec{}
		specTwoFound := models.JobSpec{}

		require.NoError(t, db.Create(&specWithPayment).Error)
		require.NoError(t, db.Create(&specNoPayment).Error)
		require.NoError(t, db.Where("id = ?", specNoPayment.ID).Find(&specOneFound).Error)
		require.Equal(t, *assets.NewLink(0), specNoPayment.MinPayment)
		require.NoError(t, db.Where("id = ?", specWithPayment.ID).Find(&specTwoFound).Error)
		require.Equal(t, *assets.NewLink(5), specWithPayment.MinPayment)
		return nil
	})
	require.NoError(t, err)
}

func TestMigrate_Migration1565210496(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	err := orm.RawDB(func(db *gorm.DB) error {
		require.NoError(t, migrations.MigrateTo(db, "0"))

		jobSpec := migration0.JobSpec{
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

		require.NoError(t, migrations.MigrateTo(db, "1565210496"))

		jobRunFound := models.JobRun{}
		require.NoError(t, db.Where("id = ?", jobRun.ID).Find(&jobRunFound).Error)
		assert.Equal(t, "115792089237316195423570985008687907853269984665640564039457584007913129639936", jobRunFound.ObservedHeight.String())
		assert.Equal(t, "0", jobRunFound.CreationHeight.String())
		return nil
	})
	require.NoError(t, err)
}

func TestMigrate_Migration1565291711(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	err := orm.RawDB(func(db *gorm.DB) error {
		require.NoError(t, migrations.MigrateTo(db, "1560881855"))

		jobSpec := migration0.JobSpec{
			ID:        utils.NewBytes32ID(),
			CreatedAt: time.Now(),
		}
		require.NoError(t, db.Create(&jobSpec).Error)
		jobRun := migration0.JobRun{
			ID:             utils.NewBytes32ID(),
			JobSpecID:      jobSpec.ID,
			CreationHeight: "0",
			ObservedHeight: "0",
		}
		require.NoError(t, db.Create(&jobRun).Error)

		require.NoError(t, migrations.MigrateTo(db, "1565291711"))

		jobRunFound := models.JobRun{}
		require.NoError(t, db.Where("id = ?", jobRun.ID).Find(&jobRunFound).Error)
		assert.Equal(t, jobRun.ID, jobRunFound.ID.String())
		return nil
	})
	require.NoError(t, err)
}

func TestMigrate_Migration1565877314(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	err := orm.RawDB(func(db *gorm.DB) error {
		require.NoError(t, migrations.MigrateTo(db, "0"))

		exi := migration0.ExternalInitiator{
			AccessKey:    "access_key",
			Salt:         "salt",
			HashedSecret: "hashed_secret",
		}
		require.NoError(t, db.Create(&exi).Error)

		require.NoError(t, migrations.MigrateTo(db, "1565877314"))

		exiFound := models.ExternalInitiator{}
		require.NoError(t, db.Where("id = ?", exi.ID).Find(&exiFound).Error)
		assert.Equal(t, "access_key", exiFound.Name)
		assert.Equal(t, "https://unset.url", exiFound.URL.String())
		return nil
	})
	require.NoError(t, err)
}

func TestMigrate_Migration1570675883(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	err := orm.RawDB(func(db *gorm.DB) error {
		require.NoError(t, migrations.MigrateTo(db, "0"))

		overrides := models.RunResult{
			Data: cltest.JSONFromString(t, `{"a": "b"}`),
		}
		require.NoError(t, db.Create(&overrides).Error)

		jobSpec := migration0.JobSpec{
			ID:        utils.NewBytes32ID(),
			CreatedAt: time.Now(),
		}
		require.NoError(t, db.Create(&jobSpec).Error)
		jobRun := migration0.JobRun{
			ID:             utils.NewBytes32ID(),
			JobSpecID:      jobSpec.ID,
			OverridesID:    overrides.ID,
			CreationHeight: "0",
			ObservedHeight: "0",
		}
		require.NoError(t, db.Create(&jobRun).Error)

		require.NoError(t, migrations.MigrateTo(db, "1570675883"))

		jobRunFound := models.JobRun{}
		require.NoError(t, db.Where("id = ?", jobRun.ID).Find(&jobRunFound).Error)
		assert.Equal(t, `{"a": "b"}`, jobRunFound.Overrides.String())
		require.Error(t, db.Where("id = ?", overrides.ID).Find(&overrides).Error)
		return nil
	})
	require.NoError(t, err)
}

func TestMigrate_NewerVersionGuard(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	err := orm.RawDB(func(db *gorm.DB) error {
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
		return nil
	})
	require.NoError(t, err)
}
