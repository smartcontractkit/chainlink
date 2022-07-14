package ocr2vrf

import (
	"context"
	"math/big"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/ocr2vrf/types"
)

// fixedBlockhashProvider returns blockhashes with fixed-sized windows relative to current block height
type fixedBlockhashProvider struct {
	client client.Client
	// start block = current head - lookbackBlocks
	lookbackBlocks uint
	// number of blocks to query in a batch
	batchSize uint
}

var _ types.Blockhashes = (*fixedBlockhashProvider)(nil)

func NewFixedBlockhashProvider(client client.Client, lookbackBlocks, batchSize uint) types.Blockhashes {
	return &fixedBlockhashProvider{
		client,
		lookbackBlocks,
		batchSize,
	}
}

func (b *fixedBlockhashProvider) OnchainVerifiableBlocks(
	ctx context.Context,
) (startHeight uint64, hashes []common.Hash, err error) {
	toBlock, err := b.CurrentHeight(ctx)
	if err != nil {
		return 0, nil, errors.Wrap(err, "current height")
	}
	fromBlock := toBlock - uint64(b.lookbackBlocks)
	var reqs []rpc.BatchElem

	for i := fromBlock; i <= toBlock; i++ {
		req := rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{hexutil.EncodeBig(big.NewInt(int64(i))), true},
			Result: &evmtypes.Head{},
		}
		reqs = append(reqs, req)
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

	var blockhashes []common.Hash
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
		if block.Hash == (common.Hash{}) {
			return 0, nil, errors.New("missing block hash")
		}
		blockhashes = append(blockhashes, block.Hash)
	}

	return fromBlock, blockhashes, nil
}

func (b *fixedBlockhashProvider) CurrentHeight(ctx context.Context) (uint64, error) {
	h, err := b.client.HeadByNumber(ctx, nil)
	if err != nil {
		return 0, errors.Wrap(err, "head by number")
	}
	if h.Number < 0 {
		return 0, errors.Errorf("unexpected head number: %d", h.Number)
	}
	return uint64(h.Number), nil
}
