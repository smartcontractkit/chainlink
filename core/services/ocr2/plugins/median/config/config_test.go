package config

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func TestValidatePluginConfig(t *testing.T) {
	type testCase struct {
		name          string
		pipeline      string
		cacheDuration models.Interval
		expectedError error
	}

	t.Run("pipeline validation", func(t *testing.T) {
		for _, tc := range []testCase{
			{"empty pipeline", "", models.Interval(time.Minute), fmt.Errorf("invalid juelsPerFeeCoinSource pipeline: empty pipeline")},
			{"blank pipeline", " ", models.Interval(time.Minute), fmt.Errorf("invalid juelsPerFeeCoinSource pipeline: empty pipeline")},
			{"foo pipeline", "foo", models.Interval(time.Minute), fmt.Errorf("invalid juelsPerFeeCoinSource pipeline: UnmarshalTaskFromMap: unknown task type: \"\"")},
		} {
			t.Run(tc.name, func(t *testing.T) {
				pc := PluginConfig{JuelsPerFeeCoinPipeline: tc.pipeline}
				assert.EqualError(t, pc.ValidatePluginConfig(), tc.expectedError.Error())
			})
		}
	})

	t.Run("cache duration validation", func(t *testing.T) {
		for _, tc := range []testCase{
			{"cache duration below minimum", `ds1 [type=bridge name=voter_turnout];`, models.Interval(time.Second * 29), fmt.Errorf("juelsPerFeeCoinSourceCache update interval: 29s is below 30 second minimum")},
			{"cache duration above maximum", `ds1 [type=bridge name=voter_turnout];`, models.Interval(time.Minute*20 + time.Second), fmt.Errorf("juelsPerFeeCoinSourceCache update interval: 20m1s is above 20 minute maximum")},
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
			{"valid duration and valid pipeline", `ds1 [type=bridge name=voter_turnout];`, models.Interval(time.Second * 30), nil},
			{"valid duration and valid pipeline", `ds1 [type=bridge name=voter_turnout];`, models.Interval(time.Minute * 20), nil},
		} {
			t.Run(s.name, func(t *testing.T) {
				pc := PluginConfig{JuelsPerFeeCoinPipeline: s.pipeline}
				assert.Nil(t, pc.ValidatePluginConfig())
			})
		}
	})
}
