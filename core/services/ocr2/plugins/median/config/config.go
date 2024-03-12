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
	JuelsPerFeeCoinPipeline      string          `json:"juelsPerFeeCoinSource"`
	JuelsPerFeeCoinCacheDuration models.Interval `json:"juelsPerFeeCoinCacheDuration"`
	JuelsPerFeeCoinCacheDisabled bool            `json:"juelsPerFeeCoinCacheDisabled"`
}

// ValidatePluginConfig validates the arguments for the Median plugin.
func ValidatePluginConfig(config PluginConfig) error {
	if _, err := pipeline.Parse(config.JuelsPerFeeCoinPipeline); err != nil {
		return errors.Wrap(err, "invalid juelsPerFeeCoinSource pipeline")
	}

	// unset duration defaults later
	if config.JuelsPerFeeCoinCacheDuration != 0 {
		if config.JuelsPerFeeCoinCacheDuration.Duration() < time.Second*30 {
			return errors.Errorf("juelsPerFeeCoinSource cache duration: %s is below 30 second minimum", config.JuelsPerFeeCoinCacheDuration.Duration().String())
		} else if config.JuelsPerFeeCoinCacheDuration.Duration() > time.Minute*20 {
			return errors.Errorf("juelsPerFeeCoinSource cache duration: %s is above 20 minute maximum", config.JuelsPerFeeCoinCacheDuration.Duration().String())
		}
	}

	return nil
}
