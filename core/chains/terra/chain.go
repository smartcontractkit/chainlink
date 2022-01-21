package terra

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/sqlx"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra"
	terraclient "github.com/smartcontractkit/chainlink-terra/pkg/terra/client"
	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/chains/terra/terratxm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// DefaultRequestTimeout is the default Terra client timeout.
const DefaultRequestTimeout = 10 * time.Second

var _ terra.Chain = (*chain)(nil)

type chain struct {
	utils.StartStopOnce
	id     string
	cfg    terra.Config
	client *terraclient.Client
	txm    *terratxm.Txm
	lggr   logger.Logger
}

// NewChain returns a new chain backed by node.
func NewChain(db *sqlx.DB, ks keystore.Terra, logCfg pg.LogConfig, eb pg.EventBroadcaster, dbchain db.Chain, lggr logger.Logger) (*chain, error) {
	if !dbchain.Enabled {
		return nil, errors.Errorf("cannot create new chain with ID %s, the chain is disabled", dbchain.ID)
	}
	if len(dbchain.Nodes) == 0 {
		return nil, errors.Errorf("no nodes for Terra chain: %s", dbchain.ID)
	}
	cfg := terra.NewConfig(dbchain.Cfg, terra.DefaultConfigSet, lggr)
	lggr = lggr.With("terraChainID", dbchain.ID)
	node := dbchain.Nodes[0] // TODO multi-node client pool https://app.shortcut.com/chainlinklabs/story/26278/terra-multi-node-client-pools
	lggr.Debugw(fmt.Sprintf("Terra chain %q has %d nodes - using %q", dbchain.ID, len(dbchain.Nodes), node.Name),
		"tendermint-url", node.TendermintURL)
	gpeFCD, err := terraclient.NewFCDGasPriceEstimator(node.FCDURL, DefaultRequestTimeout, lggr)
	if err != nil {
		return nil, err
	}
	gpe := terraclient.NewMustGasPriceEstimator([]terraclient.GasPricesEstimator{
		terraclient.NewCachingGasPriceEstimator(gpeFCD, lggr),
		terraclient.NewClosureGasPriceEstimator(func() (map[string]sdk.DecCoin, error) {
			return map[string]sdk.DecCoin{
				"uluna": sdk.NewDecCoinFromDec("uluna", cfg.FallbackGasPriceULuna()),
			}, nil
		}),
	}, lggr)
	client, err := terraclient.NewClient(dbchain.ID,
		node.TendermintURL, DefaultRequestTimeout, lggr.Named("Client"))
	if err != nil {
		return nil, err
	}
	txm, err := terratxm.NewTxm(db, client, *gpe, dbchain.ID, cfg, ks, lggr, logCfg, eb)
	if err != nil {
		return nil, err
	}
	return &chain{
		id:     dbchain.ID,
		cfg:    cfg,
		client: client,
		txm:    txm,
		lggr:   lggr.Named("Chain"),
	}, nil
}

func (c *chain) ID() string {
	return c.id
}

func (c *chain) Config() terra.Config {
	return c.cfg
}

func (c *chain) UpdateConfig(cfg db.ChainCfg) {
	c.cfg.Update(cfg)
}

func (c *chain) MsgEnqueuer() terra.MsgEnqueuer {
	return c.txm
}

func (c *chain) Reader() terraclient.Reader {
	return c.client
}

func (c *chain) Start() error {
	return c.StartOnce("Chain", func() error {
		c.lggr.Debug("Starting")
		//TODO dial client?

		c.lggr.Debug("Starting txm")
		return c.txm.Start()
	})
}

func (c *chain) Close() error {
	return c.StopOnce("Chain", func() error {
		c.lggr.Debug("Stopping")
		c.lggr.Debug("Stopping txm")
		return c.txm.Close()
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
