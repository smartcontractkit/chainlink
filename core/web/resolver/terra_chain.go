package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"
)

// TerraChainResolver resolves the TerraChain type.
type TerraChainResolver struct {
	chain db.Chain
}

func NewTerraChain(chain db.Chain) *TerraChainResolver {
	return &TerraChainResolver{chain: chain}
}

func NewTerraChains(chains []db.Chain) []*TerraChainResolver {
	var resolvers []*TerraChainResolver
	for _, c := range chains {
		resolvers = append(resolvers, NewTerraChain(c))
	}

	return resolvers
}

// ID resolves the chain's unique identifier.
func (r *TerraChainResolver) ID() graphql.ID {
	return graphql.ID(r.chain.ID)
}

// Enabled resolves the chain's enabled field.
func (r *TerraChainResolver) Enabled() bool {
	return r.chain.Enabled
}

// TODO: implement ConfigResolver .. do we really need to remap every config variable?

// Config resolves the chain's configuration field
/* func (r *TerraChainResolver) Config() *TerraChainConfigResolver {
	return NewTerraChainConfig(r.chain.Cfg)
}
*/

// CreatedAt resolves the chain's created at field.
func (r *TerraChainResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.chain.CreatedAt}
}

// UpdatedAt resolves the chain's updated at field.
func (r *TerraChainResolver) UpdatedAt() graphql.Time {
	return graphql.Time{Time: r.chain.UpdatedAt}
}

// TODO: implement Nodes resolver.. how to make GetNodesByChainID load from terra orm?

/* func (r *TerraChainResolver) Nodes(ctx context.Context) ([]*NodeResolver, error) {
	nodes, err := loader.GetNodesByChainID(ctx, r.chain.ID.String())
	if err != nil {
		return nil, err
	}

	return NewNodes(nodes), nil
} */

type TerraChainPayloadResolver struct {
	chain db.Chain
	NotFoundErrorUnionType
}

func NewTerraChainPayload(chain db.Chain, err error) *TerraChainPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "chain not found", isExpectedErrorFn: nil}

	return &TerraChainPayloadResolver{chain: chain, NotFoundErrorUnionType: e}
}

func (r *TerraChainPayloadResolver) ToTerraChain() (*TerraChainResolver, bool) {
	if r.err != nil {
		return nil, false
	}

	return NewTerraChain(r.chain), true
}

type TerraChainsPayloadResolver struct {
	chains []db.Chain
	total  int32
}

func NewTerraChainsPayload(chains []db.Chain, total int32) *TerraChainsPayloadResolver {
	return &TerraChainsPayloadResolver{chains: chains, total: total}
}

func (r *TerraChainsPayloadResolver) Results() []*TerraChainResolver {
	return NewTerraChains(r.chains)
}

func (r *TerraChainsPayloadResolver) Metadata() *PaginationMetadataResolver {
	return NewPaginationMetadata(r.total)
}
