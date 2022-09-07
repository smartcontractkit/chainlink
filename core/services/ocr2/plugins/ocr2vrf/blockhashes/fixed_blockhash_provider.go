package blockhashes

import (
	"context"
	"math/big"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/ocr2vrf/types"

	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// fixedBlockhashProvider returns blockhashes with fixed-sized windows relative to current block height
type fixedBlockhashProvider struct {
	client client.Client
	lp     logpoller.LogPoller
	logger logger.Logger
	// start block = current head - lookbackBlocks
	lookbackBlocks uint64
	// number of blocks to query in a batch
	batchSize uint64
}

var _ types.Blockhashes = (*fixedBlockhashProvider)(nil)

func NewFixedBlockhashProvider(client client.Client, logPoller logpoller.LogPoller, lggr logger.Logger, lookbackBlocks, batchSize uint64) types.Blockhashes {
	return &fixedBlockhashProvider{
		client:         client,
		lp:             logPoller,
		logger:         lggr,
		lookbackBlocks: lookbackBlocks,
		batchSize:      batchSize,
	}
}

func (b *fixedBlockhashProvider) OnchainVerifiableBlocks(
	ctx context.Context,
) (startHeight uint64, hashes []common.Hash, err error) {
	toBlock, err := b.CurrentHeight(ctx)
	if err != nil {
		return 0, nil, errors.Wrap(err, "current height")
	}

	fromBlock := uint64(0)
	if toBlock > b.lookbackBlocks {
		fromBlock = toBlock - b.lookbackBlocks
	}

	var blockHeights []uint64
	for i := fromBlock; i <= toBlock; i++ {
		blockHeights = append(blockHeights, i)
	}

	blockhashes := make([]common.Hash, len(blockHeights))

	heads, err := b.lp.GetBlocks(blockHeights, pg.WithParentCtx(ctx))
	if err != nil {
		b.logger.Warnw("error in get blocks", "err", err)
	}

	for _, h := range heads {
		blockhashes[h.BlockNumber-int64(fromBlock)] = h.BlockHash
	}

	// Fallback to RPC for blocks not found in log poller
	var reqs []rpc.BatchElem
	for i, bh := range blockhashes {
		if bh == utils.EmptyHash {
			req := rpc.BatchElem{
				Method: "eth_getBlockByNumber",
				Args:   []interface{}{hexutil.EncodeBig(big.NewInt(int64(i) + int64(fromBlock))), false},
				Result: &evmtypes.Head{},
			}
			reqs = append(reqs, req)
		}
	}

	for i := 0; i < len(reqs); i += int(b.batchSize) {
		j := i + int(b.batchSize)
		if j > len(reqs) {
			j = len(reqs)
		}
		err := b.client.BatchCallContext(ctx, reqs[i:j])
		if err != nil {
			return 0, nil, errors.Wrap(err, "batch call context eth_getBlockByNumber")
		}
	}

	for _, r := range reqs {
		if r.Error != nil {
			return 0, nil, errors.Wrap(r.Error, "error found in eth_getBlockByNumber response")
		}
		block, is := r.Result.(*evmtypes.Head)

		if !is {
			return 0, nil, errors.Errorf("expected result to be a %T, got %T", &evmtypes.Head{}, r.Result)
		}
		if block == nil {
			return 0, nil, errors.New("invariant violation: got nil block")
		}
		if block.Hash == utils.EmptyHash {
			return 0, nil, errors.New("missing block hash")
		}
		blockhashes[block.Number-int64(fromBlock)] = block.Hash
	}

	return fromBlock, blockhashes, nil
}

func (b *fixedBlockhashProvider) CurrentHeight(ctx context.Context) (uint64, error) {
	head, err := b.lp.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return 0, errors.Wrap(err, "latest block")
	}
	return uint64(head), nil
}
