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
		return newArbitrumBlockTranslator(chain)
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

type arbitrumBlockTranslator struct {
	min int64
}

func newArbitrumBlockTranslator(chain *chains.Chain) *arbitrumBlockTranslator {
	var min int64
	if chain == chains.ArbitrumMainnet {
		// sequencer contract deployed in https://etherscan.io/tx/0xd71d3c90adcce0fabc903fceed668561c92f5be5d8837295f9e46e2f6d99894e
		min = 12525700
	} else if chain == chains.ArbitrumRinkeby {
		// sequencer inbox was deployed in https://rinkeby.etherscan.io/tx/0x01eb72f978399a61549d71fe723e3d9943d1314717e0df10c902b2b2256fc974
		min = 8700589
	}
	return &arbitrumBlockTranslator{min}
}

func (a *arbitrumBlockTranslator) NumberToQueryRange(changedInBlock uint64) (fromBlock *big.Int, toBlock *big.Int) {
	// TODO: OPTIMISE THIS
	// TODO: Logic goes here that:
	// 1. Subscribes to SequencerBatchDeliveredFromOrigin on https://github.com/OffchainLabs/arbitrum/blob/next/packages/arb-bridge-eth/contracts/bridge/SequencerInbox.sol#L138
	// 2. Updates local state with mapping of L1 range -> L2 range
	// 3. Includes some sort of database persistence?
	// Currently we simply query the entire block range. This is correct, but very slow and suboptimal
	// See: https://app.clubhouse.io/chainlinklabs/story/11270/optimise-blocktranslator-ocr-for-optimism-and-arbitrum

	// NOTE: Mainnet l2
	return big.NewInt(a.min), nil
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
