package store

import (
	"context"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
)

// EthClient holds the CallerSubscriber interface for the Ethereum blockchain.
type EthClient struct {
	CallerSubscriber
}

// CallerSubscriber implements the Call and EthSubscribe functions. Call performs
// a JSON-RPC call with the given arguments and EthSubscribe registers a subscription.
type CallerSubscriber interface {
	Call(result interface{}, method string, args ...interface{}) error
	EthSubscribe(context.Context, interface{}, ...interface{}) (*rpc.ClientSubscription, error)
}

// GetNonce returns the nonce (transaction count) for a given address.
func (eth *EthClient) GetNonce(address common.Address) (uint64, error) {
	result := ""
	err := eth.Call(&result, "eth_getTransactionCount", address.Hex())
	if err != nil {
		return 0, err
	}
	return utils.HexToUint64(result)
}

func (eth *EthClient) GetWeiBalance(address common.Address) (*big.Int, error) {
	result := ""
	numWeiBigInt := new(big.Int)
	err := eth.Call(&result, "eth_getBalance", address.Hex(), "latest")
	if err != nil {
		return numWeiBigInt, err
	}
	numWeiBigInt.SetString(result, 0)
	return numWeiBigInt, nil
}

func (eth *EthClient) GetEthBalance(address common.Address) (float64, error) {
	numWei, err := eth.GetWeiBalance(address)
	if err != nil {
		return 0, err
	}
	return utils.WeiToEth(numWei), nil
}

// SendRawTx sends a signed transaction to the transaction pool.
func (eth *EthClient) SendRawTx(hex string) (common.Hash, error) {
	result := common.Hash{}
	err := eth.Call(&result, "eth_sendRawTransaction", hex)
	return result, err
}

// GetTxReceipt returns the transaction receipt for the given transaction hash.
func (eth *EthClient) GetTxReceipt(hash common.Hash) (*TxReceipt, error) {
	receipt := TxReceipt{}
	err := eth.Call(&receipt, "eth_getTransactionReceipt", hash.String())
	return &receipt, err
}

// GetBlockNumber returns the block number of the chain head.
func (eth *EthClient) GetBlockNumber() (uint64, error) {
	result := ""
	if err := eth.Call(&result, "eth_blockNumber"); err != nil {
		return 0, err
	}
	return utils.HexToUint64(result)
}

// SubscribeToLogs registers a subscription for push notifications of logs
// from a given address.
func (eth *EthClient) SubscribeToLogs(
	channel chan<- types.Log,
	addresses []common.Address,
) (*rpc.ClientSubscription, error) {
	// https://github.com/ethereum/go-ethereum/blob/762f3a48a00da02fe58063cb6ce8dc2d08821f15/ethclient/ethclient.go#L359
	ctx := context.Background()
	sub, err := eth.EthSubscribe(ctx, channel, "logs", toFilterArg(addresses))
	return sub, err
}

// SubscribeToNewHeads registers a subscription for push notifications of new blocks.
func (eth *EthClient) SubscribeToNewHeads(
	channel chan<- models.BlockHeader,
) (*rpc.ClientSubscription, error) {
	ctx := context.Background()
	sub, err := eth.EthSubscribe(ctx, channel, "newHeads")
	return sub, err
}

// https://github.com/ethereum/go-ethereum/blob/762f3a48a00da02fe58063cb6ce8dc2d08821f15/ethclient/ethclient.go#L363
// https://github.com/ethereum/go-ethereum/blob/762f3a48a00da02fe58063cb6ce8dc2d08821f15/interfaces.go#L132
func toFilterArg(addresses []common.Address) interface{} {
	withoutZeros := utils.WithoutZeroAddresses(addresses)
	if len(withoutZeros) == 0 {
		return map[string]interface{}{}
	}
	arg := map[string]interface{}{
		"address": addresses,
	}
	return arg
}

// TxReceipt holds the block number and the transaction hash of a signed
// transaction that has been written to the blockchain.
type TxReceipt struct {
	BlockNumber hexutil.Big `json:"blockNumber"`
	Hash        common.Hash `json:"transactionHash"`
}

// Unconfirmed returns true if the transaction is not confirmed.
func (txr *TxReceipt) Unconfirmed() bool {
	return common.EmptyHash(txr.Hash)
}
