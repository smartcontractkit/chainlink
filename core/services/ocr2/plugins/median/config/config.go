// config is a separate package so that we can validate
// the config in other packages, for example in job at job create time.

package config

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

// The PluginConfig struct contains the custom arguments needed for the Median plugin.
type PluginConfig struct {
	JuelsPerFeeCoinPipeline string `json:"juelsPerFeeCoinSource"`
}

// ValidatePluginConfig validates the arguments for the Median plugin.
func ValidatePluginConfig(config PluginConfig) error {
	if _, err := pipeline.Parse(config.JuelsPerFeeCoinPipeline); err != nil {
		return errors.Wrap(err, "invalid juelsPerFeeCoinSource pipeline")
	}

	return nil
}
