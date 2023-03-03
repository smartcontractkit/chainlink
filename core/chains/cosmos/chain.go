package cosmos

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/smartcontractkit/sqlx"

	cosmos "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos"
	cosmosclient "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/client"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"
	v2 "github.com/smartcontractkit/chainlink/core/config/v2"

	"github.com/smartcontractkit/chainlink/core/chains/cosmos/cosmostxm"
	"github.com/smartcontractkit/chainlink/core/chains/cosmos/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// DefaultRequestTimeout is the default Cosmos client timeout.
// Note that while the cosmos node is processing a heavy block,
// requests can be delayed significantly (https://github.com/tendermint/tendermint/issues/6899),
// however there's nothing we can do but wait until the block is processed.
// So we set a fairly high timeout here.
const DefaultRequestTimeout = 30 * time.Second

//go:generate mockery --quiet --name TxManager --srcpkg github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos --output ./mocks/ --case=underscore
//go:generate mockery --quiet --name Reader --srcpkg github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/client --output ./mocks/ --case=underscore
//go:generate mockery --quiet --name Chain --srcpkg github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos --output ./mocks/ --case=underscore
var _ cosmos.Chain = (*chain)(nil)

type chain struct {
	utils.StartStopOnce
	id           string
	cfg          cosmos.Config
	cfgImmutable bool // toml config is immutable
	txm          *cosmostxm.Txm
	orm          types.ORM
	lggr         logger.Logger
}

func newChain(id string, cfg cosmos.Config, db *sqlx.DB, ks keystore.Cosmos, logCfg pg.QConfig, eb pg.EventBroadcaster, orm types.ORM, lggr logger.Logger) (*chain, error) {
	lggr = lggr.With("cosmosChainID", id)
	var ch = chain{
		id:   id,
		cfg:  cfg,
		orm:  orm,
		lggr: lggr.Named("Chain"),
	}
	tc := func() (cosmosclient.ReaderWriter, error) {
		return ch.getClient("")
	}
	gpe := cosmosclient.NewMustGasPriceEstimator([]cosmosclient.GasPricesEstimator{
		cosmosclient.NewClosureGasPriceEstimator(func() (map[string]sdk.DecCoin, error) {
			return map[string]sdk.DecCoin{
				"uatom": sdk.NewDecCoinFromDec("uatom", cfg.FallbackGasPriceUAtom()),
			}, nil
		}),
	}, lggr)
	ch.txm = cosmostxm.NewTxm(db, tc, *gpe, ch.id, cfg, ks, lggr, logCfg, eb)

	return &ch, nil
}

func (c *chain) Name() string {
	return c.lggr.Name()
}

func (c *chain) ID() string {
	return c.id
}

func (c *chain) Config() cosmos.Config {
	return c.cfg
}

func (c *chain) UpdateConfig(cfg *db.ChainCfg) {
	if c.cfgImmutable {
		c.lggr.Criticalw("TOML configuration cannot be updated", "err", v2.ErrUnsupported)
		return
	}
	c.cfg.Update(*cfg)
}

func (c *chain) TxManager() cosmos.TxManager {
	return c.txm
}

func (c *chain) Reader(name string) (cosmosclient.Reader, error) {
	return c.getClient(name)
}

// getClient returns a client, optionally requiring a specific node by name.
func (c *chain) getClient(name string) (cosmosclient.ReaderWriter, error) {
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
		if node.CosmosChainID != c.id {
			return nil, fmt.Errorf("failed to create client for chain %s with node %s: wrong chain id %s", c.id, name, node.CosmosChainID)
		}
	}
	client, err := cosmosclient.NewClient(c.id, node.TendermintURL, DefaultRequestTimeout, c.lggr.Named("Client-"+name))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create client")
	}
	c.lggr.Debugw("Created client", "name", node.Name, "tendermint-url", node.TendermintURL)
	return client, nil
}

// Start starts cosmos chain.
func (c *chain) Start(ctx context.Context) error {
	return c.StartOnce("Chain", func() error {
		c.lggr.Debug("Starting")
		//TODO dial client?
		return c.txm.Start(ctx)
	})
}

func (c *chain) Close() error {
	return c.StopOnce("Chain", func() error {
		c.lggr.Debug("Stopping")
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

func (c *chain) HealthReport() map[string]error {
	return map[string]error{c.Name(): c.Healthy()}
}
