package client

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// ETH interface represents the connection to the Ethereum node
type ETH interface {
	SubscribeToLogs(chan<- types.Log, ethereum.FilterQuery) (Subscription, error)
	TransactionByHash(txHash common.Hash) (*types.Transaction, error)
	SubscribeToNewHeads(chan<- BlockHeader) (Subscription, error)
}

type eth struct {
	url     *url.URL
	rpc     *rpc.Client
	timeout time.Duration
}

// NewClient will return a connected ETH implementation
func NewClient(urlStr string) (ETH, error) {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, err
	}

	rpc, err := rpc.Dial(u.String())
	return &eth{
		url: u,
		rpc: rpc,
	}, err
}

// TransactionByHash calls `eth_getTransactionByHash` for a given tx hash
func (c *eth) TransactionByHash(txHash common.Hash) (*types.Transaction, error) {
	var tx types.Transaction
	return &tx, c.rpc.Call(&tx, "eth_getTransactionByHash", txHash.String())
}

// BlockHeader represents a block header in the Ethereum blockchain.
// Deliberately does not have required fields because some fields aren't
// present depending on the Ethereum node.
// i.e. Parity does not always send mixHash
type BlockHeader struct {
	ParentHash  common.Hash      `json:"parentHash"`
	UncleHash   common.Hash      `json:"sha3Uncles"`
	Coinbase    common.Address   `json:"miner"`
	Root        common.Hash      `json:"stateRoot"`
	TxHash      common.Hash      `json:"transactionsRoot"`
	ReceiptHash common.Hash      `json:"receiptsRoot"`
	Bloom       types.Bloom      `json:"logsBloom"`
	Difficulty  hexutil.Big      `json:"difficulty"`
	Number      hexutil.Big      `json:"number"`
	GasLimit    hexutil.Uint64   `json:"gasLimit"`
	GasUsed     hexutil.Uint64   `json:"gasUsed"`
	Time        hexutil.Big      `json:"timestamp"`
	Extra       hexutil.Bytes    `json:"extraData"`
	Nonce       types.BlockNonce `json:"nonce"`
	GethHash    common.Hash      `json:"mixHash"`
	ParityHash  common.Hash      `json:"hash"`
}

// SubscribeToNewHeads returns an instantiated subscription type, subscribing to heads
func (c *eth) SubscribeToNewHeads(channel chan<- BlockHeader) (Subscription, error) {
	ctx := context.Background()
	sub, err := c.rpc.EthSubscribe(ctx, channel, "newHeads")
	return sub, err
}

// Subscription is the interface for managing eth log subscriptions
type Subscription interface {
	Err() <-chan error
	Unsubscribe()
}

// SubscribeToLogs returns an instantiated subscription type, subscribing to logs based on the
// given filter query
func (c *eth) SubscribeToLogs(channel chan<- types.Log, q ethereum.FilterQuery) (Subscription, error) {
	ctx := context.Background()
	sub, err := c.rpc.EthSubscribe(ctx, channel, "logs", toFilterArg(q))
	return sub, err
}

func (c *eth) call(to common.Address, data []byte) ([]byte, error) {
	var res string
	if err := c.rpc.Call(&res, "eth_call", map[string]interface{}{
		"to":   to.String(),
		"data": fmt.Sprintf("0x%s", common.Bytes2Hex(data)),
	}, "latest"); err != nil {
		return []byte{}, err
	}
	return common.FromHex(res), nil
}

// https://github.com/ethereum/go-ethereum/blob/762f3a48a00da02fe58063cb6ce8dc2d08821f15/ethclient/ethclient.go#L359
func toFilterArg(q ethereum.FilterQuery) interface{} {
	arg := map[string]interface{}{
		"fromBlock": toBlockNumArg(q.FromBlock),
		"toBlock":   toBlockNumArg(q.ToBlock),
		"address":   q.Addresses,
		"topics":    q.Topics,
	}
	if q.FromBlock == nil {
		arg["fromBlock"] = "0x0"
	}
	return arg
}

// https://github.com/ethereum/go-ethereum/blob/762f3a48a00da02fe58063cb6ce8dc2d08821f15/ethclient/ethclient.go#L256
func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return hexutil.EncodeBig(number)
}
