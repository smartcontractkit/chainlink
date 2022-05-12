package solana

import (
	"context"
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
	orm            ORM
	lggr           logger.Logger

	// tracking node chain id for verification
	clientCache map[string]cachedClient // map URL -> {client, chainId} [mainnet/testnet/devnet/localnet]
	clientLock  sync.RWMutex
}

type cachedClient struct {
	id string
	rw solanaclient.ReaderWriter
}

// NewChain returns a new chain backed by node.
func NewChain(db *sqlx.DB, ks keystore.Solana, logCfg pg.LogConfig, eb pg.EventBroadcaster, dbchain DBChain, orm ORM, lggr logger.Logger) (*chain, error) {
	cfg := config.NewConfig(*dbchain.Cfg, lggr)
	lggr = lggr.With("chainID", dbchain.ID, "chainSet", "solana")
	var ch = chain{
		id:          dbchain.ID,
		cfg:         cfg,
		orm:         orm,
		lggr:        lggr.Named("Chain"),
		clientCache: map[string]cachedClient{},
	}
	tc := func() (solanaclient.ReaderWriter, error) {
		return ch.getClient()
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

func (c *chain) UpdateConfig(cfg *db.ChainCfg) {
	c.cfg.Update(*cfg)
}

func (c *chain) TxManager() solana.TxManager {
	return c.txm
}

func (c *chain) Reader() (solanaclient.Reader, error) {
	return c.getClient()
}

// getClient returns a client, randomly selecting one from available and valid nodes
func (c *chain) getClient() (solanaclient.ReaderWriter, error) {
	var node db.Node
	var client solanaclient.ReaderWriter
	nodes, cnt, err := c.orm.NodesForChain(c.id, 0, math.MaxInt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get nodes")
	}
	if cnt == 0 {
		return nil, errors.New("no nodes available")
	}
	rand.Seed(time.Now().Unix()) // seed randomness otherwise it will return the same each time
	// #nosec
	index := rand.Perm(len(nodes)) // list of node indexes to try
	for _, i := range index {
		node = nodes[i]
		// create client and check
		client, err = c.verifiedClient(node)
		// if error, try another node
		if err != nil {
			c.lggr.Warnw("failed to create node", "name", node.Name, "solana-url", node.SolanaURL, "error", err.Error())
			continue
		}
		// if all checks passed, mark found and break loop
		break
	}
	// if no valid node found, exit with error
	if client == nil {
		return nil, errors.New("no node valid nodes available")
	}
	c.lggr.Debugw("Created client", "name", node.Name, "solana-url", node.SolanaURL)
	return client, nil
}

// verifiedClient returns a client for node or an error if the chain id does not match.
func (c *chain) verifiedClient(node db.Node) (solanaclient.ReaderWriter, error) {
	url := node.SolanaURL
	var err error

	// check if cached client exists
	c.clientLock.RLock()
	client, exists := c.clientCache[url]
	c.clientLock.RUnlock()

	if !exists {
		// create client
		client.rw, err = solanaclient.NewClient(url, c.cfg, DefaultRequestTimeout, c.lggr.Named("Client-"+node.Name))
		if err != nil {
			return nil, errors.Wrap(err, "failed to create client")
		}

		client.id, err = client.rw.ChainID()
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch ChainID in checkClient")
		}
	}

	// check chainID matches expected chainID
	expectedID := strings.ToLower(c.id)
	if client.id != expectedID {
		return nil, errors.Errorf("client returned mismatched chain id (expected: %s, got: %s): %s", expectedID, client.id, url)
	}

	// save client if doesn't exist and checks have passed
	// if checks failed, client is not saved and can retry when a new client is requested
	if !exists {
		c.clientLock.Lock()
		// recheck when writing to prevent parallel writes (discard duplicate if exists)
		if cached, exists := c.clientCache[url]; !exists {
			c.clientCache[url] = client
		} else {
			client = cached
		}
		c.clientLock.Unlock()
	}

	return client.rw, nil
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
