package resolver

import (
	"context"

	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/web/loader"
)

// NodeResolver resolves the Node type.
type NodeResolver struct {
	node types.Node
}

func NewNode(node types.Node) *NodeResolver {
	return &NodeResolver{node: node}
}

func NewNodes(nodes []types.Node) []*NodeResolver {
	var resolvers []*NodeResolver
	for _, n := range nodes {
		resolvers = append(resolvers, NewNode(n))
	}

	return resolvers
}

// ID resolves the node's unique identifier.
func (r *NodeResolver) ID() graphql.ID {
	return graphql.ID(r.node.Name)
}

// Name resolves the node's name field.
func (r *NodeResolver) Name() string {
	return r.node.Name
}

// WSURL resolves the node's websocket url field.
func (r *NodeResolver) WSURL() string {
	return r.node.WSURL.String
}

// HTTPURL resolves the node's http url field.
func (r *NodeResolver) HTTPURL() string {
	return r.node.HTTPURL.String
}

// State resolves the node state
func (r *NodeResolver) State() string {
	return r.node.State
}

// SendOnly resolves the node's sendOnly bool
func (r *NodeResolver) SendOnly() bool {
	return r.node.SendOnly
}

// Chain resolves the node's chain object field.
func (r *NodeResolver) Chain(ctx context.Context) (*ChainResolver, error) {
	chain, err := loader.GetChainByID(ctx, r.node.EVMChainID.String())
	if err != nil {
		return nil, err
	}

	return NewChain(*chain), nil
}

// CreatedAt resolves the node's created at field.
func (r *NodeResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.node.CreatedAt}
}

// UpdatedAt resolves the node's updated at field.
func (r *NodeResolver) UpdatedAt() graphql.Time {
	return graphql.Time{Time: r.node.UpdatedAt}
}

// -- Node Query --

type NodePayloadResolver struct {
	node *types.Node
	NotFoundErrorUnionType
}

func NewNodePayloadResolver(node *types.Node, err error) *NodePayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "node not found", isExpectedErrorFn: nil}

	return &NodePayloadResolver{node: node, NotFoundErrorUnionType: e}
}

// ToNode resolves the Node object to be returned if it is found
func (r *NodePayloadResolver) ToNode() (*NodeResolver, bool) {
	if r.node != nil {
		return NewNode(*r.node), true
	}

	return nil, false
}

// -- Nodes Query --

type NodesPayloadResolver struct {
	nodes []types.Node
	total int32
}

func NewNodesPayload(nodes []types.Node, total int32) *NodesPayloadResolver {
	return &NodesPayloadResolver{nodes: nodes, total: total}
}

func (r *NodesPayloadResolver) Results() []*NodeResolver {
	return NewNodes(r.nodes)
}

func (r *NodesPayloadResolver) Metadata() *PaginationMetadataResolver {
	return NewPaginationMetadata(r.total)
}

// -- CreateNode Mutation --

type CreateNodePayloadResolver struct {
	node *types.Node
}

func NewCreateNodePayloadResolver(node *types.Node) *CreateNodePayloadResolver {
	return &CreateNodePayloadResolver{node: node}
}

func (r *CreateNodePayloadResolver) ToCreateNodeSuccess() (*CreateNodeSuccessResolve, bool) {
	if r.node != nil {
		return NewCreateNodeSuccessResolve(*r.node), true
	}

	return nil, false
}

type CreateNodeSuccessResolve struct {
	node types.Node
}

func NewCreateNodeSuccessResolve(node types.Node) *CreateNodeSuccessResolve {
	return &CreateNodeSuccessResolve{node}
}

func (r *CreateNodeSuccessResolve) Node() *NodeResolver {
	return NewNode(r.node)
}

// -- DeleteNode Mutation --

type DeleteNodePayloadResolver struct {
	node *types.Node
	NotFoundErrorUnionType
}

func NewDeleteNodePayloadResolver(node *types.Node, err error) *DeleteNodePayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "node not found", isExpectedErrorFn: nil}

	return &DeleteNodePayloadResolver{node: node, NotFoundErrorUnionType: e}
}

func (r *DeleteNodePayloadResolver) ToDeleteNodeSuccess() (*DeleteNodeSuccessResolver, bool) {
	if r.node != nil {
		return NewDeleteNodeSuccessResolver(r.node), true
	}

	return nil, false
}

type DeleteNodeSuccessResolver struct {
	node *types.Node
}

func NewDeleteNodeSuccessResolver(node *types.Node) *DeleteNodeSuccessResolver {
	return &DeleteNodeSuccessResolver{node: node}
}

func (r *DeleteNodeSuccessResolver) Node() *NodeResolver {
	return NewNode(*r.node)
}
