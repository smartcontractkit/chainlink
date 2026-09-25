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
