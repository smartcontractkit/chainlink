package starknet

import (
	"context"
	"math"
	"math/rand"
	"time"

	caigotypes "github.com/dontpanicdao/caigo/types"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/txm"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/starknet"
	starkChain "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/starknet/chain"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/starknet/db"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/starknet/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var _ starkChain.Chain = (*chain)(nil)

type chain struct {
	utils.StartStopOnce
	id   string
	cfg  starknet.Config
	orm  types.ORM
	lggr logger.Logger
	txm  txm.StarkTXM
}

func NewChain(db *sqlx.DB, ks keystore.StarkNet, dbchain types.DBChain, orm types.ORM, lggr logger.Logger) (ch *chain, err error) {
	cfg := starknet.NewConfig(*dbchain.Cfg, lggr)
	lggr = lggr.With("starknetChainID", dbchain.ID)

	ch = &chain{
		id:   dbchain.ID,
		cfg:  cfg,
		orm:  orm,
		lggr: lggr.Named("Chain"),
	}

	getClient := func() (caigotypes.Provider, error) {
		return ch.getClient()
	}

	ch.txm, err = txm.New(lggr, ks, cfg, getClient)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (c *chain) Config() starknet.Config {
	return c.cfg
}

func (c *chain) UpdateConfig(cfg *db.ChainCfg) {
	c.cfg.Update(*cfg)
}

func (c *chain) TxManager() txm.TxManager {
	return nil
}

func (c *chain) Reader() (starknet.Reader, error) {
	return nil, nil
}

// getClient returns a client, randomly selecting one from available and valid nodes
func (c *chain) getClient() (*starknet.Client, error) {
	var node db.Node
	var client *starknet.Client
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
	timeout := c.cfg.RequestTimeout()
	for _, i := range index {
		node = nodes[i]
		// create client and check
		client, err = starknet.NewClient(node.ChainID, node.URL, c.lggr, &timeout)
		// if error, try another node
		if err != nil {
			c.lggr.Warnw("failed to create node", "name", node.Name, "solana-url", node.URL, "error", err.Error())
			continue
		}
		// if all checks passed, mark found and break loop
		break
	}
	// if no valid node found, exit with error
	if client == nil {
		return nil, errors.New("no node valid nodes available")
	}
	c.lggr.Debugw("Created client", "name", node.Name, "solana-url", node.URL)
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

func (c *chain) Healthy() error {
	return c.StartStopOnce.Healthy()
}
