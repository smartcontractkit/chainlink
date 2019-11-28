package store

import (
	"context"
	"math/big"

	"chainlink/core/eth"
	"chainlink/core/store/assets"
	"chainlink/core/store/models"
	"chainlink/core/utils"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

//go:generate mockery -name EthClient -output ../internal/mocks/ -case=underscore

// EthClient is the interface supplied by EthCallerSubscriber
type EthClient interface {
	GetNonce(address common.Address) (uint64, error)
	GetEthBalance(address common.Address) (*assets.Eth, error)
	GetERC20Balance(address common.Address, contractAddress common.Address) (*big.Int, error)
	SendRawTx(hex string) (common.Hash, error)
	GetTxReceipt(hash common.Hash) (*eth.TxReceipt, error)
	GetBlockByNumber(hex string) (eth.BlockHeader, error)
	GetLogs(q ethereum.FilterQuery) ([]eth.Log, error)
	GetChainID() (*big.Int, error)
	SubscribeToLogs(channel chan<- eth.Log, q ethereum.FilterQuery) (eth.EthSubscription, error)
	SubscribeToNewHeads(channel chan<- eth.BlockHeader) (eth.EthSubscription, error)
}

// EthCallerSubscriber holds the CallerSubscriber interface for the Ethereum blockchain.
type EthCallerSubscriber struct {
	CallerSubscriber
}

// CallerSubscriber implements the Call and EthSubscribe functions. Call performs
// a JSON-RPC call with the given arguments and EthSubscribe registers a subscription.
type CallerSubscriber interface {
	Call(result interface{}, method string, args ...interface{}) error
	EthSubscribe(context.Context, interface{}, ...interface{}) (eth.EthSubscription, error)
}

// GetNonce returns the nonce (transaction count) for a given address.
func (ecs *EthCallerSubscriber) GetNonce(address common.Address) (uint64, error) {
	result := ""
	err := ecs.Call(&result, "eth_getTransactionCount", address.Hex(), "latest")
	if err != nil {
		return 0, err
	}
	return utils.HexToUint64(result)
}

// GetEthBalance returns the balance of the given addresses in Ether.
func (ecs *EthCallerSubscriber) GetEthBalance(address common.Address) (*assets.Eth, error) {
	result := ""
	amount := new(assets.Eth)
	err := ecs.Call(&result, "eth_getBalance", address.Hex(), "latest")
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
func (ecs *EthCallerSubscriber) GetERC20Balance(address common.Address, contractAddress common.Address) (*big.Int, error) {
	result := ""
	numLinkBigInt := new(big.Int)
	functionSelector := models.HexToFunctionSelector("0x70a08231") // balanceOf(address)
	data := utils.ConcatBytes(functionSelector.Bytes(), common.LeftPadBytes(address.Bytes(), utils.EVMWordByteLen))
	args := CallArgs{
		To:   contractAddress,
		Data: data,
	}
	err := ecs.Call(&result, "eth_call", args, "latest")
	if err != nil {
		return numLinkBigInt, err
	}
	numLinkBigInt.SetString(result, 0)
	return numLinkBigInt, nil
}

// SendRawTx sends a signed transaction to the transaction pool.
func (ecs *EthCallerSubscriber) SendRawTx(hex string) (common.Hash, error) {
	result := common.Hash{}
	err := ecs.Call(&result, "eth_sendRawTransaction", hex)
	return result, err
}

// GetTxReceipt returns the transaction receipt for the given transaction hash.
func (ecs *EthCallerSubscriber) GetTxReceipt(hash common.Hash) (*eth.TxReceipt, error) {
	receipt := eth.TxReceipt{}
	err := ecs.Call(&receipt, "eth_getTransactionReceipt", hash.String())
	return &receipt, err
}

// GetBlockByNumber returns the block for the passed hex, or "latest", "earliest", "pending".
func (ecs *EthCallerSubscriber) GetBlockByNumber(hex string) (eth.BlockHeader, error) {
	var header eth.BlockHeader
	err := ecs.Call(&header, "eth_getBlockByNumber", hex, false)
	return header, err
}

// GetLogs returns all logs that respect the passed filter query.
func (ecs *EthCallerSubscriber) GetLogs(q ethereum.FilterQuery) ([]eth.Log, error) {
	var results []eth.Log
	err := ecs.Call(&results, "eth_getLogs", utils.ToFilterArg(q))
	return results, err
}

// GetChainID returns the ethereum ChainID.
func (ecs *EthCallerSubscriber) GetChainID() (*big.Int, error) {
	value := new(utils.Big)
	err := ecs.Call(value, "eth_chainId")
	return value.ToInt(), err
}

// SubscribeToLogs registers a subscription for push notifications of logs
// from a given address.
func (ecs *EthCallerSubscriber) SubscribeToLogs(
	channel chan<- eth.Log,
	q ethereum.FilterQuery,
) (eth.EthSubscription, error) {
	// https://github.com/ethereum/go-ethereum/blob/762f3a48a00da02fe58063cb6ce8dc2d08821f15/ethclient/ethclient.go#L359
	ctx := context.Background()
	sub, err := ecs.EthSubscribe(ctx, channel, "logs", utils.ToFilterArg(q))
	return sub, err
}

// SubscribeToNewHeads registers a subscription for push notifications of new blocks.
func (ecs *EthCallerSubscriber) SubscribeToNewHeads(
	channel chan<- eth.BlockHeader,
) (eth.EthSubscription, error) {
	ctx := context.Background()
	sub, err := ecs.EthSubscribe(ctx, channel, "newHeads")
	return sub, err
}
