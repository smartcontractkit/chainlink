package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/hex"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func init() {
	var err error

	balanceOfABI, err = abi.JSON(strings.NewReader(balanceOfABIString))
	if err != nil {
		panic(fmt.Errorf("%w: while parsing erc20ABI", err))
	}
}

var (
	balanceOfABIString = `[
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

	balanceOfABI abi.ABI
)

// SimulatedBackendClient is an Client implementation using a simulated
// blockchain backend. Note that not all RPC methods are implemented here.
type SimulatedBackendClient struct {
	b                    evmtypes.Backend // *simulated.Backend, or something satisfying same interface
	client               simulated.Client
	t                    testing.TB
	chainID              *big.Int
	chainType            chaintype.ChainType
	headByNumberCallback func(ctx context.Context, c *SimulatedBackendClient, n *big.Int) error
}

// NewSimulatedBackendClient creates an eth client backed by a simulated backend.
func NewSimulatedBackendClient(t testing.TB, b *simulated.Backend, chainID *big.Int) *SimulatedBackendClient {
	return &SimulatedBackendClient{
		b:       b,
		client:  b.Client(),
		t:       t,
		chainID: chainID,
	}
}

// Switch to a new backend client (simulating an rpc failover event)
// If chainFamily = Optimism, the new backend will exhibit the non-geth behavior of optimism (and some other rpc clients),
// where success rather than an error code is returned when a call to FilterLogs() fails to find the block hash
// requested. This combined with a failover event can lead to the "eventual consistency" behavior that Backup LogPoller
// and other solutions were designed to recover from.
func (c *SimulatedBackendClient) SetBackend(backend evmtypes.Backend, chainType chaintype.ChainType) {
	c.chainType = chainType
	c.b = backend
	c.client = backend.Client()
}

// Dial noop for the sim.
func (c *SimulatedBackendClient) Dial(context.Context) error {
	return nil
}

// Close does nothing. We ought not close the underlying backend here since
// other simulated clients might still be using it
func (c *SimulatedBackendClient) Close() {}

// CallContext mocks the ethereum client RPC calls used by chainlink, copying the
// return value into result.
// The simulated client avoids the old block error from the simulated backend by
// passing `nil` to `CallContract` when calling `CallContext` or `BatchCallContext`
// and will not return an error when an old block is used.
func (c *SimulatedBackendClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	switch method {
	case "eth_getTransactionReceipt":
		return c.ethGetTransactionReceipt(ctx, result, args...)
	case "eth_getBlockByNumber":
		return c.ethGetBlockByNumber(ctx, result, args...)
	case "eth_call":
		return c.ethCall(ctx, result, args...)
	case "eth_getHeaderByNumber":
		return c.ethGetHeaderByNumber(ctx, result, args...)
	case "eth_estimateGas":
		return c.ethEstimateGas(ctx, result, args...)
	default:
		return fmt.Errorf("second arg to SimulatedBackendClient.Call is an RPC API method which has not yet been implemented: %s. Add processing for it here", method)
	}
}

// FilterLogs returns all logs that respect the passed filter query.
func (c *SimulatedBackendClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) (logs []types.Log, err error) {
	logs, err = c.client.FilterLogs(ctx, q)
	if c.chainType == chaintype.ChainOptimismBedrock {
		if err != nil && err.Error() == "unknown block" {
			return []types.Log{}, nil // emulate optimism behavior of returning success instead of "unknown block"
		}
	}

	return logs, err
}

// SubscribeFilterLogs registers a subscription for push notifications of logs
// from a given address.
func (c *SimulatedBackendClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, channel chan<- types.Log) (ethereum.Subscription, error) {
	return c.client.SubscribeFilterLogs(ctx, q, channel)
}

func (c *SimulatedBackendClient) TokenBalance(ctx context.Context, address common.Address, contractAddress common.Address) (balance *big.Int, err error) {
	callData, err := balanceOfABI.Pack("balanceOf", address)
	if err != nil {
		return nil, fmt.Errorf("%w: while seeking the ERC20 balance of %s on %s", err,
			address, contractAddress)
	}
	b, err := c.client.CallContract(ctx, ethereum.CallMsg{
		To: &contractAddress, Data: callData},
		big.NewInt(int64(rpc.LatestBlockNumber)))
	if err != nil {
		return nil, fmt.Errorf("%w: while calling ERC20 balanceOf method on %s "+
			"for balance of %s", err, contractAddress, address)
	}
	err = balanceOfABI.UnpackIntoInterface(balance, "balanceOf", b)
	if err != nil {
		return nil, fmt.Errorf("unable to unpack balance")
	}
	return balance, nil
}

// GetLINKBalance get link balance.
func (c *SimulatedBackendClient) LINKBalance(ctx context.Context, address common.Address, linkAddress common.Address) (*assets.Link, error) {
	panic("not implemented")
}

func (c *SimulatedBackendClient) FeeHistory(ctx context.Context, blockCount uint64, rewardPercentiles []float64) (feeHistory *ethereum.FeeHistory, err error) {
	panic("not implemented")
}

// TransactionReceipt returns the transaction receipt for the given transaction hash.
func (c *SimulatedBackendClient) TransactionReceipt(ctx context.Context, receipt common.Hash) (*types.Receipt, error) {
	return c.client.TransactionReceipt(ctx, receipt)
}

func (c *SimulatedBackendClient) TransactionByHash(ctx context.Context, txHash common.Hash) (tx *types.Transaction, err error) {
	tx, _, err = c.client.TransactionByHash(ctx, txHash)
	return
}

func (c *SimulatedBackendClient) blockNumber(ctx context.Context, number interface{}) (blockNumber *big.Int, err error) {
	switch n := number.(type) {
	case string:
		switch n {
		case "latest":
			var n uint64
			n, err = c.client.BlockNumber(ctx)
			if err != nil {
				return
			}
			blockNumber = new(big.Int)
			blockNumber.SetUint64(n)
			return
		case "earliest":
			return big.NewInt(0), nil
		case "pending":
			var h *types.Header
			h, err = c.client.HeaderByNumber(ctx, new(big.Int).SetInt64(rpc.PendingBlockNumber.Int64()))
			if err != nil {
				return
			}
			blockNumber = h.Number
			return
		case "finalized":
			var h *types.Header
			h, err = c.client.HeaderByNumber(ctx, new(big.Int).SetInt64(rpc.FinalizedBlockNumber.Int64()))
			if err != nil {
				return
			}
			blockNumber = h.Number
			return
		default:
			blockNumber, err := hexutil.DecodeBig(n)
			if err != nil {
				return nil, fmt.Errorf("%w: while parsing '%s' as hex-encoded block number", err, n)
			}
			return blockNumber, nil
		}
	case *big.Int:
		if n == nil {
			return nil, nil
		}
		if n.Sign() < 0 {
			return nil, fmt.Errorf("block number must be non-negative")
		}
		return n, nil
	default:
		return nil, fmt.Errorf("invalid type %T for block number, must be string or *big.Int", n)
	}
}

func (c *SimulatedBackendClient) RegisterHeadByNumberCallback(cb func(ctx context.Context, c *SimulatedBackendClient, n *big.Int) error) {
	c.headByNumberCallback = cb
}

// HeadByNumber returns our own header type.
func (c *SimulatedBackendClient) HeadByNumber(ctx context.Context, n *big.Int) (*evmtypes.Head, error) {
	if n == nil {
		n = big.NewInt(int64(rpc.LatestBlockNumber))
	}
	header, err := c.client.HeaderByNumber(ctx, n)
	if err != nil {
		return nil, err
	} else if header == nil {
		return nil, ethereum.NotFound
	}

	if c.headByNumberCallback != nil {
		err = c.headByNumberCallback(ctx, c, n)
		if err != nil {
			return nil, err
		}
	}

	head := &evmtypes.Head{EVMChainID: ubig.New(c.chainID)}
	head.SetFromHeader(header)
	return head, nil
}

// HeadByHash returns our own header type.
func (c *SimulatedBackendClient) HeadByHash(ctx context.Context, h common.Hash) (*evmtypes.Head, error) {
	header, err := c.client.HeaderByHash(ctx, h)
	if err != nil {
		return nil, err
	} else if header == nil {
		return nil, ethereum.NotFound
	}
	head := &evmtypes.Head{EVMChainID: ubig.NewI(c.chainID.Int64())}
	head.SetFromHeader(header)
	return head, nil
}

// BlockByNumber returns a geth block type.
func (c *SimulatedBackendClient) BlockByNumber(ctx context.Context, n *big.Int) (*types.Block, error) {
	return c.client.BlockByNumber(ctx, n)
}

// BlockByNumber returns a geth block type.
func (c *SimulatedBackendClient) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return c.client.BlockByHash(ctx, hash)
}

func (c *SimulatedBackendClient) LatestBlockHeight(ctx context.Context) (*big.Int, error) {
	header, err := c.client.HeaderByNumber(ctx, nil)
	return header.Number, err
}

// ChainID returns the ethereum ChainID.
func (c *SimulatedBackendClient) ConfiguredChainID() *big.Int {
	return c.chainID
}

// ChainID RPC call
func (c *SimulatedBackendClient) ChainID() (*big.Int, error) {
	panic("not implemented")
}

// PendingNonceAt gets pending nonce i.e. mempool nonce.
func (c *SimulatedBackendClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return c.client.PendingNonceAt(ctx, account)
}

// NonceAt gets nonce as of a specified block.
func (c *SimulatedBackendClient) SequenceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (evmtypes.Nonce, error) {
	nonce, err := c.client.NonceAt(ctx, account, blockNumber)
	return evmtypes.Nonce(nonce), err
}

// BalanceAt gets balance as of a specified block.
func (c *SimulatedBackendClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return c.client.BalanceAt(ctx, account, blockNumber)
}

type headSubscription struct {
	unSub        chan chan struct{}
	subscription ethereum.Subscription
}

var _ ethereum.Subscription = (*headSubscription)(nil)

func (h *headSubscription) Unsubscribe() {
	done := make(chan struct{})
	h.unSub <- done
	<-done
}

// Err returns err channel
func (h *headSubscription) Err() <-chan error { return h.subscription.Err() }

// SubscribeNewHead registers a subscription for push notifications of new blocks.
// Note the sim's API only accepts types.Head so we have this goroutine
// to convert those into evmtypes.Head.
func (c *SimulatedBackendClient) SubscribeNewHead(
	ctx context.Context,
	channel chan<- *evmtypes.Head,
) (ethereum.Subscription, error) {
	subscription := &headSubscription{unSub: make(chan chan struct{})}
	ch := make(chan *types.Header)

	var err error
	subscription.subscription, err = c.client.SubscribeNewHead(ctx, ch)
	if err != nil {
		return nil, fmt.Errorf("%w: could not subscribe to new heads on "+
			"simulated backend", err)
	}
	go func() {
		var lastHead *evmtypes.Head
		for {
			select {
			case h := <-ch:
				var head *evmtypes.Head
				if h != nil {
					head = &evmtypes.Head{
						Difficulty: h.Difficulty,
						Timestamp:  time.Unix(int64(h.Time), 0), //nolint:gosec
						Number:     h.Number.Int64(),
						Hash:       h.Hash(),
						ParentHash: h.ParentHash,
						EVMChainID: ubig.New(c.chainID),
					}
					head.Parent.Store(lastHead)
					lastHead = head
				}
				select {
				case channel <- head:
				case done := <-subscription.unSub:
					subscription.subscription.Unsubscribe()
					close(done)
					return
				}

			case done := <-subscription.unSub:
				subscription.subscription.Unsubscribe()
				close(done)
				return
			}
		}
	}()
	return subscription, err
}

// HeaderByNumber returns the geth header type.
func (c *SimulatedBackendClient) HeaderByNumber(ctx context.Context, n *big.Int) (*types.Header, error) {
	return c.client.HeaderByNumber(ctx, n)
}

func (c *SimulatedBackendClient) HeaderByHash(ctx context.Context, h common.Hash) (*types.Header, error) {
	return c.client.HeaderByHash(ctx, h)
}

func (c *SimulatedBackendClient) SendTransactionReturnCode(ctx context.Context, tx *types.Transaction, fromAddress common.Address) (commonclient.SendTxReturnCode, error) {
	err := c.SendTransaction(ctx, tx)
	if err == nil {
		return commonclient.Successful, nil
	}
	if strings.Contains(err.Error(), "could not fetch parent") || strings.Contains(err.Error(), "invalid transaction") {
		return commonclient.Fatal, err
	}
	// All remaining error messages returned from SendTransaction are considered Unknown.
	return commonclient.Unknown, err
}

// SendTransaction sends a transaction.
func (c *SimulatedBackendClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	var (
		sender common.Address
		err    error
	)
	// try to recover the sender from the transaction using the configured chain id
	// first. if that fails, try again with the simulated chain id (1337)
	sender, err = types.Sender(types.NewLondonSigner(c.chainID), tx)
	if err != nil {
		sender, err = types.Sender(types.NewLondonSigner(big.NewInt(1337)), tx)
		if err != nil {
			logger.Test(c.t).Panic(fmt.Errorf("invalid transaction: %v (tx: %#v)", err, tx))
		}
	}
	pendingNonce, err := c.client.PendingNonceAt(ctx, sender)
	if err != nil {
		panic(fmt.Errorf("unable to determine nonce for account %s: %v", sender.Hex(), err))
	}
	// the simulated backend does not gracefully handle tx rebroadcasts (gas bumping) so just
	// ignore the situation where nonces are reused
	// github.com/ethereum/go-ethereum/blob/fb2c79df1995b4e8dfe79f9c75464d29d23aaaf4/accounts/abi/bind/backends/simulated.go#L556
	if tx.Nonce() < pendingNonce {
		return nil
	}

	err = c.client.SendTransaction(ctx, tx)
	return err
}

type revertError struct {
	error
	reason string
}

func (e *revertError) ErrorCode() int {
	return 3
}

// ErrorData returns the hex encoded revert reason.
func (e *revertError) ErrorData() interface{} {
	return e.reason
}

var _ rpc.DataError = &revertError{}

// CallContract calls a contract.
func (c *SimulatedBackendClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	// Expected error is
	// type JsonError struct {
	//	Code    int         `json:"code"`
	//	Message string      `json:"message"`
	//	Data    interface{} `json:"data,omitempty"`
	//}
	res, err := c.client.CallContract(ctx, msg, blockNumber)
	if err != nil {
		dataErr := revertError{}
		if errors.As(err, &dataErr) {
			return nil, &JsonError{Data: dataErr.ErrorData(), Message: dataErr.Error(), Code: 3}
		}
		// Generic revert, no data
		return nil, &JsonError{Data: []byte{}, Message: err.Error(), Code: 3}
	}
	return res, nil
}

func (c *SimulatedBackendClient) PendingCallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
	// Expected error is
	// type JsonError struct {
	//	Code    int         `json:"code"`
	//	Message string      `json:"message"`
	//	Data    interface{} `json:"data,omitempty"`
	//}
	res, err := c.client.PendingCallContract(ctx, msg)
	if err != nil {
		dataErr := revertError{}
		if errors.As(err, &dataErr) {
			return nil, &JsonError{Data: dataErr.ErrorData(), Message: dataErr.Error(), Code: 3}
		}
		// Generic revert, no data
		return nil, &JsonError{Data: []byte{}, Message: err.Error(), Code: 3}
	}
	return res, nil
}

// CodeAt gets the code associated with an account as of a specified block.
func (c *SimulatedBackendClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return c.client.CodeAt(ctx, account, blockNumber)
}

// PendingCodeAt gets the latest code.
func (c *SimulatedBackendClient) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return c.client.PendingCodeAt(ctx, account)
}

// EstimateGas estimates gas for a msg.
func (c *SimulatedBackendClient) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	return c.client.EstimateGas(ctx, call)
}

// SuggestGasPrice recommends a gas price.
func (c *SimulatedBackendClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return c.client.SuggestGasPrice(ctx)
}

// BatchCallContext makes a batch rpc call.
// The simulated client avoids the old block error from the simulated backend by
// passing `nil` to `CallContract` when calling `CallContext` or `BatchCallContext`
// and will not return an error when an old block is used.
func (c *SimulatedBackendClient) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	select {
	case <-ctx.Done():
		return errors.New("context canceled")
	default:
		//do nothing
	}

	for i, elem := range b {
		var method func(context.Context, interface{}, ...interface{}) error
		switch elem.Method {
		case "eth_getTransactionReceipt":
			method = c.ethGetTransactionReceipt
		case "eth_getBlockByNumber":
			method = c.ethGetBlockByNumber
		case "eth_call":
			method = c.ethCall
		case "eth_getHeaderByNumber":
			method = c.ethGetHeaderByNumber
		case "eth_estimateGas":
			method = c.ethEstimateGas
		case "eth_getLogs":
			method = c.ethGetLogs
		default:
			return fmt.Errorf("SimulatedBackendClient got unsupported method %s", elem.Method)
		}
		b[i].Error = method(ctx, b[i].Result, b[i].Args...)
	}

	return nil
}

// BatchCallContextAll makes a batch rpc call.
func (c *SimulatedBackendClient) BatchCallContextAll(ctx context.Context, b []rpc.BatchElem) error {
	return c.BatchCallContext(ctx, b)
}

// SuggestGasTipCap suggests a gas tip cap.
func (c *SimulatedBackendClient) SuggestGasTipCap(ctx context.Context) (tipCap *big.Int, err error) {
	return c.client.SuggestGasTipCap(ctx)
}

func (c *SimulatedBackendClient) Backend() evmtypes.Backend {
	return c.b
}

// NodeStates implements evmclient.Client
func (c *SimulatedBackendClient) NodeStates() map[string]string { return nil }

// Commit imports all the pending transactions as a single block and starts a
// fresh new state.
func (c *SimulatedBackendClient) Commit() common.Hash {
	return c.b.Commit()
}

func (c *SimulatedBackendClient) IsL2() bool {
	return false
}

func (c *SimulatedBackendClient) fetchHeader(ctx context.Context, blockNumOrTag string) (*types.Header, error) {
	switch blockNumOrTag {
	case rpc.SafeBlockNumber.String():
		return c.client.HeaderByNumber(ctx, big.NewInt(int64(rpc.SafeBlockNumber)))
	case rpc.LatestBlockNumber.String():
		return c.client.HeaderByNumber(ctx, big.NewInt(int64(rpc.LatestBlockNumber)))
	case rpc.FinalizedBlockNumber.String():
		return c.client.HeaderByNumber(ctx, big.NewInt(int64(rpc.FinalizedBlockNumber)))
	default:
		blockNum, ok := new(big.Int).SetString(blockNumOrTag, 0)
		if !ok {
			return nil, fmt.Errorf("error while converting block number string: %s to big.Int ", blockNumOrTag)
		}
		return c.client.HeaderByNumber(ctx, blockNum)
	}
}

func (c *SimulatedBackendClient) ethGetTransactionReceipt(ctx context.Context, result interface{}, args ...interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("SimulatedBackendClient expected 1 arg, got %d for eth_getTransactionReceipt", len(args))
	}

	hash, is := args[0].(common.Hash)
	if !is {
		return fmt.Errorf("SimulatedBackendClient expected arg to be a hash, got: %T", args[0])
	}

	receipt, err := c.client.TransactionReceipt(ctx, hash)
	if err != nil {
		return err
	}

	// strongly typing the result here has the consequence of not being flexible in
	// custom types where a real-world RPC client would allow for custom types with
	// custom marshalling.
	switch typed := result.(type) {
	case *types.Receipt:
		*typed = *receipt
	case *evmtypes.Receipt:
		*typed = *evmtypes.FromGethReceipt(receipt)
	default:
		return fmt.Errorf("SimulatedBackendClient expected return type of *evmtypes.Receipt for eth_getTransactionReceipt, got type %T", result)
	}

	return nil
}

func (c *SimulatedBackendClient) ethGetBlockByNumber(ctx context.Context, result interface{}, args ...interface{}) error {
	if len(args) != 2 {
		return fmt.Errorf("SimulatedBackendClient expected 2 args, got %d for eth_getBlockByNumber", len(args))
	}

	blockNumOrTag, is := args[0].(string)
	if !is {
		return fmt.Errorf("SimulatedBackendClient expected first arg to be a string for eth_getBlockByNumber, got: %T", args[0])
	}

	_, is = args[1].(bool)
	if !is {
		return fmt.Errorf("SimulatedBackendClient expected second arg to be a boolean for eth_getBlockByNumber, got: %T", args[1])
	}

	header, err := c.fetchHeader(ctx, blockNumOrTag)
	if err != nil {
		return err
	}

	switch res := result.(type) {
	case *evmtypes.Head:
		res.SetFromHeader(header)
	case *evmtypes.Block:
		res.Number = header.Number.Int64()
		res.Hash = header.Hash()
		res.ParentHash = header.ParentHash
		res.Timestamp = time.Unix(int64(header.Time), 0).UTC()
	default:
		return fmt.Errorf("SimulatedBackendClient Unexpected Type %T", res)
	}

	return nil
}
func (c *SimulatedBackendClient) ethEstimateGas(ctx context.Context, result interface{}, args ...interface{}) error {
	if len(args) != 2 {
		return fmt.Errorf("SimulatedBackendClient expected 2 args, got %d for eth_estimateGas", len(args))
	}

	params, ok := args[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("SimulatedBackendClient expected first arg to be map[string]interface{} for eth_call, got: %T", args[0])
	}

	_, err := c.blockNumber(ctx, args[1])
	if err != nil {
		return fmt.Errorf("SimulatedBackendClient expected second arg to be the string 'latest' or a *big.Int for eth_call, got: %T", args[1])
	}

	resp, err := c.client.EstimateGas(ctx, toCallMsg(params))
	if err != nil {
		return err
	}

	switch typedResult := result.(type) {
	case *uint64:
		*typedResult = resp
	case *hexutil.Uint64:
		*typedResult = hexutil.Uint64(resp)
	default:
		return fmt.Errorf("SimulatedBackendClient unexpected type %T", result)
	}

	return nil
}

func (c *SimulatedBackendClient) ethCall(ctx context.Context, result interface{}, args ...interface{}) error {
	if len(args) != 2 {
		return fmt.Errorf("SimulatedBackendClient expected 2 args, got %d for eth_call", len(args))
	}

	params, ok := args[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("SimulatedBackendClient expected first arg to be map[string]interface{} for eth_call, got: %T", args[0])
	}

	if _, err := c.blockNumber(ctx, args[1]); err != nil {
		return fmt.Errorf("SimulatedBackendClient expected second arg to be the string 'latest' or a *big.Int for eth_call, got: %T", args[1])
	}

	resp, err := c.client.CallContract(ctx, toCallMsg(params), nil /* always latest block on simulated backend */)
	if err != nil {
		return err
	}

	switch typedResult := result.(type) {
	case *hexutil.Bytes:
		*typedResult = append(*typedResult, resp...)

		if !bytes.Equal(*typedResult, resp) {
			return fmt.Errorf("SimulatedBackendClient was passed a non-empty array, or failed to copy answer. Expected %x = %x", *typedResult, resp)
		}
	case *string:
		*typedResult = hexutil.Encode(resp)
	default:
		return fmt.Errorf("SimulatedBackendClient unexpected type %T", result)
	}

	return nil
}

func (c *SimulatedBackendClient) ethGetHeaderByNumber(ctx context.Context, result interface{}, args ...interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("SimulatedBackendClient expected 1 arg, got %d for eth_getHeaderByNumber", len(args))
	}

	blockNumber, err := c.blockNumber(ctx, args[0])
	if err != nil {
		return fmt.Errorf("SimulatedBackendClient expected first arg to be a string for eth_getHeaderByNumber: %w", err)
	}

	header, err := c.client.HeaderByNumber(ctx, blockNumber)
	if err != nil {
		return err
	}

	switch typedResult := result.(type) {
	case *types.Header:
		*typedResult = *header
	default:
		return fmt.Errorf("SimulatedBackendClient unexpected Type %T", typedResult)
	}

	return nil
}

func (c *SimulatedBackendClient) LatestFinalizedBlock(ctx context.Context) (*evmtypes.Head, error) {
	h, err := c.client.HeaderByNumber(ctx, big.NewInt(rpc.FinalizedBlockNumber.Int64()))
	if err != nil {
		return nil, err
	}
	head := &evmtypes.Head{EVMChainID: ubig.New(c.chainID)}
	head.SetFromHeader(h)
	return head, nil
}

func (c *SimulatedBackendClient) ethGetLogs(ctx context.Context, result interface{}, args ...interface{}) error {
	var from, to *big.Int
	var hash *common.Hash
	var err error
	var addresses []common.Address
	var topics [][]common.Hash

	params := args[0].(map[string]interface{})
	if blockHash, ok := params["blockHash"]; ok {
		hash, err = interfaceToHash(blockHash)
		if err != nil {
			return fmt.Errorf("SimultaedBackendClient received unexpected 'blockhash' param: %w", err)
		}
	}

	if fromBlock, ok := params["fromBlock"]; ok {
		from, err = c.blockNumber(ctx, fromBlock)
		if err != nil {
			return fmt.Errorf("SimulatedBackendClient expected 'fromBlock' to be a string: %w", err)
		}
	}

	if toBlock, ok := params["toBlock"]; ok {
		to, err = c.blockNumber(ctx, toBlock)
		if err != nil {
			return fmt.Errorf("SimulatedBackendClient expected 'toBlock' to be a string: %w", err)
		}
	}

	if a, ok := params["addresses"]; ok {
		addresses = a.([]common.Address)
	}

	if t, ok := params["topics"]; ok {
		tt := t.([][]common.Hash)
		lastTopic := len(tt) - 1
		for lastTopic >= 0 {
			if tt[lastTopic] != nil {
				break
			}
			lastTopic--
		}
		// lastTopic is the topic index of the last non-nil topic slice
		//  We have to drop any nil values in the topics slice after that due to a quirk in FilterLogs(),
		//  which will only use nil as a wildcard if there are non-nil values after it in the slice
		for i := 0; i < lastTopic; i++ {
			topics = append(topics, tt[i])
		}
	}

	query := ethereum.FilterQuery{
		BlockHash: hash,
		FromBlock: from,
		ToBlock:   to,
		Addresses: addresses,
		Topics:    topics,
	}
	logs, err := c.FilterLogs(ctx, query)
	if err != nil {
		return err
	}
	switch r := result.(type) {
	case *[]types.Log:
		*r = logs
		return nil
	default:
		return fmt.Errorf("SimulatedBackendClient unexpected Type %T", r)
	}
}

func (c *SimulatedBackendClient) CheckTxValidity(ctx context.Context, from common.Address, to common.Address, data []byte) *SendError {
	return nil
}

func toCallMsg(params map[string]interface{}) ethereum.CallMsg {
	var callMsg ethereum.CallMsg
	toAddr, err := interfaceToAddress(params["to"])
	if err != nil {
		panic(fmt.Errorf("unexpected 'to' parameter: %s", err))
	}

	callMsg.To = &toAddr

	// from is optional in the standard client; default to 0x when missing
	if value, ok := params["from"]; ok {
		addr, err := interfaceToAddress(value)
		if err != nil {
			panic(fmt.Errorf("unexpected 'from' parameter: %s", err))
		}

		callMsg.From = addr
	} else {
		callMsg.From = common.HexToAddress("0x")
	}

	if params["data"] != nil && params["input"] != nil {
		panic("cannot have both 'data' and 'input' parameters")
	}

	switch data := params["data"].(type) {
	case nil:
		// This parameter is not required so nil is acceptable
	case hexutil.Bytes:
		callMsg.Data = data
	case []byte:
		callMsg.Data = data
	default:
		panic("unexpected type of 'data' parameter; try hexutil.Bytes, []byte, or nil")
	}

	switch input := params["input"].(type) {
	case nil:
		// This parameter is not required so nil is acceptable
	case hexutil.Bytes:
		callMsg.Data = input
	case []byte:
		callMsg.Data = input
	default:
		panic("unexpected type of 'input' parameter; try hexutil.Bytes, []byte, or nil")
	}

	if value, ok := params["value"].(*big.Int); ok {
		callMsg.Value = value
	}

	switch gas := params["gas"].(type) {
	case nil:
		// This parameter is not required so nil is acceptable
	case uint64:
		callMsg.Gas = gas
	case hexutil.Uint64:
		callMsg.Gas = uint64(gas)
	default:
		panic("unexpected type of 'gas' parameter; try hexutil.Uint64, or uint64")
	}

	switch gasPrice := params["gasPrice"].(type) {
	case nil:
		// This parameter is not required so nil is acceptable
	case *big.Int:
		callMsg.GasPrice = gasPrice
	case *hexutil.Big:
		callMsg.GasPrice = gasPrice.ToInt()
	default:
		panic("unexpected type of 'gasPrice' parameter; try *big.Int, or *hexutil.Big")
	}

	return callMsg
}

func interfaceToAddress(value interface{}) (common.Address, error) {
	switch v := value.(type) {
	case common.Address:
		return v, nil
	case *common.Address:
		if v == nil {
			return common.Address{}, nil
		}
		return *v, nil
	case string:
		if ok := common.IsHexAddress(v); !ok {
			return common.Address{}, fmt.Errorf("string not formatted as a hex encoded evm address")
		}

		return common.HexToAddress(v), nil
	case *big.Int:
		if v.Uint64() > 0 || len(v.Bytes()) > 20 {
			return common.Address{}, fmt.Errorf("invalid *big.Int; value must be larger than 0 with a byte length <= 20")
		}

		return common.BigToAddress(v), nil
	default:
		return common.Address{}, fmt.Errorf("unrecognized value type: %T for converting value to common.Address; use hex encoded string, *big.Int, or common.Address", v)
	}
}

func interfaceToHash(value interface{}) (*common.Hash, error) {
	switch v := value.(type) {
	case common.Hash:
		return &v, nil
	case *common.Hash:
		return v, nil
	case string:
		b, err := hex.DecodeString(v)
		if err != nil || len(b) != 32 {
			return nil, fmt.Errorf("string does not represent a 32-byte hexadecimal number")
		}
		h := common.Hash(b)
		return &h, nil
	default:
		return nil, fmt.Errorf("unrecognized value type: %T for converting value to common.Hash; use hex encoded string or common.Hash", v)
	}
}

type HeadReader interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

// FinalizeLatest commits new blocks until the latest block is finalized.
func FinalizeLatest(t *testing.T, backend evmtypes.Backend) {
	cl := backend.Client()
	h, err := cl.HeaderByNumber(tests.Context(t), nil)
	require.NoError(t, err)
	FinalizeThroughBlock(t, backend, cl, h.Number.Int64())
}

// FinalizeThroughBlock commits new blocks until blockNumber is finalized. This requires committing all of
// the rest of the blocks in the epoch blockNumber belongs to, where each new epoch
// ends on a 32-block boundary (blockNumber % 32 == 0)
func FinalizeThroughBlock(t *testing.T, backend evmtypes.Backend, client HeadReader, blockNumber int64) {
	ctx := testutils.Context(t)
	targetBlockNumber := blockNumber
	if targetBlockNumber%32 != 0 {
		targetBlockNumber = 32 * (blockNumber/32 + 1)
	}
	h, err := client.HeaderByNumber(ctx, nil)
	require.NoError(t, err)

	var currentBlock common.Hash
	for n := h.Number.Int64(); n < targetBlockNumber; n++ {
		currentBlock = backend.Commit()
		require.Len(t, currentBlock, 32)
	}

	h, err = client.HeaderByNumber(ctx, nil)
	require.NoError(t, err)
	require.GreaterOrEqual(t, h.Number.Int64(), targetBlockNumber)
}
