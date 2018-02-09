package store

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink/utils"
)

type EthClient struct {
	CallerSubscriber
}

type CallerSubscriber interface {
	Call(result interface{}, method string, args ...interface{}) error
	EthSubscribe(context.Context, interface{}, ...interface{}) (*rpc.ClientSubscription, error)
}

func (eth *EthClient) GetNonce(account accounts.Account) (uint64, error) {
	var result string
	err := eth.Call(&result, "eth_getTransactionCount", account.Address.Hex())
	if err != nil {
		return 0, err
	}
	return utils.HexToUint64(result)
}

func (eth *EthClient) SendRawTx(hex string) (common.Hash, error) {
	result := common.Hash{}
	err := eth.Call(&result, "eth_sendRawTransaction", hex)
	return result, err
}

func (eth *EthClient) GetTxReceipt(hash common.Hash) (*TxReceipt, error) {
	receipt := TxReceipt{}
	err := eth.Call(&receipt, "eth_getTransactionReceipt", hash.String())
	return &receipt, err
}

func (eth *EthClient) BlockNumber() (uint64, error) {
	result := ""
	if err := eth.Call(&result, "eth_blockNumber"); err != nil {
		return 0, err
	}
	return utils.HexToUint64(result)
}

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

type TxReceipt struct {
	BlockNumber uint64      `json:"blockNumber"`
	Hash        common.Hash `json:"transactionHash"`
}

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
	if txr.Hash, err = utils.StringToHash(rcpt.Hash); err != nil {
		return err
	}
	return nil
}

func (txr *TxReceipt) Unconfirmed() bool {
	return common.EmptyHash(txr.Hash)
}

// why this over geth types.Log?
//type EventLog struct {
//Address     common.Address  `json:"address"`
//BlockHash   common.Hash     `json:"blockHash"`
//BlockNumber hexutil.Uint64  `json:"blockNumber"`
//Data        hexutil.Bytes   `json:"data"`
//LogIndex    hexutil.Uint64  `json:"logIndex"`
//Topics      []hexutil.Bytes `json:"topics"`
//TxHash      common.Hash     `json:"transactionHash"`
//TxIndex     hexutil.Uint64  `json:"transactionIndex"`
//}

//type EthNotification struct {
//Params json.RawMessage `json:"params"`
//}

//func (en EthNotification) UnmarshalLog() (EventLog, error) {
//var el EventLog
//var rval map[string]json.RawMessage

//if err := json.Unmarshal(en.Params, &rval); err != nil {
//return el, err
//}

//if err := json.Unmarshal(rval["result"], &el); err != nil {
//return el, err
//}

//if el.Address == utils.ZeroAddress {
//return el, errors.New("Cannot unmarshal a log with a zero address")
//}

//return el, nil
//}
