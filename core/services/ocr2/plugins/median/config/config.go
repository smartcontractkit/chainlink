// config is a separate package so that we can validate
// the config in other packages, for example in job at job create time.

package config

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

// The PluginConfig struct contains the custom arguments needed for the Median plugin.
type PluginConfig struct {
	JuelsPerFeeCoinPipeline string `json:"juelsPerFeeCoinSource"`
	GasPricePipeline        string `json:"gasPriceSource"`
}

// ValidatePluginConfig validates the arguments for the Median plugin.
func (config *PluginConfig) ValidatePluginConfig() error {
	if _, err := pipeline.Parse(config.JuelsPerFeeCoinPipeline); err != nil {
		return errors.Wrap(err, "invalid juelsPerFeeCoinSource pipeline")
	}

	// Gas price pipeline is optional
	if !config.GasPricePipelineExists() {
		return nil
	} else if _, err := pipeline.Parse(config.GasPricePipeline); err != nil {
		return errors.Wrap(err, "invalid gasPriceSource pipeline")
	}

	return nil
}

func (config *PluginConfig) GasPricePipelineExists() bool {
	return !(strings.TrimSpace(config.GasPricePipeline) == "")
}
