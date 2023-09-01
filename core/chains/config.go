package chains

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

var (
	// ErrChainIDEmpty is returned when chain is required but was empty.
	ErrChainIDEmpty = errors.New("chain id empty")
	ErrNotFound     = errors.New("not found")
)

type ChainConfigs interface {
	Chains(offset, limit int, ids ...string) ([]types.ChainStatus, int, error)
}

type NodeConfigs[I ID, N Node] interface {
	Node(name string) (N, error)
	Nodes(chainID I) (nodes []N, err error)

	NodeStatus(name string) (types.NodeStatus, error)
}

// Configs holds chain and node configurations.
// TODO: BCF-2605 audit the usage of this interface and potentially remove it
type Configs[I ID, N Node] interface {
	ChainConfigs
	NodeConfigs[I, N]
}

// ChainStatuser is a generic interface for chain configuration.
type ChainStatuser interface {
	// must return [ErrNotFound] if the id is not found
	ChainStatus(ctx context.Context, id string) (types.ChainStatus, error)
	ChainStatuses(ctx context.Context, offset, limit int) ([]types.ChainStatus, int, error)
}

// NodesStatuser is an interface for node configuration and state.
// TODO BCF-2440, BCF-2511 may need Node(ctx,name) to get a node status by name
type NodesStatuser interface {
	NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []types.NodeStatus, count int, err error)
}

// ChainOpts holds options for configuring a Chain
type ChainOpts[I ID, N Node] interface {
	Validate() error
	ConfigsAndLogger() (Configs[I, N], logger.Logger)
}
