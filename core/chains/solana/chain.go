package solana

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"sync"
	"time"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/txm"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/chains/internal"
	"github.com/smartcontractkit/chainlink/v2/core/chains/solana/monitor"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// DefaultRequestTimeout is the default Solana client timeout.
const DefaultRequestTimeout = 30 * time.Second

// ChainOpts holds options for configuring a Chain.
type ChainOpts struct {
	Logger   logger.Logger
	KeyStore loop.Keystore
}

func (o *ChainOpts) Validate() (err error) {
	required := func(s string) error {
		return errors.Errorf("%s is required", s)
	}
	if o.Logger == nil {
		err = multierr.Append(err, required("Logger"))
	}
	if o.KeyStore == nil {
		err = multierr.Append(err, required("KeyStore"))
	}
	return
}

func (o *ChainOpts) GetLogger() logger.Logger {
	return o.Logger
}

func NewChain(cfg *SolanaConfig, opts ChainOpts) (solana.Chain, error) {
	if !cfg.IsEnabled() {
		return nil, fmt.Errorf("cannot create new chain with ID %s: %w", *cfg.ChainID, chains.ErrChainDisabled)
	}
	c, err := newChain(*cfg.ChainID, cfg, opts.KeyStore, opts.Logger)
	if err != nil {
		return nil, err
	}
	return c, nil
}

var _ solana.Chain = (*chain)(nil)

type chain struct {
	utils.StartStopOnce
	id             string
	cfg            *SolanaConfig
	txm            *txm.Txm
	balanceMonitor services.ServiceCtx
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
		return v.chainIDVerified, errors.Wrap(err, "failed to fetch ChainID in verifiedCachedClient")
	}

	// check chainID matches expected chainID
	expectedChainID := strings.ToLower(v.expectedChainID)
	if v.chainID != expectedChainID {
		v.chainIDVerified = false
		return v.chainIDVerified, errors.Errorf("client returned mismatched chain id (expected: %s, got: %s): %s", expectedChainID, v.chainID, v.nodeURL)
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

func newChain(id string, cfg *SolanaConfig, ks loop.Keystore, lggr logger.Logger) (*chain, error) {
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
	ch.balanceMonitor = monitor.NewBalanceMonitor(ch.id, cfg, lggr, ks, ch.Reader)
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
	return internal.ListNodeStatuses(int(pageSize), pageToken, c.listNodeStatuses)
}

func (c *chain) Transact(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	return c.sendTx(ctx, from, to, amount, balanceCheck)
}

func (c *chain) listNodeStatuses(start, end int) ([]relaytypes.NodeStatus, int, error) {
	stats := make([]relaytypes.NodeStatus, 0)
	total := len(c.cfg.Nodes)
	if start >= total {
		return stats, total, internal.ErrOutOfRange
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

func (c *chain) Name() string {
	return c.lggr.Name()
}

func (c *chain) ID() string {
	return c.id
}

func (c *chain) Config() config.Config {
	return c.cfg
}

func (c *chain) TxManager() solana.TxManager {
	return c.txm
}

func (c *chain) Reader() (client.Reader, error) {
	return c.getClient()
}

func (c *chain) ChainID() relay.ChainID {
	return relay.ChainID(c.id)
}

// getClient returns a client, randomly selecting one from available and valid nodes
func (c *chain) getClient() (client.ReaderWriter, error) {
	var node db.Node
	var client client.ReaderWriter
	nodes, err := c.cfg.ListNodes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get nodes")
	}
	if len(nodes) == 0 {
		return nil, errors.New("no nodes available")
	}
	// #nosec
	index := rand.Perm(len(nodes)) // list of node indexes to try
	for _, i := range index {
		node = nodes[i]
		// create client and check
		client, err = c.verifiedClient(node)
		// if error, try another node
		if err != nil {
			c.lggr.Warnw("failed to create node", "name", node.Name, "solana-url", node.SolanaURL, "err", err.Error())
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

// verifiedClient returns a client for node or an error if fails to create the client.
// The client will still be returned if the nodes are not valid, or the chain id doesn't match.
// Further client calls will try and verify the client, and fail if the client is still not valid.
func (c *chain) verifiedClient(node db.Node) (client.ReaderWriter, error) {
	url := node.SolanaURL
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
		cl.ReaderWriter, err = client.NewClient(url, c.cfg, DefaultRequestTimeout, logger.Named(c.lggr, "Client."+node.Name))
		if err != nil {
			return nil, errors.Wrap(err, "failed to create client")
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
	return multierr.Combine(
		c.StartStopOnce.Ready(),
		c.txm.Ready(),
	)
}

func (c *chain) HealthReport() map[string]error {
	report := map[string]error{c.Name(): c.StartStopOnce.Healthy()}
	maps.Copy(report, c.txm.HealthReport())
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
