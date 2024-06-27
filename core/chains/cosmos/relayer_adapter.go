package cosmos

import (
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/adapters"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"

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

func NewLoopRelayerChain(r *pkgcosmos.Relayer, s adapters.Chain) *LoopRelayerChain {
	ra := relay.NewServerAdapter(r, s)
	return &LoopRelayerChain{
		Relayer: ra,
		chain:   s,
	}
}
func (r *LoopRelayerChain) Chain() adapters.Chain {
	return r.chain
}

var _ LoopRelayerChainer = &LoopRelayerChain{}
