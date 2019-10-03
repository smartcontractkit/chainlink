package store

import (
	"context"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"chainlink/core/store/assets"
	"chainlink/core/store/models"
	"chainlink/core/utils"
)

//go:generate mockgen -package=mocks -destination=../internal/mocks/eth_client_mocks.go chainlink/core/store EthClient

// EthClient is the interface supplied by EthCallerSubscriber
type EthClient interface {
	GetNonce(address common.Address) (uint64, error)
	GetEthBalance(address common.Address) (*assets.Eth, error)
	GetERC20Balance(address common.Address, contractAddress common.Address) (*big.Int, error)
	SendRawTx(hex string) (common.Hash, error)
	GetTxReceipt(hash common.Hash) (*models.TxReceipt, error)
	GetBlockByNumber(hex string) (models.BlockHeader, error)
	GetLogs(q ethereum.FilterQuery) ([]models.Log, error)
	GetChainID() (*big.Int, error)
	SubscribeToLogs(channel chan<- models.Log, q ethereum.FilterQuery) (models.EthSubscription, error)
	SubscribeToNewHeads(channel chan<- models.BlockHeader) (models.EthSubscription, error)
}

// EthCallerSubscriber holds the CallerSubscriber interface for the Ethereum blockchain.
type EthCallerSubscriber struct {
	CallerSubscriber
}

// CallerSubscriber implements the Call and EthSubscribe functions. Call performs
// a JSON-RPC call with the given arguments and EthSubscribe registers a subscription.
type CallerSubscriber interface {
	Call(result interface{}, method string, args ...interface{}) error
	EthSubscribe(context.Context, interface{}, ...interface{}) (models.EthSubscription, error)
}

// GetNonce returns the nonce (transaction count) for a given address.
func (eth *EthCallerSubscriber) GetNonce(address common.Address) (uint64, error) {
	result := ""
	err := eth.Call(&result, "eth_getTransactionCount", address.Hex(), "latest")
	if err != nil {
		return 0, err
	}
	return utils.HexToUint64(result)
}

// GetEthBalance returns the balance of the given addresses in Ether.
func (eth *EthCallerSubscriber) GetEthBalance(address common.Address) (*assets.Eth, error) {
	result := ""
	amount := new(assets.Eth)
	err := eth.Call(&result, "eth_getBalance", address.Hex(), "latest")
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
func (eth *EthCallerSubscriber) GetERC20Balance(address common.Address, contractAddress common.Address) (*big.Int, error) {
	result := ""
	numLinkBigInt := new(big.Int)
	functionSelector := models.HexToFunctionSelector("0x70a08231") // balanceOf(address)
	data := utils.ConcatBytes(functionSelector.Bytes(), common.LeftPadBytes(address.Bytes(), utils.EVMWordByteLen))
	args := CallArgs{
		To:   contractAddress,
		Data: data,
	}
	err := eth.Call(&result, "eth_call", args, "latest")
	if err != nil {
		return numLinkBigInt, err
	}
	numLinkBigInt.SetString(result, 0)
	return numLinkBigInt, nil
}

// SendRawTx sends a signed transaction to the transaction pool.
func (eth *EthCallerSubscriber) SendRawTx(hex string) (common.Hash, error) {
	result := common.Hash{}
	err := eth.Call(&result, "eth_sendRawTransaction", hex)
	return result, err
}

// GetTxReceipt returns the transaction receipt for the given transaction hash.
func (eth *EthCallerSubscriber) GetTxReceipt(hash common.Hash) (*models.TxReceipt, error) {
	receipt := models.TxReceipt{}
	err := eth.Call(&receipt, "eth_getTransactionReceipt", hash.String())
	return &receipt, err
}

// GetBlockByNumber returns the block for the passed hex, or "latest", "earliest", "pending".
func (eth *EthCallerSubscriber) GetBlockByNumber(hex string) (models.BlockHeader, error) {
	var header models.BlockHeader
	err := eth.Call(&header, "eth_getBlockByNumber", hex, false)
	return header, err
}

// GetLogs returns all logs that respect the passed filter query.
func (eth *EthCallerSubscriber) GetLogs(q ethereum.FilterQuery) ([]models.Log, error) {
	var results []models.Log
	err := eth.Call(&results, "eth_getLogs", utils.ToFilterArg(q))
	return results, err
}

// GetChainID returns the ethereum ChainID.
func (eth *EthCallerSubscriber) GetChainID() (*big.Int, error) {
	value := new(models.Big)
	err := eth.Call(value, "eth_chainId")
	return value.ToInt(), err
}

// SubscribeToLogs registers a subscription for push notifications of logs
// from a given address.
func (eth *EthCallerSubscriber) SubscribeToLogs(
	channel chan<- models.Log,
	q ethereum.FilterQuery,
) (models.EthSubscription, error) {
	// https://github.com/ethereum/go-ethereum/blob/762f3a48a00da02fe58063cb6ce8dc2d08821f15/ethclient/ethclient.go#L359
	ctx := context.Background()
	sub, err := eth.EthSubscribe(ctx, channel, "logs", utils.ToFilterArg(q))
	return sub, err
}

// SubscribeToNewHeads registers a subscription for push notifications of new blocks.
func (eth *EthCallerSubscriber) SubscribeToNewHeads(
	channel chan<- models.BlockHeader,
) (models.EthSubscription, error) {
	ctx := context.Background()
	sub, err := eth.EthSubscribe(ctx, channel, "newHeads")
	return sub, err
}
