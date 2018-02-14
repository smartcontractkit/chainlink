package store

import (
	"context"
	"encoding/json"
	"strconv"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
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

// Subscribe registers a subscription for the given address.
// https://github.com/ethereum/go-ethereum/blob/762f3a48a00da02fe58063cb6ce8dc2d08821f15/ethclient/ethclient.go#L359
func (eth *EthClient) Subscribe(channel chan<- []types.Log, addresses []common.Address) (*rpc.ClientSubscription, error) {
	ctx := context.Background()
	sub, err := eth.EthSubscribe(ctx, channel, "logs", toFilterArg(addresses))
	return sub, err
}

// https://github.com/ethereum/go-ethereum/blob/762f3a48a00da02fe58063cb6ce8dc2d08821f15/ethclient/ethclient.go#L363
// https://github.com/ethereum/go-ethereum/blob/762f3a48a00da02fe58063cb6ce8dc2d08821f15/interfaces.go#L132
func toFilterArg(addresses []common.Address) interface{} {
	arg := map[string]interface{}{
		"address": addresses,
	}
	return arg
}

// TxReceipt holds the block number and the transaction hash of a signed
// transaction that has been written to the blockchain.
type TxReceipt struct {
	BlockNumber uint64      `json:"blockNumber"`
	Hash        common.Hash `json:"transactionHash"`
}

// UnmarshalJSON parses the given JSON-encoded data and updates the TxReceipt
// with the values.
func (txr *TxReceipt) UnmarshalJSON(b []byte) error {
	type Rcpt struct {
		BlockNumber string `json:"blockNumber"`
		Hash        string `json:"transactionHash"`
	}
	var rcpt Rcpt
	if err := json.Unmarshal(b, &rcpt); err != nil {
		return err
	}
	block, err := strconv.ParseUint(rcpt.BlockNumber[2:], 16, 64)
	if err != nil {
		return err
	}
	txr.BlockNumber = block
	txr.Hash = common.HexToHash(rcpt.Hash)
	return nil
}

// Unconfirmed returns true if the transaction is not confirmed.
func (txr *TxReceipt) Unconfirmed() bool {
	return common.EmptyHash(txr.Hash)
}
