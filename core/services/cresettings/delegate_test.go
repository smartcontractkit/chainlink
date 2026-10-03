package cresettings

import (
	"context"
	"fmt"
	"sync"
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
	d, _ := newTestDelegateWithCommitted(t)
	return d
}

// committedStore stands in for the database: it holds the committed capabilities_registry spec
// the projector reads.
type committedStore struct {
	mu   sync.Mutex
	spec *job.CRESettingsSpec
	err  error
}

func (c *committedStore) set(spec *job.CRESettingsSpec) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.spec = spec
}

func (c *committedStore) load(context.Context) (*job.CRESettingsSpec, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.spec, c.err
}

func newTestDelegateWithCommitted(t *testing.T) (*delegate, *committedStore) {
	t.Helper()
	committed := &committedStore{}
	p := newCapRegistryProjector(logger.TestLogger(t), committed.load, globalconfig.New())
	return NewDelegate(logger.TestLogger(t), &loop.AtomicSettings{}, &loop.AtomicSettings{}, p), committed
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

// TestDelegate_CapRegistry_AppliesOnlyCommittedState: ServicesForSpec may run inside the
// uncommitted create transaction, so it must not apply the spec it is handed; it applies
// whatever is committed.
func TestDelegate_CapRegistry_AppliesOnlyCommittedState(t *testing.T) {
	t.Parallel()

	d, committed := newTestDelegateWithCommitted(t)
	gc := d.capRegistry.GlobalConfig()

	// Not committed yet (e.g. inside a feeds approval transaction): nothing is applied.
	_, err := d.ServicesForSpec(t.Context(), capRegistryJob(1, `{"version":7}`, "h7"))
	require.NoError(t, err)
	reg, v := gc.LoadParsed()
	assert.Nil(t, reg)
	assert.Equal(t, uint64(0), v)

	// Committed (API create, or boot): applied synchronously.
	committed.set(&job.CRESettingsSpec{ConfigType: ConfigTypeCapRegistry, OffchainConfig: `{"version":7}`, Hash: "h7"})
	_, err = d.ServicesForSpec(t.Context(), capRegistryJob(1, `{"version":7}`, "h7"))
	require.NoError(t, err)
	raw, v := gc.Load()
	assert.Equal(t, `{"version":7}`, raw)
	assert.Equal(t, uint64(7), v)
}

// TestDelegate_CapRegistry_DeleteDoesNotClearBeforeCommit: OnDeleteJob runs inside the delete
// transaction. The payload stays applied until the delete is committed (and stays if it rolls
// back).
func TestDelegate_CapRegistry_DeleteDoesNotClearBeforeCommit(t *testing.T) {
	t.Parallel()

	d, committed := newTestDelegateWithCommitted(t)
	gc := d.capRegistry.GlobalConfig()
	committed.set(&job.CRESettingsSpec{ConfigType: ConfigTypeCapRegistry, OffchainConfig: `{"version":3}`, Hash: "h3"})
	_, err := d.ServicesForSpec(t.Context(), capRegistryJob(1, `{"version":3}`, "h3"))
	require.NoError(t, err)

	require.NoError(t, d.OnDeleteJob(t.Context(), capRegistryJob(1, `{"version":3}`, "h3")))
	_, v := gc.Load()
	assert.Equal(t, uint64(3), v, "uncommitted delete must not clear")

	// Rolled back: a refresh still sees the job.
	require.NoError(t, d.capRegistry.Refresh(t.Context()))
	_, v = gc.Load()
	assert.Equal(t, uint64(3), v)

	// Committed: the next refresh clears it.
	committed.set(nil)
	require.NoError(t, d.capRegistry.Refresh(t.Context()))
	reg, v := gc.LoadParsed()
	assert.Nil(t, reg)
	assert.Equal(t, uint64(0), v)
}

// TestDelegate_CapRegistry_NoInMemorySlot: capabilities_registry uniqueness is enforced by the
// database, so a rolled-back create cannot leave an in-memory reservation behind that would
// block the next replacement.
func TestDelegate_CapRegistry_NoInMemorySlot(t *testing.T) {
	t.Parallel()

	d, _ := newTestDelegateWithCommitted(t)
	// Job 1 was "created" in a transaction that then rolled back (never deleted).
	_, err := d.ServicesForSpec(t.Context(), capRegistryJob(1, `{"version":5}`, "h5"))
	require.NoError(t, err)
	_, err = d.ServicesForSpec(t.Context(), capRegistryJob(2, `{"version":6}`, "h6"))
	require.NoError(t, err)
	_, loaded := d.activeJobIDs.Load(ConfigTypeCapRegistry)
	assert.False(t, loaded)
}

func TestDelegate_CapRegistry_Unconfigured(t *testing.T) {
	t.Parallel()

	d := NewDelegate(logger.TestLogger(t), &loop.AtomicSettings{}, &loop.AtomicSettings{}, nil)
	_, err := d.ServicesForSpec(t.Context(), capRegistryJob(1, `{"version":1}`, "h1"))
	require.ErrorContains(t, err, "no offchain capabilities registry configured")
	require.NoError(t, d.OnDeleteJob(t.Context(), capRegistryJob(1, `{"version":1}`, "h1")))
}

func TestDelegate_CapRegistry_RefreshErrorIsNotFatal(t *testing.T) {
	t.Parallel()

	d, committed := newTestDelegateWithCommitted(t)
	committed.err = assert.AnError
	_, err := d.ServicesForSpec(t.Context(), capRegistryJob(1, `{"version":1}`, "h1"))
	require.NoError(t, err, "the projector retries in the background")
}

func TestDelegate_DifferentConfigTypesCoexist(t *testing.T) {
	t.Parallel()

	d := newTestDelegate(t)

	settingsJob := job.Job{ID: 10, Type: job.CRESettings, CRESettingsSpec: &job.CRESettingsSpec{Settings: `Foo = "bar"`, Hash: "hs"}}
	_, err := d.ServicesForSpec(t.Context(), settingsJob)
	require.NoError(t, err)
	_, err = d.ServicesForSpec(t.Context(), capRegistryJob(11, `{"version":1}`, "h1"))
	require.NoError(t, err)
	_, err = d.ServicesForSpec(t.Context(), cresettingsJob(12, shardAssignmentToml))
	require.NoError(t, err)

	// Deleting the capabilities_registry job does not free the settings slot.
	require.NoError(t, d.OnDeleteJob(t.Context(), capRegistryJob(11, `{"version":1}`, "h1")))
	_, err = d.ServicesForSpec(t.Context(), cresettingsJob(13, settingsToml))
	require.ErrorContains(t, err, "already active: 10")
}

func TestDelegate_FailedSettingsStoreReleasesSlot(t *testing.T) {
	t.Parallel()

	d := newTestDelegate(t)
	_, err := d.ServicesForSpec(t.Context(), cresettingsJob(1, `not = [valid toml`))
	require.Error(t, err)
	_, loaded := d.activeJobIDs.Load(ConfigTypeSettings)
	assert.False(t, loaded, "a rejected settings job must not hold the slot")
	_, err = d.ServicesForSpec(t.Context(), cresettingsJob(2, settingsToml))
	require.NoError(t, err)
}

func TestDelegate_EmbeddedCapRegistryConfigTypeRejected(t *testing.T) {
	t.Parallel()

	d := newTestDelegate(t)
	_, err := d.ServicesForSpec(t.Context(), cresettingsJob(1, `config_type = "capabilities_registry"
`))
	require.ErrorContains(t, err, "unknown config_type")
	require.ErrorContains(t, err, "top-level config_type")
}
