package chains

import (
	"context"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

var (
	// ErrChainIDEmpty is returned when chain is required but was empty.
	ErrChainIDEmpty = errors.New("chain id empty")
	ErrNotFound     = errors.New("not found")
)

// ChainStatuser is a generic interface for chain configuration.
type ChainStatuser interface {
	// must return [ErrNotFound] if the id is not found
	ChainStatus(ctx context.Context, id string) (types.ChainStatus, error)
	ChainStatuses(ctx context.Context, offset, limit int) ([]types.ChainStatus, int, error)
}

// NodesStatuser is an interface for node configuration and state.
// TODO BCF2440, BCF-2511 may need Node(ctx,name) to get a node status by name
type NodesStatuser interface {
	NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []types.NodeStatus, count int, err error)
}

// ChainSetOpts holds options for configuring a ChainSet via NewChainSet.
type ChainSetOpts[I ID, N Node] interface {
	Validate() error
	ConfigsAndLogger() (Configs[I, N], logger.Logger)
}
