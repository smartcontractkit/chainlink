package evm

import (
	"github.com/smartcontractkit/chainlink-common/pkg/loop"

	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

type LOOPRelayAdapter interface {
	loop.Relayer
	Chain() legacyevm.Chain
}
type loopRelayAdapter struct {
	loop.Relayer
	chain legacyevm.Chain
}

var _ LOOPRelayAdapter = &loopRelayAdapter{}

func NewLOOPRelayAdapter(r *Relayer) *loopRelayAdapter {
	return &loopRelayAdapter{
		Relayer: relay.NewServerAdapter(r),
		chain:   r.chain,
	}
}

func (la *loopRelayAdapter) Chain() legacyevm.Chain {
	return la.chain
}
