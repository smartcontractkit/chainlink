package median

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

// The PluginConfig struct contains the custom arguments needed for the Median plugin.
type PluginConfig struct {
	JuelsPerFeeCoinPipeline string `json:"juelsPerFeeCoinSource"`
}

// validatePluginConfig validates the arguments for the Median plugin.
func validatePluginConfig(config PluginConfig) error {
	if _, err := pipeline.Parse(config.JuelsPerFeeCoinPipeline); err != nil {
		return errors.Wrap(err, "invalid juelsPerFeeCoinSource pipeline")
	}

	return nil
}
