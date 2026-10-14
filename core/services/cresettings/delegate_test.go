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

const capRegPayload = `{"domain":"cre","env":"test","version":1,"dons":{"don-a":{"capabilities":{"cron@1.0.0":{"spec_config":{"fields":{"interval":{"stringValue":"30"}}}}}}}}`

func capRegSettings(payload string) string {
	return "config_type = \"capabilities_registry\"\noffchain_config = '''" + payload + "'''\n"
}

func TestServicesForSpecCapRegistry(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	gc := globalconfig.New()
	d := NewDelegate(logger.TestLogger(t), &loop.AtomicSettings{}, &loop.AtomicSettings{}, gc)
	capRegJob := func(id int32, payload string) job.Job {
		jb := cresettingsJob(id, capRegSettings(payload))
		jb.CRESettingsSpec.Hash = payload
		return jb
	}
	jb := capRegJob(1, capRegPayload)

	_, err := d.ServicesForSpec(ctx, jb)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), gc.LoadParsed().GetVersion())

	// A second capabilities_registry job is rejected; other config types are independent.
	_, err = d.ServicesForSpec(ctx, capRegJob(2, capRegPayload))
	require.Error(t, err)
	_, err = d.ServicesForSpec(ctx, cresettingsJob(3, settingsToml))
	require.NoError(t, err)

	// Deleting the job frees the slot but keeps the payload, so a JD replace (delete + create)
	// whose new payload is rejected leaves the node on the old one.
	require.NoError(t, d.OnDeleteJob(ctx, jb))
	stale := `{"domain":"cre","env":"test","version":1,"dons":{}}`
	_, err = d.ServicesForSpec(ctx, capRegJob(4, stale))
	require.ErrorContains(t, err, "not newer")
	assert.Contains(t, gc.LoadParsed().GetDons(), "don-a")

	_, err = d.ServicesForSpec(ctx, capRegJob(5, capRegPayload))
	require.NoError(t, err)
}

func TestServicesForSpecCapRegistryInvalidPayloadFreesSlot(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	d := NewDelegate(logger.TestLogger(t), &loop.AtomicSettings{}, &loop.AtomicSettings{}, globalconfig.New())

	_, err := d.ServicesForSpec(ctx, cresettingsJob(1, "config_type = \"capabilities_registry\"\n"))
	require.Error(t, err)
	// The failed job must not hold the slot.
	_, err = d.ServicesForSpec(ctx, cresettingsJob(2, capRegSettings(capRegPayload)))
	require.NoError(t, err)
}
