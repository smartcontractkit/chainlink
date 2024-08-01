package cosmos

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink-common/pkg/chains"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/adapters"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/client"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/txm"
)

// defaultRequestTimeout is the default Cosmos client timeout.
// Note that while the cosmos node is processing a heavy block,
// requests can be delayed significantly (https://github.com/tendermint/tendermint/issues/6899),
// however there's nothing we can do but wait until the block is processed.
// So we set a fairly high timeout here.
// TODO(BCI-979): Remove this, or make this configurable with the updated client.
const defaultRequestTimeout = 30 * time.Second

// Chain is a wrap for easy use in other places in the core node
type Chain = adapters.Chain

// ChainOpts holds options for configuring a Chain.
type ChainOpts struct {
	Logger   logger.Logger
	DS       sqlutil.DataSource
	KeyStore loop.Keystore
}

func (o *ChainOpts) Validate() (err error) {
	required := func(s string) error {
		return fmt.Errorf("%s is required", s)
	}
	if o.Logger == nil {
		err = errors.Join(err, required("Logger'"))
	}
	if o.DS == nil {
		err = errors.Join(err, required("DataSource"))
	}
	if o.KeyStore == nil {
		err = errors.Join(err, required("KeyStore"))
	}
	return
}

func NewChain(cfg *config.TOMLConfig, opts ChainOpts) (adapters.Chain, error) {
	if !cfg.IsEnabled() {
		return nil, fmt.Errorf("cannot create new chain with ID %s, the chain is disabled", *cfg.ChainID)
	}
	c, err := newChain(*cfg.ChainID, cfg, opts.DS, opts.KeyStore, opts.Logger)
	if err != nil {
		return nil, err
	}
	return c, nil
}

var _ adapters.Chain = (*chain)(nil)

type chain struct {
	services.StateMachine
	id   string
	cfg  *config.TOMLConfig
	txm  *txm.Txm
	lggr logger.Logger
}

func newChain(id string, cfg *config.TOMLConfig, ds sqlutil.DataSource, ks loop.Keystore, lggr logger.Logger) (*chain, error) {
	lggr = logger.With(lggr, "cosmosChainID", id)
	var ch = chain{
		id:   id,
		cfg:  cfg,
		lggr: logger.Named(lggr, "Chain"),
	}
	tc := func() (client.ReaderWriter, error) {
		return ch.getClient("")
	}
	gpe := client.NewMustGasPriceEstimator([]client.GasPricesEstimator{
		client.NewClosureGasPriceEstimator(func() (map[string]sdk.DecCoin, error) {
			return map[string]sdk.DecCoin{
				cfg.GasToken(): sdk.NewDecCoinFromDec(cfg.GasToken(), cfg.FallbackGasPrice()),
			}, nil
		}),
	}, lggr)
	ch.txm = txm.NewTxm(ds, tc, *gpe, ch.id, cfg, ks, lggr)

	return &ch, nil
}

func (c *chain) Name() string {
	return c.lggr.Name()
}

func (c *chain) ID() string {
	return c.id
}

func (c *chain) ChainID() string {
	return c.id
}

func (c *chain) Config() config.Config {
	return c.cfg
}

func (c *chain) TxManager() adapters.TxManager {
	return c.txm
}

func (c *chain) Reader(name string) (client.Reader, error) {
	return c.getClient(name)
}

// getClient returns a client, optionally requiring a specific node by name.
func (c *chain) getClient(name string) (client.ReaderWriter, error) {
	var node db.Node
	if name == "" { // Any node
		nodes, err := c.cfg.ListNodes()
		if err != nil {
			return nil, fmt.Errorf("failed to list nodes: %w", err)
		}
		if len(nodes) == 0 {
			return nil, errors.New("no nodes available")
		}
		nodeIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(nodes))))
		if err != nil {
			return nil, fmt.Errorf("could not generate a random node index: %w", err)
		}
		node = nodes[nodeIndex.Int64()]
	} else { // Named node
		var err error
		node, err = c.cfg.GetNode(name)
		if err != nil {
			return nil, fmt.Errorf("failed to get node named %s: %w", name, err)
		}
		if node.CosmosChainID != c.id {
			return nil, fmt.Errorf("failed to create client for chain %s with node %s: wrong chain id %s", c.id, name, node.CosmosChainID)
		}
	}
	client, err := client.NewClient(c.id, node.TendermintURL, defaultRequestTimeout, logger.Named(c.lggr, "Client."+name))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
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
	return errors.Join(
		c.StateMachine.Ready(),
		c.txm.Ready(),
	)
}

func (c *chain) HealthReport() map[string]error {
	m := map[string]error{c.Name(): c.Healthy()}
	services.CopyHealth(m, c.txm.HealthReport())
	return m
}

// ChainService interface
func (c *chain) GetChainStatus(ctx context.Context) (types.ChainStatus, error) {
	toml, err := c.cfg.TOMLString()
	if err != nil {
		return types.ChainStatus{}, err
	}
	return types.ChainStatus{
		ID:      c.id,
		Enabled: *c.cfg.Enabled,
		Config:  toml,
	}, nil
}
func (c *chain) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) (stats []types.NodeStatus, nextPageToken string, total int, err error) {
	return chains.ListNodeStatuses(int(pageSize), pageToken, c.listNodeStatuses)
}

func (c *chain) Transact(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	fromAcc, err := sdk.AccAddressFromBech32(from)
	if err != nil {
		return fmt.Errorf("failed to parse from account: %s", fromAcc)
	}
	toAcc, err := sdk.AccAddressFromBech32(to)
	if err != nil {
		return fmt.Errorf("failed to parse from account: %s", toAcc)
	}
	coin := sdk.Coin{Amount: sdk.NewIntFromBigInt(amount), Denom: c.Config().GasToken()}

	txm := c.TxManager()

	if balanceCheck {
		var reader client.Reader
		reader, err = c.Reader("")
		if err != nil {
			return fmt.Errorf("chain unreachable: %v", err)
		}
		gasPrice, err2 := txm.GasPrice()
		if err2 != nil {
			return fmt.Errorf("gas price unavailable: %v", err2)
		}

		err = validateBalance(reader, gasPrice, fromAcc, coin)
		if err != nil {
			return fmt.Errorf("failed to validate balance: %v", err)
		}
	}

	sendMsg := bank.NewMsgSend(fromAcc, toAcc, sdk.Coins{coin})
	_, err = txm.Enqueue(ctx, "", sendMsg)
	if err != nil {
		return fmt.Errorf("failed to enqueue tx: %w", err)
	}
	return nil
}

// TODO BCF-2602 statuses are static for non-evm chain and should be dynamic
func (c *chain) listNodeStatuses(start, end int) ([]types.NodeStatus, int, error) {
	stats := make([]types.NodeStatus, 0)
	total := len(c.cfg.Nodes)
	if start >= total {
		return stats, total, chains.ErrOutOfRange
	}
	if end > total {
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

// maxGasUsedTransfer is an upper bound on how much gas we expect a MsgSend for a single coin to use.
const maxGasUsedTransfer = 100_000

// validateBalance validates that fromAddr's balance can cover coin, including fees at gasPrice.
func validateBalance(reader client.Reader, gasPrice sdk.DecCoin, fromAddr sdk.AccAddress, coin sdk.Coin) error {
	balance, err := reader.Balance(fromAddr, coin.GetDenom())
	if err != nil {
		return err
	}

	fee := gasPrice.Amount.MulInt64(maxGasUsedTransfer).RoundInt()
	need := coin.Amount.Add(fee)

	if balance.Amount.LT(need) {
		return fmt.Errorf("balance %q is too low for this transaction to be executed: need %s total, including %s fee", balance, need, fee)
	}
	return nil
}
