package headtracker

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	evmclient "github.com/smartcontractkit/chainlink-integrations/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink-integrations/evm/types"
)

// simulatedHeadTracker - simplified version of HeadTracker that works with simulated backed
type simulatedHeadTracker struct {
	ec             evmclient.Client
	useFinalityTag bool
	finalityDepth  int64
}

func NewSimulatedHeadTracker(ec evmclient.Client, useFinalityTag bool, finalityDepth int64) *simulatedHeadTracker {
	return &simulatedHeadTracker{
		ec:             ec,
		useFinalityTag: useFinalityTag,
		finalityDepth:  finalityDepth,
	}
}

func (ht *simulatedHeadTracker) LatestAndFinalizedBlock(ctx context.Context) (*evmtypes.Head, *evmtypes.Head, error) {
	latest, err := ht.ec.HeadByNumber(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	if latest == nil {
		return nil, nil, fmt.Errorf("expected latest block to be valid")
	}

	var finalizedBlock *evmtypes.Head
	if ht.useFinalityTag {
		finalizedBlock, err = ht.ec.LatestFinalizedBlock(ctx)
	} else {
		finalizedBlock, err = ht.ec.HeadByNumber(ctx, big.NewInt(max(latest.Number-ht.finalityDepth, 0)))
	}

	if err != nil {
		return nil, nil, fmt.Errorf("simulatedHeadTracker failed to get finalized block")
	}

	if finalizedBlock == nil {
		return nil, nil, fmt.Errorf("expected finalized block to be valid")
	}

	return latest, finalizedBlock, nil
}

func (ht *simulatedHeadTracker) LatestChain() *evmtypes.Head {
	return nil
}

func (ht *simulatedHeadTracker) HealthReport() map[string]error {
	return nil
}

func (ht *simulatedHeadTracker) Start(_ context.Context) error {
	return nil
}

func (ht *simulatedHeadTracker) Close() error {
	return nil
}

func (ht *simulatedHeadTracker) Backfill(_ context.Context, _ *evmtypes.Head) error {
	return errors.New("unimplemented")
}

func (ht *simulatedHeadTracker) Name() string {
	return "SimulatedHeadTracker"
}

func (ht *simulatedHeadTracker) Ready() error {
	return nil
}
