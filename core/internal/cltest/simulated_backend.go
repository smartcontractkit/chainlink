package cltest

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func NewSimulatedBackend(t *testing.T, alloc core.GenesisAlloc, gasLimit uint64) *backends.SimulatedBackend {
	backend := backends.NewSimulatedBackend(alloc, gasLimit)
	// NOTE: Make sure to finish closing any application/client before
	// backend.Close or they can hang
	t.Cleanup(func() {
		logger.TestLogger(t).ErrorIfClosing(backend, "simulated backend")
	})
	return backend
}

const SimulatedBackendEVMChainID int64 = 1337

// newIdentity returns a go-ethereum abstraction of an ethereum account for
// interacting with contract golang wrappers
func NewSimulatedBackendIdentity(t *testing.T) *bind.TransactOpts {
	key, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate ethereum identity")
	return MustNewSimulatedBackendKeyedTransactor(t, key)
}

func NewApplicationWithConfigAndKeyOnSimulatedBlockchain(
	t testing.TB,
	cfg *configtest.TestGeneralConfig,
	backend *backends.SimulatedBackend,
	flagsAndDeps ...interface{},
) *TestApplication {
	chainId := backend.Blockchain().Config().ChainID
	cfg.Overrides.DefaultChainID = chainId

	client := &SimulatedBackendClient{b: backend, t: t, chainId: chainId}
	eventBroadcaster := pg.NewEventBroadcaster(cfg.DatabaseURL(), 0, 0, logger.TestLogger(t), uuid.NewV4())

	zero := models.MustMakeDuration(0 * time.Millisecond)
	reaperThreshold := models.MustMakeDuration(100 * time.Millisecond)
	simulatedBackendChain := evmtypes.Chain{
		ID: *utils.NewBigI(SimulatedBackendEVMChainID),
		Cfg: evmtypes.ChainCfg{
			GasEstimatorMode:                 null.StringFrom("FixedPrice"),
			EvmHeadTrackerMaxBufferSize:      null.IntFrom(100),
			EvmHeadTrackerSamplingInterval:   &zero, // Head sampling disabled
			EthTxResendAfterThreshold:        &zero,
			EvmFinalityDepth:                 null.IntFrom(15),
			EthTxReaperThreshold:             &reaperThreshold,
			MinIncomingConfirmations:         null.IntFrom(1),
			MinRequiredOutgoingConfirmations: null.IntFrom(1),
			MinimumContractPayment:           assets.NewLinkFromJuels(100),
		},
		Enabled: true,
	}

	flagsAndDeps = append(flagsAndDeps, client, eventBroadcaster, simulatedBackendChain)

	//  app.Stop() will call client.Close on the simulated backend
	return NewApplicationWithConfigAndKey(t, cfg, flagsAndDeps...)
}

func MustNewSimulatedBackendKeyedTransactor(t *testing.T, key *ecdsa.PrivateKey) *bind.TransactOpts {
	t.Helper()
	return MustNewKeyedTransactor(t, key, SimulatedBackendEVMChainID)
}

func MustNewKeyedTransactor(t *testing.T, key *ecdsa.PrivateKey, chainID int64) *bind.TransactOpts {
	t.Helper()
	transactor, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(chainID))
	require.NoError(t, err)
	return transactor
}

// SimulatedBackendClient is an eth.Client implementation using a simulated
// blockchain backend. Note that not all RPC methods are implemented here.
type SimulatedBackendClient struct {
	b       *backends.SimulatedBackend
	t       testing.TB
	chainId *big.Int
}

var _ eth.Client = (*SimulatedBackendClient)(nil)

func (c *SimulatedBackendClient) Dial(context.Context) error {
	return nil
}

// Close does nothing. We ought not close the underlying backend here since
// other simulated clients might still be using it
func (c *SimulatedBackendClient) Close() {}

// checkEthCallArgs extracts and verifies the arguments for an eth_call RPC
func (c *SimulatedBackendClient) checkEthCallArgs(
	args []interface{}) (*eth.CallArgs, *big.Int, error) {
	if len(args) != 2 {
		return nil, nil, fmt.Errorf(
			"should have two arguments after \"eth_call\", got %d", len(args))
	}
	callArgs, ok := args[0].(eth.CallArgs)
	if !ok {
		return nil, nil, fmt.Errorf("third arg to SimulatedBackendClient.Call "+
			"must be an eth.CallArgs, got %+#v", args[0])
	}
	blockNumber, err := c.blockNumber(args[1])
	if err != nil || blockNumber.Cmp(c.currentBlockNumber()) != 0 {
		return nil, nil, fmt.Errorf("fourth arg to SimulatedBackendClient.Call "+
			"must be the string \"latest\", or a *big.Int equal to current "+
			"blocknumber, got %#+v", args[1])
	}
	return &callArgs, blockNumber, nil
}

// Call mocks the ethereum client RPC calls used by chainlink, copying the
// return value into result.
func (c *SimulatedBackendClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	switch method {
	case "eth_call":
		callArgs, _, err := c.checkEthCallArgs(args)
		if err != nil {
			return err
		}
		callMsg := ethereum.CallMsg{To: &callArgs.To, Data: callArgs.Data}
		b, err := c.b.CallContract(ctx, callMsg, nil /* always latest block */)
		if err != nil {
			return errors.Wrapf(err, "while calling contract at address %x with "+
				"data %x", callArgs.To, callArgs.Data)
		}
		switch r := result.(type) {
		case *hexutil.Bytes:
			*r = append(*r, b...)
			if !bytes.Equal(*r, b) {
				return fmt.Errorf("was passed a non-empty array, or failed to copy "+
					"answer. Expected %x = %x", *r, b)
			}
			return nil
		default:
			return fmt.Errorf("first arg to SimulatedBackendClient.Call is an "+
				"unrecognized type: %T; add processing logic for it here", result)
		}
	default:
		return fmt.Errorf("second arg to SimulatedBackendClient.Call is an RPC "+
			"API method which has not yet been implemented: %s. Add processing for "+
			"it here", method)
	}
}

// FilterLogs returns all logs that respect the passed filter query.
func (c *SimulatedBackendClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) (logs []types.Log, err error) {
	return c.b.FilterLogs(ctx, q)
}

// SubscribeToLogs registers a subscription for push notifications of logs
// from a given address.
func (c *SimulatedBackendClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, channel chan<- types.Log) (ethereum.Subscription, error) {
	return c.b.SubscribeFilterLogs(ctx, q, channel)
}

func (c *SimulatedBackendClient) GetEthBalance(ctx context.Context, account common.Address, blockNumber *big.Int) (*assets.Eth, error) {
	panic("not implemented")
}

// currentBlockNumber returns index of *pending* block in simulated blockchain
func (c *SimulatedBackendClient) currentBlockNumber() *big.Int {
	return c.b.Blockchain().CurrentBlock().Number()
}

var balanceOfABIString = `[
  {
    "constant": true,
    "inputs": [
      {
        "name": "_owner",
        "type": "address"
      }
    ],
    "name": "balanceOf",
    "outputs": [
      {
        "name": "balance",
        "type": "uint256"
      }
    ],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  }
]`

var balanceOfABI abi.ABI

func init() {
	var err error
	balanceOfABI, err = abi.JSON(strings.NewReader(balanceOfABIString))
	if err != nil {
		panic(errors.Wrapf(err, "while parsing erc20ABI"))
	}
}

// GetERC20Balance returns the balance of the given address for the token
// contract address.
func (c *SimulatedBackendClient) GetERC20Balance(address common.Address, contractAddress common.Address) (balance *big.Int, err error) {
	callData, err := balanceOfABI.Pack("balanceOf", address)
	if err != nil {
		return nil, errors.Wrapf(err, "while seeking the ERC20 balance of %s on %s",
			address, contractAddress)
	}
	b, err := c.b.CallContract(context.Background(), ethereum.CallMsg{
		To: &contractAddress, Data: callData},
		c.currentBlockNumber())
	if err != nil {
		return nil, errors.Wrapf(err, "while calling ERC20 balanceOf method on %s "+
			"for balance of %s", contractAddress, address)
	}
	err = balanceOfABI.UnpackIntoInterface(balance, "balanceOf", b)
	if err != nil {
		return nil, errors.New("unable to unpack balance")
	}
	return balance, nil
}

func (c *SimulatedBackendClient) GetLINKBalance(linkAddress common.Address, address common.Address) (*assets.Link, error) {
	panic("not implemented")
}

// TransactionReceipt returns the transaction receipt for the given transaction hash.
func (c *SimulatedBackendClient) TransactionReceipt(ctx context.Context, receipt common.Hash) (*types.Receipt, error) {
	return c.b.TransactionReceipt(ctx, receipt)
}

func (c *SimulatedBackendClient) blockNumber(number interface{}) (blockNumber *big.Int, err error) {
	switch n := number.(type) {
	case string:
		switch n {
		case "latest":
			return c.currentBlockNumber(), nil
		case "earliest":
			return big.NewInt(0), nil
		case "pending":
			panic("not implemented") // I don't understand the semantics of this.
			// return big.NewInt(0).Add(c.currentBlockNumber(), big.NewInt(1)), nil
		default:
			blockNumber, err = utils.HexToUint256(n)
			if err != nil {
				return nil, errors.Wrapf(err, "while parsing '%s' as hex-encoded"+
					"block number", n)
			}
			return blockNumber, nil
		}
	case *big.Int:
		if n.Sign() < 0 {
			return nil, fmt.Errorf("block number must be non-negative")
		}
		return n, nil
	}
	panic("can never reach here")
}

func (c *SimulatedBackendClient) HeadByNumber(ctx context.Context, n *big.Int) (*eth.Head, error) {
	if n == nil {
		n = c.currentBlockNumber()
	}
	header, err := c.b.HeaderByNumber(ctx, n)
	if err != nil {
		return nil, err
	} else if header == nil {
		return nil, ethereum.NotFound
	}
	return &eth.Head{
		EVMChainID: utils.NewBigI(SimulatedBackendEVMChainID),
		Hash:       header.Hash(),
		Number:     header.Number.Int64(),
		ParentHash: header.ParentHash,
	}, nil
}

func (c *SimulatedBackendClient) BlockByNumber(ctx context.Context, n *big.Int) (*types.Block, error) {
	return c.b.BlockByNumber(ctx, n)
}

// GetChainID returns the ethereum ChainID.
func (c *SimulatedBackendClient) ChainID() *big.Int {
	return c.chainId
}

func (c *SimulatedBackendClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return c.b.PendingNonceAt(ctx, account)
}

func (c *SimulatedBackendClient) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return c.b.NonceAt(ctx, account, blockNumber)
}

func (c *SimulatedBackendClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return c.b.BalanceAt(ctx, account, blockNumber)
}

type headSubscription struct {
	close        chan struct{}
	subscription ethereum.Subscription
}

var _ ethereum.Subscription = (*headSubscription)(nil)

func (h *headSubscription) Unsubscribe() {
	h.subscription.Unsubscribe()
	close(h.close)
}

func (h *headSubscription) Err() <-chan error { return h.subscription.Err() }

// SubscribeToNewHeads registers a subscription for push notifications of new
// blocks.
// Note the sim's API only accepts types.Head so we have this goroutine
// to convert those into eth.Head.
func (c *SimulatedBackendClient) SubscribeNewHead(
	ctx context.Context,
	channel chan<- *eth.Head,
) (ethereum.Subscription, error) {
	subscription := &headSubscription{close: make(chan struct{})}
	ch := make(chan *types.Header)
	go func() {
		var lastHead *eth.Head

		for {
			select {
			case h := <-ch:
				switch h {
				case nil:
					channel <- nil
				default:
					head := &eth.Head{Number: h.Number.Int64(), Hash: h.Hash(), ParentHash: h.ParentHash, Parent: lastHead}
					lastHead = head
					select {
					// In head tracker shutdown the heads reader is closed, so the channel <- head write
					// may hang.
					case channel <- head:
					case <-subscription.close:
						return
					}
				}
			case <-subscription.close:
				return
			}
		}
	}()
	var err error
	subscription.subscription, err = c.b.SubscribeNewHead(ctx, ch)
	if err != nil {
		return nil, errors.Wrapf(err, "could not subscribe to new heads on "+
			"simulated backend")
	}
	return subscription, err
}

func (c *SimulatedBackendClient) HeaderByNumber(ctx context.Context, n *big.Int) (*types.Header, error) {
	return c.b.HeaderByNumber(ctx, n)
}

func (c *SimulatedBackendClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	sender, err := types.Sender(types.NewLondonSigner(c.chainId), tx)
	if err != nil {
		logger.TestLogger(c.t).Panic(fmt.Errorf("invalid transaction: %v (tx: %#v)", err, tx))
	}
	pendingNonce, err := c.b.PendingNonceAt(ctx, sender)
	if err != nil {
		panic(fmt.Errorf("unable to determine nonce for account %s: %v", sender.Hex(), err))
	}
	// the simulated backend does not gracefully handle tx rebroadcasts (gas bumping) so just
	// ignore the situation where nonces are reused
	// github.com/ethereum/go-ethereum/blob/fb2c79df1995b4e8dfe79f9c75464d29d23aaaf4/accounts/abi/bind/backends/simulated.go#L556
	if tx.Nonce() < pendingNonce {
		return nil
	}

	err = c.b.SendTransaction(ctx, tx)
	return err
}

func (c *SimulatedBackendClient) Call(result interface{}, method string, args ...interface{}) error {
	return c.CallContext(context.Background(), result, method, args)
}

func (c *SimulatedBackendClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return c.b.CallContract(ctx, msg, blockNumber)
}

func (c *SimulatedBackendClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return c.b.CodeAt(ctx, account, blockNumber)
}

func (c *SimulatedBackendClient) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return c.b.PendingCodeAt(ctx, account)
}

func (c *SimulatedBackendClient) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	return c.b.EstimateGas(ctx, call)
}

func (c *SimulatedBackendClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	panic("unimplemented")
}

func (c *SimulatedBackendClient) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	for i, elem := range b {
		if elem.Method != "eth_getTransactionReceipt" || len(elem.Args) != 1 {
			return errors.New("SimulatedBackendClient BatchCallContext only supports eth_getTransactionReceipt")
		}
		switch v := elem.Result.(type) {
		case *bulletprooftxmanager.Receipt:
			hash, is := elem.Args[0].(common.Hash)
			if !is {
				return errors.Errorf("SimulatedBackendClient expected arg to be a hash, got: %T", elem.Args[0])
			}
			receipt, err := c.b.TransactionReceipt(ctx, hash)
			b[i].Result = bulletprooftxmanager.FromGethReceipt(receipt)
			b[i].Error = err
		default:
			return errors.Errorf("SimulatedBackendClient unsupported elem.Result type %T", v)
		}
	}
	return nil
}

func (c *SimulatedBackendClient) SuggestGasTipCap(ctx context.Context) (tipCap *big.Int, err error) {
	return nil, nil
}

// Mine forces the simulated backend to produce a new block every 2 seconds
func Mine(backend *backends.SimulatedBackend, blockTime time.Duration) (stopMining func()) {
	timer := time.NewTicker(blockTime)
	chStop := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case <-timer.C:
				backend.Commit()
			case <-chStop:
				wg.Done()
				return
			}
		}
	}()
	return func() { close(chStop); timer.Stop(); wg.Wait() }
}
