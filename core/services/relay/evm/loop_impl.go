package evm

import (
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

// RelayAdapter extends loop.Relayer with a method for accessing the internal legacy chain type.
// Only avaialable in embedded mode, not LOOPP mode.
type RelayAdapter interface {
	loop.Relayer
	Chain() types.ChainService
}
type relayAdapter struct {
	loop.Relayer
	chain types.ChainService
}

var _ RelayAdapter = &relayAdapter{}

func NewLOOPAdapter(r loop.Relayer) *relayAdapter {
	return &relayAdapter{Relayer: r, chain: r}
}

func NewLegacyAdapter(r *Relayer) *relayAdapter {
	return &relayAdapter{
		Relayer: relay.NewServerAdapter(r),
		chain:   r.chain,
	}
}

func (la *relayAdapter) Chain() types.ChainService {
	return la.chain
}
