package cresettings

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

const (
	settingsToml = `Foo = "bar"
`
	shardAssignmentToml = `config_type = "shard_assignment"
static_default_assignment = [0, 1]
`
)

func newTestDelegate(t *testing.T) *delegate {
	t.Helper()
	return NewDelegate(logger.TestLogger(t), &loop.AtomicSettings{}, &loop.AtomicSettings{}, globalconfig.New())
}

func cresettingsJob(id int32, settings string) job.Job {
	return job.Job{
		ID:              id,
		Type:            job.CRESettings,
		CRESettingsSpec: &job.CRESettingsSpec{Settings: settings},
	}
}

func capRegistryJob(id int32, raw, hash string) job.Job {
	return job.Job{
		ID:   id,
		Type: job.CRESettings,
		CRESettingsSpec: &job.CRESettingsSpec{
			ConfigType:     ConfigTypeCapRegistry,
			OffchainConfig: raw,
			Hash:           hash,
		},
	}
}

func TestServicesForSpecOneJobPerConfigType(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	d := newTestDelegate(t)

	_, err := d.ServicesForSpec(ctx, cresettingsJob(1, settingsToml))
	require.NoError(t, err)

	_, err = d.ServicesForSpec(ctx, cresettingsJob(2, shardAssignmentToml))
	require.NoError(t, err)

	_, err = d.ServicesForSpec(ctx, cresettingsJob(3, settingsToml))
	require.ErrorContains(t, err, fmt.Sprintf("another %s job with config_type %q is already active: 1", job.CRESettings, ConfigTypeSettings))

	_, err = d.ServicesForSpec(ctx, cresettingsJob(4, shardAssignmentToml))
	require.ErrorContains(t, err, fmt.Sprintf("another %s job with config_type %q is already active: 2", job.CRESettings, ConfigTypeShardAssignment))
}

func TestServicesForSpecDefaultConfigTypeSharesSettingsSlot(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	d := newTestDelegate(t)

	_, err := d.ServicesForSpec(ctx, cresettingsJob(1, shardAssignmentToml))
	require.NoError(t, err)

	_, err = d.ServicesForSpec(ctx, cresettingsJob(2, settingsToml))
	require.NoError(t, err)

	// spec with no config_type defaults to settings, so it collides with job 2
	_, err = d.ServicesForSpec(ctx, cresettingsJob(3, `Foo = "baz"
`))
	require.ErrorContains(t, err, "already active: 2")
}

func TestServicesForSpecUnknownConfigType(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	d := newTestDelegate(t)

	_, err := d.ServicesForSpec(ctx, cresettingsJob(1, `config_type = "unknown"
Foo = "bar"
`))
	require.ErrorContains(t, err, `unknown config_type "unknown"`)

	// unknown config_type must not claim a slot
	_, err = d.ServicesForSpec(ctx, cresettingsJob(2, settingsToml))
	require.NoError(t, err)
}

func TestOnDeleteJobFreesSlotPerConfigType(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	d := newTestDelegate(t)

	_, err := d.ServicesForSpec(ctx, cresettingsJob(1, settingsToml))
	require.NoError(t, err)
	_, err = d.ServicesForSpec(ctx, cresettingsJob(2, shardAssignmentToml))
	require.NoError(t, err)

	// deleting a rejected duplicate must not free the active job's slot
	_, err = d.ServicesForSpec(ctx, cresettingsJob(3, settingsToml))
	require.ErrorContains(t, err, "already active: 1")
	require.NoError(t, d.OnDeleteJob(ctx, cresettingsJob(3, settingsToml)))
	_, err = d.ServicesForSpec(ctx, cresettingsJob(4, settingsToml))
	require.ErrorContains(t, err, "already active: 1")

	// deleting the active settings job frees only the settings slot
	require.NoError(t, d.OnDeleteJob(ctx, cresettingsJob(1, settingsToml)))
	_, err = d.ServicesForSpec(ctx, cresettingsJob(5, settingsToml))
	require.NoError(t, err)

	// shard assignment slot is unaffected
	_, err = d.ServicesForSpec(ctx, cresettingsJob(6, shardAssignmentToml))
	require.ErrorContains(t, err, "already active: 2")
	require.NoError(t, d.OnDeleteJob(ctx, cresettingsJob(2, shardAssignmentToml)))
	_, err = d.ServicesForSpec(ctx, cresettingsJob(7, shardAssignmentToml))
	require.NoError(t, err)
}

func TestDelegate_CapabilitiesRegistry_StoresIntoGlobalConfig(t *testing.T) {
	t.Parallel()

	d := newTestDelegate(t)

	_, err := d.ServicesForSpec(t.Context(), capRegistryJob(1, `{"version":7}`, "h7"))
	require.NoError(t, err)

	raw, v := d.globalConfig.Load()
	assert.Equal(t, `{"version":7}`, raw)
	assert.Equal(t, uint64(7), v)
}

func TestDelegate_RejectsSecondJobOfSameConfigType(t *testing.T) {
	t.Parallel()

	d := newTestDelegate(t)

	_, err := d.ServicesForSpec(t.Context(), capRegistryJob(1, `{"version":1}`, "h1"))
	require.NoError(t, err)

	_, err = d.ServicesForSpec(t.Context(), capRegistryJob(2, `{"version":2}`, "h2"))
	require.ErrorContains(t, err, "already active")
}

func TestDelegate_DifferentConfigTypesCoexist(t *testing.T) {
	t.Parallel()

	d := newTestDelegate(t)

	// settings job
	settingsJob := job.Job{ID: 10, Type: job.CRESettings, CRESettingsSpec: &job.CRESettingsSpec{Settings: `Foo = "bar"`, Hash: "hs"}}
	_, err := d.ServicesForSpec(t.Context(), settingsJob)
	require.NoError(t, err)

	// capabilities_registry job coexists
	_, err = d.ServicesForSpec(t.Context(), capRegistryJob(11, `{"version":1}`, "h1"))
	require.NoError(t, err)
}

func TestDelegate_OnDeleteJobClearsConfigType(t *testing.T) {
	t.Parallel()

	d := newTestDelegate(t)

	_, err := d.ServicesForSpec(t.Context(), capRegistryJob(1, `{"version":1}`, "h1"))
	require.NoError(t, err)

	require.NoError(t, d.OnDeleteJob(t.Context(), capRegistryJob(1, `{"version":1}`, "h1")))

	// A new job of the same config_type is now accepted.
	_, err = d.ServicesForSpec(t.Context(), capRegistryJob(2, `{"version":2}`, "h2"))
	require.NoError(t, err)
}

// TestDelegate_FailedStoreReleasesSlot covers the feeds-manager replacement flow: the old job is
// deleted, then its replacement is created. If the replacement's payload is rejected, its slot
// must be released so a later valid job of the same config_type can still be created.
func TestDelegate_FailedStoreReleasesSlot(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	d := newTestDelegate(t)

	_, err := d.ServicesForSpec(ctx, capRegistryJob(1, `{"version":5}`, "h5"))
	require.NoError(t, err)
	require.NoError(t, d.OnDeleteJob(ctx, capRegistryJob(1, `{"version":5}`, "h5")))

	// Stale version: rejected by GlobalConfig.
	_, err = d.ServicesForSpec(ctx, capRegistryJob(2, `{"version":4}`, "h4"))
	require.ErrorContains(t, err, "not newer than applied version 5")
	// Invalid payload: rejected by GlobalConfig.
	_, err = d.ServicesForSpec(ctx, capRegistryJob(3, `{"version":0}`, "h0"))
	require.ErrorContains(t, err, "version must be >= 1")

	// Neither rejected job holds the slot.
	_, err = d.ServicesForSpec(ctx, capRegistryJob(4, `{"version":6}`, "h6"))
	require.NoError(t, err)
	_, v := d.globalConfig.Load()
	assert.Equal(t, uint64(6), v)

	// The active job still owns the slot after a rejected attempt.
	_, err = d.ServicesForSpec(ctx, capRegistryJob(5, `{"version":7}`, "h7"))
	require.ErrorContains(t, err, "already active: 4")
}

func TestDelegate_FailedSettingsStoreReleasesSlot(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	d := NewDelegate(logger.TestLogger(t), &loop.AtomicSettings{}, &loop.AtomicSettings{}, nil)

	// No GlobalConfig wired: capabilities_registry fails, and must not hold its slot.
	_, err := d.ServicesForSpec(ctx, capRegistryJob(1, `{"version":1}`, "h1"))
	require.ErrorContains(t, err, "no global config store")
	_, loaded := d.activeJobIDs.Load(ConfigTypeCapRegistry)
	assert.False(t, loaded)
}

// TestDelegate_RestartFromPersistedSpecs simulates node restart: specs as produced by the
// validator (and round-tripped through the DB) are replayed into a fresh delegate. The
// capabilities_registry job must route to GlobalConfig, not collide with the settings job or
// overwrite the CRE settings with its empty Settings field.
func TestDelegate_RestartFromPersistedSpecs(t *testing.T) {
	t.Parallel()

	settingsJob, err := ValidatedCRESettingsSpec(`type = "cresettings"
schemaVersion = 1
settings = '''Foo = "bar"'''`)
	require.NoError(t, err)
	settingsJob.ID = 1
	capRegJob, err := ValidatedCRESettingsSpec(`type = "cresettings"
schemaVersion = 1
config_type = "capabilities_registry"
offchain_config = '''{"version":4}'''`)
	require.NoError(t, err)
	capRegJob.ID = 2

	atomicSettings := &loop.AtomicSettings{}
	gc := globalconfig.New()
	d := NewDelegate(logger.TestLogger(t), atomicSettings, &loop.AtomicSettings{}, gc)

	for _, j := range []job.Job{capRegJob, settingsJob} {
		_, err = d.ServicesForSpec(t.Context(), j)
		require.NoError(t, err)
	}

	_, v := gc.Load()
	assert.Equal(t, uint64(4), v)
	applied, err := atomicSettings.Load()
	require.NoError(t, err)
	assert.Equal(t, settingsJob.CRESettingsSpec.Hash, applied.Hash)
}
