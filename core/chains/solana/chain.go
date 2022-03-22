package solana

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	solanaclient "github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/solana/monitor"
	"github.com/smartcontractkit/chainlink/core/chains/solana/soltxm"
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
	id             string
	cfg            config.Config
	txm            *soltxm.Txm
	balanceMonitor services.ServiceCtx
	orm            types.ORM
	lggr           logger.Logger

	// tracking node chain id for verification
	clientChainID map[string]string
	chainIDLock   *sync.RWMutex
}

// NewChain returns a new chain backed by node.
func NewChain(db *sqlx.DB, ks keystore.Solana, logCfg pg.LogConfig, eb pg.EventBroadcaster, dbchain db.Chain, orm types.ORM, lggr logger.Logger) (*chain, error) {
	cfg := config.NewConfig(dbchain.Cfg, lggr)
	lggr = lggr.With("solanaChainID", dbchain.ID)
	var ch = chain{
		id:            dbchain.ID,
		cfg:           cfg,
		orm:           orm,
		lggr:          lggr.Named("Chain"),
		clientChainID: map[string]string{},
		chainIDLock:   &sync.RWMutex{},
	}
	tc := func() (solanaclient.ReaderWriter, error) {
		return ch.getClient("")
	}
	ch.txm = soltxm.NewTxm(tc, cfg, lggr)
	ch.balanceMonitor = monitor.NewBalanceMonitor(ch.id, cfg, lggr, ks, ch.Reader)
	return &ch, nil
}

func (c *chain) ID() string {
	return c.id
}

func (c *chain) Config() config.Config {
	return c.cfg
}

func (c *chain) UpdateConfig(cfg db.ChainCfg) {
	c.cfg.Update(cfg)
}

func (c *chain) TxManager() solana.TxManager {
	return c.txm
}

func (c *chain) Reader(name string) (solanaclient.Reader, error) {
	return c.getClient(name)
}

// getClient returns a client, optionally requiring a specific node by name.
func (c *chain) getClient(name string) (solanaclient.ReaderWriter, error) {
	var node db.Node
	var client solanaclient.ReaderWriter
	if name == "" { // Any node
		nodes, cnt, err := c.orm.NodesForChain(c.id, 0, math.MaxInt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get nodes")
		}
		if cnt == 0 {
			return nil, errors.New("no nodes available")
		}
		// #nosec
		index := rand.Perm(len(nodes)) // list of node indexes to try
		found := false
		for _, i := range index {
			node = nodes[i]
			// create client and check
			client, err = c.getOrCreate(node)
			// if error, try another node
			if err != nil {
				c.lggr.Warnw("failed to create node", "name", node.Name, "solana-url", node.SolanaURL, "error", err.Error())
				continue
			}
			// if all checks passed, mark found and break loop
			found = true
			break
		}
		// if no valid node found, exit with error
		if !found {
			return nil, errors.New("no node valid nodes available")
		}
	} else { // Named node
		var err error
		node, err = c.orm.NodeNamed(name)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get node named %s", name)
		}
		if node.SolanaChainID != c.id {
			return nil, fmt.Errorf("failed to create client for chain %s with node %s: wrong chain id %s", c.id, name, node.SolanaChainID)
		}

		// create client and check
		client, err = c.getOrCreate(node)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create node %s", name)
		}
	}
	c.lggr.Debugw("Created client", "name", node.Name, "solana-url", node.SolanaURL)
	return client, nil
}

func (c *chain) getOrCreate(node db.Node) (solanaclient.ReaderWriter, error) {
	// create client
	client, err := solanaclient.NewClient(node.SolanaURL, c.cfg, DefaultRequestTimeout, c.lggr.Named("Client-"+node.Name))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create client")
	}

	// if endpoint has not been checked, fetch chainID
	url := node.SolanaURL
	c.chainIDLock.RLock()
	_, exists := c.clientChainID[url]
	c.chainIDLock.RUnlock()
	if !exists {
		id, err := client.ChainID()
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch ChainID in checkClient")
		}
		c.chainIDLock.Lock()
		c.clientChainID[url] = id
		c.chainIDLock.Unlock()
	}

	// check chainID matches expected chainID
	expectedID := strings.ToLower(c.id)
	c.chainIDLock.RLock()
	defer c.chainIDLock.RUnlock()
	if c.clientChainID[url] != expectedID {
		return nil, fmt.Errorf("client returned mismatched chain id (expected: %s, got: %s): %s", expectedID, c.clientChainID[url], url)
	}
	return client, nil
}

func (c *chain) Start(ctx context.Context) error {
	return c.StartOnce("Chain", func() error {
		c.lggr.Debug("Starting")
		c.lggr.Debug("Starting txm")
		c.lggr.Debug("Starting balance monitor")
		return multierr.Combine(
			c.txm.Start(ctx),
			c.balanceMonitor.Start(ctx))
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
