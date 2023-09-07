package chains

import (
	"errors"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

var (
	// ErrChainIDEmpty is returned when chain is required but was empty.
	ErrChainIDEmpty = errors.New("chain id empty")
	ErrNotFound     = errors.New("not found")
)

type ChainConfig[N Node] interface {
	GetChainStatus() (types.ChainStatus, error)
	GetNodeStatus(name string) (types.NodeStatus, error)
	NodeConfig[N]
}

type NodeConfig[N Node] interface {
	ListNodes() (nodes []N, err error)
	Node(name string) (N, error)
}

// ChainOpts holds options for configuring a Chain
type ChainOpts[N Node] interface {
	Validate() error
	Logger() logger.Logger
	ChainConfig() ChainConfig[N]
}
