package headtracker

import (
	"context"
	"fmt"
	"math/big"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// simulatedHeadTracker - simplified version of HeadTracker that works with simulated backed
type simulatedHeadTracker struct {
	ec             evmclient.Client
	useFinalityTag bool
	finalityDepth  int64
	ctx            context.Context
}

func NewSimulatedHeadTracker(ctx context.Context, ec evmclient.Client, useFinalityTag bool, finalityDepth int64) *simulatedHeadTracker {
	return &simulatedHeadTracker{
		ec:             ec,
		useFinalityTag: useFinalityTag,
		finalityDepth:  finalityDepth,
		ctx:            ctx,
	}
}

func (ht *simulatedHeadTracker) ChainWithLatestFinalized() (*evmtypes.Head, error) {
	latestHead, err := ht.ec.HeadByNumber(ht.ctx, nil)
	if err != nil {
		return nil, err
	}

	var finalizedBlock *evmtypes.Head
	if ht.useFinalityTag {
		finalizedBlock, err = ht.ec.LatestFinalizedBlock(ht.ctx)
	} else {
		finalizedBlock, err = ht.ec.HeadByNumber(ht.ctx, big.NewInt(max(latestHead.Number-ht.finalityDepth, 0)))
	}

	if err != nil {
		return nil, fmt.Errorf("simulatedHeadTracker failed to get finalized block")
	}

	finalizedBlock.IsFinalized = true
	latestHead.Parent = finalizedBlock
	return latestHead, nil
}
