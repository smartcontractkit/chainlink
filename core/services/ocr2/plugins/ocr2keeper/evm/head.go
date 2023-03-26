package evm

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/ocr2keepers/pkg/chain"
	"github.com/smartcontractkit/ocr2keepers/pkg/types"

	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
)

type HeadProvider struct {
	ht         httypes.HeadTracker
	hb         httypes.HeadBroadcaster
	chHead     chan types.BlockKey
	subscribed bool
}

// HeadTicker provides external access to the heads channel
func (hw *HeadProvider) HeadTicker() chan types.BlockKey {
	if !hw.subscribed {
		_, _ = hw.hb.Subscribe(&headWrapper{c: hw.chHead})
		hw.subscribed = true
	}
	return hw.chHead
}

func (hw *HeadProvider) LatestBlock() int64 {
	lc := hw.ht.LatestChain()
	if lc == nil {
		return 0
	}
	return lc.Number
}

// send does a non-blocking send of the key on c.
func send(c chan types.BlockKey, k types.BlockKey) {
	select {
	case c <- k:
	default:
	}
}

type headWrapper struct {
	c chan types.BlockKey
}

func (w *headWrapper) OnNewLongestChain(ctx context.Context, head *evmtypes.Head) {
	var bl int64
	if head != nil {
		bl = head.Number
	}

	send(w.c, chain.BlockKey(fmt.Sprintf("%d", bl)))
}
