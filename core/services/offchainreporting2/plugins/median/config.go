package median

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

type PluginConfig struct {
	JuelsPerFeeCoinPipeline string `json:"juelsPerFeeCoinSource"`
}

func validatePluginConfig(config PluginConfig) error {
	if _, err := pipeline.Parse(config.JuelsPerFeeCoinPipeline); err != nil {
		return errors.Wrap(err, "invalid juelsPerFeeCoinSource pipeline")
	}

	return nil
}
