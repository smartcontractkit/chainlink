package cosmos

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/adapters"
	cosmosclient "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/client"
	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"
	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos/cosmostxm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// DefaultRequestTimeout is the default Cosmos client timeout.
// Note that while the cosmos node is processing a heavy block,
// requests can be delayed significantly (https://github.com/tendermint/tendermint/issues/6899),
// however there's nothing we can do but wait until the block is processed.
// So we set a fairly high timeout here.
// TODO(BCI-979): Remove this, or make this configurable with the updated client.
const DefaultRequestTimeout = 30 * time.Second

var _ adapters.Chain = (*chain)(nil)

type chain struct {
	utils.StartStopOnce
	id   string
	cfg  coscfg.Config
	txm  *cosmostxm.Txm
	cfgs types.Configs
	lggr logger.Logger
}

func newChain(id string, cfg coscfg.Config, db *sqlx.DB, ks keystore.Cosmos, logCfg pg.QConfig, eb pg.EventBroadcaster, cfgs types.Configs, lggr logger.Logger) (*chain, error) {
	lggr = lggr.With("cosmosChainID", id)
	var ch = chain{
		id:   id,
		cfg:  cfg,
		cfgs: cfgs,
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

func (c *chain) Config() coscfg.Config {
	return c.cfg
}

func (c *chain) TxManager() adapters.TxManager {
	return c.txm
}

func (c *chain) Reader(name string) (cosmosclient.Reader, error) {
	return c.getClient(name)
}

// getClient returns a client, optionally requiring a specific node by name.
func (c *chain) getClient(name string) (cosmosclient.ReaderWriter, error) {
	var node db.Node
	if name == "" { // Any node
		nodes, err := c.cfgs.Nodes(c.id)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get nodes")
		}
		if len(nodes) == 0 {
			return nil, errors.New("no nodes available")
		}
		nodeIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(nodes))))
		if err != nil {
			return nil, errors.Wrap(err, "could not generate a random node index")
		}
		node = nodes[nodeIndex.Int64()]
	} else { // Named node
		var err error
		node, err = c.cfgs.Node(name)
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

func (c *chain) HealthReport() map[string]error {
	return map[string]error{
		c.Name(): multierr.Combine(
			c.StartStopOnce.Healthy(),
			c.txm.Healthy()),
	}
}

func (c *chain) SendTx(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	return errors.New("unsupported") //TODO
}
