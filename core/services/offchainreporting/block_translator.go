package offchainreporting

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
)

// BlockTranslator converts emitted block numbers (from block.number) into a
// block number range suitable for query in FilterLogs
type BlockTranslator interface {
	NumberToQueryRange(ctx context.Context, changedInL1Block uint64) (fromBlock *big.Int, toBlock *big.Int)
}

// NewBlockTranslator returns the block translator for the given chain
func NewBlockTranslator(cfg Config, client eth.Client, lggr logger.Logger) BlockTranslator {
	switch cfg.ChainType() {
	case chains.Arbitrum:
		return NewArbitrumBlockTranslator(client, lggr)
	case chains.Optimism:
		return newOptimismBlockTranslator()

	case chains.XDai, chains.ExChain:
		fallthrough
	default:
		return &l1BlockTranslator{}
	}
}

type l1BlockTranslator struct{}

func (*l1BlockTranslator) NumberToQueryRange(_ context.Context, changedInL1Block uint64) (fromBlock *big.Int, toBlock *big.Int) {
	return big.NewInt(int64(changedInL1Block)), big.NewInt(int64(changedInL1Block))
}

func (*l1BlockTranslator) OnNewLongestChain(context.Context, eth.Head) {}

type optimismBlockTranslator struct{}

func newOptimismBlockTranslator() *optimismBlockTranslator {
	return &optimismBlockTranslator{}
}

func (*optimismBlockTranslator) NumberToQueryRange(_ context.Context, changedInL1Block uint64) (fromBlock *big.Int, toBlock *big.Int) {
	// TODO: OPTIMISE THIS
	// Currently we simply query the entire block range. This is correct, but very slow and suboptimal
	// https://app.clubhouse.io/chainlinklabs/story/11524/optimise-blocktranslator-ocr-for-optimism
	return big.NewInt(0), nil
}
