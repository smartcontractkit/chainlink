package cresettings

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func newTestDelegate(t *testing.T, gc *globalconfig.GlobalConfig) *delegate {
	return NewDelegate(logger.TestLogger(t), &loop.AtomicSettings{}, &loop.AtomicSettings{}, gc)
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

func TestDelegate_CapabilitiesRegistry_StoresIntoGlobalConfig(t *testing.T) {
	t.Parallel()

	gc := globalconfig.New()
	d := newTestDelegate(t, gc)

	_, err := d.ServicesForSpec(t.Context(), capRegistryJob(1, `{"version":7}`, "h7"))
	require.NoError(t, err)

	raw, v := gc.Load()
	assert.Equal(t, `{"version":7}`, raw)
	assert.Equal(t, uint64(7), v)
}

func TestDelegate_RejectsSecondJobOfSameConfigType(t *testing.T) {
	t.Parallel()

	d := newTestDelegate(t, globalconfig.New())

	_, err := d.ServicesForSpec(t.Context(), capRegistryJob(1, `{"version":1}`, "h1"))
	require.NoError(t, err)

	_, err = d.ServicesForSpec(t.Context(), capRegistryJob(2, `{"version":2}`, "h2"))
	require.ErrorContains(t, err, "already active")

	// Same job ID re-applying is allowed (idempotent).
	_, err = d.ServicesForSpec(t.Context(), capRegistryJob(1, `{"version":1}`, "h1"))
	require.NoError(t, err)
}

func TestDelegate_DifferentConfigTypesCoexist(t *testing.T) {
	t.Parallel()

	d := newTestDelegate(t, globalconfig.New())

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

	d := newTestDelegate(t, globalconfig.New())

	_, err := d.ServicesForSpec(t.Context(), capRegistryJob(1, `{"version":1}`, "h1"))
	require.NoError(t, err)

	require.NoError(t, d.OnDeleteJob(t.Context(), capRegistryJob(1, `{"version":1}`, "h1")))

	// A new job of the same config_type is now accepted.
	_, err = d.ServicesForSpec(t.Context(), capRegistryJob(2, `{"version":2}`, "h2"))
	require.NoError(t, err)
}
