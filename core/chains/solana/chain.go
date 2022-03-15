package solana

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	solanaclient "github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/solana/monitor"
	"github.com/smartcontractkit/chainlink/core/chains/solana/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// DefaultRequestTimeout is the default Solana client timeout.
const DefaultRequestTimeout = 30 * time.Second

//go:generate mockery --name TxManager --srcpkg github.com/smartcontractkit/chainlink-solana/pkg/solana --output ./mocks/ --case=underscore
//go:generate mockery --name Reader --srcpkg github.com/smartcontractkit/chainlink-solana/pkg/solana/client --output ./mocks/ --case=underscore
//go:generate mockery --name Chain --srcpkg github.com/smartcontractkit/chainlink-solana/pkg/solana --output ./mocks/ --case=underscore
var _ solana.Chain = (*chain)(nil)

type chain struct {
	utils.StartStopOnce
	id  string
	cfg solana.Config
	// txm            *terratxm.Txm
	balanceMonitor services.Service
	orm            types.ORM
	lggr           logger.Logger
}

// NewChain returns a new chain backed by node.
func NewChain(db *sqlx.DB, ks keystore.Solana, logCfg pg.LogConfig, eb pg.EventBroadcaster, dbchain db.Chain, orm types.ORM, lggr logger.Logger) (*chain, error) {
	cfg := solana.NewConfig(dbchain.Cfg, lggr)
	lggr = lggr.With("solanaChainID", dbchain.ID)
	var ch = chain{
		id:   dbchain.ID,
		cfg:  cfg,
		orm:  orm,
		lggr: lggr.Named("Chain"),
	}
	tc := func() (solanaclient.ReaderWriter, error) {
		return ch.getClient("")
	}
	// ch.txm = terratxm.NewTxm(db, tc, *gpe, ch.id, cfg, ks, lggr, logCfg, eb)
	ch.balanceMonitor = monitor.NewBalanceMonitor(ch.id, cfg, lggr, ks, ch.Reader)
	return &ch, nil
}

func (c *chain) ID() string {
	return c.id
}

func (c *chain) Config() solana.Config {
	return c.cfg
}

func (c *chain) UpdateConfig(cfg db.ChainCfg) {
	c.cfg.Update(cfg)
}

// func (c *chain) TxManager() solana.TxManager {
// 	return c.txm
// }

func (c *chain) Reader(name string) (solanaclient.Reader, error) {
	return c.getClient(name)
}

// getClient returns a client, optionally requiring a specific node by name.
func (c *chain) getClient(name string) (solanaclient.ReaderWriter, error) {
	//TODO cache clients?
	var node db.Node
	if name == "" { // Any node
		nodes, cnt, err := c.orm.NodesForChain(c.id, 0, math.MaxInt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get nodes")
		}
		if cnt == 0 {
			return nil, errors.New("no nodes available")
		}
		// #nosec
		node = nodes[rand.Intn(len(nodes))]
	} else { // Named node
		var err error
		node, err = c.orm.NodeNamed(name)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get node named %s", name)
		}
		if node.SolanaChainID != c.id {
			return nil, fmt.Errorf("failed to create client for chain %s with node %s: wrong chain id %s", c.id, name, node.SolanaChainID)
		}
	}
	client, err := solanaclient.NewClient(node.SolanaURL, c.cfg, DefaultRequestTimeout, c.lggr.Named("Client-"+node.Name))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create client")
	}
	c.lggr.Debugw("Created client", "name", node.Name, "solana-url", node.SolanaURL)
	return client, nil
}

func (c *chain) Start() error {
	return c.StartOnce("Chain", func() error {
		c.lggr.Debug("Starting")
		c.lggr.Debug("Starting txm")
		c.lggr.Debug("Starting balance monitor")
		return multierr.Combine(
			c.txm.Start(),
			c.balanceMonitor.Start())
	})
}

func (c *chain) Close() error {
	return c.StopOnce("Chain", func() error {
		c.lggr.Debug("Stopping")
		c.lggr.Debug("Stopping txm")
		c.lggr.Debug("Stopping balance monitor")
		return multierr.Combine(c.txm.Close(),
			c.balanceMonitor.Close())
	})
}

func (c *chain) Ready() error {
	return multierr.Combine(
		c.StartStopOnce.Ready(),
		c.txm.Ready(),
	)
}

func (c *chain) Healthy() error {
	return multierr.Combine(
		c.StartStopOnce.Healthy(),
		c.txm.Healthy(),
	)
}
