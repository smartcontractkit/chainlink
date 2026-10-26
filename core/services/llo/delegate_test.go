package llo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	lloconfig "github.com/smartcontractkit/chainlink-data-streams/llo/pluginconfig"
	llov30 "github.com/smartcontractkit/chainlink-data-streams/llo/v30"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/telem"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

// newTestDelegate returns a delegate carrying only what v31FactoryParams reads.
func newTestDelegate(cfg DelegateConfig) *delegate {
	return &delegate{cfg: cfg, telem: telem.NewTelemeterService(telem.TelemeterParams{})}
}

func Test_delegate_v31FactoryParams(t *testing.T) {
	t.Parallel()

	t.Run("forwards the V31Config knobs", func(t *testing.T) {
		t.Parallel()

		params := newTestDelegate(DelegateConfig{
			DonID: 42,
			V31Config: lloconfig.V31Config{
				MaxSnapshotRounds:          7,
				BlobLifetimeRounds:         11,
				MaxDurationBlobObservation: lloconfig.Duration(3 * time.Second),
				BlobInFlightWaitFactor:     5,
				MaxBlobSnapshotAge:         lloconfig.Duration(9 * time.Second),
				MaxRoundPeriod:             lloconfig.Duration(time.Minute),
				CaptureStagingTelemetry:    true,
			},
		}).v31FactoryParams(logger.Test(t), nil)

		assert.True(t, params.CaptureStagingTelemetry)
		assert.Equal(t, uint64(7), params.MaxSnapshotRounds)
		assert.Equal(t, uint64(11), params.BlobLifetimeRounds)
		assert.Equal(t, 3*time.Second, params.MaxDurationBlobObservation)
		assert.Equal(t, uint64(5), params.BlobInFlightWaitFactor)
		assert.Equal(t, 9*time.Second, params.MaxBlobSnapshotAge)
		assert.Equal(t, time.Minute, params.MaxRoundPeriod)
		assert.Equal(t, uint32(42), params.DonID)
	})

	t.Run("leaves unset knobs zero so the plugin applies its defaults", func(t *testing.T) {
		t.Parallel()

		params := newTestDelegate(DelegateConfig{}).v31FactoryParams(logger.Test(t), nil)

		assert.Zero(t, params.MaxSnapshotRounds)
		assert.Zero(t, params.BlobLifetimeRounds)
		assert.Zero(t, params.MaxDurationBlobObservation)
		assert.Zero(t, params.BlobInFlightWaitFactor)
		assert.Zero(t, params.MaxBlobSnapshotAge)
		assert.Zero(t, params.MaxRoundPeriod)
		assert.Nil(t, params.AttributedObservationTelemetryCh)
	})

	t.Run("forwards a negative MaxBlobSnapshotAge, which disables the age check", func(t *testing.T) {
		t.Parallel()

		params := newTestDelegate(DelegateConfig{
			V31Config: lloconfig.V31Config{MaxBlobSnapshotAge: lloconfig.Duration(-time.Second)},
		}).v31FactoryParams(logger.Test(t), nil)

		assert.Equal(t, -time.Second, params.MaxBlobSnapshotAge)
	})

	t.Run("enables VerboseLogging from either the node or the job knob", func(t *testing.T) {
		t.Parallel()

		for _, tc := range []struct {
			name string
			node bool
			job  bool
			want bool
		}{
			{name: "neither", node: false, job: false, want: false},
			{name: "node only", node: true, job: false, want: true},
			{name: "job only", node: false, job: true, want: true},
			{name: "both", node: true, job: true, want: true},
		} {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				params := newTestDelegate(DelegateConfig{
					ReportingPluginConfig: llov30.Config{VerboseLogging: tc.node},
					V31Config:             lloconfig.V31Config{VerboseLogging: tc.job},
				}).v31FactoryParams(logger.Test(t), nil)

				assert.Equal(t, tc.want, params.VerboseLogging)
			})
		}
	})
}

// Test_delegate_pluginTelemetry checks that each plugin version only gets the
// telemetry channels its own config enables: v30 from the node driven flags, v31
// from V31Config alone.
func Test_delegate_pluginTelemetry(t *testing.T) {
	t.Parallel()

	type channels struct{ outcome, report, attributed bool }
	for _, tc := range []struct {
		name     string
		cfg      DelegateConfig
		v30, v31 channels
	}{
		{
			name: "all off",
		},
		{
			name: "node flags enable v30 only",
			cfg: DelegateConfig{
				CaptureEATelemetry:      true,
				CaptureOutcomeTelemetry: true,
				CaptureReportTelemetry:  true,
			},
			v30: channels{outcome: true, report: true},
		},
		{
			name: "v31 flags enable v31 only, without the node flag",
			cfg: DelegateConfig{
				V31Config: lloconfig.V31Config{
					CaptureOutcomeTelemetry:               true,
					CaptureReportTelemetry:                true,
					CaptureAttributedObservationTelemetry: true,
				},
			},
			v31: channels{outcome: true, report: true, attributed: true},
		},
		{
			name: "each flag independently",
			cfg: DelegateConfig{
				CaptureEATelemetry:      true,
				CaptureOutcomeTelemetry: true,
				V31Config:               lloconfig.V31Config{CaptureReportTelemetry: true},
			},
			v30: channels{outcome: true},
			v31: channels{report: true},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.cfg.PluginMonitoringEndpoint = nopTypedEndpoint{}
			d := &delegate{cfg: tc.cfg, telem: telem.NewTelemeterService(tc.cfg.telemeterParams(logger.Test(t)))}

			v30 := d.v30FactoryParams(logger.Test(t), nil)
			assert.Equal(t, tc.v30.outcome, v30.OutcomeTelemetryCh != nil, "v30 outcome")
			assert.Equal(t, tc.v30.report, v30.ReportTelemetryCh != nil, "v30 report")

			v31 := d.v31FactoryParams(logger.Test(t), nil)
			assert.Equal(t, tc.v31.outcome, v31.OutcomeTelemetryCh != nil, "v31 outcome")
			assert.Equal(t, tc.v31.report, v31.ReportTelemetryCh != nil, "v31 report")
			assert.Equal(t, tc.v31.attributed, v31.AttributedObservationTelemetryCh != nil, "v31 attributed observation")
		})
	}
}

type nopTypedEndpoint struct{}

func (nopTypedEndpoint) SendTypedLog(synchronization.TelemetryType, []byte) {}

// stubKeyValueDatabaseFactory and stubBinaryNetworkEndpoint2Factory stand in
// for the OCR3.1-only dependencies; validateInstances only checks that they
// are present.
type stubKeyValueDatabaseFactory struct{}

func (stubKeyValueDatabaseFactory) NewKeyValueDatabase(ocr2types.ConfigDigest) (ocr3_1types.KeyValueDatabase, error) {
	panic("not implemented")
}

func (stubKeyValueDatabaseFactory) NewKeyValueDatabaseIfExists(ocr2types.ConfigDigest) (ocr3_1types.KeyValueDatabase, error) {
	panic("not implemented")
}

type stubBinaryNetworkEndpoint2Factory struct{}

func (stubBinaryNetworkEndpoint2Factory) NewEndpoint(ocr2types.ConfigDigest, []string, []ocrcommontypes.BootstrapperLocator, ocr2types.BinaryNetworkEndpoint2Config, ocr2types.BinaryNetworkEndpoint2Config) (ocr2types.BinaryNetworkEndpoint2, error) {
	panic("not implemented")
}

func (stubBinaryNetworkEndpoint2Factory) PeerID() string { return "" }

func Test_DelegateConfig_validateInstances(t *testing.T) {
	t.Parallel()

	// trackers of length n; their contents are never read here.
	trackers := func(n int) []ocr2types.ContractConfigTracker {
		return make([]ocr2types.ContractConfigTracker, n)
	}

	t.Run("accepts one entry per tracker", func(t *testing.T) {
		t.Parallel()

		cfg := DelegateConfig{
			ContractConfigTrackers: trackers(2),
			PluginVersions:         []lloconfig.PluginVersion{lloconfig.PluginVersionV30, lloconfig.PluginVersionV30},
		}
		assert.NoError(t, cfg.validateInstances())
	})

	t.Run("rejects a length mismatch with the trackers", func(t *testing.T) {
		t.Parallel()

		cfg := DelegateConfig{
			ContractConfigTrackers: trackers(2),
			PluginVersions:         []lloconfig.PluginVersion{lloconfig.PluginVersionV30},
		}
		assert.ErrorContains(t, cfg.validateInstances(), "got 1 entries for 2 trackers")
	})

	t.Run("rejects an unsupported version", func(t *testing.T) {
		t.Parallel()

		cfg := DelegateConfig{
			ContractConfigTrackers: trackers(2),
			PluginVersions:         []lloconfig.PluginVersion{lloconfig.PluginVersionV30, "v32"},
		}
		assert.ErrorContains(t, cfg.validateInstances(), `unsupported plugin version for instance 1: "v32"`)
	})

	t.Run("rejects an empty version, which the config normalizes before it gets here", func(t *testing.T) {
		t.Parallel()

		cfg := DelegateConfig{
			ContractConfigTrackers: trackers(1),
			PluginVersions:         []lloconfig.PluginVersion{""},
		}
		assert.ErrorContains(t, cfg.validateInstances(), "unsupported plugin version for instance 0")
	})

	t.Run("requires the OCR3.1 dependencies when any instance is v31", func(t *testing.T) {
		t.Parallel()

		// A v30 production instance handing over to a v31 staging instance: the
		// dependencies are per job, so instance 1 alone makes them required.
		cfg := DelegateConfig{
			ContractConfigTrackers: trackers(2),
			PluginVersions:         []lloconfig.PluginVersion{lloconfig.PluginVersionV30, lloconfig.PluginVersionV31},
		}
		require.ErrorContains(t, cfg.validateInstances(), "KeyValueDatabaseFactory must not be nil")

		cfg.KeyValueDatabaseFactory = stubKeyValueDatabaseFactory{}
		require.ErrorContains(t, cfg.validateInstances(), "BinaryNetworkEndpoint2Factory must not be nil")

		cfg.BinaryNetworkEndpoint2Factory = stubBinaryNetworkEndpoint2Factory{}
		assert.NoError(t, cfg.validateInstances())
	})

	t.Run("does not require them when every instance is v30", func(t *testing.T) {
		t.Parallel()

		cfg := DelegateConfig{
			ContractConfigTrackers: trackers(2),
			PluginVersions:         []lloconfig.PluginVersion{lloconfig.PluginVersionV30, lloconfig.PluginVersionV30},
		}
		assert.NoError(t, cfg.validateInstances())
	})
}
