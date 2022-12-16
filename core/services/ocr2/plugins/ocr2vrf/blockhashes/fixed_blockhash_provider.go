package blockhashes

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ocr2vrf/types"

	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils/mathutil"
)

// fixedBlockhashProvider returns blockhashes with fixed-sized windows relative to current block height
type fixedBlockhashProvider struct {
	lp     logpoller.LogPoller
	logger logger.Logger
	// start block = current head - lookbackBlocks
	lookbackBlocks uint64
}

var _ types.Blockhashes = (*fixedBlockhashProvider)(nil)

func NewFixedBlockhashProvider(logPoller logpoller.LogPoller, lggr logger.Logger, lookbackBlocks uint64) types.Blockhashes {
	return &fixedBlockhashProvider{
		lp:             logPoller,
		logger:         lggr,
		lookbackBlocks: mathutil.Min(uint64(256), lookbackBlocks),
	}
}

func (b *fixedBlockhashProvider) OnchainVerifiableBlocks(
	ctx context.Context,
) (startHeight uint64, hashes []common.Hash, err error) {
	toBlock, err := b.CurrentHeight(ctx)
	if err != nil {
		return 0, nil, err
	}

	fromBlock := uint64(0)
	if toBlock > b.lookbackBlocks {
		fromBlock = toBlock - b.lookbackBlocks + 1
	}

	var blockHeights []uint64
	for i := fromBlock; i <= toBlock; i++ {
		blockHeights = append(blockHeights, i)
	}

	var blockhashes []common.Hash

	heads, err := b.lp.GetBlocksRange(ctx, blockHeights, pg.WithParentCtx(ctx))
	if err != nil {
		return 0, nil, err
	}

	for _, h := range heads {
		blockhashes = append(blockhashes, h.BlockHash)
	}

	return fromBlock, blockhashes, nil
}

func (b *fixedBlockhashProvider) CurrentHeight(ctx context.Context) (uint64, error) {
	head, err := b.lp.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return 0, err
	}
	return uint64(head), nil
}
