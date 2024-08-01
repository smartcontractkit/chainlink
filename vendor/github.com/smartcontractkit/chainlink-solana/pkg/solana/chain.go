package solana

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"sync"
	"time"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/smartcontractkit/chainlink-common/pkg/chains"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	relaytypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/monitor"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/txm"
)

type Chain interface {
	types.ChainService

	ID() string
	Config() config.Config
	TxManager() TxManager
	// Reader returns a new Reader from the available list of nodes (if there are multiple, it will randomly select one)
	Reader() (client.Reader, error)
}

// DefaultRequestTimeout is the default Solana client timeout.
const DefaultRequestTimeout = 30 * time.Second

// ChainOpts holds options for configuring a Chain.
type ChainOpts struct {
	Logger   logger.Logger
	KeyStore core.Keystore
}

func (o *ChainOpts) Validate() (err error) {
	required := func(s string) error {
		return fmt.Errorf("%s is required", s)
	}
	if o.Logger == nil {
		err = errors.Join(err, required("Logger"))
	}
	if o.KeyStore == nil {
		err = errors.Join(err, required("KeyStore"))
	}
	return
}

func (o *ChainOpts) GetLogger() logger.Logger {
	return o.Logger
}

func NewChain(cfg *config.TOMLConfig, opts ChainOpts) (Chain, error) {
	if !cfg.IsEnabled() {
		return nil, fmt.Errorf("cannot create new chain with ID %s: chain is disabled", *cfg.ChainID)
	}
	c, err := newChain(*cfg.ChainID, cfg, opts.KeyStore, opts.Logger)
	if err != nil {
		return nil, err
	}
	return c, nil
}

var _ Chain = (*chain)(nil)

type chain struct {
	services.StateMachine
	id             string
	cfg            *config.TOMLConfig
	txm            *txm.Txm
	balanceMonitor services.Service
	lggr           logger.Logger

	// tracking node chain id for verification
	clientCache map[string]*verifiedCachedClient // map URL -> {client, chainId} [mainnet/testnet/devnet/localnet]
	clientLock  sync.RWMutex
}

type verifiedCachedClient struct {
	chainID         string
	expectedChainID string
	nodeURL         string

	chainIDVerified     bool
	chainIDVerifiedLock sync.RWMutex

	client.ReaderWriter
}

func (v *verifiedCachedClient) verifyChainID() (bool, error) {
	v.chainIDVerifiedLock.RLock()
	if v.chainIDVerified {
		v.chainIDVerifiedLock.RUnlock()
		return true, nil
	}
	v.chainIDVerifiedLock.RUnlock()

	var err error

	v.chainIDVerifiedLock.Lock()
	defer v.chainIDVerifiedLock.Unlock()

	v.chainID, err = v.ReaderWriter.ChainID()
	if err != nil {
		v.chainIDVerified = false
		return v.chainIDVerified, fmt.Errorf("failed to fetch ChainID in verifiedCachedClient: %w", err)
	}

	// check chainID matches expected chainID
	expectedChainID := strings.ToLower(v.expectedChainID)
	if v.chainID != expectedChainID {
		v.chainIDVerified = false
		return v.chainIDVerified, fmt.Errorf("client returned mismatched chain id (expected: %s, got: %s): %s", expectedChainID, v.chainID, v.nodeURL)
	}

	v.chainIDVerified = true

	return v.chainIDVerified, nil
}

func (v *verifiedCachedClient) SendTx(ctx context.Context, tx *solanago.Transaction) (solanago.Signature, error) {
	verified, err := v.verifyChainID()
	if !verified {
		return [64]byte{}, err
	}

	return v.ReaderWriter.SendTx(ctx, tx)
}

func (v *verifiedCachedClient) SimulateTx(ctx context.Context, tx *solanago.Transaction, opts *rpc.SimulateTransactionOpts) (*rpc.SimulateTransactionResult, error) {
	verified, err := v.verifyChainID()
	if !verified {
		return nil, err
	}

	return v.ReaderWriter.SimulateTx(ctx, tx, opts)
}

func (v *verifiedCachedClient) SignatureStatuses(ctx context.Context, sigs []solanago.Signature) ([]*rpc.SignatureStatusesResult, error) {
	verified, err := v.verifyChainID()
	if !verified {
		return nil, err
	}

	return v.ReaderWriter.SignatureStatuses(ctx, sigs)
}

func (v *verifiedCachedClient) Balance(addr solanago.PublicKey) (uint64, error) {
	verified, err := v.verifyChainID()
	if !verified {
		return 0, err
	}

	return v.ReaderWriter.Balance(addr)
}

func (v *verifiedCachedClient) SlotHeight() (uint64, error) {
	verified, err := v.verifyChainID()
	if !verified {
		return 0, err
	}

	return v.ReaderWriter.SlotHeight()
}

func (v *verifiedCachedClient) LatestBlockhash() (*rpc.GetLatestBlockhashResult, error) {
	verified, err := v.verifyChainID()
	if !verified {
		return nil, err
	}

	return v.ReaderWriter.LatestBlockhash()
}

func (v *verifiedCachedClient) ChainID() (string, error) {
	verified, err := v.verifyChainID()
	if !verified {
		return "", err
	}

	return v.chainID, nil
}

func (v *verifiedCachedClient) GetFeeForMessage(msg string) (uint64, error) {
	verified, err := v.verifyChainID()
	if !verified {
		return 0, err
	}

	return v.ReaderWriter.GetFeeForMessage(msg)
}

func (v *verifiedCachedClient) GetAccountInfoWithOpts(ctx context.Context, addr solanago.PublicKey, opts *rpc.GetAccountInfoOpts) (*rpc.GetAccountInfoResult, error) {
	verified, err := v.verifyChainID()
	if !verified {
		return nil, err
	}

	return v.ReaderWriter.GetAccountInfoWithOpts(ctx, addr, opts)
}

func newChain(id string, cfg *config.TOMLConfig, ks loop.Keystore, lggr logger.Logger) (*chain, error) {
	lggr = logger.With(lggr, "chainID", id, "chain", "solana")
	var ch = chain{
		id:          id,
		cfg:         cfg,
		lggr:        logger.Named(lggr, "Chain"),
		clientCache: map[string]*verifiedCachedClient{},
	}
	tc := func() (client.ReaderWriter, error) {
		return ch.getClient()
	}
	ch.txm = txm.NewTxm(ch.id, tc, cfg, ks, lggr)
	bc := func() (monitor.BalanceClient, error) {
		return ch.getClient()
	}
	ch.balanceMonitor = monitor.NewBalanceMonitor(ch.id, cfg, lggr, ks, bc)
	return &ch, nil
}

// ChainService interface
func (c *chain) GetChainStatus(ctx context.Context) (relaytypes.ChainStatus, error) {
	toml, err := c.cfg.TOMLString()
	if err != nil {
		return relaytypes.ChainStatus{}, err
	}
	return relaytypes.ChainStatus{
		ID:      c.id,
		Enabled: c.cfg.IsEnabled(),
		Config:  toml,
	}, nil
}

func (c *chain) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) (stats []relaytypes.NodeStatus, nextPageToken string, total int, err error) {
	return chains.ListNodeStatuses(int(pageSize), pageToken, c.listNodeStatuses)
}

func (c *chain) Transact(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	return c.sendTx(ctx, from, to, amount, balanceCheck)
}

func (c *chain) listNodeStatuses(start, end int) ([]relaytypes.NodeStatus, int, error) {
	stats := make([]relaytypes.NodeStatus, 0)
	total := len(c.cfg.Nodes)
	if start >= total {
		return stats, total, chains.ErrOutOfRange
	}
	if end > total {
		end = total
	}
	nodes := c.cfg.Nodes[start:end]
	for _, node := range nodes {
		stat, err := config.NodeStatus(node, c.ChainID())
		if err != nil {
			return stats, total, err
		}
		stats = append(stats, stat)
	}
	return stats, total, nil
}

func (c *chain) Name() string {
	return c.lggr.Name()
}

func (c *chain) ID() string {
	return c.id
}

func (c *chain) Config() config.Config {
	return c.cfg
}

func (c *chain) TxManager() TxManager {
	return c.txm
}

func (c *chain) Reader() (client.Reader, error) {
	return c.getClient()
}

func (c *chain) ChainID() string {
	return c.id
}

// getClient returns a client, randomly selecting one from available and valid nodes
func (c *chain) getClient() (client.ReaderWriter, error) {
	var node *config.Node
	var client client.ReaderWriter
	nodes := c.cfg.ListNodes()
	if len(nodes) == 0 {
		return nil, errors.New("no nodes available")
	}
	// #nosec
	index := rand.Perm(len(nodes)) // list of node indexes to try
	for _, i := range index {
		node = nodes[i]
		// create client and check
		var err error
		client, err = c.verifiedClient(node)
		// if error, try another node
		if err != nil {
			c.lggr.Warnw("failed to create node", "name", node.Name, "solana-url", node.URL, "err", err.Error())
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

// verifiedClient returns a client for node or an error if fails to create the client.
// The client will still be returned if the nodes are not valid, or the chain id doesn't match.
// Further client calls will try and verify the client, and fail if the client is still not valid.
func (c *chain) verifiedClient(node *config.Node) (client.ReaderWriter, error) {
	if node == nil {
		return nil, fmt.Errorf("nil node")
	}

	if node.Name == nil || node.URL == nil {
		return nil, fmt.Errorf("node config contains nil: %+v", node)
	}

	url := node.URL.String()
	var err error

	// check if cached client exists
	c.clientLock.RLock()
	cl, exists := c.clientCache[url]
	c.clientLock.RUnlock()

	if !exists {
		cl = &verifiedCachedClient{
			nodeURL:         url,
			expectedChainID: c.id,
		}
		// create client
		cl.ReaderWriter, err = client.NewClient(url, c.cfg, DefaultRequestTimeout, logger.Named(c.lggr, "Client."+*node.Name))
		if err != nil {
			return nil, fmt.Errorf("failed to create client: %w", err)
		}

		c.clientLock.Lock()
		// recheck when writing to prevent parallel writes (discard duplicate if exists)
		if cached, exists := c.clientCache[url]; !exists {
			c.clientCache[url] = cl
		} else {
			cl = cached
		}
		c.clientLock.Unlock()
	}

	return cl, nil
}

func (c *chain) Start(ctx context.Context) error {
	return c.StartOnce("Chain", func() error {
		c.lggr.Debug("Starting")
		c.lggr.Debug("Starting txm")
		c.lggr.Debug("Starting balance monitor")
		var ms services.MultiStart
		return ms.Start(ctx, c.txm, c.balanceMonitor)
	})
}

func (c *chain) Close() error {
	return c.StopOnce("Chain", func() error {
		c.lggr.Debug("Stopping")
		c.lggr.Debug("Stopping txm")
		c.lggr.Debug("Stopping balance monitor")
		return services.CloseAll(c.txm, c.balanceMonitor)
	})
}

func (c *chain) Ready() error {
	return errors.Join(
		c.StateMachine.Ready(),
		c.txm.Ready(),
	)
}

func (c *chain) HealthReport() map[string]error {
	report := map[string]error{c.Name(): c.Healthy()}
	services.CopyHealth(report, c.txm.HealthReport())
	return report
}

func (c *chain) sendTx(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	reader, err := c.Reader()
	if err != nil {
		return fmt.Errorf("chain unreachable: %w", err)
	}

	fromKey, err := solanago.PublicKeyFromBase58(from)
	if err != nil {
		return fmt.Errorf("failed to parse from key: %w", err)
	}
	toKey, err := solanago.PublicKeyFromBase58(to)
	if err != nil {
		return fmt.Errorf("failed to parse to key: %w", err)
	}
	if !amount.IsUint64() {
		return fmt.Errorf("amount %s overflows uint64", amount)
	}
	amountI := amount.Uint64()

	blockhash, err := reader.LatestBlockhash()
	if err != nil {
		return fmt.Errorf("failed to get latest block hash: %w", err)
	}
	tx, err := solanago.NewTransaction(
		[]solanago.Instruction{
			system.NewTransferInstruction(
				amountI,
				fromKey,
				toKey,
			).Build(),
		},
		blockhash.Value.Blockhash,
		solanago.TransactionPayer(fromKey),
	)
	if err != nil {
		return fmt.Errorf("failed to create tx: %w", err)
	}

	if balanceCheck {
		if err = solanaValidateBalance(reader, fromKey, amountI, tx.Message.ToBase64()); err != nil {
			return fmt.Errorf("failed to validate balance: %w", err)
		}
	}

	txm := c.TxManager()
	err = txm.Enqueue("", tx)
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	return nil
}

func solanaValidateBalance(reader client.Reader, from solanago.PublicKey, amount uint64, msg string) error {
	balance, err := reader.Balance(from)
	if err != nil {
		return err
	}

	fee, err := reader.GetFeeForMessage(msg)
	if err != nil {
		return err
	}

	if balance < (amount + fee) {
		return fmt.Errorf("balance %d is too low for this transaction to be executed: amount %d + fee %d", balance, amount, fee)
	}
	return nil
}
