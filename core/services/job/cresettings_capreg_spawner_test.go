package job_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/cresettings"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
)

func capRegistrySpec(t *testing.T, offchainConfig string) job.Job {
	t.Helper()
	jb, err := cresettings.ValidatedCRESettingsSpec(fmt.Sprintf(`type = "cresettings"
schemaVersion = 1
externalJobID = "%s"
config_type = "capabilities_registry"
offchain_config = '''%s'''`, uuid.New(), offchainConfig))
	require.NoError(t, err)
	return jb
}

// TestSpawner_CRESettingsCapabilitiesRegistry exercises the node-side path end to end with a
// real database, job spawner and cresettings delegate: create -> applied to GlobalConfig; node
// restart (fresh spawner + GlobalConfig booting from the DB) re-applies the same payload
// alongside a settings job; deleting the job withdraws it; a replacement applies; a duplicate
// is rejected before it is persisted.
func TestSpawner_CRESettingsCapabilitiesRegistry(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	lggr := logger.TestLogger(t)
	orm := NewTestORM(t, db, pipeline.NewORM(db, lggr, config.JobPipeline().MaxSuccessfulRuns()), bridges.NewORM(db), keyStore)

	type node struct {
		spawner   job.Spawner
		projector *cresettings.CapRegistryProjector
		gc        *globalconfig.GlobalConfig
		settings  *loop.AtomicSettings
	}
	boot := func() node {
		n := node{gc: globalconfig.New(), settings: &loop.AtomicSettings{}}
		// pgtest never commits, so here "committed" means visible in the test transaction.
		n.projector = cresettings.NewCapRegistryProjector(lggr, db, n.gc)
		d := cresettings.NewDelegate(lggr, n.settings, &loop.AtomicSettings{}, n.projector)
		n.spawner = job.NewSpawner(orm, config.Database(), noopChecker{}, map[job.Type]job.Delegate{job.CRESettings: d}, lggr, nil)
		require.NoError(t, n.spawner.Start(ctx))
		return n
	}

	const payload = `{"version":3,"dons":{"7":{"capabilityConfigs":{"cron@1.0.0":{"specConfig":{"fields":{"interval":{"stringValue":"30"}}}}}}}}`

	first := boot()
	capRegJob := capRegistrySpec(t, payload)
	require.NoError(t, first.spawner.CreateJob(ctx, nil, &capRegJob))
	settingsJob, err := cresettings.ValidatedCRESettingsSpec(testspecs.GetCRESettingsSpec())
	require.NoError(t, err)
	require.NoError(t, first.spawner.CreateJob(ctx, nil, &settingsJob))

	raw, v := first.gc.Load()
	assert.JSONEq(t, payload, raw)
	assert.Equal(t, uint64(3), v)

	require.NoError(t, first.spawner.Close())

	// Restart: a fresh node boots from the DB and re-applies the same payload, routed to
	// GlobalConfig (not to the CRE settings).
	second := boot()
	raw, v = second.gc.Load()
	assert.JSONEq(t, payload, raw)
	assert.Equal(t, uint64(3), v)
	applied, err := second.settings.Load()
	require.NoError(t, err)
	assert.Equal(t, settingsJob.CRESettingsSpec.Hash, applied.Hash)
	reg, _ := second.gc.LoadParsed()
	require.NotNil(t, reg)
	assert.Equal(t, "30", reg.GetDons()[7].GetCapabilityConfigs()["cron@1.0.0"].GetSpecConfig().GetFields()["interval"].GetStringValue())

	// Deleting the job withdraws the payload (revert to on-chain/TOML) once the projector
	// observes the committed delete.
	require.NoError(t, second.spawner.DeleteJob(ctx, nil, capRegJob.ID))
	require.NoError(t, second.projector.Refresh(ctx))
	reg, v = second.gc.LoadParsed()
	assert.Nil(t, reg)
	assert.Equal(t, uint64(0), v)

	// Replacement with a newer payload is accepted once the old job is gone.
	next := capRegistrySpec(t, `{"version":4}`)
	require.NoError(t, second.spawner.CreateJob(ctx, nil, &next))
	_, v = second.gc.Load()
	assert.Equal(t, uint64(4), v)

	// A second capabilities_registry job is rejected at insert time, before anything is
	// persisted, so it can never displace the applied job on a later boot (jobs start
	// newest-first). This must be the last DB statement: pgtest runs the whole test in one
	// transaction without savepoints, so the rejected INSERT aborts it.
	dup := capRegistrySpec(t, `{"version":9}`)
	require.ErrorIs(t, second.spawner.CreateJob(ctx, nil, &dup), job.ErrCRESettingsCapRegistryExists)
	_, v = second.gc.Load()
	assert.Equal(t, uint64(4), v)
	require.NoError(t, second.spawner.Close())
}
