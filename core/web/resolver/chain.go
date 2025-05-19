package resolver

import (
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

// ChainResolver resolves the Chain type.
type ChainResolver struct {
	chain chainlink.NetworkChainStatus
}

func NewChain(chain chainlink.NetworkChainStatus) *ChainResolver {
	return &ChainResolver{chain: chain}
}

func NewChains(chains []chainlink.NetworkChainStatus) []*ChainResolver {
	var resolvers []*ChainResolver
	for _, c := range chains {
		resolvers = append(resolvers, NewChain(c))
	}

	return resolvers
}

// ID resolves the chain's unique identifier.
func (r *ChainResolver) ID() graphql.ID {
	return graphql.ID(r.chain.ID)
}

// Enabled resolves the chain's enabled field.
func (r *ChainResolver) Enabled() bool {
	return r.chain.Enabled
}

// Config resolves the chain's configuration field
func (r *ChainResolver) Config() string {
	return r.chain.Config
}

// Network resolves the chain's network field
func (r *ChainResolver) Network() string {
	return r.chain.Network
}

type ChainPayloadResolver struct {
	chain chainlink.NetworkChainStatus
	NotFoundErrorUnionType
}

func NewChainPayload(chain chainlink.NetworkChainStatus, err error) *ChainPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "chain not found", isExpectedErrorFn: nil}

	return &ChainPayloadResolver{chain: chain, NotFoundErrorUnionType: e}
}

func (r *ChainPayloadResolver) ToChain() (*ChainResolver, bool) {
	if r.err != nil {
		return nil, false
	}

	return NewChain(r.chain), true
}

type ChainsPayloadResolver struct {
	chains []chainlink.NetworkChainStatus
	total  int32
}

func NewChainsPayload(chains []chainlink.NetworkChainStatus, total int32) *ChainsPayloadResolver {
	return &ChainsPayloadResolver{chains: chains, total: total}
}

func (r *ChainsPayloadResolver) Results() []*ChainResolver {
	return NewChains(r.chains)
}

func (r *ChainsPayloadResolver) Metadata() *PaginationMetadataResolver {
	return NewPaginationMetadata(r.total)
}
