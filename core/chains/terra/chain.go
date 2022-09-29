package terra

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra"
	terraclient "github.com/smartcontractkit/chainlink-terra/pkg/terra/client"
	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/terra/monitor"
	"github.com/smartcontractkit/chainlink/core/chains/terra/terratxm"
	"github.com/smartcontractkit/chainlink/core/chains/terra/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// DefaultRequestTimeout is the default Terra client timeout.
// Note that while the terra node is processing a heavy block,
// requests can be delayed significantly (https://github.com/tendermint/tendermint/issues/6899),
// however there's nothing we can do but wait until the block is processed.
// So we set a fairly high timeout here.
const DefaultRequestTimeout = 30 * time.Second

//go:generate mockery --name TxManager --srcpkg github.com/smartcontractkit/chainlink-terra/pkg/terra --output ./mocks/ --case=underscore
//go:generate mockery --name Reader --srcpkg github.com/smartcontractkit/chainlink-terra/pkg/terra/client --output ./mocks/ --case=underscore
//go:generate mockery --name Chain --srcpkg github.com/smartcontractkit/chainlink-terra/pkg/terra --output ./mocks/ --case=underscore
var _ terra.Chain = (*chain)(nil)

type chain struct {
	utils.StartStopOnce
	id             string
	cfg            terra.Config
	txm            *terratxm.Txm
	balanceMonitor services.ServiceCtx
	orm            types.ORM
	lggr           logger.Logger
}

// NewChain returns a new chain backed by node.
func NewChain(db *sqlx.DB, ks keystore.Terra, logCfg pg.LogConfig, eb pg.EventBroadcaster, dbchain types.DBChain, orm types.ORM, lggr logger.Logger) (*chain, error) {
	cfg := terra.NewConfig(*dbchain.Cfg, lggr)
	lggr = lggr.With("terraChainID", dbchain.ID)
	var ch = chain{
		id:   dbchain.ID,
		cfg:  cfg,
		orm:  orm,
		lggr: lggr.Named("Chain"),
	}
	tc := func() (terraclient.ReaderWriter, error) {
		return ch.getClient("")
	}
	gpeFCD := terraclient.NewFCDGasPriceEstimator(cfg, DefaultRequestTimeout, lggr)
	gpe := terraclient.NewMustGasPriceEstimator([]terraclient.GasPricesEstimator{
		terraclient.NewCachingGasPriceEstimator(gpeFCD, lggr),
		terraclient.NewClosureGasPriceEstimator(func() (map[string]sdk.DecCoin, error) {
			return map[string]sdk.DecCoin{
				"uluna": sdk.NewDecCoinFromDec("uluna", cfg.FallbackGasPriceULuna()),
			}, nil
		}),
	}, lggr)
	ch.txm = terratxm.NewTxm(db, tc, *gpe, ch.id, cfg, ks, lggr, logCfg, eb)
	ch.balanceMonitor = monitor.NewBalanceMonitor(ch.id, cfg, lggr, ks, ch.Reader)

	return &ch, nil
}

func (c *chain) ID() string {
	return c.id
}

func (c *chain) Config() terra.Config {
	return c.cfg
}

func (c *chain) UpdateConfig(cfg *db.ChainCfg) {
	c.cfg.Update(*cfg)
}

func (c *chain) TxManager() terra.TxManager {
	return c.txm
}

func (c *chain) Reader(name string) (terraclient.Reader, error) {
	return c.getClient(name)
}

// getClient returns a client, optionally requiring a specific node by name.
func (c *chain) getClient(name string) (terraclient.ReaderWriter, error) {
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
		if node.TerraChainID != c.id {
			return nil, fmt.Errorf("failed to create client for chain %s with node %s: wrong chain id %s", c.id, name, node.TerraChainID)
		}
	}
	client, err := terraclient.NewClient(c.id, node.TendermintURL, DefaultRequestTimeout, c.lggr.Named("Client-"+name))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create client")
	}
	c.lggr.Debugw("Created client", "name", node.Name, "tendermint-url", node.TendermintURL)
	return client, nil
}

// Start starts terra chain.
func (c *chain) Start(ctx context.Context) error {
	return c.StartOnce("Chain", func() error {
		c.lggr.Debug("Starting")
		//TODO dial client?

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
