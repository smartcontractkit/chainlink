package ocrcommon

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/common/config"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// BlockTranslator converts emitted block numbers (from block.number) into a
// block number range suitable for query in FilterLogs
type BlockTranslator interface {
	NumberToQueryRange(ctx context.Context, changedInL1Block uint64) (fromBlock *big.Int, toBlock *big.Int)
}

// NewBlockTranslator returns the block translator for the given chain
func NewBlockTranslator(cfg Config, client evmclient.Client, lggr logger.Logger) BlockTranslator {
	switch cfg.ChainType() {
	case config.ChainArbitrum:
		return NewArbitrumBlockTranslator(client, lggr)
	case "", config.ChainCelo, config.ChainGnosis, config.ChainKroma, config.ChainMetis, config.ChainOptimismBedrock, config.ChainScroll, config.ChainWeMix, config.ChainXDai, config.ChainXLayer, config.ChainZkSync:
		fallthrough
	default:
		return &l1BlockTranslator{}
	}
}

type l1BlockTranslator struct{}

func (*l1BlockTranslator) NumberToQueryRange(_ context.Context, changedInL1Block uint64) (fromBlock *big.Int, toBlock *big.Int) {
	return big.NewInt(int64(changedInL1Block)), big.NewInt(int64(changedInL1Block))
}

func (*l1BlockTranslator) OnNewLongestChain(context.Context, *evmtypes.Head) {}
