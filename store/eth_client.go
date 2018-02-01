package store

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

func (eth *EthClient) Subscribe(channel chan EventLog, address string) error {
	ctx := context.Background()
	_, err := eth.EthSubscribe(ctx, channel, "logs", address)
	return err
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

type EventLog struct {
	Address   common.Address  `json:"address"`
	BlockHash common.Hash     `json:"blockHash"`
	TxHash    common.Hash     `json:"transactionHash"`
	Data      hexutil.Bytes   `json:"data"`
	Topics    []hexutil.Bytes `json:"topics"`
}

type EthNotification struct {
	Params json.RawMessage `json:"params"`
}

func (en EthNotification) UnmarshalLog() (EventLog, error) {
	var el EventLog
	var rval map[string]json.RawMessage

	if err := json.Unmarshal(en.Params, &rval); err != nil {
		return el, err
	}

	if err := json.Unmarshal(rval["result"], &el); err != nil {
		return el, err
	}

	if el.Address == utils.ZeroAddress {
		return el, errors.New("Cannot unmarshal a log with a zero address")
	}

	return el, nil
}
