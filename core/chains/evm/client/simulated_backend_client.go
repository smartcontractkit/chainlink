package client

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/assets"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// SimulatedBackendClient is an Client implementation using a simulated
// blockchain backend. Note that not all RPC methods are implemented here.
type SimulatedBackendClient struct {
	b       *backends.SimulatedBackend
	t       testing.TB
	chainId *big.Int
}

// NewSimulatedBackendClient creates an eth client backed by a simulated backend.
func NewSimulatedBackendClient(t testing.TB, b *backends.SimulatedBackend, chainId *big.Int) *SimulatedBackendClient {
	return &SimulatedBackendClient{
		b:       b,
		t:       t,
		chainId: chainId,
	}
}

// Dial noop for the sim.
func (c *SimulatedBackendClient) Dial(context.Context) error {
	return nil
}

// Close does nothing. We ought not close the underlying backend here since
// other simulated clients might still be using it
func (c *SimulatedBackendClient) Close() {}

// checkEthCallArgs extracts and verifies the arguments for an eth_call RPC
func (c *SimulatedBackendClient) checkEthCallArgs(
	args []interface{}) (*CallArgs, *big.Int, error) {
	if len(args) != 2 {
		return nil, nil, fmt.Errorf(
			"should have two arguments after \"eth_call\", got %d", len(args))
	}
	callArgs, ok := args[0].(map[string]interface{})
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
	ca := CallArgs{
		From: callArgs["from"].(common.Address),
		To:   *callArgs["to"].(*common.Address),
		Data: callArgs["data"].(hexutil.Bytes),
	}
	return &ca, blockNumber, nil
}

// CallContext mocks the ethereum client RPC calls used by chainlink, copying the
// return value into result.
func (c *SimulatedBackendClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	switch method {
	case "eth_call":
		callArgs, _, err := c.checkEthCallArgs(args)
		if err != nil {
			return err
		}
		callMsg := ethereum.CallMsg{From: callArgs.From, To: &callArgs.To, Data: callArgs.Data}
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

// SubscribeFilterLogs registers a subscription for push notifications of logs
// from a given address.
func (c *SimulatedBackendClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, channel chan<- types.Log) (ethereum.Subscription, error) {
	return c.b.SubscribeFilterLogs(ctx, q, channel)
}

// GetEthBalance helper to get eth balance
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
func (c *SimulatedBackendClient) GetERC20Balance(ctx context.Context, address common.Address, contractAddress common.Address) (balance *big.Int, err error) {
	callData, err := balanceOfABI.Pack("balanceOf", address)
	if err != nil {
		return nil, errors.Wrapf(err, "while seeking the ERC20 balance of %s on %s",
			address, contractAddress)
	}
	b, err := c.b.CallContract(ctx, ethereum.CallMsg{
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

// GetLINKBalance get link balance.
func (c *SimulatedBackendClient) GetLINKBalance(ctx context.Context, linkAddress common.Address, address common.Address) (*assets.Link, error) {
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

// HeadByNumber returns our own header type.
func (c *SimulatedBackendClient) HeadByNumber(ctx context.Context, n *big.Int) (*evmtypes.Head, error) {
	if n == nil {
		n = c.currentBlockNumber()
	}
	header, err := c.b.HeaderByNumber(ctx, n)
	if err != nil {
		return nil, err
	} else if header == nil {
		return nil, ethereum.NotFound
	}
	return &evmtypes.Head{
		EVMChainID: utils.NewBigI(c.chainId.Int64()),
		Hash:       header.Hash(),
		Number:     header.Number.Int64(),
		ParentHash: header.ParentHash,
	}, nil
}

// HeadByHash returns our own header type.
func (c *SimulatedBackendClient) HeadByHash(ctx context.Context, h common.Hash) (*evmtypes.Head, error) {
	header, err := c.b.HeaderByHash(ctx, h)
	if err != nil {
		return nil, err
	} else if header == nil {
		return nil, ethereum.NotFound
	}
	return &evmtypes.Head{
		EVMChainID: utils.NewBigI(c.chainId.Int64()),
		Hash:       header.Hash(),
		Number:     header.Number.Int64(),
		ParentHash: header.ParentHash,
	}, nil
}

// BlockByNumber returns a geth block type.
func (c *SimulatedBackendClient) BlockByNumber(ctx context.Context, n *big.Int) (*types.Block, error) {
	return c.b.BlockByNumber(ctx, n)
}

// BlockByNumber returns a geth block type.
func (c *SimulatedBackendClient) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return c.b.BlockByHash(ctx, hash)
}

// ChainID returns the ethereum ChainID.
func (c *SimulatedBackendClient) ChainID() *big.Int {
	return c.chainId
}

// PendingNonceAt gets pending nonce i.e. mempool nonce.
func (c *SimulatedBackendClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return c.b.PendingNonceAt(ctx, account)
}

// NonceAt gets nonce as of a specified block.
func (c *SimulatedBackendClient) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return c.b.NonceAt(ctx, account, blockNumber)
}

// BalanceAt gets balance as of a specified block.
func (c *SimulatedBackendClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return c.b.BalanceAt(ctx, account, blockNumber)
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
	subscription.subscription, err = c.b.SubscribeNewHead(ctx, ch)
	if err != nil {
		return nil, errors.Wrapf(err, "could not subscribe to new heads on "+
			"simulated backend")
	}
	go func() {
		var lastHead *evmtypes.Head
		for {
			select {
			case h := <-ch:
				var head *evmtypes.Head
				if h != nil {
					head = &evmtypes.Head{Number: h.Number.Int64(), Hash: h.Hash(), ParentHash: h.ParentHash, Parent: lastHead, EVMChainID: utils.NewBig(c.chainId)}
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
	return c.b.HeaderByNumber(ctx, n)
}

func (c *SimulatedBackendClient) HeaderByHash(ctx context.Context, h common.Hash) (*types.Header, error) {
	return c.b.HeaderByHash(ctx, h)
}

// SendTransaction sends a transaction.
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
	res, err := c.b.CallContract(ctx, msg, blockNumber)
	if err != nil {
		dataErr := revertError{}
		isCustomRevert := errors.As(err, &dataErr)
		if isCustomRevert {
			return nil, &JsonError{Data: dataErr.ErrorData(), Message: dataErr.Error(), Code: 3}
		}
		// Generic revert, no data
		return nil, &JsonError{Data: []byte{}, Message: err.Error(), Code: 3}
	}
	return res, nil
}

// CodeAt gets the code associated with an account as of a specified block.
func (c *SimulatedBackendClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return c.b.CodeAt(ctx, account, blockNumber)
}

// PendingCodeAt gets the latest code.
func (c *SimulatedBackendClient) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return c.b.PendingCodeAt(ctx, account)
}

// EstimateGas estimates gas for a msg.
func (c *SimulatedBackendClient) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	return c.b.EstimateGas(ctx, call)
}

// SuggestGasPrice recommends a gas price.
func (c *SimulatedBackendClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	panic("unimplemented")
}

// BatchCallContext makes a batch rpc call.
func (c *SimulatedBackendClient) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	select {
	case <-ctx.Done():
		return errors.New("context canceled")
	default:
		//do nothing
	}

	for i, elem := range b {
		switch elem.Method {
		case "eth_getTransactionReceipt":
			if _, ok := elem.Result.(*evmtypes.Receipt); !ok {
				return errors.Errorf("SimulatedBackendClient expected return type of *evmtypes.Receipt for eth_getTransactionReceipt, got type %T", elem.Result)
			}
			if len(elem.Args) != 1 {
				return errors.Errorf("SimulatedBackendClient expected 1 arg, got %d for eth_getTransactionReceipt", len(elem.Args))
			}
			hash, is := elem.Args[0].(common.Hash)
			if !is {
				return errors.Errorf("SimulatedBackendClient expected arg to be a hash, got: %T", elem.Args[0])
			}
			receipt, err := c.b.TransactionReceipt(ctx, hash)
			b[i].Result = evmtypes.FromGethReceipt(receipt)
			b[i].Error = err
		case "eth_getBlockByNumber":
			switch v := elem.Result.(type) {
			case *evmtypes.Head:
			case *evmtypes.Block:
			default:
				return errors.Errorf("SimulatedBackendClient expected return type of [*evmtypes.Head] or [*evmtypes.Block] for eth_getBlockByNumber, got type %T", v)
			}
			if len(elem.Args) != 2 {
				return errors.Errorf("SimulatedBackendClient expected 2 args, got %d for eth_getBlockByNumber", len(elem.Args))
			}
			blockNum, is := elem.Args[0].(string)
			if !is {
				return errors.Errorf("SimulatedBackendClient expected first arg to be a string for eth_getBlockByNumber, got: %T", elem.Args[0])
			}
			_, is = elem.Args[1].(bool)
			if !is {
				return errors.Errorf("SimulatedBackendClient expected second arg to be a boolean for eth_getBlockByNumber, got: %T", elem.Args[1])
			}
			n, ok := new(big.Int).SetString(blockNum, 0)
			if !ok {
				return errors.Errorf("error while converting block number string: %s to big.Int ", blockNum)
			}
			header, err := c.b.HeaderByNumber(ctx, n)
			if err != nil {
				return err
			}
			switch v := elem.Result.(type) {
			case *evmtypes.Head:
				b[i].Result = &evmtypes.Head{
					Number: header.Number.Int64(),
					Hash:   header.Hash(),
				}
			case *evmtypes.Block:
				b[i].Result = &evmtypes.Block{
					Number: header.Number.Int64(),
					Hash:   header.Hash(),
				}
			default:
				return errors.Errorf("SimulatedBackendClient Unexpected Type %T", v)
			}

			b[i].Error = err
		case "eth_call":
			if len(elem.Args) != 2 {
				return errors.Errorf("SimulatedBackendClient expected 2 args, got %d for eth_call", len(elem.Args))
			}

			_, ok := elem.Result.(*string)
			if !ok {
				return errors.Errorf("SimulatedBackendClient expected result to be *string for eth_call, got: %T", elem.Result)
			}

			params, ok := elem.Args[0].(map[string]interface{})
			if !ok {
				return errors.Errorf("SimulatedBackendClient expected first arg to be map[string]interface{} for eth_call, got: %T", elem.Args[0])
			}

			blockNum, ok := elem.Args[1].(string)
			if !ok {
				return errors.Errorf("SimulatedBackendClient expected second arg to be a string for eth_call, got: %T", elem.Args[1])
			}

			if blockNum != "" {
				if _, ok = new(big.Int).SetString(blockNum, 0); !ok {
					return errors.Errorf("error while converting block number string: %s to big.Int ", blockNum)
				}
			}

			callMsg := toCallMsg(params)
			resp, err := c.b.CallContract(ctx, callMsg, nil)
			*(b[i].Result.(*string)) = hexutil.Encode(resp)
			b[i].Error = err
		default:
			return errors.Errorf("SimulatedBackendClient got unsupported method %s", elem.Method)
		}
	}
	return nil
}

// BatchCallContextAll makes a batch rpc call.
func (c *SimulatedBackendClient) BatchCallContextAll(ctx context.Context, b []rpc.BatchElem) error {
	return c.BatchCallContext(ctx, b)
}

// SuggestGasTipCap suggests a gas tip cap.
func (c *SimulatedBackendClient) SuggestGasTipCap(ctx context.Context) (tipCap *big.Int, err error) {
	return c.b.SuggestGasTipCap(ctx)
}

func (c *SimulatedBackendClient) Backend() *backends.SimulatedBackend {
	return c.b
}

// NodeStates implements evmclient.Client
func (c *SimulatedBackendClient) NodeStates() map[string]string { return nil }

// Commit imports all the pending transactions as a single block and starts a
// fresh new state.
func (c *SimulatedBackendClient) Commit() common.Hash {
	return c.b.Commit()
}

func toCallMsg(params map[string]interface{}) ethereum.CallMsg {
	var callMsg ethereum.CallMsg

	switch to := params["to"].(type) {
	case string:
		toAddr := common.HexToAddress(to)
		callMsg.To = &toAddr
	case common.Address:
		callMsg.To = &to
	case *common.Address:
		callMsg.To = to
	default:
		panic("unexpected type of 'to' parameter")
	}

	switch from := params["from"].(type) {
	case nil:
		// This parameter is not required so nil is acceptable
	case string:
		callMsg.From = common.HexToAddress(from)
	case common.Address:
		callMsg.From = from
	case *common.Address:
		callMsg.From = *from
	default:
		panic("unexpected type of 'from' parameter")
	}

	switch data := params["data"].(type) {
	case nil:
		// This parameter is not required so nil is acceptable
	case hexutil.Bytes:
		callMsg.Data = data
	case []byte:
		callMsg.Data = data
	default:
		panic("unexpected type of 'data' parameter")
	}

	if value, ok := params["value"].(*big.Int); ok {
		callMsg.Value = value
	}

	if gas, ok := params["gas"].(uint64); ok {
		callMsg.Gas = gas
	}

	if gasPrice, ok := params["gasPrice"].(*big.Int); ok {
		callMsg.GasPrice = gasPrice
	}

	return callMsg
}
