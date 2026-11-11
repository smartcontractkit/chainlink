package cresettings

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	wfmonitoring "github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
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
	return NewDelegate(logger.TestLogger(t), &settings.Atomic{}, &settings.Atomic{}, globalconfig.New())
}

var shardAssignmentReader *sdkmetric.ManualReader

func TestMain(m *testing.M) {
	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	client := beholder.NoopClientConfig{}.New()
	client.Meter = mp.Meter("cresettings-test")
	client.MeterProvider = mp
	beholder.SetClient(client)
	shardAssignmentReader = reader
	code := m.Run()
	_ = mp.Shutdown(context.Background())
	os.Exit(code)
}

var arrivalTestHashCounter atomic.Uint32

// TestServicesForSpecShardAssignmentArrivalRecorded proves the delegate
// records the arrival of the exact raw blob it Stores, before Store: applying
// that blob through the same global metrics afterwards produces exactly one
// arrival→in-effect staleness sample.
func TestServicesForSpecShardAssignmentArrivalRecorded(t *testing.T) { //nolint:paralleltest // uses the process-global metrics singleton
	ctx := context.Background()
	shardSettings := &settings.Atomic{}
	d := NewDelegate(logger.TestLogger(t), &settings.Atomic{}, shardSettings, globalconfig.New())

	_, err := d.ServicesForSpec(ctx, cresettingsJob(1, shardAssignmentToml))
	require.NoError(t, err)

	u, err := shardSettings.Load()
	require.NoError(t, err)
	assert.Equal(t, shardAssignmentToml, u.Settings)

	m, err := wfmonitoring.GlobalShardAssignmentMetrics()
	require.NoError(t, err)
	// A unique parsed hash per run keeps -count reruns a real change; the
	// arrival slot matches on the raw blob alone.
	before := stalenessSampleCount(t)
	require.True(t, m.RecordSettingsApplied(ctx, arrivalTestHashCounter.Add(1), shardAssignmentToml))
	require.Equal(t, before+1, stalenessSampleCount(t))
}

func stalenessSampleCount(t *testing.T) uint64 {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, shardAssignmentReader.Collect(t.Context(), &rm))
	for _, sm := range rm.ScopeMetrics {
		for _, met := range sm.Metrics {
			if met.Name != "platform_shard_assignment_config_staleness_seconds" {
				continue
			}
			hist, ok := met.Data.(metricdata.Histogram[float64])
			require.True(t, ok)
			if len(hist.DataPoints) == 0 {
				return 0
			}
			return hist.DataPoints[0].Count
		}
	}
	return 0
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

func TestServicesForSpecExplicitSettingsConfigType(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	d := newTestDelegate(t)

	explicit := "config_type = \"settings\"\n" + settingsToml
	_, err := d.ServicesForSpec(ctx, cresettingsJob(1, explicit))
	require.NoError(t, err)

	got, err := d.atomicSettings.Load()
	require.NoError(t, err)
	assert.Equal(t, explicit, got.Settings)

	// explicit and implicit settings jobs share the same slot
	_, err = d.ServicesForSpec(ctx, cresettingsJob(2, settingsToml))
	require.ErrorContains(t, err, fmt.Sprintf("another %s job with config_type %q is already active: 1", job.CRESettings, ConfigTypeSettings))

	_, err = d.ServicesForSpec(ctx, cresettingsJob(3, shardAssignmentToml))
	require.NoError(t, err)
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
	d := NewDelegate(logger.TestLogger(t), &settings.Atomic{}, &settings.Atomic{}, gc)
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
	d := NewDelegate(logger.TestLogger(t), &settings.Atomic{}, &settings.Atomic{}, globalconfig.New())

	_, err := d.ServicesForSpec(ctx, cresettingsJob(1, "config_type = \"capabilities_registry\"\n"))
	require.Error(t, err)
	// The failed job must not hold the slot.
	_, err = d.ServicesForSpec(ctx, cresettingsJob(2, capRegSettings(capRegPayload)))
	require.NoError(t, err)
}
