package mocks

import (
	"context"
	"slices"

	services2 "github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// FakeRelayerChainInteroperators is a fake chainlink.RelayerChainInteroperators.
// This exists because mockery generation doesn't understand how to produce an alias instead of the underlying type (which is not exported in this case).
type FakeRelayerChainInteroperators struct {
	Relayers  map[types.RelayID]loop.Relayer
	EVMChains legacyevm.LegacyChainContainer
	Nodes     []types.NodeStatus
	NodesErr  error
}

func (f *FakeRelayerChainInteroperators) LegacyEVMChains() legacyevm.LegacyChainContainer {
	return f.EVMChains
}

func (f *FakeRelayerChainInteroperators) NodeStatuses(ctx context.Context, offset, limit int, relayIDs ...types.RelayID) (nodes []types.NodeStatus, count int, err error) {
	return slices.Clone(f.Nodes), len(f.Nodes), f.NodesErr
}

func (f *FakeRelayerChainInteroperators) Services() []services2.ServiceCtx {
	panic("unimplemented")
}

func (f *FakeRelayerChainInteroperators) List(filter chainlink.FilterFn) chainlink.RelayerChainInteroperators {
	panic("unimplemented")
}

func (f *FakeRelayerChainInteroperators) Get(id types.RelayID) (loop.Relayer, error) {
	r, ok := f.Relayers[id]
	if !ok {
		return nil, chainlink.ErrNoSuchRelayer
	}
	return r, nil
}

func (f *FakeRelayerChainInteroperators) GetIDToRelayerMap() map[types.RelayID]loop.Relayer {
	return f.Relayers
}

func (f *FakeRelayerChainInteroperators) Slice() []loop.Relayer {
	var relayers []loop.Relayer

	for _, value := range f.Relayers {
		relayers = append(relayers, value)
	}

	return relayers
}

func (f *FakeRelayerChainInteroperators) ChainStatus(ctx context.Context, id types.RelayID) (types.ChainStatus, error) {
	panic("unimplemented")
}

func (f *FakeRelayerChainInteroperators) ChainStatuses(ctx context.Context, offset, limit int) ([]chainlink.NetworkChainStatus, int, error) {
	panic("unimplemented")
}
