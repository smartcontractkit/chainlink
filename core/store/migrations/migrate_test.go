package migrations_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/migrations"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560881855"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormigrate "gopkg.in/gormigrate.v1"
)

func TestMigrate_Migrations(t *testing.T) {
	_, orm, cleanup := cltest.BootstrapThrowawayORM(t, "migrations", false)
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

func TestMigrate_Migration1560881855(t *testing.T) {
	_, orm, cleanup := cltest.BootstrapThrowawayORM(t, "migrations", false)
	defer cleanup()

	err := orm.RawDB(func(db *gorm.DB) error {
		require.NoError(t, migrations.MigrateTo(db, "1560924400"))

		befCreation := time.Now()
		jobSpecID := uuid.NewV4()
		jobID := uuid.NewV4()
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
	_, orm, cleanup := cltest.BootstrapThrowawayORM(t, "migrations", false)
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

func TestMigrate_Migration1565210496(t *testing.T) {
	_, orm, cleanup := cltest.BootstrapThrowawayORM(t, "migrations", false)
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
	_, orm, cleanup := cltest.BootstrapThrowawayORM(t, "migrations", false)
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
	_, orm, cleanup := cltest.BootstrapThrowawayORM(t, "migrations", false)
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

func TestMigrate_Migration1586369235(t *testing.T) {
	// Make sure that the data still reads OK afterward
	_, orm, cleanup := cltest.BootstrapThrowawayORM(t, "migrations", false)
	defer cleanup()

	err := orm.RawDB(func(db *gorm.DB) error {
		require.NoError(t, migrations.MigrateTo(db, "1586163842"))
		hexEncodedData := "0x3162323831336636383832373462366261623565663264366135343866323038"
		binaryData := hexutil.MustDecode(hexEncodedData)
		address := hexutil.MustDecode("0xa0788FC17B1dEe36f057c42B6F373A34B014687e")
		bigInt := "42000000000000000000" // 42 LINK

		require.NoError(t, db.Exec(`INSERT INTO encumbrances (payment, aggregator, agg_initiate_job_selector, agg_fulfill_selector) VALUES (?, 'a', E'\\xDEADBEEF', E'\\xDEADBEEF')`, bigInt).Error)
		require.NoError(t, db.Exec(`INSERT INTO run_requests (request_id) VALUES (?::text)`, hexEncodedData).Error)
		require.NoError(t, db.Exec(`INSERT INTO txes (signed_raw_tx, "from", "to", data, nonce, value, gas_limit, hash, gas_price, confirmed, sent_at) VALUES (?::text, ?, ?, E'\\xDEADBEEF', 42, ?, 42, ?, ?, false, 42)`, hexEncodedData, address, address, bigInt, binaryData, bigInt).Error)
		require.NoError(t, db.Exec(`INSERT INTO tx_attempts (signed_raw_tx, created_at, hash, gas_price, confirmed, sent_at) VALUES (?::text, NOW(), ?, ?, false, 42)`, hexEncodedData, binaryData, bigInt).Error)

		require.NoError(t, migrations.MigrateTo(db, "1586369235"))

		var e models.Encumbrance
		require.NoError(t, db.First(&e, "true").Error)
		assert.Equal(t, e.Payment.ToInt().String(), bigInt)

		var rr models.RunRequest
		require.NoError(t, db.First(&rr, "true").Error)
		assert.Equal(t, rr.RequestID.Bytes(), binaryData)

		var tx models.Tx
		require.NoError(t, db.First(&tx, "true").Error)
		assert.Equal(t, tx.SignedRawTx, binaryData)
		assert.Equal(t, tx.GasPrice.String(), bigInt)
		assert.Equal(t, tx.Value.String(), bigInt)

		var txa models.TxAttempt
		require.NoError(t, db.First(&txa, "true").Error)
		assert.Equal(t, txa.SignedRawTx, binaryData)
		assert.Equal(t, txa.GasPrice.String(), bigInt)

		return nil
	})
	require.NoError(t, err)
}

func TestMigrate_NewerVersionGuard(t *testing.T) {
	_, orm, cleanup := cltest.BootstrapThrowawayORM(t, "migrations", false)
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
