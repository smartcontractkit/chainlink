package cosmos

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/smartcontractkit/sqlx"

	relaychains "github.com/smartcontractkit/chainlink-relay/pkg/chains"
	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-relay/pkg/services"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/adapters"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/client"
	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos/cosmostxm"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

// defaultRequestTimeout is the default Cosmos client timeout.
// Note that while the cosmos node is processing a heavy block,
// requests can be delayed significantly (https://github.com/tendermint/tendermint/issues/6899),
// however there's nothing we can do but wait until the block is processed.
// So we set a fairly high timeout here.
// TODO(BCI-979): Remove this, or make this configurable with the updated client.
const defaultRequestTimeout = 30 * time.Second

var (
	// ErrChainIDEmpty is returned when chain is required but was empty.
	ErrChainIDEmpty = errors.New("chain id empty")
	// ErrChainIDInvalid is returned when a chain id does not match any configured chains.
	ErrChainIDInvalid = errors.New("chain id does not match any local chains")
)

// Chain is a wrap for easy use in other places in the core node
type Chain = adapters.Chain

// ChainOpts holds options for configuring a Chain.
type ChainOpts struct {
	QueryConfig      pg.QConfig
	Logger           logger.Logger
	DB               *sqlx.DB
	KeyStore         loop.Keystore
	EventBroadcaster pg.EventBroadcaster
}

func (o *ChainOpts) Validate() (err error) {
	required := func(s string) error {
		return fmt.Errorf("%s is required", s)
	}
	if o.QueryConfig == nil {
		err = multierr.Append(err, required("Config"))
	}
	if o.Logger == nil {
		err = multierr.Append(err, required("Logger'"))
	}
	if o.DB == nil {
		err = multierr.Append(err, required("DB"))
	}
	if o.KeyStore == nil {
		err = multierr.Append(err, required("KeyStore"))
	}
	if o.EventBroadcaster == nil {
		err = multierr.Append(err, required("EventBroadcaster"))
	}
	return
}

func NewChain(cfg *coscfg.TOMLConfig, opts ChainOpts) (adapters.Chain, error) {
	if !cfg.IsEnabled() {
		return nil, fmt.Errorf("cannot create new chain with ID %s, the chain is disabled", *cfg.ChainID)
	}
	c, err := newChain(*cfg.ChainID, cfg, opts.DB, opts.KeyStore, opts.QueryConfig, opts.EventBroadcaster, opts.Logger)
	if err != nil {
		return nil, err
	}
	return c, nil
}

var _ adapters.Chain = (*chain)(nil)

type chain struct {
	services.StateMachine
	id   string
	cfg  *coscfg.TOMLConfig
	txm  *cosmostxm.Txm
	lggr logger.Logger
}

func newChain(id string, cfg *coscfg.TOMLConfig, db *sqlx.DB, ks loop.Keystore, logCfg pg.QConfig, eb pg.EventBroadcaster, lggr logger.Logger) (*chain, error) {
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
	ch.txm = cosmostxm.NewTxm(db, tc, *gpe, ch.id, cfg, ks, lggr, logCfg, eb)

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

func (c *chain) Config() coscfg.Config {
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
	return multierr.Combine(
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
func (c *chain) GetChainStatus(ctx context.Context) (relaytypes.ChainStatus, error) {
	toml, err := c.cfg.TOMLString()
	if err != nil {
		return relaytypes.ChainStatus{}, err
	}
	return relaytypes.ChainStatus{
		ID:      c.id,
		Enabled: *c.cfg.Enabled,
		Config:  toml,
	}, nil
}
func (c *chain) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) (stats []relaytypes.NodeStatus, nextPageToken string, total int, err error) {
	return relaychains.ListNodeStatuses(int(pageSize), pageToken, c.listNodeStatuses)
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
func (c *chain) listNodeStatuses(start, end int) ([]relaytypes.NodeStatus, int, error) {
	stats := make([]relaytypes.NodeStatus, 0)
	total := len(c.cfg.Nodes)
	if start >= total {
		return stats, total, relaychains.ErrOutOfRange
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

func nodeStatus(n *coscfg.Node, id relay.ChainID) (relaytypes.NodeStatus, error) {
	var s relaytypes.NodeStatus
	s.ChainID = id
	s.Name = *n.Name
	b, err := toml.Marshal(n)
	if err != nil {
		return relaytypes.NodeStatus{}, err
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
		return errors.Errorf("balance %q is too low for this transaction to be executed: need %s total, including %s fee", balance, need, fee)
	}
	return nil
}
