package resolver

import (
	"context"
	"errors"

	"github.com/graph-gophers/graphql-go"
	"github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"

	v2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/web/loader"
)

// NodeResolver resolves the Node type.
type NodeResolver struct {
	node   v2.Node
	status types.NodeStatus
}

func NewNode(status types.NodeStatus) (nr *NodeResolver, warn error) {
	nr = &NodeResolver{status: status}
	warn = toml.Unmarshal([]byte(status.Config), &nr.node)
	return
}

func NewNodes(nodes []types.NodeStatus) (resolvers []*NodeResolver, warns error) {
	for _, n := range nodes {
		nr, warn := NewNode(n)
		if warn != nil {
			warns = errors.Join(warns, warn)
		}
		resolvers = append(resolvers, nr)
	}

	return
}

func orZero[P any](s *P) P {
	if s == nil {
		var zero P
		return zero
	}
	return *s
}

// ID resolves the node's unique identifier.
func (r *NodeResolver) ID() graphql.ID {
	return graphql.ID(r.Name())
}

// Name resolves the node's name field.
func (r *NodeResolver) Name() string {
	return orZero(r.node.Name)
}

// WSURL resolves the node's websocket url field.
func (r *NodeResolver) WSURL() string {
	if r.node.WSURL == nil {
		return ""
	}
	return r.node.WSURL.String()
}

// HTTPURL resolves the node's http url field.
func (r *NodeResolver) HTTPURL() string {
	if r.node.HTTPURL == nil {
		return ""
	}
	return r.node.HTTPURL.String()
}

// State resolves the node state
func (r *NodeResolver) State() string {
	return r.status.State
}

// SendOnly resolves the node's sendOnly bool
func (r *NodeResolver) SendOnly() bool {
	return orZero(r.node.SendOnly)
}

// Chain resolves the node's chain object field.
func (r *NodeResolver) Chain(ctx context.Context) (*ChainResolver, error) {
	chain, err := loader.GetChainByID(ctx, r.status.ChainID)
	if err != nil {
		return nil, err
	}

	return NewChain(*chain), nil
}

// -- Node Query --

type NodePayloadResolver struct {
	nr *NodeResolver
	NotFoundErrorUnionType
}

func NewNodePayloadResolver(node *types.NodeStatus, err error) (npr *NodePayloadResolver, warn error) {
	e := NotFoundErrorUnionType{err: err, message: "node not found", isExpectedErrorFn: nil}
	npr = &NodePayloadResolver{NotFoundErrorUnionType: e}
	if node != nil {
		npr.nr, warn = NewNode(*node)
	}
	return
}

// ToNode resolves the Node object to be returned if it is found
func (r *NodePayloadResolver) ToNode() (*NodeResolver, bool) {
	return r.nr, r.nr != nil
}

// -- Nodes Query --

type NodesPayloadResolver struct {
	nrs   []*NodeResolver
	total int32
}

func NewNodesPayload(nodes []types.NodeStatus, total int32) (npr *NodesPayloadResolver, warn error) {
	npr = &NodesPayloadResolver{total: total}
	npr.nrs, warn = NewNodes(nodes)
	return
}

func (r *NodesPayloadResolver) Results() []*NodeResolver {
	return r.nrs
}

func (r *NodesPayloadResolver) Metadata() *PaginationMetadataResolver {
	return NewPaginationMetadata(r.total)
}
