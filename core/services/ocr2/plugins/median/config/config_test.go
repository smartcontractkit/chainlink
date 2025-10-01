package config

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

func TestValidatePluginConfig(t *testing.T) {
	type testCase struct {
		name          string
		pipeline      string
		cacheDuration sqlutil.Interval
		expectedError error
	}

	t.Run("pipeline validation", func(t *testing.T) {
		for _, tc := range []testCase{
			{"empty pipeline", "", sqlutil.Interval(time.Minute), errors.New("invalid juelsPerFeeCoinSource pipeline: empty pipeline")},
			{"blank pipeline", " ", sqlutil.Interval(time.Minute), errors.New("invalid juelsPerFeeCoinSource pipeline: empty pipeline")},
			{"foo pipeline", "foo", sqlutil.Interval(time.Minute), errors.New("invalid juelsPerFeeCoinSource pipeline: UnmarshalTaskFromMap: unknown task type: \"\"")},
		} {
			t.Run(tc.name, func(t *testing.T) {
				pc := PluginConfig{JuelsPerFeeCoinPipeline: tc.pipeline}
				assert.EqualError(t, pc.ValidatePluginConfig(), tc.expectedError.Error())
			})
		}
	})

	t.Run("cache duration validation", func(t *testing.T) {
		for _, tc := range []testCase{
			{"cache duration below minimum", `ds1 [type=bridge name=voter_turnout];`, sqlutil.Interval(time.Second * 29), errors.New("juelsPerFeeCoinSourceCache update interval: 29s is below 30 second minimum")},
			{"cache duration above maximum", `ds1 [type=bridge name=voter_turnout];`, sqlutil.Interval(time.Minute*20 + time.Second), errors.New("juelsPerFeeCoinSourceCache update interval: 20m1s is above 20 minute maximum")},
		} {
			t.Run(tc.name, func(t *testing.T) {
				pc := PluginConfig{JuelsPerFeeCoinPipeline: tc.pipeline, JuelsPerFeeCoinCache: &JuelsPerFeeCoinCache{UpdateInterval: tc.cacheDuration}}
				assert.EqualError(t, pc.ValidatePluginConfig(), tc.expectedError.Error())
			})
		}
	})

	t.Run("valid values", func(t *testing.T) {
		for _, s := range []testCase{
			{"valid 0 cache duration and valid pipeline", `ds1 [type=bridge name=voter_turnout];`, 0, nil},
			{"valid duration and valid pipeline", `ds1 [type=bridge name=voter_turnout];`, sqlutil.Interval(time.Second * 30), nil},
			{"valid duration and valid pipeline", `ds1 [type=bridge name=voter_turnout];`, sqlutil.Interval(time.Minute * 20), nil},
		} {
			t.Run(s.name, func(t *testing.T) {
				pc := PluginConfig{JuelsPerFeeCoinPipeline: s.pipeline}
				assert.NoError(t, pc.ValidatePluginConfig())
			})
		}
	})
}
