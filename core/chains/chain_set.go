package chains

import (
	"context"
	"math/big"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	// ErrChainIDEmpty is returned when chain is required but was empty.
	ErrChainIDEmpty = errors.New("chain id empty")
	ErrNotFound     = errors.New("not found")
)

// ChainStatuser is a generic interface for chain configuration.
type ChainStatuser interface {
	ChainStatus(ctx context.Context) (types.ChainStatus, error)
	// ChainStatuses(ctx context.Context, offset, limit int) ([]types.ChainStatus, int, error)
}

// NodesStatuser is an interface for node configuration and state.
type NodesStatuser interface {
	NodeStatuses(ctx context.Context, offset, limit int) (nodes []types.NodeStatus, count int, err error)
}

// ChainService is a live, runtime chain instance, with supporting services.
type ChainService interface {
	services.ServiceCtx
	ChainStatus(ctx context.Context) (types.ChainStatus, error)
	NodeStatuses(ctx context.Context, offset, limit int) (nodes []types.NodeStatus, count int, err error)
	SendTx(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error
}

// ChainSetOpts holds options for configuring a ChainSet via NewChainSet.
type ChainSetOpts[I ID, N Node] interface {
	Validate() error
	ConfigsAndLogger() (Statuser[N], logger.Logger)
}

type chainSet[N Node, S ChainService] struct {
	utils.StartStopOnce
	opts    ChainSetOpts[string, N]
	configs Statuser[N]
	lggr    logger.Logger
	chains  map[string]S
}

// NewChainSet returns a new immutable ChainSet for the given ChainSetOpts.
func NewChainSet[N Node, S ChainService](
	chains map[string]S,
	opts ChainSetOpts[string, N],
) (types.ChainSet[string, S], error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	cfgs, lggr := opts.ConfigsAndLogger()
	cs := chainSet[N, S]{
		opts:    opts,
		configs: cfgs,
		lggr:    logger.Named(lggr, "ChainSet"),
		chains:  chains,
	}

	return &cs, nil
}

func (c *chainSet[N, S]) Chain(ctx context.Context, id string) (s S, err error) {
	if err = c.StartStopOnce.Ready(); err != nil {
		return
	}
	ch, ok := c.chains[id]
	if !ok {
		err = ErrNotFound
		return
	}
	return ch, nil
}

func (c *chainSet[N, S]) ChainStatus(ctx context.Context, id string) (cfg types.ChainStatus, err error) {
	panic("unimplemented")
	return
}

func (c *chainSet[N, S]) ChainStatuses(ctx context.Context, offset, limit int) (x []types.ChainStatus, y int, z error) {
	panic("unimplemented")
	return
}

func (c *chainSet[N, S]) NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []types.NodeStatus, count int, err error) {
	panic("unimplemented")
	return
}

func (c *chainSet[N, S]) SendTx(ctx context.Context, chainID, from, to string, amount *big.Int, balanceCheck bool) error {
	panic("unimplemented")
	return nil
}

func (c *chainSet[N, S]) Start(ctx context.Context) error {
	panic("unimplemented")
}

func (c *chainSet[N, S]) Close() error {
	panic("unimplemented")
}

func (c *chainSet[N, S]) Ready() (err error) {
	panic("unimplemented")
}

func (c *chainSet[N, S]) Name() string {
	panic("unimplemented")
}

func (c *chainSet[N, S]) HealthReport() map[string]error {
	panic("unimplemented")
}
