package eth

import (
	"context"
	"math/big"

	"chainlink/core/assets"
	"chainlink/core/utils"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

//go:generate mockery -name Client -output ../internal/mocks/ -case=underscore

// Client is the interface used to interact with an ethereum node.
type Client interface {
	LogSubscriber
	GetNonce(address common.Address) (uint64, error)
	GetEthBalance(address common.Address) (*assets.Eth, error)
	GetERC20Balance(address common.Address, contractAddress common.Address) (*big.Int, error)
	SendRawTx(hex string) (common.Hash, error)
	GetTxReceipt(hash common.Hash) (*TxReceipt, error)
	GetBlockByNumber(hex string) (BlockHeader, error)
	GetChainID() (*big.Int, error)
	SubscribeToNewHeads(channel chan<- BlockHeader) (Subscription, error)
}

// LogSubscriber encapsulates only the methods needed for subscribing to ethereum log events.
type LogSubscriber interface {
	GetLogs(q ethereum.FilterQuery) ([]Log, error)
	SubscribeToLogs(channel chan<- Log, q ethereum.FilterQuery) (Subscription, error)
}

// Subscription holds the methods for an ethereum log subscription.
type Subscription interface {
	Err() <-chan error
	Unsubscribe()
}

// CallerSubscriberClient implements the ethereum Client interface using a
// CallerSubscriber instance.
type CallerSubscriberClient struct {
	CallerSubscriber
}

// CallerSubscriber implements the Call and Subscribe functions. Call performs
// a JSON-RPC call with the given arguments and Subscribe registers a subscription,
// using an open stream to receive updates from ethereum node.
type CallerSubscriber interface {
	Call(result interface{}, method string, args ...interface{}) error
	Subscribe(context.Context, interface{}, ...interface{}) (Subscription, error)
}

// GetNonce returns the nonce (transaction count) for a given address.
func (client *CallerSubscriberClient) GetNonce(address common.Address) (uint64, error) {
	result := ""
	err := client.Call(&result, "eth_getTransactionCount", address.Hex(), "latest")
	if err != nil {
		return 0, err
	}
	return utils.HexToUint64(result)
}

// GetEthBalance returns the balance of the given addresses in Ether.
func (client *CallerSubscriberClient) GetEthBalance(address common.Address) (*assets.Eth, error) {
	result := ""
	amount := new(assets.Eth)
	err := client.Call(&result, "eth_getBalance", address.Hex(), "latest")
	if err != nil {
		return amount, err
	}
	amount.SetString(result, 0)
	return amount, nil
}

// CallArgs represents the data used to call the balance method of an ERC
// contract. "To" is the address of the ERC contract. "Data" is the message sent
// to the contract.
type CallArgs struct {
	To   common.Address `json:"to"`
	Data hexutil.Bytes  `json:"data"`
}

// GetERC20Balance returns the balance of the given address for the token contract address.
func (client *CallerSubscriberClient) GetERC20Balance(address common.Address, contractAddress common.Address) (*big.Int, error) {
	result := ""
	numLinkBigInt := new(big.Int)
	functionSelector := HexToFunctionSelector("0x70a08231") // balanceOf(address)
	data := utils.ConcatBytes(functionSelector.Bytes(), common.LeftPadBytes(address.Bytes(), utils.EVMWordByteLen))
	args := CallArgs{
		To:   contractAddress,
		Data: data,
	}
	err := client.Call(&result, "eth_call", args, "latest")
	if err != nil {
		return numLinkBigInt, err
	}
	numLinkBigInt.SetString(result, 0)
	return numLinkBigInt, nil
}

// SendRawTx sends a signed transaction to the transaction pool.
func (client *CallerSubscriberClient) SendRawTx(hex string) (common.Hash, error) {
	result := common.Hash{}
	err := client.Call(&result, "eth_sendRawTransaction", hex)
	return result, err
}

// GetTxReceipt returns the transaction receipt for the given transaction hash.
func (client *CallerSubscriberClient) GetTxReceipt(hash common.Hash) (*TxReceipt, error) {
	receipt := TxReceipt{}
	err := client.Call(&receipt, "eth_getTransactionReceipt", hash.String())
	return &receipt, err
}

// GetBlockByNumber returns the block for the passed hex, or "latest", "earliest", "pending".
func (client *CallerSubscriberClient) GetBlockByNumber(hex string) (BlockHeader, error) {
	var header BlockHeader
	err := client.Call(&header, "eth_getBlockByNumber", hex, false)
	return header, err
}

// GetLogs returns all logs that respect the passed filter query.
func (client *CallerSubscriberClient) GetLogs(q ethereum.FilterQuery) ([]Log, error) {
	var results []Log
	err := client.Call(&results, "eth_getLogs", utils.ToFilterArg(q))
	return results, err
}

// GetChainID returns the ethereum ChainID.
func (client *CallerSubscriberClient) GetChainID() (*big.Int, error) {
	value := new(utils.Big)
	err := client.Call(value, "eth_chainId")
	return value.ToInt(), err
}

// SubscribeToLogs registers a subscription for push notifications of logs
// from a given address.
func (client *CallerSubscriberClient) SubscribeToLogs(
	channel chan<- Log,
	q ethereum.FilterQuery,
) (Subscription, error) {
	// https://github.com/ethereum/go-ethereum/blob/762f3a48a00da02fe58063cb6ce8dc2d08821f15/ethclient/ethclient.go#L359
	ctx := context.Background()
	sub, err := client.Subscribe(ctx, channel, "logs", utils.ToFilterArg(q))
	return sub, err
}

// SubscribeToNewHeads registers a subscription for push notifications of new blocks.
func (client *CallerSubscriberClient) SubscribeToNewHeads(
	channel chan<- BlockHeader,
) (Subscription, error) {
	ctx := context.Background()
	sub, err := client.Subscribe(ctx, channel, "newHeads")
	return sub, err
}
