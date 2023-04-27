package starknet

import (
	"context"
	"math/big"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"

	starkChain "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/chain"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/txm"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/starknet"
	"github.com/smartcontractkit/chainlink/v2/core/chains/starknet/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var _ starkChain.Chain = (*chain)(nil)

type chain struct {
	utils.StartStopOnce
	id   string
	cfg  config.Config
	cfgs types.Configs
	lggr logger.Logger
	txm  txm.StarkTXM
}

func newChain(id string, cfg config.Config, ks keystore.StarkNet, cfgs types.Configs, lggr logger.Logger) (ch *chain, err error) {
	lggr = lggr.With("starknetChainID", id)
	ch = &chain{
		id:   id,
		cfg:  cfg,
		cfgs: cfgs,
		lggr: lggr.Named("Chain"),
	}

	getClient := func() (*starknet.Client, error) {
		return ch.getClient()
	}

	ch.txm, err = txm.New(lggr, ks, cfg, getClient)
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

// getClient returns a client, randomly selecting one from available and valid nodes
func (c *chain) getClient() (*starknet.Client, error) {
	var node db.Node
	var client *starknet.Client
	nodes, err := c.cfgs.Nodes(c.id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get nodes")
	}
	if len(nodes) == 0 {
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
			c.lggr.Warnw("failed to create node", "name", node.Name, "starknet-url", node.URL, "error", err.Error())
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
	report := map[string]error{c.Name(): c.StartStopOnce.Healthy()}
	maps.Copy(report, c.txm.HealthReport())
	return report
}

func (c *chain) SendTx(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	return errors.New("unsupported") //TODO
}
