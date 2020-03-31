package client

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gobuffalo/packr/v2"
)

// ETH interface represents the connection to the Ethereum node
type ETH interface {
	ABI(string) (abi.ABI, error)
	Call(
		to common.Address,
		abi *abi.ABI,
		sig string,
		res interface{},
		args ...interface{},
	) error
	SubscribeToLogs(chan<- types.Log, ethereum.FilterQuery) (Subscription, error)
	TransactionByHash(txHash common.Hash) (*types.Transaction, error)
	SubscribeToNewHeads(chan<- types.Header) (Subscription, error)
}

type eth struct {
	url     *url.URL
	rpc     *rpc.Client
	timeout time.Duration
	abiBox  *packr.Box
}

// NewClient will return a connected ETH implementation
func NewClient(urlStr string) (ETH, error) {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, err
	}

	rpc, err := rpc.Dial(u.String())
	box := packr.New("contract-abi", "./abi")
	return &eth{
		url:    u,
		rpc:    rpc,
		abiBox: box,
	}, err
}

// ABI will return the ABI instance for a given filename
func (c *eth) ABI(filename string) (abi.ABI, error) {
	b, err := c.abiBox.Find(filename)
	if err != nil {
		return abi.ABI{}, err
	}
	return abi.JSON(bytes.NewBuffer(b))
}

func (c *eth) Call(
	to common.Address,
	abi *abi.ABI,
	sig string,
	res interface{},
	args ...interface{},
) error {
	if data, err := abi.Pack(sig, args...); err != nil {
		return err
	} else if resp, err := c.call(to, data); err != nil {
		return err
	} else if err := abi.Unpack(res, sig, resp); err != nil {
		return err
	}
	return nil
}

// TransactionByHash calls `eth_getTransactionByHash` for a given tx hash
func (c *eth) TransactionByHash(txHash common.Hash) (*types.Transaction, error) {
	var tx types.Transaction
	return &tx, c.rpc.Call(&tx, "eth_getTransactionByHash", txHash.String())
}

// SubscribeToNewHeads returns an instantiated subscription type, subscribing to heads
func (c *eth) SubscribeToNewHeads(channel chan<- types.Header) (Subscription, error) {
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
