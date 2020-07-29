package eth

import (
	"context"
	"math/big"
	"strings"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

//go:generate mockery --name Client --output ../../internal/mocks/ --case=underscore
//go:generate mockery --name GethClient --output ../../internal/mocks/ --case=underscore
//go:generate mockery --name RPCClient --output ../../internal/mocks/ --case=underscore
//go:generate mockery --name Subscription --output ../../internal/mocks/ --case=underscore

// Client is the interface used to interact with an ethereum node.
type Client interface {
	GethClient
	RPCClient

	Dial(ctx context.Context) error
	GetERC20Balance(address common.Address, contractAddress common.Address) (*big.Int, error)
	SendRawTx(bytes []byte) (common.Hash, error)
}

// GethClient is an interface that represents go-ethereum's own ethclient
// https://github.com/ethereum/go-ethereum/blob/master/ethclient/ethclient.go
type GethClient interface {
	ChainID(ctx context.Context) (*big.Int, error)
	SendTransaction(context.Context, *types.Transaction) error
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	HeaderByNumber(ctx context.Context, n *big.Int) (*types.Header, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
}

// RPCClient is an interface that represents go-ethereum's own rpc.Client
// https://github.com/ethereum/go-ethereum/blob/master/rpc/client.go
type RPCClient interface {
	Call(result interface{}, method string, args ...interface{}) error
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	Close()
}

// This interface only exists so that we can generate a mock for it.  It is
// identical to `ethereum.Subscription`.
type Subscription interface {
	Err() <-chan error
	Unsubscribe()
}

// client implements the ethereum Client interface using a
// CallerSubscriber instance.
type client struct {
	GethClient
	RPCClient
	url    string // For reestablishing the connection after a disconnect
	mocked bool
}

var _ Client = (*client)(nil)

func NewClient(url string) *client {
	return &client{url: url}
}

// This alternate constructor exists for testing purposes.
func NewClientWith(rpcClient RPCClient, gethClient GethClient) *client {
	return &client{
		GethClient: gethClient,
		RPCClient:  rpcClient,
		mocked:     true,
	}
}

func (client *client) Dial(ctx context.Context) error {
	logger.Debugw("eth.Client#Dial(...)")
	if client.mocked {
		return nil
	} else if client.RPCClient != nil || client.GethClient != nil {
		panic("eth.Client.Dial(...) should only be called once during the application's lifetime.")
	}

	rpcClient, err := rpc.DialContext(ctx, client.url)
	if err != nil {
		return err
	}
	client.RPCClient = rpcClient
	client.GethClient = ethclient.NewClient(rpcClient)
	return nil
}

// CallArgs represents the data used to call the balance method of a contract.
// "To" is the address of the ERC contract. "Data" is the message sent
// to the contract.
type CallArgs struct {
	To   common.Address `json:"to"`
	Data hexutil.Bytes  `json:"data"`
}

// GetERC20Balance returns the balance of the given address for the token contract address.
func (client *client) GetERC20Balance(address common.Address, contractAddress common.Address) (*big.Int, error) {
	logger.Debugw("eth.Client#GetERC20Balance(...)",
		"address", address,
		"contractAddress", contractAddress,
	)
	result := ""
	numLinkBigInt := new(big.Int)
	functionSelector := models.HexToFunctionSelector("0x70a08231") // balanceOf(address)
	data := utils.ConcatBytes(functionSelector.Bytes(), common.LeftPadBytes(address.Bytes(), utils.EVMWordByteLen))
	args := CallArgs{
		To:   contractAddress,
		Data: data,
	}
	err := client.RPCClient.Call(&result, "eth_call", args, "latest")
	if err != nil {
		return numLinkBigInt, err
	}
	numLinkBigInt.SetString(result, 0)
	return numLinkBigInt, nil
}

// SendRawTx sends a signed transaction to the transaction pool.
func (client *client) SendRawTx(bytes []byte) (common.Hash, error) {
	logger.Debugw("eth.Client#SendRawTx(...)",
		"bytes", bytes,
	)
	result := common.Hash{}
	err := client.RPCClient.Call(&result, "eth_sendRawTransaction", hexutil.Encode(bytes))
	return result, err
}

// We wrap the GethClient's `TransactionReceipt` method so that we can ignore the error that arises
// when we're talking to a Parity node that has no receipt yet.
func (client *client) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	logger.Debugw("eth.Client#TransactionReceipt(...)",
		"txHash", txHash,
	)
	receipt, err := client.GethClient.TransactionReceipt(ctx, txHash)
	if err != nil && strings.Contains(err.Error(), "missing required field 'transactionHash'") {
		return nil, ethereum.NotFound
	}
	return receipt, err
}

func (client *client) ChainID(ctx context.Context) (*big.Int, error) {
	logger.Debugw("eth.Client#ChainID(...)")
	return client.GethClient.ChainID(ctx)
}

func (client *client) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	logger.Debugw("eth.Client#SendTransaction(...)",
		"tx", tx,
	)
	return client.GethClient.SendTransaction(ctx, tx)
}

func (client *client) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	logger.Debugw("eth.Client#PendingNonceAt(...)",
		"account", account,
	)
	return client.GethClient.PendingNonceAt(ctx, account)
}

func (client *client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	logger.Debugw("eth.Client#BlockByNumber(...)",
		"number", number,
	)
	return client.GethClient.BlockByNumber(ctx, number)
}

func (client *client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	logger.Debugw("eth.Client#HeaderByNumber(...)",
		"number", number,
	)
	return client.GethClient.HeaderByNumber(ctx, number)
}

func (client *client) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	logger.Debugw("eth.Client#BalanceAt(...)",
		"account", account,
		"blockNumber", blockNumber,
	)
	return client.GethClient.BalanceAt(ctx, account, blockNumber)
}

func (client *client) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	logger.Debugw("eth.Client#FilterLogs(...)",
		"q", q,
	)
	return client.GethClient.FilterLogs(ctx, q)
}

func (client *client) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	logger.Debugw("eth.Client#SubscribeFilterLogs(...)",
		"q", q,
	)
	return client.GethClient.SubscribeFilterLogs(ctx, q, ch)
}

func (client *client) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	logger.Debugw("eth.Client#SubscribeNewHead(...)")
	return client.GethClient.SubscribeNewHead(ctx, ch)
}

func (client *client) Call(result interface{}, method string, args ...interface{}) error {
	logger.Debugw("eth.Client#Call(...)",
		"method", method,
		"args", args,
	)
	return client.RPCClient.Call(result, method, args...)
}

func (client *client) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	logger.Debugw("eth.Client#Call(...)",
		"method", method,
		"args", args,
	)
	return client.RPCClient.CallContext(ctx, result, method, args...)
}
