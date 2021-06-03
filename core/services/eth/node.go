package eth

import (
	"context"
	"fmt"
	"math/big"
	"net/url"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
)

type rawclient struct {
	rpc  *rpc.Client
	geth *ethclient.Client
	uri  url.URL
}

// node represents one ethereum node.
// It must have a ws url and may have a http url
type node struct {
	ws     rawclient
	http   *rawclient
	log    *logger.Logger
	dialed bool
}

func newNode(wsuri url.URL, httpuri *url.URL, name string) (n *node) {
	n = new(node)
	n.log = logger.CreateLogger(logger.Default.With(
		"nodeName", name,
		"nodeTier", "primary",
	))
	n.ws.uri = wsuri
	if httpuri != nil {
		n.http = &rawclient{uri: *httpuri}
	}
	return
}

func (n *node) Dial(ctx context.Context) error {
	if n.dialed {
		panic("eth.Client.Dial(...) should only be called once during the node's lifetime.")
	}

	{
		var httpuri string
		if n.http != nil {
			httpuri = n.http.uri.String()
		}
		n.log.Debugw("eth.Client#Dial(...)", "wsuri", n.ws.uri.String(), "httpuri", httpuri)
	}

	{
		uri := n.ws.uri.String()
		rpc, err := rpc.DialWebsocket(ctx, uri, "")
		if err != nil {
			return err
		}
		n.dialed = true
		n.ws.rpc = rpc
		n.ws.geth = ethclient.NewClient(rpc)
	}

	if n.http != nil {
		uri := n.http.uri.String()
		rpc, err := rpc.DialHTTP(uri)
		if err != nil {
			return err
		}
		n.http.rpc = rpc
		n.http.geth = ethclient.NewClient(rpc)
	}

	return nil
}

// RPC wrappers

func (n node) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	n.log.Debugw("eth.Client#Call(...)",
		"method", method,
		"args", args,
		"mode", switching(n),
	)
	if n.http != nil {
		return n.wrapHTTP(n.http.rpc.CallContext(ctx, result, method, args...))
	}
	return n.wrapWS(n.ws.rpc.CallContext(ctx, result, method, args...))
}

func (n node) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	n.log.Debugw("eth.Client#BatchCall(...)",
		"nBatchElems", len(b),
		"mode", switching(n),
	)
	if n.http != nil {
		return n.wrapHTTP(n.http.rpc.BatchCallContext(ctx, b))
	}
	return n.wrapWS(n.ws.rpc.BatchCallContext(ctx, b))
}

func (n node) EthSubscribe(ctx context.Context, channel interface{}, args ...interface{}) (ethereum.Subscription, error) {
	n.log.Debugw("eth.Client#EthSubscribe", "mode", "websocket")
	return n.ws.rpc.EthSubscribe(ctx, channel, args...)
}

func (n node) Close() {
	n.ws.rpc.Close()
}

// GethClient wrappers

func (n node) TransactionReceipt(ctx context.Context, txHash common.Hash) (receipt *types.Receipt, err error) {
	n.log.Debugw("eth.Client#TransactionReceipt(...)",
		"txHash", txHash,
		"mode", switching(n),
	)

	if n.http != nil {
		receipt, err = n.http.geth.TransactionReceipt(ctx, txHash)
		err = n.wrapHTTP(err)
	} else {
		receipt, err = n.ws.geth.TransactionReceipt(ctx, txHash)
		err = n.wrapWS(err)
	}

	return
}

// NOTE: ChainID may need a bit of rethinking if we implement multiple clients since in theory they could have different ChainIDs
func (n node) ChainID(ctx context.Context) (chainID *big.Int, err error) {
	n.log.Debugw("eth.Client#ChainID(...)", "mode", "websocket")
	chainID, err = n.ws.geth.ChainID(ctx)
	err = n.wrapWS(err)
	return
}

func (n node) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	n.log.Debugw("eth.Client#SendTransaction(...)",
		"tx", tx,
		"mode", switching(n),
	)
	if n.http != nil {
		return n.wrapHTTP(n.http.geth.SendTransaction(ctx, tx))
	}
	return n.wrapWS(n.ws.geth.SendTransaction(ctx, tx))
}

func (n node) PendingNonceAt(ctx context.Context, account common.Address) (nonce uint64, err error) {
	n.log.Debugw("eth.Client#PendingNonceAt(...)",
		"account", account,
		"mode", switching(n),
	)
	if n.http != nil {
		nonce, err = n.http.geth.PendingNonceAt(ctx, account)
		err = n.wrapHTTP(err)
	} else {
		nonce, err = n.ws.geth.PendingNonceAt(ctx, account)
		err = n.wrapWS(err)
	}
	return
}

func (n node) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (nonce uint64, err error) {
	n.log.Debugw("eth.Client#NonceAt(...)",
		"account", account,
		"blockNumber", blockNumber,
		"mode", switching(n),
	)
	if n.http != nil {
		nonce, err = n.http.geth.NonceAt(ctx, account, blockNumber)
		err = n.wrapHTTP(err)
	} else {
		nonce, err = n.ws.geth.NonceAt(ctx, account, blockNumber)
		err = n.wrapWS(err)
	}
	return
}

func (n node) PendingCodeAt(ctx context.Context, account common.Address) (code []byte, err error) {
	n.log.Debugw("eth.Client#PendingCodeAt(...)",
		"account", account,
		"mode", switching(n),
	)
	if n.http != nil {
		code, err = n.http.geth.PendingCodeAt(ctx, account)
		err = n.wrapHTTP(err)
	} else {
		code, err = n.ws.geth.PendingCodeAt(ctx, account)
		err = n.wrapWS(err)
	}
	return
}

func (n node) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) (code []byte, err error) {
	n.log.Debugw("eth.Client#CodeAt(...)",
		"account", account,
		"blockNumber", blockNumber,
		"mode", switching(n),
	)
	if n.http != nil {
		code, err = n.http.geth.CodeAt(ctx, account, blockNumber)
		err = n.wrapHTTP(err)
	} else {
		code, err = n.ws.geth.CodeAt(ctx, account, blockNumber)
		err = n.wrapWS(err)
	}
	return
}

func (n node) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	n.log.Debugw("eth.Client#EstimateGas(...)",
		"call", call,
		"mode", switching(n),
	)
	if n.http != nil {
		gas, err = n.http.geth.EstimateGas(ctx, call)
		err = n.wrapHTTP(err)
	} else {
		gas, err = n.ws.geth.EstimateGas(ctx, call)
		err = n.wrapWS(err)
	}
	return
}

func (n node) SuggestGasPrice(ctx context.Context) (price *big.Int, err error) {
	n.log.Debugw("eth.Client#SuggestGasPrice()", "mode", "websocket")
	price, err = n.ws.geth.SuggestGasPrice(ctx)
	err = n.wrapWS(err)
	return
}

func (n node) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) (val []byte, err error) {
	n.log.Debugw("eth.Client#CallContract()",
		"mode", switching(n),
	)
	if n.http != nil {
		val, err = n.http.geth.CallContract(ctx, msg, blockNumber)
		err = n.wrapHTTP(err)
	} else {
		val, err = n.ws.geth.CallContract(ctx, msg, blockNumber)
		err = n.wrapWS(err)
	}
	return

}

func (n node) BlockByNumber(ctx context.Context, number *big.Int) (b *types.Block, err error) {
	n.log.Debugw("eth.Client#BlockByNumber(...)",
		"number", number,
		"mode", switching(n),
	)
	if n.http != nil {
		b, err = n.http.geth.BlockByNumber(ctx, number)
		err = n.wrapHTTP(err)
	} else {
		b, err = n.ws.geth.BlockByNumber(ctx, number)
		err = n.wrapWS(err)
	}
	return
}

func (n node) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (balance *big.Int, err error) {
	n.log.Debugw("eth.Client#BalanceAt(...)",
		"account", account,
		"blockNumber", blockNumber,
		"mode", switching(n),
	)
	if n.http != nil {
		balance, err = n.http.geth.BalanceAt(ctx, account, blockNumber)
		err = n.wrapHTTP(err)
	} else {
		balance, err = n.ws.geth.BalanceAt(ctx, account, blockNumber)
		err = n.wrapWS(err)
	}
	return
}

func (n node) FilterLogs(ctx context.Context, q ethereum.FilterQuery) (l []types.Log, err error) {
	n.log.Debugw("eth.Client#FilterLogs(...)",
		"q", q,
		"mode", switching(n),
	)
	if n.http != nil {
		l, err = n.http.geth.FilterLogs(ctx, q)
		err = n.wrapHTTP(err)
	} else {
		l, err = n.ws.geth.FilterLogs(ctx, q)
		err = n.wrapWS(err)
	}
	return
}

func (n node) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (sub ethereum.Subscription, err error) {
	n.log.Debugw("eth.Client#SubscribeFilterLogs(...)", "q", q, "mode", "websocket")
	sub, err = n.ws.geth.SubscribeFilterLogs(ctx, q, ch)
	err = n.wrapWS(err)
	return
}

func (n node) wrapWS(err error) error {
	return wrap(err, fmt.Sprintf("primary websocket (%s)", n.ws.uri.String()))
}

func (n node) wrapHTTP(err error) error {
	return wrap(err, fmt.Sprintf("primary http (%s)", n.http.uri.String()))
}

func wrap(err error, tp string) error {
	if err == nil {
		return nil
	}
	if errors.Cause(err).Error() == "context deadline exceeded" {
		err = errors.Wrap(err, "remote eth node timed out")
	}
	return errors.Wrapf(err, "%s call failed", tp)
}

func switching(n node) string {
	if n.http != nil {
		return "http"
	}
	return "websocket"
}
