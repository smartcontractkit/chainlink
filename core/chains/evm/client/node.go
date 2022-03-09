package client

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"sync"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	uuid "github.com/satori/go.uuid"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	promEVMPoolRPCNodeDials = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_dials_total",
		Help: "The total number of dials for the given RPC node",
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodeDialsFailed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_dials_failed",
		Help: "The total number of failed dials for the given RPC node",
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodeDialsSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_dials_success",
		Help: "The total number of successful dials for the given RPC node",
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodeVerifies = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_verifies",
		Help: "The total number of chain ID verifications for the given RPC node",
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodeVerifiesFailed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_verifies_failed",
		Help: "The total number of failed chain ID verifications for the given RPC node",
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodeVerifiesSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_verifies_success",
		Help: "The total number of successful chain ID verifications for the given RPC node",
	}, []string{"evmChainID", "nodeName"})

	promEVMPoolRPCNodeCalls = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_calls_total",
		Help: "The approximate total number of RPC calls for the given RPC node",
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodeCallsFailed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_calls_failed",
		Help: "The approximate total number of failed RPC calls for the given RPC node",
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodeCallsSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_calls_success",
		Help: "The approximate total number of successful RPC calls for the given RPC node",
	}, []string{"evmChainID", "nodeName"})
)

//go:generate mockery --name Node --output ../mocks/ --case=underscore

// Node represents a client that connects to an ethereum-compatible RPC node
type Node interface {
	Start(ctx context.Context) error
	Close()

	State() NodeState
	// Unique identifier for node
	ID() int32
	ChainID() *big.Int

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
	EthSubscribe(ctx context.Context, channel chan<- *evmtypes.Head, args ...interface{}) (ethereum.Subscription, error)

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
	utils.StartStopOnce
	ws      rawclient
	http    *rawclient
	log     logger.Logger
	name    string
	id      int32
	chainID *big.Int
	cfg     NodeConfig

	state   NodeState
	stateMu sync.RWMutex

	// ctx can be cancelled to immediately cancel all in-flight requests on
	// this node. All access to ctx and cancel should be serialized through
	// stateMu since the context can be cancelled and replaced on state
	// transitions as well as node Close.
	ctx    context.Context
	cancel context.CancelFunc
	// wg waits for subsidiary goroutines
	wg sync.WaitGroup

	// nLiveNodes is a passed in function that allows this node to
	// query a parent object to see how many live nodes there are in total.
	// This is done so we can prevent the last alive node in a pool from being
	// moved to out-of-sync state. It is better to have one out-of-sync node
	// than no nodes at all.
	nLiveNodes func() int
}

// NodeConfig allows configuration of the node
type NodeConfig interface {
	NodeNoNewHeadsThreshold() time.Duration
	NodePollFailureThreshold() uint32
	NodePollInterval() time.Duration
}

// NewNode returns a new *node as Node
func NewNode(nodeCfg NodeConfig, lggr logger.Logger, wsuri url.URL, httpuri *url.URL, name string, id int32, chainID *big.Int) Node {
	n := new(node)
	n.name = name
	n.id = id
	n.chainID = chainID
	n.cfg = nodeCfg
	n.ws.uri = wsuri
	if httpuri != nil {
		n.http = &rawclient{uri: *httpuri}
	}
	n.ctx, n.cancel = context.WithCancel(context.Background())
	n.log = lggr.Named("Node").With(
		"nodeTier", "primary",
		"nodeName", name,
		"node", n.String(),
		"evmChainID", chainID,
	)
	return n
}

// Start dials and verifies the node
// Should only be called once in a node's lifecycle
// Return value is necessary to conform to interface but this will never
// actually return an error.
func (n *node) Start(startCtx context.Context) error {
	return n.StartOnce(n.name, func() error {
		n.start(startCtx)
		return nil
	})
}

// start initially dials the node and verifies chain ID
// This spins off lifecycle goroutines.
// Not thread-safe.
// Node lifecycle is synchronous: only one goroutine should be running at a
// time.
func (n *node) start(startCtx context.Context) {
	if n.state != NodeStateUndialed {
		panic(fmt.Sprintf("cannot dial node with state %v", n.state))
	}

	dialCtx, cancel := n.wrapCtx(startCtx)
	defer cancel()
	if err := n.dial(dialCtx); err != nil {
		n.log.Errorw("Dial failed: EVM Node is unreachable", "err", err)
		n.declareUnreachable()
		return
	}
	n.setState(NodeStateDialed)

	verifyCtx, cancel := n.wrapCtx(startCtx)
	defer cancel()
	if err := n.verify(verifyCtx); errors.Is(err, errInvalidChainID) {
		n.log.Errorw("Verify failed: EVM Node has the wrong chain ID", "err", err)
		n.declareInvalidChainID()
		return
	} else if err != nil {
		n.log.Errorw(fmt.Sprintf("Verify failed: %v", err), "err", err)
		n.declareUnreachable()
		return
	}

	n.declareAlive()
}

// Not thread-safe
// Pure dial: does not mutate node "state" field.
func (n *node) dial(ctx context.Context) error {
	promEVMPoolRPCNodeDials.WithLabelValues(n.chainID.String(), n.name).Inc()
	var httpuri string
	if n.http != nil {
		httpuri = n.http.uri.String()
	}
	n.log.Debugw("RPC dial: evmclient.Client#dial", "wsuri", n.ws.uri.String(), "httpuri", httpuri)

	uri := n.ws.uri.String()
	wsrpc, err := rpc.DialWebsocket(ctx, uri, "")
	if err != nil {
		promEVMPoolRPCNodeDialsFailed.WithLabelValues(n.chainID.String(), n.name).Inc()
		return errors.Wrapf(err, "error while dialing websocket: %v", uri)
	}

	var httprpc *rpc.Client
	if n.http != nil {
		httprpc, err = rpc.DialHTTP(httpuri)
		if err != nil {
			promEVMPoolRPCNodeDialsFailed.WithLabelValues(n.chainID.String(), n.name).Inc()
			return errors.Wrapf(err, "error while dialing HTTP: %v", uri)
		}
	}

	n.ws.rpc = wsrpc
	n.ws.geth = ethclient.NewClient(wsrpc)

	if n.http != nil {
		n.http.rpc = httprpc
		n.http.geth = ethclient.NewClient(httprpc)
	}

	n.log.Debugw("RPC dial: success", "wsuri", n.ws.uri.String(), "httpuri", httpuri)
	promEVMPoolRPCNodeDialsSuccess.WithLabelValues(n.chainID.String(), n.name).Inc()

	return nil
}

var errInvalidChainID = errors.New("invalid chain id")

// verify checks that all connections to eth nodes match the given chain ID
// Not thread-safe
// Pure verify: does not mutate node "state" field.
func (n *node) verify(ctx context.Context) (err error) {
	promEVMPoolRPCNodeVerifies.WithLabelValues(n.chainID.String(), n.name).Inc()
	ctx, cancel := n.wrapCtx(ctx)
	defer cancel()

	promFailed := func() {
		promEVMPoolRPCNodeVerifiesFailed.WithLabelValues(n.chainID.String(), n.name).Inc()
	}

	switch n.state {
	case NodeStateDialed, NodeStateOutOfSync, NodeStateInvalidChainID:
	default:
		panic(fmt.Sprintf("cannot verify node in state %v", n.state))
	}

	var chainID *big.Int
	if chainID, err = n.ws.geth.ChainID(ctx); err != nil {
		promFailed()
		return errors.Wrapf(err, "failed to verify chain ID for node %s", n.name)
	} else if chainID.Cmp(n.chainID) != 0 {
		promFailed()
		return errors.Wrapf(
			errInvalidChainID,
			"websocket rpc ChainID doesn't match local chain ID: RPC ID=%s, local ID=%s, node name=%s",
			chainID.String(),
			n.chainID.String(),
			n.name,
		)
	}
	if n.http != nil {
		if chainID, err = n.http.geth.ChainID(ctx); err != nil {
			promFailed()
			return errors.Wrapf(err, "failed to verify chain ID for node %s", n.name)
		} else if chainID.Cmp(n.chainID) != 0 {
			promFailed()
			return errors.Wrapf(
				errInvalidChainID,
				"http rpc ChainID doesn't match local chain ID: RPC ID=%s, local ID=%s, node name=%s",
				chainID.String(),
				n.chainID.String(),
				n.name,
			)
		}
	}

	promEVMPoolRPCNodeVerifiesSuccess.WithLabelValues(n.chainID.String(), n.name).Inc()

	return nil
}

func (n *node) Close() {
	err := n.StopOnce(n.name, func() error {
		defer n.wg.Wait()

		n.stateMu.Lock()
		defer n.stateMu.Unlock()

		n.cancel()
		n.state = NodeStateClosed
		if n.ws.rpc != nil {
			n.ws.rpc.Close()
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

// RPC wrappers

// CallContext implementation
func (n *node) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	ctx, cancel, err := n.wrapLiveCtx(ctx)
	if err != nil {
		return err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With(
		"method", method,
		"args", args,
	)

	lggr.Debug("RPC call: evmclient.Client#CallContext")
	if n.http != nil {
		err = n.wrapHTTP(n.http.rpc.CallContext(ctx, result, method, args...))
	} else {
		err = n.wrapWS(n.ws.rpc.CallContext(ctx, result, method, args...))
	}

	n.logResult(lggr, err)

	return err
}

func (n *node) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	ctx, cancel, err := n.wrapLiveCtx(ctx)
	if err != nil {
		return err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("nBatchElems", len(b))

	lggr.Debug("RPC call: evmclient.Client#BatchCallContext")
	if n.http != nil {
		err = n.wrapHTTP(n.http.rpc.BatchCallContext(ctx, b))
	} else {
		err = n.wrapWS(n.ws.rpc.BatchCallContext(ctx, b))
	}

	n.logResult(lggr, err)

	return err
}

func (n *node) EthSubscribe(ctx context.Context, channel chan<- *evmtypes.Head, args ...interface{}) (ethereum.Subscription, error) {
	ctx, cancel, err := n.wrapLiveCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr("websocket").With("args", args)

	lggr.Debug("RPC call: evmclient.Client#EthSubscribe")
	sub, err := n.ws.rpc.EthSubscribe(ctx, channel, args...)

	n.logResult(lggr, err)

	return sub, err
}

// GethClient wrappers

func (n *node) TransactionReceipt(ctx context.Context, txHash common.Hash) (receipt *types.Receipt, err error) {
	ctx, cancel, err := n.wrapLiveCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("txHash", txHash)

	lggr.Debug("RPC call: evmclient.Client#TransactionReceipt")

	if n.http != nil {
		receipt, err = n.http.geth.TransactionReceipt(ctx, txHash)
		err = n.wrapHTTP(err)
	} else {
		receipt, err = n.ws.geth.TransactionReceipt(ctx, txHash)
		err = n.wrapWS(err)
	}

	n.logResult(lggr, err, "receipt", receipt)

	return
}

func (n *node) HeaderByNumber(ctx context.Context, number *big.Int) (header *types.Header, err error) {
	ctx, cancel, err := n.wrapLiveCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("number", number)

	lggr.Debug("RPC call: evmclient.Client#HeaderByNumber")
	if n.http != nil {
		header, err = n.http.geth.HeaderByNumber(ctx, number)
		err = n.wrapHTTP(err)
	} else {
		header, err = n.ws.geth.HeaderByNumber(ctx, number)
		err = n.wrapWS(err)
	}

	n.logResult(lggr, err, "header", header)

	return
}

func (n *node) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	ctx, cancel, err := n.wrapLiveCtx(ctx)
	if err != nil {
		return err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("tx", tx)

	lggr.Debug("RPC call: evmclient.Client#SendTransaction")
	if n.http != nil {
		err = n.wrapHTTP(n.http.geth.SendTransaction(ctx, tx))
	} else {
		err = n.wrapWS(n.ws.geth.SendTransaction(ctx, tx))
	}

	n.logResult(lggr, err)

	return err
}

func (n *node) PendingNonceAt(ctx context.Context, account common.Address) (nonce uint64, err error) {
	ctx, cancel, err := n.wrapLiveCtx(ctx)
	if err != nil {
		return 0, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("account", account)

	lggr.Debug("RPC call: evmclient.Client#PendingNonceAt")
	if n.http != nil {
		nonce, err = n.http.geth.PendingNonceAt(ctx, account)
		err = n.wrapHTTP(err)
	} else {
		nonce, err = n.ws.geth.PendingNonceAt(ctx, account)
		err = n.wrapWS(err)
	}

	n.logResult(lggr, err, "nonce", nonce)

	return
}

// NonceAt is a bit of a misnomer. You might expect it to return the highest
// mined nonce at the given block number, but it actually returns the total
// transaction count which is the highest mined nonce + 1
func (n *node) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (nonce uint64, err error) {
	ctx, cancel, err := n.wrapLiveCtx(ctx)
	if err != nil {
		return 0, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("account", account, "blockNumber", blockNumber)

	lggr.Debug("RPC call: evmclient.Client#NonceAt")
	if n.http != nil {
		nonce, err = n.http.geth.NonceAt(ctx, account, blockNumber)
		err = n.wrapHTTP(err)
	} else {
		nonce, err = n.ws.geth.NonceAt(ctx, account, blockNumber)
		err = n.wrapWS(err)
	}

	n.logResult(lggr, err, "nonce", nonce)

	return
}

func (n *node) PendingCodeAt(ctx context.Context, account common.Address) (code []byte, err error) {
	ctx, cancel, err := n.wrapLiveCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("account", account)

	lggr.Debug("RPC call: evmclient.Client#PendingCodeAt")
	if n.http != nil {
		code, err = n.http.geth.PendingCodeAt(ctx, account)
		err = n.wrapHTTP(err)
	} else {
		code, err = n.ws.geth.PendingCodeAt(ctx, account)
		err = n.wrapWS(err)
	}

	n.logResult(lggr, err, "code", code)

	return
}

func (n *node) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) (code []byte, err error) {
	ctx, cancel, err := n.wrapLiveCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("account", account, "blockNumber", blockNumber)

	lggr.Debug("RPC call: evmclient.Client#CodeAt")
	if n.http != nil {
		code, err = n.http.geth.CodeAt(ctx, account, blockNumber)
		err = n.wrapHTTP(err)
	} else {
		code, err = n.ws.geth.CodeAt(ctx, account, blockNumber)
		err = n.wrapWS(err)
	}

	n.logResult(lggr, err, "code", code)

	return
}

func (n *node) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	ctx, cancel, err := n.wrapLiveCtx(ctx)
	if err != nil {
		return 0, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("call", call)

	lggr.Debug("RPC call: evmclient.Client#EstimateGas")
	if n.http != nil {
		gas, err = n.http.geth.EstimateGas(ctx, call)
		err = n.wrapHTTP(err)
	} else {
		gas, err = n.ws.geth.EstimateGas(ctx, call)
		err = n.wrapWS(err)
	}

	n.logResult(lggr, err, "gas", gas)

	return
}

func (n *node) SuggestGasPrice(ctx context.Context) (price *big.Int, err error) {
	ctx, cancel, err := n.wrapLiveCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n))

	lggr.Debug("RPC call: evmclient.Client#SuggestGasPrice")
	if n.http != nil {
		price, err = n.http.geth.SuggestGasPrice(ctx)
		err = n.wrapHTTP(err)
	} else {
		price, err = n.ws.geth.SuggestGasPrice(ctx)
		err = n.wrapWS(err)
	}

	n.logResult(lggr, err, "price", price)

	return
}

func (n *node) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) (val []byte, err error) {
	ctx, cancel, err := n.wrapLiveCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("msg", msg, "blockNumber", blockNumber)

	lggr.Debug("RPC call: evmclient.Client#CallContract")
	if n.http != nil {
		val, err = n.http.geth.CallContract(ctx, msg, blockNumber)
		err = n.wrapHTTP(err)
	} else {
		val, err = n.ws.geth.CallContract(ctx, msg, blockNumber)
		err = n.wrapWS(err)
	}

	n.logResult(lggr, err, "val", val)

	return

}

func (n *node) BlockByNumber(ctx context.Context, number *big.Int) (b *types.Block, err error) {
	ctx, cancel, err := n.wrapLiveCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("number", number)

	lggr.Debug("RPC call: evmclient.Client#BlockByNumber")
	if n.http != nil {
		b, err = n.http.geth.BlockByNumber(ctx, number)
		err = n.wrapHTTP(err)
	} else {
		b, err = n.ws.geth.BlockByNumber(ctx, number)
		err = n.wrapWS(err)
	}

	n.logResult(lggr, err, "block", b)

	return
}

func (n *node) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (balance *big.Int, err error) {
	ctx, cancel, err := n.wrapLiveCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("account", account.Hex(), "blockNumber", blockNumber)

	lggr.Debug("RPC call: evmclient.Client#BalanceAt")
	if n.http != nil {
		balance, err = n.http.geth.BalanceAt(ctx, account, blockNumber)
		err = n.wrapHTTP(err)
	} else {
		balance, err = n.ws.geth.BalanceAt(ctx, account, blockNumber)
		err = n.wrapWS(err)
	}

	n.logResult(lggr, err, "balance", balance)

	return
}

func (n *node) FilterLogs(ctx context.Context, q ethereum.FilterQuery) (l []types.Log, err error) {
	ctx, cancel, err := n.wrapLiveCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("q", q)

	lggr.Debug("RPC call: evmclient.Client#FilterLogs")
	if n.http != nil {
		l, err = n.http.geth.FilterLogs(ctx, q)
		err = n.wrapHTTP(err)
	} else {
		l, err = n.ws.geth.FilterLogs(ctx, q)
		err = n.wrapWS(err)
	}

	n.logResult(lggr, err, "log", l)

	return
}

func (n *node) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (sub ethereum.Subscription, err error) {
	ctx, cancel, err := n.wrapLiveCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr("websocket").With("q", q)

	lggr.Debug("RPC call: evmclient.Client#SubscribeFilterLogs")
	sub, err = n.ws.geth.SubscribeFilterLogs(ctx, q, ch)
	err = n.wrapWS(err)

	n.logResult(lggr, err)

	return
}

func (n *node) SuggestGasTipCap(ctx context.Context) (tipCap *big.Int, err error) {
	ctx, cancel, err := n.wrapLiveCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n))

	lggr.Debug("RPC call: evmclient.Client#SuggestGasTipCap")
	if n.http != nil {
		tipCap, err = n.http.geth.SuggestGasTipCap(ctx)
		err = n.wrapHTTP(err)
	} else {
		tipCap, err = n.ws.geth.SuggestGasTipCap(ctx)
		err = n.wrapWS(err)
	}

	n.logResult(lggr, err, "tipCap", tipCap)

	return
}

func (n *node) ChainID() (chainID *big.Int) { return n.chainID }

// newRqLggr generates a new logger with a unique request ID
func (n *node) newRqLggr(mode string) logger.Logger {
	return n.log.With(
		"requestID", uuid.NewV4(),
		"mode", mode,
	)
}

func (n *node) logResult(lggr logger.Logger, err error, results ...interface{}) {
	promEVMPoolRPCNodeCalls.WithLabelValues(n.chainID.String(), n.name).Inc()
	if err == nil {
		promEVMPoolRPCNodeCallsSuccess.WithLabelValues(n.chainID.String(), n.name).Inc()
		lggr.Debugw("RPC call success", results...)
	} else {
		promEVMPoolRPCNodeCallsFailed.WithLabelValues(n.chainID.String(), n.name).Inc()
		lggr.Debugw("RPC call failure", "err", err)
	}
}

func (n *node) wrapWS(err error) error {
	err = wrap(err, fmt.Sprintf("primary websocket (%s)", n.ws.uri.String()))
	if err != nil {
		n.log.Debugw("Call failed", "err", err)
	} else {
		n.log.Trace("Call succeeded")
	}
	return err
}

func (n *node) wrapHTTP(err error) error {
	err = wrap(err, fmt.Sprintf("primary http (%s)", n.http.uri.String()))
	if err != nil {
		n.log.Debugw("Call failed", "err", err)
	} else {
		n.log.Trace("Call succeeded")
	}
	return err
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

// wrapLiveCtx adds a default timeout and combines with the node's master context
// which can be cancelled early.
//
// Returns error if node is not "alive".
//
// NOTE: We don't want to wrap all calls in a mutex lock via IfAlive or
// something similar because rpc calls can be slow and there are a lot
// of them: we might end up holding a read lock for a really long time
// through various overlapping requests, preventing state transitions.
//
// Instead, we check if the node is alive, and if so, copy the context pointer.
// If the node is marked dead during the request, the master context will be
// cancelled and the request will exit early.
func (n *node) wrapLiveCtx(parentCtx context.Context) (combinedCtx context.Context, cancel context.CancelFunc, err error) {
	// Need to wrap in mutex because state transition can cancel and replace the
	// context
	n.stateMu.RLock()
	if n.state != NodeStateAlive {
		err = errors.Errorf("cannot execute RPC call on node with state: %s", n.state)
		n.stateMu.RUnlock()
		return
	}
	nodeCtx := n.ctx
	n.stateMu.RUnlock()
	combinedCtx, cancel = wrapCtx(parentCtx, nodeCtx)
	return
}

func (n *node) wrapCtx(parentCtx context.Context) (combinedCtx context.Context, cancel context.CancelFunc) {
	nodeCtx := n.getCtx()
	combinedCtx, cancel = wrapCtx(parentCtx, nodeCtx)
	return
}

// getCtx wraps context access in the stateMu since state transitions can
// cancel and replace contexts
func (n *node) getCtx() context.Context {
	n.stateMu.RLock()
	defer n.stateMu.RUnlock()
	return n.ctx
}

func wrapCtx(parentCtx, nodeCtx context.Context) (combinedCtx context.Context, cancel context.CancelFunc) {
	combinedCtx, cancel = utils.CombinedContext(parentCtx, nodeCtx, queryTimeout)
	return
}

func (n *node) ctxWithDefaultTimeout() (context.Context, context.CancelFunc) {
	return DefaultQueryCtx(n.getCtx())
}

func switching(n *node) string {
	if n.http != nil {
		return "http"
	}
	return "websocket"
}

func (n *node) String() string {
	s := fmt.Sprintf("(primary)%s:%s", n.name, n.ws.uri.String())
	if n.http != nil {
		s = s + fmt.Sprintf(":%s", n.http.uri.String())
	}
	return s
}

func (n *node) ID() int32 {
	return n.id
}
