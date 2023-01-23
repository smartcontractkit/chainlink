package evm

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/ocr2keepers/pkg/types"

	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
)

type HeadProvider struct {
	ht httypes.HeadTracker
	hb httypes.HeadBroadcaster
}

// OnNewHead should continue running until the context ends
func (hw *HeadProvider) OnNewHead(ctx context.Context, f func(blockKey types.BlockKey)) error {
	_, _ = hw.hb.Subscribe(&headWrapper{f: f})
	<-ctx.Done()
	return nil
}

func (hw *HeadProvider) LatestBlock() int64 {
	lc := hw.ht.LatestChain()
	if lc == nil {
		return 0
	}
	return lc.Number
}

type headWrapper struct {
	f func(blockKey types.BlockKey)
}

func (w *headWrapper) OnNewLongestChain(ctx context.Context, head *evmtypes.Head) {
	var bl int64
	if head != nil {
		bl = head.Number
	}
	w.f(types.BlockKey(fmt.Sprintf("%d", bl)))
}
