package block_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/config"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	evmbig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"

	"github.com/smartcontractkit/chainlink/v2/block"
)

func Test_BlockTranslator(t *testing.T) {
	t.Parallel()

	ethClient := clienttest.NewClient(t)
	ctx := testutils.Context(t)
	lggr := logger.Test(t)

	t.Run("for L1 chains, returns the block changed argument", func(t *testing.T) {
		bt := block.NewBlockTranslator(ChainEthMainnet(t).EVM().ChainType(), ethClient, lggr)

		from, to := bt.NumberToQueryRange(ctx, 42)

		assert.Equal(t, big.NewInt(42), from)
		assert.Equal(t, big.NewInt(42), to)
	})

	t.Run("for optimism, uses the default translator", func(t *testing.T) {
		bt := block.NewBlockTranslator(ChainOptimismMainnet(t).EVM().ChainType(), ethClient, lggr)
		from, to := bt.NumberToQueryRange(ctx, 42)
		assert.Equal(t, big.NewInt(42), from)
		assert.Equal(t, big.NewInt(42), to)
	})

	t.Run("for arbitrum, returns the ArbitrumBlockTranslator", func(t *testing.T) {
		bt := block.NewBlockTranslator(ChainArbitrumMainnet(t).EVM().ChainType(), ethClient, lggr)
		assert.IsType(t, &block.ArbitrumBlockTranslator{}, bt)

		bt = block.NewBlockTranslator(ChainArbitrumRinkeby(t).EVM().ChainType(), ethClient, lggr)
		assert.IsType(t, &block.ArbitrumBlockTranslator{}, bt)
	})
}

func ChainEthMainnet(t *testing.T) config.ChainScopedConfig      { return scopedConfig(t, 1) }
func ChainOptimismMainnet(t *testing.T) config.ChainScopedConfig { return scopedConfig(t, 10) }
func ChainArbitrumMainnet(t *testing.T) config.ChainScopedConfig { return scopedConfig(t, 42161) }
func ChainArbitrumRinkeby(t *testing.T) config.ChainScopedConfig { return scopedConfig(t, 421611) }

func scopedConfig(t *testing.T, chainID int64) config.ChainScopedConfig {
	id := evmbig.NewI(chainID)
	evmCfg := toml.EVMConfig{ChainID: id, Chain: toml.Defaults(id)}
	return config.NewTOMLChainScopedConfig(&evmCfg)
}
