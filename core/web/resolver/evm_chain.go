package resolver

import (
	"context"

	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/web/loader"
)

// ChainResolver resolves the Chain type.
type ChainResolver struct {
	chain types.DBChain
}

func NewChain(chain types.DBChain) *ChainResolver {
	return &ChainResolver{chain: chain}
}

func NewChains(chains []types.DBChain) []*ChainResolver {
	var resolvers []*ChainResolver
	for _, c := range chains {
		resolvers = append(resolvers, NewChain(c))
	}

	return resolvers
}

// ID resolves the chain's unique identifier.
func (r *ChainResolver) ID() graphql.ID {
	return graphql.ID(r.chain.ID.String())
}

// Enabled resolves the chain's enabled field.
func (r *ChainResolver) Enabled() bool {
	return r.chain.Enabled
}

// Config resolves the chain's configuration field
func (r *ChainResolver) Config() *ChainConfigResolver {
	return NewChainConfig(*r.chain.Cfg)
}

// CreatedAt resolves the chain's created at field.
func (r *ChainResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.chain.CreatedAt}
}

// UpdatedAt resolves the chain's updated at field.
func (r *ChainResolver) UpdatedAt() graphql.Time {
	return graphql.Time{Time: r.chain.UpdatedAt}
}

func (r *ChainResolver) Nodes(ctx context.Context) ([]*NodeResolver, error) {
	nodes, err := loader.GetNodesByChainID(ctx, r.chain.ID.String())
	if err != nil {
		return nil, err
	}

	return NewNodes(nodes), nil
}

type ChainPayloadResolver struct {
	chain types.DBChain
	NotFoundErrorUnionType
}

func NewChainPayload(chain types.DBChain, err error) *ChainPayloadResolver {
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
	chains []types.DBChain
	total  int32
}

func NewChainsPayload(chains []types.DBChain, total int32) *ChainsPayloadResolver {
	return &ChainsPayloadResolver{chains: chains, total: total}
}

func (r *ChainsPayloadResolver) Results() []*ChainResolver {
	return NewChains(r.chains)
}

func (r *ChainsPayloadResolver) Metadata() *PaginationMetadataResolver {
	return NewPaginationMetadata(r.total)
}

// -- CreateChain Mutation --

type CreateChainPayloadResolver struct {
	chain     *types.DBChain
	inputErrs map[string]string
}

func NewCreateChainPayload(chain *types.DBChain, inputErrs map[string]string) *CreateChainPayloadResolver {
	return &CreateChainPayloadResolver{chain: chain, inputErrs: inputErrs}
}

func (r *CreateChainPayloadResolver) ToCreateChainSuccess() (*CreateChainSuccessResolver, bool) {
	if r.chain == nil {
		return nil, false
	}

	return NewCreateChainSuccess(r.chain), true
}

func (r *CreateChainPayloadResolver) ToInputErrors() (*InputErrorsResolver, bool) {
	if r.inputErrs != nil {
		var errs []*InputErrorResolver

		for path, message := range r.inputErrs {
			errs = append(errs, NewInputError(path, message))
		}

		return NewInputErrors(errs), true
	}

	return nil, false
}

type CreateChainSuccessResolver struct {
	chain *types.DBChain
}

func NewCreateChainSuccess(chain *types.DBChain) *CreateChainSuccessResolver {
	return &CreateChainSuccessResolver{chain: chain}
}

func (r *CreateChainSuccessResolver) Chain() *ChainResolver {
	return NewChain(*r.chain)
}

type UpdateChainPayloadResolver struct {
	chain     *types.DBChain
	inputErrs map[string]string
	NotFoundErrorUnionType
}

func NewUpdateChainPayload(chain *types.DBChain, inputErrs map[string]string, err error) *UpdateChainPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "chain not found", isExpectedErrorFn: nil}

	return &UpdateChainPayloadResolver{chain: chain, inputErrs: inputErrs, NotFoundErrorUnionType: e}
}

func (r *UpdateChainPayloadResolver) ToUpdateChainSuccess() (*UpdateChainSuccessResolver, bool) {
	if r.chain == nil {
		return nil, false
	}

	return NewUpdateChainSuccess(*r.chain), true
}

func (r *UpdateChainPayloadResolver) ToInputErrors() (*InputErrorsResolver, bool) {
	if r.inputErrs != nil {
		var errs []*InputErrorResolver

		for path, message := range r.inputErrs {
			errs = append(errs, NewInputError(path, message))
		}

		return NewInputErrors(errs), true
	}

	return nil, false
}

type UpdateChainSuccessResolver struct {
	chain types.DBChain
}

func NewUpdateChainSuccess(chain types.DBChain) *UpdateChainSuccessResolver {
	return &UpdateChainSuccessResolver{chain: chain}
}

func (r *UpdateChainSuccessResolver) Chain() *ChainResolver {
	return NewChain(r.chain)
}

type DeleteChainPayloadResolver struct {
	chain *types.DBChain
	NotFoundErrorUnionType
}

func NewDeleteChainPayload(chain *types.DBChain, err error) *DeleteChainPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "chain not found", isExpectedErrorFn: nil}

	return &DeleteChainPayloadResolver{chain: chain, NotFoundErrorUnionType: e}
}

func (r *DeleteChainPayloadResolver) ToDeleteChainSuccess() (*DeleteChainSuccessResolver, bool) {
	if r.chain == nil {
		return nil, false
	}

	return NewDeleteChainSuccess(*r.chain), true
}

type DeleteChainSuccessResolver struct {
	chain types.DBChain
}

func NewDeleteChainSuccess(chain types.DBChain) *DeleteChainSuccessResolver {
	return &DeleteChainSuccessResolver{chain: chain}
}

func (r *DeleteChainSuccessResolver) Chain() *ChainResolver {
	return NewChain(r.chain)
}
