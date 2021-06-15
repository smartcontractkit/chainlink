package offchainreporting

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/core/chains"
)

// BlockTranslator converts emitted block numbers (from block.number) into a
// block number range suitable for query in FilterLogs
type BlockTranslator interface {
	NumberToQueryRange(changedInBlock uint64) (fromBlock *big.Int, toBlock *big.Int)
	Start()
	Close()
}

// NewBlockTranslator returns the block translator for the given chain
func NewBlockTranslator(chain *chains.Chain) BlockTranslator {
	if chain.IsArbitrum() {
		return newArbitrumBlockTranslator()
	} else if chain.IsOptimism() {
		return newOptimismBlockTranslator()
	}
	return &l1BlockTranslator{}
}

type l1BlockTranslator struct{}

func (*l1BlockTranslator) NumberToQueryRange(changedInBlock uint64) (fromBlock *big.Int, toBlock *big.Int) {
	return big.NewInt(int64(changedInBlock)), big.NewInt(int64(changedInBlock))
}
func (*l1BlockTranslator) Start() {}
func (*l1BlockTranslator) Close() {}

type arbitrumBlockTranslator struct{}

func newArbitrumBlockTranslator() *arbitrumBlockTranslator {
	return &arbitrumBlockTranslator{}
}

func (*arbitrumBlockTranslator) NumberToQueryRange(changedInBlock uint64) (fromBlock *big.Int, toBlock *big.Int) {
	// TODO: OPTIMISE THIS
	// TODO: Logic goes here that:
	// 1. Subscribes to SequencerBatchDeliveredFromOrigin on https://github.com/OffchainLabs/arbitrum/blob/next/packages/arb-bridge-eth/contracts/bridge/SequencerInbox.sol#L138
	// 2. Updates local state with mapping of L1 range -> L2 range
	// 3. Includes some sort of database persistence?
	// Currently we simply query the entire block range. This is correct, but very slow and suboptimal
	// See: https://app.clubhouse.io/chainlinklabs/story/11270/optimise-blocktranslator-ocr-for-optimism-and-arbitrum
	return big.NewInt(0), nil
}
func (a *arbitrumBlockTranslator) Start() {}
func (a *arbitrumBlockTranslator) Close() {}

type optimismBlockTranslator struct{}

func newOptimismBlockTranslator() *optimismBlockTranslator {
	return &optimismBlockTranslator{}
}

func (*optimismBlockTranslator) NumberToQueryRange(changedInBlock uint64) (fromBlock *big.Int, toBlock *big.Int) {
	// TODO: OPTIMISE THIS
	// Currently we simply query the entire block range. This is correct, but very slow and suboptimal
	// See: https://app.clubhouse.io/chainlinklabs/story/11270/optimise-blocktranslator-ocr-for-optimism-and-arbitrum
	return big.NewInt(0), nil
}
func (a *optimismBlockTranslator) Start() {}
func (a *optimismBlockTranslator) Close() {}
