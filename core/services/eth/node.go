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

type Node interface {
	Dial(ctx context.Context) error
	Close()
	Verify(ctx context.Context, expectedChainID *big.Int) (err error)

	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
	EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
	HeaderByNumber(context.Context, *big.Int) (*types.Header, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
	EthSubscribe(ctx context.Context, channel interface{}, args ...interface{}) (ethereum.Subscription, error)

	String() string
}

type rawclient struct {
	rpc  *rpc.Client
	geth *ethclient.Client
	uri  url.URL
}

// Node represents one ethereum node.
// It must have a ws url and may have a http url
type node struct {
	ws     rawclient
	http   *rawclient
	log    logger.Logger
	name   string
	dialed bool
}

func NewNode(lggr logger.Logger, wsuri url.URL, httpuri *url.URL, name string) Node {
	n := new(node)
	n.name = name
	n.log = lggr.With(
		"nodeName", name,
		"nodeTier", "primary",
	)
	n.ws.uri = wsuri
	if httpuri != nil {
		n.http = &rawclient{uri: *httpuri}
	}
	return n
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
			return errors.Wrapf(err, "Error while dialing websocket: %v", uri)
		}
		n.dialed = true
		n.ws.rpc = rpc
		n.ws.geth = ethclient.NewClient(rpc)
	}

	if n.http != nil {
		uri := n.http.uri.String()
		rpc, err := rpc.DialHTTP(uri)
		if err != nil {
			return errors.Wrapf(err, "Error while dialing HTTP: %v", uri)
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

func (n node) HeaderByNumber(ctx context.Context, number *big.Int) (header *types.Header, err error) {
	n.log.Debugw("eth.Client#HeaderByNumber(...)",
		"number", n,
		"mode", switching(n),
	)
	if n.http != nil {
		header, err = n.http.geth.HeaderByNumber(ctx, number)
		err = n.wrapHTTP(err)
	} else {
		header, err = n.ws.geth.HeaderByNumber(ctx, number)
		err = n.wrapWS(err)
	}
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

func (n node) SuggestGasTipCap(ctx context.Context) (tipCap *big.Int, err error) {
	n.log.Debugw("eth.Client#SuggestGasTipCap(...)",
		"mode", switching(n),
	)
	if n.http != nil {
		tipCap, err = n.http.geth.SuggestGasTipCap(ctx)
		err = n.wrapHTTP(err)
	} else {
		tipCap, err = n.ws.geth.SuggestGasTipCap(ctx)
		err = n.wrapWS(err)
	}
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

func (n node) String() string {
	s := fmt.Sprintf("(primary)%s:%s", n.name, n.ws.uri.String())
	if n.http != nil {
		s = s + fmt.Sprintf(":%s", n.http.uri.String())
	}
	return s
}

// Verify checks that all connections to eth nodes match the given chain ID
func (n node) Verify(ctx context.Context, expectedChainID *big.Int) (err error) {
	var chainID *big.Int
	if chainID, err = n.ws.geth.ChainID(ctx); err != nil {
		return errors.Wrapf(err, "failed to verify chain ID for node %s", n.name)
	} else if chainID.Cmp(expectedChainID) != 0 {
		return errors.Errorf(
			"websocket rpc ChainID doesn't match local chain ID: RPC ID=%s, local ID=%s, node name=%s",
			chainID.String(),
			expectedChainID.String(),
			n.name,
		)
	}
	if n.http != nil {
		if chainID, err = n.http.geth.ChainID(ctx); err != nil {
			return errors.Wrapf(err, "failed to verify chain ID for node %s", n.name)
		} else if chainID.Cmp(expectedChainID) != 0 {
			return errors.Errorf(
				"http rpc ChainID doesn't match local chain ID: RPC ID=%s, local ID=%s, node name=%s",
				chainID.String(),
				expectedChainID.String(),
				n.name,
			)
		}
	}
	return nil
}
