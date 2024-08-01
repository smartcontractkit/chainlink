package starknet

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"math/rand"

	"github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink-common/pkg/chains"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/txm"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/starknet"
)

type Chain interface {
	types.ChainService

	ID() string
	Config() config.Config

	TxManager() txm.TxManager
	Reader() (starknet.Reader, error)
}

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
		return fmt.Errorf("%s is required", s)
	}
	if o.Logger == nil {
		err = errors.Join(err, required("Logger"))
	}
	if o.KeyStore == nil {
		err = errors.Join(err, required("KeyStore"))
	}
	return
}

var _ Chain = (*chain)(nil)

type chain struct {
	utils.StartStopOnce
	id   string
	cfg  *config.TOMLConfig
	lggr logger.Logger
	txm  txm.StarkTXM
}

func NewChain(cfg *config.TOMLConfig, opts ChainOpts) (Chain, error) {
	if !cfg.IsEnabled() {
		return nil, fmt.Errorf("cannot create new chain with ID %s: chain is disabled", *cfg.ChainID)
	}
	c, err := newChain(*cfg.ChainID, cfg, opts.KeyStore, opts.Logger)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func newChain(id string, cfg *config.TOMLConfig, loopKs loop.Keystore, lggr logger.Logger) (*chain, error) {
	lggr = logger.With(lggr, "starknetChainID", id)
	ch := &chain{
		id:   id,
		cfg:  cfg,
		lggr: logger.Named(lggr, "Chain"),
	}

	getClient := func() (*starknet.Client, error) {
		return ch.getClient()
	}

	getFeederClient := func() (*starknet.FeederClient, error) {
		return ch.getFeederClient(), nil
	}

	var err error
	ch.txm, err = txm.New(lggr, loopKs, cfg, getClient, getFeederClient)
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

func (c *chain) getFeederClient() *starknet.FeederClient {
	return starknet.NewFeederClient(c.cfg.FeederURL.String())
}

// getClient returns a client, randomly selecting one from available and valid nodes
func (c *chain) getClient() (*starknet.Client, error) {
	var node db.Node
	var client *starknet.Client
	nodes, err := c.cfg.ListNodes()
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
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
		client, err = starknet.NewClient(node.ChainID, node.URL, node.APIKey, c.lggr, &timeout)
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
func (c *chain) GetChainStatus(ctx context.Context) (types.ChainStatus, error) {
	toml, err := c.cfg.TOMLString()
	if err != nil {
		return types.ChainStatus{}, err
	}
	return types.ChainStatus{
		ID:      c.id,
		Enabled: c.cfg.IsEnabled(),
		Config:  toml,
	}, nil
}

func (c *chain) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) (stats []types.NodeStatus, nextPageToken string, total int, err error) {
	return chains.ListNodeStatuses(int(pageSize), pageToken, c.listNodeStatuses)
}

func (c *chain) Transact(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	return errors.ErrUnsupported
}

func (c *chain) SendTx(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	return c.Transact(ctx, from, to, amount, balanceCheck)
}

// TODO BCF-2602 statuses are static for non-evm chain and should be dynamic
func (c *chain) listNodeStatuses(start, end int) ([]types.NodeStatus, int, error) {
	stats := make([]types.NodeStatus, 0)
	total := len(c.cfg.Nodes)
	if start >= total {
		return stats, total, chains.ErrOutOfRange
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

func nodeStatus(n *config.Node, id string) (types.NodeStatus, error) {
	var s types.NodeStatus
	s.ChainID = id
	s.Name = *n.Name
	b, err := toml.Marshal(n)
	if err != nil {
		return types.NodeStatus{}, err
	}
	s.Config = string(b)
	return s, nil
}
