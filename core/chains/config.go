package chains

import (
	"errors"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

var (
	// ErrChainIDEmpty is returned when chain is required but was empty.
	ErrChainIDEmpty = errors.New("chain id empty")
	ErrNotFound     = errors.New("not found")
)

type ChainConfigs interface {
	Chains(offset, limit int, ids ...relay.ChainID) ([]types.ChainStatus, int, error)
}

type NodeConfigs[N Node] interface {
	Node(name string) (N, error)
	Nodes(chainID relay.ChainID) (nodes []N, err error)

	NodeStatus(name string) (types.NodeStatus, error)
}

// Configs holds chain and node configurations.
// TODO: BCF-2605 audit the usage of this interface and potentially remove it
type Configs[N Node] interface {
	ChainConfigs
	NodeConfigs[N]
}

// ChainOpts holds options for configuring a Chain
type ChainOpts[N Node] interface {
	Validate() error
	ConfigsAndLogger() (Configs[N], logger.Logger)
}
