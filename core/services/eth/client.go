package eth

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/chainlink/core/assets"
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

// Client is the interface used to interact with an ethereum node.
type Client interface {
	GethClient
	RPCClient

	Dial(ctx context.Context) error
	GetEthBalance(address common.Address) (*assets.Eth, error)
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

// client implements the ethereum Client interface using a
// CallerSubscriber instance.
type client struct {
	GethClient
	RPCClient
	url string // For reestablishing the connection after a disconnect
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
	}
}

func (client *client) Dial(ctx context.Context) error {
	if client.RPCClient != nil {
		return nil
	}

	rpcClient, err := rpc.DialContext(ctx, client.url)
	if err != nil {
		return err
	}
	client.RPCClient = rpcClient
	client.GethClient = ethclient.NewClient(rpcClient)
	return nil
}

// GetEthBalance returns the balance of the given addresses in Ether.
func (client *client) GetEthBalance(address common.Address) (*assets.Eth, error) {
	result := ""
	amount := new(assets.Eth)
	err := client.logCall(&result, "eth_getBalance", address.Hex(), "latest")
	if err != nil {
		return amount, err
	}
	amount.SetString(result, 0)
	return amount, nil
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
	result := ""
	numLinkBigInt := new(big.Int)
	functionSelector := models.HexToFunctionSelector("0x70a08231") // balanceOf(address)
	data := utils.ConcatBytes(functionSelector.Bytes(), common.LeftPadBytes(address.Bytes(), utils.EVMWordByteLen))
	args := CallArgs{
		To:   contractAddress,
		Data: data,
	}
	err := client.logCall(&result, "eth_call", args, "latest")
	if err != nil {
		return numLinkBigInt, err
	}
	numLinkBigInt.SetString(result, 0)
	return numLinkBigInt, nil
}

// SendRawTx sends a signed transaction to the transaction pool.
func (client *client) SendRawTx(bytes []byte) (common.Hash, error) {
	result := common.Hash{}
	err := client.logCall(&result, "eth_sendRawTransaction", hexutil.Encode(bytes))
	return result, err
}

// We wrap the GethClient's `TransactionReceipt` method so that we can ignore the error that arises
// when we're talking to a Parity node that has no receipt yet.
func (client *client) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	receipt, err := client.GethClient.TransactionReceipt(ctx, txHash)
	if err != nil && strings.Contains(err.Error(), "missing required field 'transactionHash' for Log") {
		return nil, ethereum.NotFound
	}
	return receipt, err
}

// logCall logs an RPC call's method and arguments, and then calls the method
func (client *client) logCall(result interface{}, method string, args ...interface{}) error {
	logger.Debugw(
		fmt.Sprintf(`Calling eth client RPC method "%s"`, method),
		"args", args,
	)
	return client.RPCClient.Call(result, method, args...)
}
