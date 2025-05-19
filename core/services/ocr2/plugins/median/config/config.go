// config is a separate package so that we can validate
// the config in other packages, for example in job at job create time.

package config

import (
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type DeviationFunctionDefinition map[string]any

// The PluginConfig struct contains the custom arguments needed for the Median plugin.
// To avoid a catastrophic libocr codec error, you must make sure that either all nodes in the same DON
// (1) have no GasPriceSubunitsPipeline or all nodes in the same DON (2) have a GasPriceSubunitsPipeline
type PluginConfig struct {
	GasPriceSubunitsPipeline string `json:"gasPriceSubunitsSource"`
	JuelsPerFeeCoinPipeline  string `json:"juelsPerFeeCoinSource"`
	// JuelsPerFeeCoinCache is disabled when nil
	JuelsPerFeeCoinCache        *JuelsPerFeeCoinCache       `json:"juelsPerFeeCoinCache"`
	DeviationFunctionDefinition DeviationFunctionDefinition `json:"deviationFunc"`
}

type JuelsPerFeeCoinCache struct {
	Disable                 bool            `json:"disable"`
	UpdateInterval          models.Interval `json:"updateInterval"`
	StalenessAlertThreshold models.Interval `json:"stalenessAlertThreshold"`
}

// ValidatePluginConfig validates the arguments for the Median plugin.
func (config *PluginConfig) ValidatePluginConfig() error {
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

	// Gas price pipeline is optional
	if !config.HasGasPriceSubunitsPipeline() {
		return nil
	} else if _, err := pipeline.Parse(config.GasPriceSubunitsPipeline); err != nil {
		return errors.Wrap(err, "invalid gasPriceSubunitsSource pipeline")
	}

	return nil
}

func (config *PluginConfig) HasGasPriceSubunitsPipeline() bool {
	return strings.TrimSpace(config.GasPriceSubunitsPipeline) != ""
}
