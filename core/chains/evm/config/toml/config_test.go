package toml_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

func TestEVMConfig_ValidateConfig(t *testing.T) {
	var err error
	for _, id := range toml.DefaultIDs {
		t.Run(fmt.Sprintf("chainID-%s", id), func(t *testing.T) {
			evmCfg := &toml.EVMConfig{
				ChainID: id,
				Chain:   toml.Defaults(id),
			}

			err = multierr.Append(err, evmCfg.ValidateConfig())
			err = multierr.Append(err, evmCfg.Chain.ValidateConfig())
			err = multierr.Append(err, evmCfg.Chain.GasEstimator.ValidateConfig())

			assert.NoError(t, err)
		})
	}
}
