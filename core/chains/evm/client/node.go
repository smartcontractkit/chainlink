package client

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
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
	promEVMPoolRPCCallTiming = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "evm_pool_rpc_node_rpc_call_time",
		Help: "The duration of an RPC call in nanoseconds",
		Buckets: []float64{
			float64(50 * time.Millisecond),
			float64(100 * time.Millisecond),
			float64(200 * time.Millisecond),
			float64(500 * time.Millisecond),
			float64(1 * time.Second),
			float64(2 * time.Second),
			float64(4 * time.Second),
			float64(8 * time.Second),
		},
	}, []string{"evmChainID", "nodeName", "rpcHost", "isSendOnly", "success", "rpcCallName"})
)

//go:generate mockery --quiet --name Node --output ../mocks/ --case=underscore

// Node represents a client that connects to an ethereum-compatible RPC node
type Node interface {
	Start(ctx context.Context) error
	Close() error

	// State returns NodeState
	State() NodeState
	// StateAndLatest returns NodeState with the latest received block number & total difficulty.
	StateAndLatest() (state NodeState, blockNum int64, totalDifficulty *utils.Big)
	// Name is a unique identifier for this node.
	Name() string
	ChainID() *big.Int

	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
	EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
	HeaderByNumber(context.Context, *big.Int) (*types.Header, error)
	HeaderByHash(context.Context, common.Hash) (*types.Header, error)
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
	lfcLog  logger.Logger
	rpcLog  logger.Logger
	name    string
	id      int32
	chainID *big.Int
	cfg     NodeConfig

	ws   rawclient
	http *rawclient

	stateMu sync.RWMutex // protects state* fields
	state   NodeState
	// Each node is tracking the last received head number and total difficulty
	stateLatestBlockNumber     int64
	stateLatestTotalDifficulty *utils.Big

	// Need to track subscriptions because closing the RPC does not (always?)
	// close the underlying subscription
	subs []ethereum.Subscription

	// chStopInFlight can be closed to immediately cancel all in-flight requests on
	// this node. Closing and replacing should be serialized through
	// stateMu since it can happen on state transitions as well as node Close.
	chStopInFlight chan struct{}
	// nodeCtx is the node lifetime's context
	nodeCtx context.Context
	// cancelNodeCtx cancels nodeCtx when stopping the node
	cancelNodeCtx context.CancelFunc
	// wg waits for subsidiary goroutines
	wg sync.WaitGroup

	// nLiveNodes is a passed in function that allows this node to:
	//  1. see how many live nodes there are in total, so we can prevent the last alive node in a pool from being
	//  moved to out-of-sync state. It is better to have one out-of-sync node than no nodes at all.
	//  2. compare against the highest head (by number or difficulty) to ensure we don't fall behind too far.
	nLiveNodes func() (count int, blockNumber int64, totalDifficulty *utils.Big)
}

// NodeConfig allows configuration of the node
type NodeConfig interface {
	NodeNoNewHeadsThreshold() time.Duration
	NodePollFailureThreshold() uint32
	NodePollInterval() time.Duration
	NodeSelectionMode() string
	NodeSyncThreshold() uint32
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
	n.chStopInFlight = make(chan struct{})
	n.nodeCtx, n.cancelNodeCtx = context.WithCancel(context.Background())
	lggr = lggr.Named("Node").With(
		"nodeTier", "primary",
		"nodeName", name,
		"node", n.String(),
		"evmChainID", chainID,
	)
	n.lfcLog = lggr.Named("Lifecycle")
	n.rpcLog = lggr.Named("RPC")
	n.stateLatestBlockNumber = -1
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

	dialCtx, dialCancel := n.makeQueryCtx(startCtx)
	defer dialCancel()
	if err := n.dial(dialCtx); err != nil {
		n.lfcLog.Errorw("Dial failed: EVM Node is unreachable", "err", err)
		n.declareUnreachable()
		return
	}
	n.setState(NodeStateDialed)

	verifyCtx, verifyCancel := n.makeQueryCtx(startCtx)
	defer verifyCancel()
	if err := n.verify(verifyCtx); errors.Is(err, errInvalidChainID) {
		n.lfcLog.Errorw("Verify failed: EVM Node has the wrong chain ID", "err", err)
		n.declareInvalidChainID()
		return
	} else if err != nil {
		n.lfcLog.Errorw(fmt.Sprintf("Verify failed: %v", err), "err", err)
		n.declareUnreachable()
		return
	}

	n.declareAlive()
}

// Not thread-safe
// Pure dial: does not mutate node "state" field.
func (n *node) dial(callerCtx context.Context) error {
	ctx, cancel := n.makeQueryCtx(callerCtx)
	defer cancel()

	promEVMPoolRPCNodeDials.WithLabelValues(n.chainID.String(), n.name).Inc()
	lggr := n.lfcLog.With("wsuri", n.ws.uri.Redacted())
	if n.http != nil {
		lggr = lggr.With("httpuri", n.http.uri.Redacted())
	}
	lggr.Debugw("RPC dial: evmclient.Client#dial")

	wsrpc, err := rpc.DialWebsocket(ctx, n.ws.uri.String(), "")
	if err != nil {
		promEVMPoolRPCNodeDialsFailed.WithLabelValues(n.chainID.String(), n.name).Inc()
		return errors.Wrapf(err, "error while dialing websocket: %v", n.ws.uri.Redacted())
	}

	var httprpc *rpc.Client
	if n.http != nil {
		httprpc, err = rpc.DialHTTP(n.http.uri.String())
		if err != nil {
			promEVMPoolRPCNodeDialsFailed.WithLabelValues(n.chainID.String(), n.name).Inc()
			return errors.Wrapf(err, "error while dialing HTTP: %v", n.http.uri.Redacted())
		}
	}

	n.ws.rpc = wsrpc
	n.ws.geth = ethclient.NewClient(wsrpc)

	if n.http != nil {
		n.http.rpc = httprpc
		n.http.geth = ethclient.NewClient(httprpc)
	}

	promEVMPoolRPCNodeDialsSuccess.WithLabelValues(n.chainID.String(), n.name).Inc()

	return nil
}

var errInvalidChainID = errors.New("invalid chain id")

// verify checks that all connections to eth nodes match the given chain ID
// Not thread-safe
// Pure verify: does not mutate node "state" field.
func (n *node) verify(callerCtx context.Context) (err error) {
	ctx, cancel := n.makeQueryCtx(callerCtx)
	defer cancel()

	promEVMPoolRPCNodeVerifies.WithLabelValues(n.chainID.String(), n.name).Inc()
	promFailed := func() {
		promEVMPoolRPCNodeVerifiesFailed.WithLabelValues(n.chainID.String(), n.name).Inc()
	}

	st := n.State()
	switch st {
	case NodeStateDialed, NodeStateOutOfSync, NodeStateInvalidChainID:
	default:
		panic(fmt.Sprintf("cannot verify node in state %v", st))
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

func (n *node) Close() error {
	return n.StopOnce(n.name, func() error {
		defer func() {
			n.wg.Wait()
			if n.ws.rpc != nil {
				n.ws.rpc.Close()
			}
		}()

		n.stateMu.Lock()
		defer n.stateMu.Unlock()

		n.cancelNodeCtx()
		n.cancelInflightRequests()
		n.state = NodeStateClosed
		return nil
	})
}

// registerSub adds the sub to the node list
func (n *node) registerSub(sub ethereum.Subscription) {
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	n.subs = append(n.subs, sub)
}

// disconnectAll disconnects all clients connected to the node
// WARNING: NOT THREAD-SAFE
// This must be called from within the n.stateMu lock
func (n *node) disconnectAll() {
	if n.ws.rpc != nil {
		n.ws.rpc.Close()
	}
	n.cancelInflightRequests()
	n.unsubscribeAll()
}

// cancelInflightRequests closes and replaces the chStopInFlight
// WARNING: NOT THREAD-SAFE
// This must be called from within the n.stateMu lock
func (n *node) cancelInflightRequests() {
	close(n.chStopInFlight)
	n.chStopInFlight = make(chan struct{})
}

// unsubscribeAll unsubscribes all subscriptions
// WARNING: NOT THREAD-SAFE
// This must be called from within the n.stateMu lock
func (n *node) unsubscribeAll() {
	for _, sub := range n.subs {
		sub.Unsubscribe()
	}
	n.subs = nil
}

// getChStopInflight provides a convenience helper that mutex wraps a
// read to the chStopInFlight
func (n *node) getChStopInflight() chan struct{} {
	n.stateMu.RLock()
	defer n.stateMu.RUnlock()
	return n.chStopInFlight
}

func (n *node) getRPCDomain() string {
	if n.http != nil {
		return n.http.uri.Host
	}
	return n.ws.uri.Host
}

// RPC wrappers

// CallContext implementation
func (n *node) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	ctx, cancel, ws, http, err := n.makeLiveQueryCtx(ctx)
	if err != nil {
		return err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With(
		"method", method,
		"args", args,
	)

	lggr.Debug("RPC call: evmclient.Client#CallContext")
	start := time.Now()
	if http != nil {
		err = n.wrapHTTP(http.rpc.CallContext(ctx, result, method, args...))
	} else {
		err = n.wrapWS(ws.rpc.CallContext(ctx, result, method, args...))
	}
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "CallContext")

	return err
}

func (n *node) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	ctx, cancel, ws, http, err := n.makeLiveQueryCtx(ctx)
	if err != nil {
		return err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("nBatchElems", len(b), "batchElems", b)

	lggr.Debug("RPC call: evmclient.Client#BatchCallContext")
	start := time.Now()
	if http != nil {
		err = n.wrapHTTP(http.rpc.BatchCallContext(ctx, b))
	} else {
		err = n.wrapWS(ws.rpc.BatchCallContext(ctx, b))
	}
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "BatchCallContext")

	return err
}

func (n *node) EthSubscribe(ctx context.Context, channel chan<- *evmtypes.Head, args ...interface{}) (ethereum.Subscription, error) {
	ctx, cancel, ws, _, err := n.makeLiveQueryCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr("websocket").With("args", args)

	lggr.Debug("RPC call: evmclient.Client#EthSubscribe")
	start := time.Now()
	sub, err := ws.rpc.EthSubscribe(ctx, channel, args...)
	if err == nil {
		n.registerSub(sub)
	}
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "EthSubscribe")

	return sub, err
}

// GethClient wrappers

func (n *node) TransactionReceipt(ctx context.Context, txHash common.Hash) (receipt *types.Receipt, err error) {
	ctx, cancel, ws, http, err := n.makeLiveQueryCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("txHash", txHash)

	lggr.Debug("RPC call: evmclient.Client#TransactionReceipt")

	start := time.Now()
	if http != nil {
		receipt, err = http.geth.TransactionReceipt(ctx, txHash)
		err = n.wrapHTTP(err)
	} else {
		receipt, err = ws.geth.TransactionReceipt(ctx, txHash)
		err = n.wrapWS(err)
	}
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "TransactionReceipt",
		"receipt", receipt,
	)

	return
}

func (n *node) HeaderByNumber(ctx context.Context, number *big.Int) (header *types.Header, err error) {
	ctx, cancel, ws, http, err := n.makeLiveQueryCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("number", number)

	lggr.Debug("RPC call: evmclient.Client#HeaderByNumber")
	start := time.Now()
	if http != nil {
		header, err = http.geth.HeaderByNumber(ctx, number)
		err = n.wrapHTTP(err)
	} else {
		header, err = ws.geth.HeaderByNumber(ctx, number)
		err = n.wrapWS(err)
	}
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "HeaderByNumber", "header", header)

	return
}

func (n *node) HeaderByHash(ctx context.Context, hash common.Hash) (header *types.Header, err error) {
	ctx, cancel, ws, http, err := n.makeLiveQueryCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("hash", hash)

	lggr.Debug("RPC call: evmclient.Client#HeaderByHash")
	start := time.Now()
	if http != nil {
		header, err = http.geth.HeaderByHash(ctx, hash)
		err = n.wrapHTTP(err)
	} else {
		header, err = ws.geth.HeaderByHash(ctx, hash)
		err = n.wrapWS(err)
	}
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "HeaderByHash",
		"header", header,
	)

	return
}

func (n *node) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	ctx, cancel, ws, http, err := n.makeLiveQueryCtx(ctx)
	if err != nil {
		return err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("tx", tx)

	lggr.Debug("RPC call: evmclient.Client#SendTransaction")
	start := time.Now()
	if http != nil {
		err = n.wrapHTTP(http.geth.SendTransaction(ctx, tx))
	} else {
		err = n.wrapWS(ws.geth.SendTransaction(ctx, tx))
	}
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "SendTransaction")

	return err
}

// PendingNonceAt returns one higher than the highest nonce from both mempool and mined transactions
func (n *node) PendingNonceAt(ctx context.Context, account common.Address) (nonce uint64, err error) {
	ctx, cancel, ws, http, err := n.makeLiveQueryCtx(ctx)
	if err != nil {
		return 0, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("account", account)

	lggr.Debug("RPC call: evmclient.Client#PendingNonceAt")
	start := time.Now()
	if http != nil {
		nonce, err = http.geth.PendingNonceAt(ctx, account)
		err = n.wrapHTTP(err)
	} else {
		nonce, err = ws.geth.PendingNonceAt(ctx, account)
		err = n.wrapWS(err)
	}
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "PendingNonceAt",
		"nonce", nonce,
	)

	return
}

// NonceAt is a bit of a misnomer. You might expect it to return the highest
// mined nonce at the given block number, but it actually returns the total
// transaction count which is the highest mined nonce + 1
func (n *node) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (nonce uint64, err error) {
	ctx, cancel, ws, http, err := n.makeLiveQueryCtx(ctx)
	if err != nil {
		return 0, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("account", account, "blockNumber", blockNumber)

	lggr.Debug("RPC call: evmclient.Client#NonceAt")
	start := time.Now()
	if http != nil {
		nonce, err = http.geth.NonceAt(ctx, account, blockNumber)
		err = n.wrapHTTP(err)
	} else {
		nonce, err = ws.geth.NonceAt(ctx, account, blockNumber)
		err = n.wrapWS(err)
	}
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "NonceAt",
		"nonce", nonce,
	)

	return
}

func (n *node) PendingCodeAt(ctx context.Context, account common.Address) (code []byte, err error) {
	ctx, cancel, ws, http, err := n.makeLiveQueryCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("account", account)

	lggr.Debug("RPC call: evmclient.Client#PendingCodeAt")
	start := time.Now()
	if http != nil {
		code, err = http.geth.PendingCodeAt(ctx, account)
		err = n.wrapHTTP(err)
	} else {
		code, err = ws.geth.PendingCodeAt(ctx, account)
		err = n.wrapWS(err)
	}
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "PendingCodeAt",
		"code", code,
	)

	return
}

func (n *node) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) (code []byte, err error) {
	ctx, cancel, ws, http, err := n.makeLiveQueryCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("account", account, "blockNumber", blockNumber)

	lggr.Debug("RPC call: evmclient.Client#CodeAt")
	start := time.Now()
	if http != nil {
		code, err = http.geth.CodeAt(ctx, account, blockNumber)
		err = n.wrapHTTP(err)
	} else {
		code, err = ws.geth.CodeAt(ctx, account, blockNumber)
		err = n.wrapWS(err)
	}
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "CodeAt",
		"code", code,
	)

	return
}

func (n *node) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	ctx, cancel, ws, http, err := n.makeLiveQueryCtx(ctx)
	if err != nil {
		return 0, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("call", call)

	lggr.Debug("RPC call: evmclient.Client#EstimateGas")
	start := time.Now()
	if http != nil {
		gas, err = http.geth.EstimateGas(ctx, call)
		err = n.wrapHTTP(err)
	} else {
		gas, err = ws.geth.EstimateGas(ctx, call)
		err = n.wrapWS(err)
	}
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "EstimateGas",
		"gas", gas,
	)

	return
}

func (n *node) SuggestGasPrice(ctx context.Context) (price *big.Int, err error) {
	ctx, cancel, ws, http, err := n.makeLiveQueryCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n))

	lggr.Debug("RPC call: evmclient.Client#SuggestGasPrice")
	start := time.Now()
	if http != nil {
		price, err = http.geth.SuggestGasPrice(ctx)
		err = n.wrapHTTP(err)
	} else {
		price, err = ws.geth.SuggestGasPrice(ctx)
		err = n.wrapWS(err)
	}
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "SuggestGasPrice",
		"price", price,
	)

	return
}

func (n *node) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) (val []byte, err error) {
	ctx, cancel, ws, http, err := n.makeLiveQueryCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("callMsg", msg, "blockNumber", blockNumber)

	lggr.Debug("RPC call: evmclient.Client#CallContract")
	start := time.Now()
	if http != nil {
		val, err = http.geth.CallContract(ctx, msg, blockNumber)
		err = n.wrapHTTP(err)
	} else {
		val, err = ws.geth.CallContract(ctx, msg, blockNumber)
		err = n.wrapWS(err)
	}
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "CallContract",
		"val", val,
	)

	return

}

func (n *node) BlockByNumber(ctx context.Context, number *big.Int) (b *types.Block, err error) {
	ctx, cancel, ws, http, err := n.makeLiveQueryCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("number", number)

	lggr.Debug("RPC call: evmclient.Client#BlockByNumber")
	start := time.Now()
	if http != nil {
		b, err = http.geth.BlockByNumber(ctx, number)
		err = n.wrapHTTP(err)
	} else {
		b, err = ws.geth.BlockByNumber(ctx, number)
		err = n.wrapWS(err)
	}
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "BlockByNumber",
		"block", b,
	)

	return
}

func (n *node) BlockByHash(ctx context.Context, hash common.Hash) (b *types.Block, err error) {
	ctx, cancel, ws, http, err := n.makeLiveQueryCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("hash", hash)

	lggr.Debug("RPC call: evmclient.Client#BlockByHash")
	start := time.Now()
	if http != nil {
		b, err = http.geth.BlockByHash(ctx, hash)
		err = n.wrapHTTP(err)
	} else {
		b, err = ws.geth.BlockByHash(ctx, hash)
		err = n.wrapWS(err)
	}
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "BlockByHash",
		"block", b,
	)

	return
}

func (n *node) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (balance *big.Int, err error) {
	ctx, cancel, ws, http, err := n.makeLiveQueryCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("account", account.Hex(), "blockNumber", blockNumber)

	lggr.Debug("RPC call: evmclient.Client#BalanceAt")
	start := time.Now()
	if http != nil {
		balance, err = http.geth.BalanceAt(ctx, account, blockNumber)
		err = n.wrapHTTP(err)
	} else {
		balance, err = ws.geth.BalanceAt(ctx, account, blockNumber)
		err = n.wrapWS(err)
	}
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "BalanceAt",
		"balance", balance,
	)

	return
}

func (n *node) FilterLogs(ctx context.Context, q ethereum.FilterQuery) (l []types.Log, err error) {
	ctx, cancel, ws, http, err := n.makeLiveQueryCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With("q", q)

	lggr.Debug("RPC call: evmclient.Client#FilterLogs")
	start := time.Now()
	if http != nil {
		l, err = http.geth.FilterLogs(ctx, q)
		err = n.wrapHTTP(err)
	} else {
		l, err = ws.geth.FilterLogs(ctx, q)
		err = n.wrapWS(err)
	}
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "FilterLogs",
		"log", l,
	)

	return
}

func (n *node) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (sub ethereum.Subscription, err error) {
	ctx, cancel, ws, _, err := n.makeLiveQueryCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr("websocket").With("q", q)

	lggr.Debug("RPC call: evmclient.Client#SubscribeFilterLogs")
	start := time.Now()
	sub, err = ws.geth.SubscribeFilterLogs(ctx, q, ch)
	if err == nil {
		n.registerSub(sub)
	}
	err = n.wrapWS(err)
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "SubscribeFilterLogs")

	return
}

func (n *node) SuggestGasTipCap(ctx context.Context) (tipCap *big.Int, err error) {
	ctx, cancel, ws, http, err := n.makeLiveQueryCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n))

	lggr.Debug("RPC call: evmclient.Client#SuggestGasTipCap")
	start := time.Now()
	if http != nil {
		tipCap, err = http.geth.SuggestGasTipCap(ctx)
		err = n.wrapHTTP(err)
	} else {
		tipCap, err = ws.geth.SuggestGasTipCap(ctx)
		err = n.wrapWS(err)
	}
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "SuggestGasTipCap",
		"tipCap", tipCap,
	)

	return
}

func (n *node) ChainID() (chainID *big.Int) { return n.chainID }

// newRqLggr generates a new logger with a unique request ID
func (n *node) newRqLggr(mode string) logger.Logger {
	return n.rpcLog.With(
		"requestID", uuid.NewV4(),
		"mode", mode,
	)
}

func (n *node) logResult(
	lggr logger.Logger,
	err error,
	callDuration time.Duration,
	rpcDomain,
	callName string,
	results ...interface{},
) {
	lggr = lggr.With("duration", callDuration, "rpcDomain", rpcDomain, "callName", callName)
	promEVMPoolRPCNodeCalls.WithLabelValues(n.chainID.String(), n.name).Inc()
	if err == nil {
		promEVMPoolRPCNodeCallsSuccess.WithLabelValues(n.chainID.String(), n.name).Inc()
		lggr.Debugw(
			fmt.Sprintf("evmclient.Client#%s RPC call success", callName),
			results...,
		)
	} else {
		promEVMPoolRPCNodeCallsFailed.WithLabelValues(n.chainID.String(), n.name).Inc()
		lggr.Debugw(
			fmt.Sprintf("evmclient.Client#%s RPC call failure", callName),
			append(results, "err", err)...,
		)
	}
	promEVMPoolRPCCallTiming.
		WithLabelValues(
			n.chainID.String(),             // chain id
			n.name,                         // node name
			rpcDomain,                      // rpc domain
			"false",                        // is send only
			strconv.FormatBool(err == nil), // is successful
			callName,                       // rpc call name
		).
		Observe(float64(callDuration))
}

func (n *node) wrapWS(err error) error {
	err = wrap(err, fmt.Sprintf("primary websocket (%s)", n.ws.uri.Redacted()))
	return err
}

func (n *node) wrapHTTP(err error) error {
	err = wrap(err, fmt.Sprintf("primary http (%s)", n.http.uri.Redacted()))
	if err != nil {
		n.rpcLog.Debugw("Call failed", "err", err)
	} else {
		n.rpcLog.Trace("Call succeeded")
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

// makeLiveQueryCtx wraps makeQueryCtx but returns error if node is not NodeStateAlive.
func (n *node) makeLiveQueryCtx(parentCtx context.Context) (ctx context.Context, cancel context.CancelFunc, ws rawclient, http *rawclient, err error) {
	// Need to wrap in mutex because state transition can cancel and replace the
	// context
	n.stateMu.RLock()
	if n.state != NodeStateAlive {
		err = errors.Errorf("cannot execute RPC call on node with state: %s", n.state)
		n.stateMu.RUnlock()
		return
	}
	cancelCh := n.chStopInFlight
	ws = n.ws
	if n.http != nil {
		cp := *n.http
		http = &cp
	}
	n.stateMu.RUnlock()
	ctx, cancel = makeQueryCtx(parentCtx, cancelCh)
	return
}

func (n *node) makeQueryCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	return makeQueryCtx(ctx, n.getChStopInflight())
}

// makeQueryCtx returns a context that cancels if:
// 1. Passed in ctx cancels
// 2. Passed in channel is closed
// 3. Default timeout is reached (queryTimeout)
func makeQueryCtx(ctx context.Context, ch chan struct{}) (context.Context, context.CancelFunc) {
	var chCancel, timeoutCancel context.CancelFunc
	ctx, chCancel = utils.WithCloseChan(ctx, ch)
	ctx, timeoutCancel = context.WithTimeout(ctx, queryTimeout)
	cancel := func() {
		chCancel()
		timeoutCancel()
	}
	return ctx, cancel
}

func switching(n *node) string {
	if n.http != nil {
		return "http"
	}
	return "websocket"
}

func (n *node) String() string {
	s := fmt.Sprintf("(primary)%s:%s", n.name, n.ws.uri.Redacted())
	if n.http != nil {
		s = s + fmt.Sprintf(":%s", n.http.uri.Redacted())
	}
	return s
}

func (n *node) Name() string {
	return n.name
}
