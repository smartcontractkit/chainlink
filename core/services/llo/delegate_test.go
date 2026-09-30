package llo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	lloconfig "github.com/smartcontractkit/chainlink-data-streams/llo/pluginconfig"
	llov30 "github.com/smartcontractkit/chainlink-data-streams/llo/v30"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/telem"
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
			},
		}).v31FactoryParams(logger.Test(t), nil)

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
