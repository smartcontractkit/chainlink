package starknet

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	relaychains "github.com/smartcontractkit/chainlink-relay/pkg/chains"
	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-relay/pkg/services"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"
	starkChain "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/chain"
	starkchain "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/chain"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/txm"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/starknet"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type ChainOpts struct {
	Logger logger.Logger
	// the implementation used here needs to be co-ordinated with the starknet transaction manager keystore adapter
	KeyStore loop.Keystore
}

func (o *ChainOpts) Name() string {
	return o.Logger.Name()
}

func (o *ChainOpts) Validate() (err error) {
	required := func(s string) error {
		return errors.Errorf("%s is required", s)
	}
	if o.Logger == nil {
		err = multierr.Append(err, required("Logger'"))
	}
	if o.KeyStore == nil {
		err = multierr.Append(err, required("KeyStore"))
	}
	return
}

var _ starkChain.Chain = (*chain)(nil)

type chain struct {
	utils.StartStopOnce
	id   string
	cfg  *StarknetConfig
	lggr logger.Logger
	txm  txm.StarkTXM
}

func NewChain(cfg *StarknetConfig, opts ChainOpts) (starkchain.Chain, error) {
	if !cfg.IsEnabled() {
		return nil, fmt.Errorf("cannot create new chain with ID %s: %w", *cfg.ChainID, chains.ErrChainDisabled)
	}
	c, err := newChain(*cfg.ChainID, cfg, opts.KeyStore, opts.Logger)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func newChain(id string, cfg *StarknetConfig, loopKs loop.Keystore, lggr logger.Logger) (*chain, error) {
	lggr = logger.With(lggr, "starknetChainID", id)
	ch := &chain{
		id:   id,
		cfg:  cfg,
		lggr: logger.Named(lggr, "Chain"),
	}

	getClient := func() (*starknet.Client, error) {
		return ch.getClient()
	}

	var err error
	ch.txm, err = txm.New(lggr, loopKs, cfg, getClient)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (c *chain) Name() string {
	return c.lggr.Name()
}

func (c *chain) Config() config.Config {
	return c.cfg
}

func (c *chain) TxManager() txm.TxManager {
	return c.txm
}

func (c *chain) Reader() (starknet.Reader, error) {
	return c.getClient()
}

func (c *chain) ChainID() string {
	return c.id
}

// getClient returns a client, randomly selecting one from available and valid nodes
func (c *chain) getClient() (*starknet.Client, error) {
	var node db.Node
	var client *starknet.Client
	nodes, err := c.cfg.ListNodes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get nodes")
	}
	if len(nodes) == 0 {
		return nil, errors.New("no nodes available")
	}
	// #nosec
	index := rand.Perm(len(nodes)) // list of node indexes to try
	timeout := c.cfg.RequestTimeout()
	for _, i := range index {
		node = nodes[i]
		// create client and check
		client, err = starknet.NewClient(node.ChainID, node.URL, c.lggr, &timeout)
		// if error, try another node
		if err != nil {
			c.lggr.Warnw("failed to create node", "name", node.Name, "starknet-url", node.URL, "err", err.Error())
			continue
		}
		// if all checks passed, mark found and break loop
		break
	}
	// if no valid node found, exit with error
	if client == nil {
		return nil, errors.New("no node valid nodes available")
	}
	c.lggr.Debugw("Created client", "name", node.Name, "starknet-url", node.URL)
	return client, nil
}

func (c *chain) Start(ctx context.Context) error {
	return c.StartOnce("Chain", func() error {
		return c.txm.Start(ctx)
	})
}

func (c *chain) Close() error {
	return c.StopOnce("Chain", func() error {
		return c.txm.Close()
	})
}

func (c *chain) Ready() error {
	return c.StartStopOnce.Ready()
}

func (c *chain) HealthReport() map[string]error {
	report := map[string]error{c.Name(): c.Healthy()}
	services.CopyHealth(report, c.txm.HealthReport())
	return report
}

func (c *chain) ID() string {
	return c.id
}

// ChainService interface
func (c *chain) GetChainStatus(ctx context.Context) (relaytypes.ChainStatus, error) {
	toml, err := c.cfg.TOMLString()
	if err != nil {
		return relaytypes.ChainStatus{}, err
	}
	return relaytypes.ChainStatus{
		ID:      c.id,
		Enabled: c.cfg.IsEnabled(),
		Config:  toml,
	}, nil
}

func (c *chain) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) (stats []relaytypes.NodeStatus, nextPageToken string, total int, err error) {
	return relaychains.ListNodeStatuses(int(pageSize), pageToken, c.listNodeStatuses)
}

func (c *chain) Transact(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	return chains.ErrLOOPPUnsupported
}

func (c *chain) SendTx(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	return c.Transact(ctx, from, to, amount, balanceCheck)
}

// TODO BCF-2602 statuses are static for non-evm chain and should be dynamic
func (c *chain) listNodeStatuses(start, end int) ([]relaytypes.NodeStatus, int, error) {
	stats := make([]relaytypes.NodeStatus, 0)
	total := len(c.cfg.Nodes)
	if start >= total {
		return stats, total, relaychains.ErrOutOfRange
	}
	if end <= 0 || end > total {
		end = total
	}
	nodes := c.cfg.Nodes[start:end]
	for _, node := range nodes {
		stat, err := nodeStatus(node, c.ChainID())
		if err != nil {
			return stats, total, err
		}
		stats = append(stats, stat)
	}
	return stats, total, nil
}
