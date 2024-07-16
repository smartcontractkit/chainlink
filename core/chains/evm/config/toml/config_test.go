package toml_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

func TestEVMConfig_ValidateConfig(t *testing.T) {
	name := "fake"
	for _, id := range toml.DefaultIDs {
		t.Run(fmt.Sprintf("chainID-%s", id), func(t *testing.T) {
			evmCfg := &toml.EVMConfig{
				ChainID: id,
				Chain:   toml.Defaults(id),
				Nodes: toml.EVMNodes{{
					Name:    &name,
					WSURL:   config.MustParseURL("wss://foo.test/ws"),
					HTTPURL: config.MustParseURL("http://foo.test"),
				}},
			}

			assert.NoError(t, config.Validate(evmCfg))
		})
	}
}
