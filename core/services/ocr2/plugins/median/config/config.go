// config is a separate package so that we can validate
// the config in other packages, for example in job at job create time.

package config

import (
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// The PluginConfig struct contains the custom arguments needed for the Median plugin.
type PluginConfig struct {
	JuelsPerFeeCoinPipeline string `json:"juelsPerFeeCoinSource"`
	// JuelsPerFeeCoinCache is disabled when nil
	JuelsPerFeeCoinCache *JuelsPerFeeCoinCache `json:"juelsPerFeeCoinCache"`
}

type JuelsPerFeeCoinCache struct {
	UpdateInterval          models.Interval `json:"updateInterval"`
	StalenessAlertThreshold models.Interval `json:"stalenessAlertThreshold"`
}

// ValidatePluginConfig validates the arguments for the Median plugin.
func ValidatePluginConfig(config PluginConfig) error {
	if _, err := pipeline.Parse(config.JuelsPerFeeCoinPipeline); err != nil {
		return errors.Wrap(err, "invalid juelsPerFeeCoinSource pipeline")
	}

	// unset durations have a default set late
	if config.JuelsPerFeeCoinCache != nil {
		updateInterval := config.JuelsPerFeeCoinCache.UpdateInterval.Duration()
		if updateInterval != 0 && updateInterval < time.Second*30 {
			return errors.Errorf("juelsPerFeeCoinSourceCache update interval: %s is below 30 second minimum", updateInterval.String())
		} else if updateInterval > time.Minute*20 {
			return errors.Errorf("juelsPerFeeCoinSourceCache update interval: %s is above 20 minute maximum", updateInterval.String())
		}
	}

	return nil
}
