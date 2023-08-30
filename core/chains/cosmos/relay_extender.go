package cosmos

import (
	"context"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/adapters"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	pkgcosmos "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos"
	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

// LegacyChainContainer is container interface for Cosmos chains
type LegacyChainContainer interface {
	Get(id string) (adapters.Chain, error)
	Len() int
	List(ids ...string) ([]adapters.Chain, error)
	Slice() []adapters.Chain
}

type LegacyChains = chains.ChainsKV[adapters.Chain]

var _ LegacyChainContainer = &LegacyChains{}

func NewLegacyChains(m map[string]adapters.Chain) *LegacyChains {
	return chains.NewChainsKV[adapters.Chain](m)
}

type LoopRelayerChainer interface {
	loop.Relayer
	Chain() adapters.Chain
}

type LoopRelayerChain struct {
	loop.Relayer
	chain adapters.Chain
}

func NewLoopRelayerChain(r *pkgcosmos.Relayer, s *RelayExtender) *LoopRelayerChain {

	ra := relay.NewRelayerAdapter(r, s)
	return &LoopRelayerChain{
		Relayer: ra,
		chain:   s,
	}
}
func (r *LoopRelayerChain) Chain() adapters.Chain {
	return r.chain
}

var _ LoopRelayerChainer = &LoopRelayerChain{}

// TODO remove these wrappers after BCF-2441
type RelayExtender struct {
	adapters.Chain
	chainImpl *chain
}

var _ relay.RelayerExt = &RelayExtender{}

func NewRelayExtender(cfg *CosmosConfig, opts ChainOpts) (*RelayExtender, error) {
	c, err := NewChain(cfg, opts)
	if err != nil {
		return nil, err
	}
	chainImpl, ok := (c).(*chain)
	if !ok {
		return nil, fmt.Errorf("internal error: cosmos relay extender got wrong type %t", c)
	}
	return &RelayExtender{Chain: chainImpl, chainImpl: chainImpl}, nil
}
func (r *RelayExtender) GetChainStatus(ctx context.Context) (relaytypes.ChainStatus, error) {
	return r.chainImpl.GetChainStatus(ctx)
}
func (r *RelayExtender) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) (stats []relaytypes.NodeStatus, nextPageToken string, total int, err error) {
	return r.chainImpl.ListNodeStatuses(ctx, pageSize, pageToken)
}
func (r *RelayExtender) Transact(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	return chains.ErrLOOPPUnsupported
}
