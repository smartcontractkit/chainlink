package migrations_test

import (
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/migrations"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1551895034"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1551895034/old"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1552418531"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1554131520"
	"github.com/smartcontractkit/chainlink/core/store/models"
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

func TestMigrate_Upgrade(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()
	db := orm.DB

	// Create an old migration schema table
	err := db.Exec(`
		CREATE TABLE IF NOT EXISTS "migration_timestamps" (
			"timestamp" varchar(12),
			PRIMARY KEY ("timestamp")
		);
		INSERT INTO migration_timestamps VALUES('0');
		INSERT INTO migration_timestamps VALUES('1549496047');
		INSERT INTO migration_timestamps VALUES('1551816486');
		INSERT INTO migration_timestamps VALUES('1551895034');
		INSERT INTO migration_timestamps VALUES('1552418531');
		INSERT INTO migration_timestamps VALUES('1553029703');
		INSERT INTO migration_timestamps VALUES('1554131520');
	`).Error
	require.NoError(t, err)

	require.NoError(t, migrations.Migrate(db))
	assert.False(t, db.HasTable("migration_timestamps"))
	assert.True(t, db.HasTable("migrations"))
}

func TestMigrationFromExistingDB(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	fixtureDBPath := fmt.Sprintf("../../../internal/fixtures/migrations/1554131520_dump.%s.sql", orm.DialectName())
	loadSqlDump(t, orm, fixtureDBPath)

	require.NoError(t, migrations.Migrate(orm.DB))
}

func loadSqlDump(t *testing.T, orm *orm.ORM, sqldump string) error {
	return orm.DB.Exec(string(cltest.MustReadFile(t, sqldump))).Error
}

func TestMigrate_Migration0(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	db := orm.DB
	tm := &migration0.Migration{}

	require.NoError(t, tm.Migrate(db))

	assert.True(t, db.HasTable("job_specs"))
	assert.True(t, db.HasTable("task_specs"))
	assert.True(t, db.HasTable("job_runs"))
	assert.True(t, db.HasTable("task_runs"))
	assert.True(t, db.HasTable("run_results"))
	assert.True(t, db.HasTable("initiators"))
	assert.True(t, db.HasTable("txes"))
	assert.True(t, db.HasTable("tx_attempts"))
	assert.True(t, db.HasTable("bridge_types"))
	assert.True(t, db.HasTable("heads"))
	assert.True(t, db.HasTable("users"))
	assert.True(t, db.HasTable("sessions"))
	assert.True(t, db.HasTable("encumbrances"))
	assert.True(t, db.HasTable("service_agreements"))
}

func TestMigrate1551816486(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	// seed db w old table
	err := orm.DB.Exec(`
		CREATE TABLE "bridge_types" (
			"name" varchar(255),
			"url" varchar(255),
			"confirmations" bigint,
			"incoming_token" varchar(255),
			"outgoing_token" varchar(255),
			"minimum_contract_payment" varchar(255),
			UNIQUE (name));
	`).Error

	require.NoError(t, err)

	initial := models.BridgeType{
		Name: "someUniqueName",
		URL:  cltest.WebURL("http://someurl.com"),
	}

	require.NoError(t, orm.DB.Save(&initial).Error)
	require.NoError(t, migration0.Migration{}.Migrate(orm.DB))

	migratedbt, err := orm.FindBridge(initial.Name.String())
	require.NoError(t, err)
	require.Equal(t, initial, migratedbt)
}

func TestMigrate1551895034(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	tm := &migration1551895034.Migration{}

	height := models.NewBig(big.NewInt(1337))
	hash := common.HexToHash("0xde3fb1df888c6c7f77f3a8e9c2582f87e7ad5277d98bd06cfd17cd2d7ea49f42")

	previous := old.IndexableBlockNumber{
		Number: *height,
		Digits: 4,
		Hash:   hash,
	}
	// seed w old schema and data
	err := orm.DB.AutoMigrate(old.IndexableBlockNumber{}).Error
	require.NoError(t, err)
	err = orm.DB.Save(&previous).Error
	require.NoError(t, err)

	// migrate
	require.NoError(t, tm.Migrate(orm.DB))

	retrieved := models.Head{}
	err = orm.DB.First(&retrieved).Error
	require.NoError(t, err)

	require.Equal(t, height.ToInt(), retrieved.ToInt())
	require.Equal(
		t,
		hash.String(),
		retrieved.Hash().Hex())
}

func TestMigrate1552418531(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	// seed w old schema
	err := orm.DB.Exec(`
		CREATE TABLE "job_specs" ("id" varchar(255) NOT NULL,"created_at" timestamp,"start_at" timestamp,"end_at" timestamp, PRIMARY KEY ("id"));
		INSERT INTO "job_specs" VALUES ('testjobspec', CURRENT_TIMESTAMP, NULL, NULL);
	`).Error
	require.NoError(t, err)

	// migrate
	tm := &migration1552418531.Migration{}

	require.NoError(t, tm.Migrate(orm.DB))

	retrieved := models.JobSpec{}
	err = orm.DB.First(&retrieved).Error
	require.NoError(t, err)

	require.Equal(t, false, retrieved.DeletedAt.Valid)

	err = orm.DB.Delete(&retrieved).Error
	require.NoError(t, err)
	err = orm.DB.First(&retrieved).Error
	require.Error(t, err)
	err = orm.DB.Unscoped().First(&retrieved).Error
	require.NoError(t, err)
	require.Equal(t, true, retrieved.DeletedAt.Valid)
}

func TestMigrate1554131520(t *testing.T) {
	orm, cleanup := bootstrapORM(t)
	defer cleanup()

	// seed w old schema
	tm0 := &migration0.Migration{}
	require.NoError(t, tm0.Migrate(orm.DB))

	j := cltest.NewJob()
	j.Initiators = []models.Initiator{
		{
			JobSpecID: j.ID,
			Type:      models.InitiatorCron,
			InitiatorParams: models.InitiatorParams{
				Schedule: models.Cron("* * * * *"),
			},
		},
		{
			JobSpecID: j.ID,
			Type:      models.InitiatorWeb,
		},
		{
			JobSpecID: j.ID,
			Type:      models.InitiatorEthLog,
			InitiatorParams: models.InitiatorParams{
				Address: cltest.NewAddress(),
			},
		},
		{
			JobSpecID: j.ID,
			Type:      models.InitiatorRunLog,
			InitiatorParams: models.InitiatorParams{
				Address: cltest.NewAddress(),
			},
		},
	}

	require.NoError(t, orm.CreateJob(&j))

	cronjr := newRunWithoutRunRequest(j, j.Initiators[0])
	webjr := newRunWithoutRunRequest(j, j.Initiators[1])
	ethlogjr := newRunWithoutRunRequest(j, j.Initiators[2])
	runlogjr := newRunWithoutRunRequest(j, j.Initiators[3])

	require.NoError(t, orm.CreateJobRun(cronjr))
	require.NoError(t, orm.CreateJobRun(webjr))
	require.NoError(t, orm.CreateJobRun(ethlogjr))
	require.NoError(t, orm.CreateJobRun(runlogjr))

	orm.DB.Exec(`
		UPDATE job_runs SET run_request_id = NULL;
	`)

	// migrate
	tm := &migration1554131520.Migration{}
	require.NoError(t, tm.Migrate(orm.DB))

	// check run request backfill
	retrieved := models.JobRun{}
	require.NoError(t, orm.DB.Where("ID = ?", cronjr.ID).Preload("RunRequest").First(&retrieved).Error)
	assert.NotEqual(t, time.Time{}, retrieved.RunRequest.CreatedAt)
	assert.Nil(t, retrieved.RunRequest.RequestID)
	assert.Nil(t, retrieved.RunRequest.Requester)
	assert.Nil(t, retrieved.RunRequest.TxHash)

	retrieved = models.JobRun{}
	require.NoError(t, orm.DB.Where("ID = ?", webjr.ID).Preload("RunRequest").First(&retrieved).Error)
	assert.NotEqual(t, time.Time{}, retrieved.RunRequest.CreatedAt)
	assert.Nil(t, retrieved.RunRequest.RequestID)
	assert.Nil(t, retrieved.RunRequest.Requester)
	assert.Nil(t, retrieved.RunRequest.TxHash)

	retrieved = models.JobRun{}
	require.NoError(t, orm.DB.Where("ID = ?", ethlogjr.ID).Preload("RunRequest").First(&retrieved).Error)
	assert.NotEqual(t, time.Time{}, retrieved.RunRequest.CreatedAt)
	assert.NotNil(t, retrieved.RunRequest.TxHash)
	assert.Nil(t, retrieved.RunRequest.RequestID)
	assert.Nil(t, retrieved.RunRequest.Requester)

	retrieved = models.JobRun{}
	require.NoError(t, orm.DB.Where("ID = ?", runlogjr.ID).Preload("RunRequest").First(&retrieved).Error)
	assert.NotEqual(t, time.Time{}, retrieved.RunRequest.CreatedAt)
	assert.NotNil(t, retrieved.RunRequest.TxHash)
	assert.NotNil(t, retrieved.RunRequest.Requester)
	assert.Equal(t, "BACKFILLED_FAKE", *retrieved.RunRequest.RequestID)
}

func newRunWithoutRunRequest(j models.JobSpec, i models.Initiator) *models.JobRun {
	jr := j.NewRun(i)
	jr.RunRequest = models.RunRequest{}
	return &jr
}
