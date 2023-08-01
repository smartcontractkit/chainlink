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
	"github.com/google/uuid"
	"github.com/pkg/errors"
	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type rpcClient struct {
	utils.StartStopOnce
	lfcLog              logger.Logger
	rpcLog              logger.Logger
	name                string
	id                  int32
	chainID             *big.Int
	nodePoolCfg         config.NodePool
	noNewHeadsThreshold time.Duration
	order               int32

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

// NewRPCCLient returns a new *rpcClient as multinodeclient.RPCClient
func NewRPCClient(nodeCfg config.NodePool, noNewHeadsThreshold time.Duration, lggr logger.Logger, wsuri url.URL, httpuri *url.URL, name string, id int32, chainID *big.Int, nodeOrder int32) commonclient.ChainRPCClient[
	*big.Int,
	evmtypes.Nonce,
	common.Address,
	types.Block,
	common.Hash,
	evmtypes.Transaction,
	common.Hash,
	types.Log,
	ethereum.FilterQuery,
	*evmtypes.Receipt,
	evmtypes.EvmFee,
	*evmtypes.Head,
] {
	r := new(rpcClient)
	r.name = name
	r.id = id
	r.chainID = chainID
	r.nodePoolCfg = nodeCfg
	r.noNewHeadsThreshold = noNewHeadsThreshold
	r.ws.uri = wsuri
	r.order = nodeOrder
	if httpuri != nil {
		r.http = &rawclient{uri: *httpuri}
	}
	r.chStopInFlight = make(chan struct{})
	r.nodeCtx, r.cancelNodeCtx = context.WithCancel(context.Background())
	lggr = lggr.Named("Client").With(
		"clientTier", "primary",
		"clientName", name,
		"client", r.String(),
		"evmChainID", chainID,
		"clientOrder", r.order,
	)
	r.lfcLog = lggr.Named("Lifecycle")
	r.rpcLog = lggr.Named("RPC")
	r.stateLatestBlockNumber = -1

	return r
}

// makeLiveQueryCtx wraps makeQueryCtx but returns error if node is not NodeStateAlive.
func (r *rpcClient) makeLiveQueryCtx(parentCtx context.Context) (ctx context.Context, cancel context.CancelFunc, ws rawclient, http *rawclient, err error) {
	// Need to wrap in mutex because state transition can cancel and replace the
	// context
	r.stateMu.RLock()
	if r.state != NodeStateAlive {
		err = errors.Errorf("cannot execute RPC call on node with state: %s", r.state)
		r.stateMu.RUnlock()
		return
	}
	cancelCh := r.chStopInFlight
	ws = r.ws
	if r.http != nil {
		cp := *r.http
		http = &cp
	}
	r.stateMu.RUnlock()
	ctx, cancel = makeQueryCtx(parentCtx, cancelCh)
	return
}

// Not thread-safe
// Pure dial: does not mutate node "state" field.
func (r *rpcClient) Dial(callerCtx context.Context) error {
	ctx, cancel := r.makeQueryCtx(callerCtx)
	defer cancel()

	promEVMPoolRPCNodeDials.WithLabelValues(r.chainID.String(), r.name).Inc()
	lggr := r.lfcLog.With("wsuri", r.ws.uri.Redacted())
	if r.http != nil {
		lggr = lggr.With("httpuri", r.http.uri.Redacted())
	}
	lggr.Debugw("RPC dial: evmclient.Client#dial")

	wsrpc, err := rpc.DialWebsocket(ctx, r.ws.uri.String(), "")
	if err != nil {
		promEVMPoolRPCNodeDialsFailed.WithLabelValues(r.chainID.String(), r.name).Inc()
		return errors.Wrapf(err, "error while dialing websocket: %v", r.ws.uri.Redacted())
	}

	var httprpc *rpc.Client
	if r.http != nil {
		httprpc, err = rpc.DialHTTP(r.http.uri.String())
		if err != nil {
			promEVMPoolRPCNodeDialsFailed.WithLabelValues(r.chainID.String(), r.name).Inc()
			return errors.Wrapf(err, "error while dialing HTTP: %v", r.http.uri.Redacted())
		}
	}

	r.ws.rpc = wsrpc
	r.ws.geth = ethclient.NewClient(wsrpc)

	if r.http != nil {
		r.http.rpc = httprpc
		r.http.geth = ethclient.NewClient(httprpc)
	}

	promEVMPoolRPCNodeDialsSuccess.WithLabelValues(r.chainID.String(), r.name).Inc()

	return nil
}

func (r *rpcClient) Close() error {
	return r.StopOnce(r.name, func() error {
		defer func() {
			r.wg.Wait()
			if r.ws.rpc != nil {
				r.ws.rpc.Close()
			}
		}()

		r.stateMu.Lock()
		defer r.stateMu.Unlock()

		r.cancelNodeCtx()
		r.cancelInflightRequests()
		r.state = NodeStateClosed
		return nil
	})
}

// cancelInflightRequests closes and replaces the chStopInFlight
// WARNING: NOT THREAD-SAFE
// This must be called from within the r.stateMu lock
func (r *rpcClient) cancelInflightRequests() {
	close(r.chStopInFlight)
	r.chStopInFlight = make(chan struct{})
}

func (r *rpcClient) String() string {
	s := fmt.Sprintf("(primary)%s:%s", r.name, r.ws.uri.Redacted())
	if r.http != nil {
		s = s + fmt.Sprintf(":%s", r.http.uri.Redacted())
	}
	return s
}

func (r *rpcClient) logResult(
	lggr logger.Logger,
	err error,
	callDuration time.Duration,
	rpcDomain,
	callName string,
	results ...interface{},
) {
	lggr = lggr.With("duration", callDuration, "rpcDomain", rpcDomain, "callName", callName)
	promEVMPoolRPCNodeCalls.WithLabelValues(r.chainID.String(), r.name).Inc()
	if err == nil {
		promEVMPoolRPCNodeCallsSuccess.WithLabelValues(r.chainID.String(), r.name).Inc()
		lggr.Tracew(
			fmt.Sprintf("evmclient.Client#%s RPC call success", callName),
			results...,
		)
	} else {
		promEVMPoolRPCNodeCallsFailed.WithLabelValues(r.chainID.String(), r.name).Inc()
		lggr.Debugw(
			fmt.Sprintf("evmclient.Client#%s RPC call failure", callName),
			append(results, "err", err)...,
		)
	}
	promEVMPoolRPCCallTiming.
		WithLabelValues(
			r.chainID.String(),             // chain id
			r.name,                         // node name
			rpcDomain,                      // rpc domain
			"false",                        // is send only
			strconv.FormatBool(err == nil), // is successful
			callName,                       // rpc call name
		).
		Observe(float64(callDuration))
}

func switching(r *rpcClient) string {
	if r.http != nil {
		return "http"
	}
	return "websocket"
}

func (r *rpcClient) getRPCDomain() string {
	if r.http != nil {
		return r.http.uri.Host
	}
	return r.ws.uri.Host
}

// registerSub adds the sub to the node list
func (r *rpcClient) registerSub(sub ethereum.Subscription) {
	r.stateMu.Lock()
	defer r.stateMu.Unlock()
	r.subs = append(r.subs, sub)
}

// disconnectAll disconnects all clients connected to the node
// WARNING: NOT THREAD-SAFE
// This must be called from within the r.stateMu lock
func (r *rpcClient) DisconnectAll() {
	if r.ws.rpc != nil {
		r.ws.rpc.Close()
	}
	r.cancelInflightRequests()
	r.unsubscribeAll()
}

// unsubscribeAll unsubscribes all subscriptions
// WARNING: NOT THREAD-SAFE
// This must be called from within the r.stateMu lock
func (r *rpcClient) unsubscribeAll() {
	for _, sub := range r.subs {
		sub.Unsubscribe()
	}
	r.subs = nil
}

// RPC wrappers

// CallContext implementation
func (r *rpcClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return err
	}
	defer cancel()
	lggr := r.newRqLggr().With(
		"method", method,
		"args", args,
	)

	lggr.Debug("RPC call: evmclient.Client#CallContext")
	start := time.Now()
	if http != nil {
		err = r.wrapHTTP(http.rpc.CallContext(ctx, result, method, args...))
	} else {
		err = r.wrapWS(ws.rpc.CallContext(ctx, result, method, args...))
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "CallContext")

	return err
}

func (r *rpcClient) BatchCallContext(ctx context.Context, b []interface{}) error {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return err
	}
	// Kind of hacky: not sure if this is the best solution.
	batch := make([]rpc.BatchElem, len(b))
	for i, arg := range b {
		batch[i] = arg.(rpc.BatchElem)
	}
	defer cancel()
	lggr := r.newRqLggr().With("nBatchElems", len(b), "batchElems", b)

	lggr.Trace("RPC call: evmclient.Client#BatchCallContext")
	start := time.Now()
	if http != nil {
		err = r.wrapHTTP(http.rpc.BatchCallContext(ctx, batch))
	} else {
		err = r.wrapWS(ws.rpc.BatchCallContext(ctx, batch))
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "BatchCallContext")

	return err
}

func (r *rpcClient) Subscribe(ctx context.Context, channel chan<- *evmtypes.Head, args ...interface{}) (ethereum.Subscription, error) {
	ctx, cancel, ws, _, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := r.newRqLggr().With("args", args)

	lggr.Debug("RPC call: evmclient.Client#EthSubscribe")
	start := time.Now()
	sub, err := ws.rpc.EthSubscribe(ctx, channel, args...)
	if err == nil {
		r.registerSub(sub)
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "EthSubscribe")

	return sub, err
}

// GethClient wrappers

func (r *rpcClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (receipt *types.Receipt, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := r.newRqLggr().With("txHash", txHash)

	lggr.Debug("RPC call: evmclient.Client#TransactionReceipt")

	start := time.Now()
	if http != nil {
		receipt, err = http.geth.TransactionReceipt(ctx, txHash)
		err = r.wrapHTTP(err)
	} else {
		receipt, err = ws.geth.TransactionReceipt(ctx, txHash)
		err = r.wrapWS(err)
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "TransactionReceipt",
		"receipt", receipt,
	)

	return
}

func (r *rpcClient) TransactionByHash(ctx context.Context, txHash common.Hash) (tx *types.Transaction, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := r.newRqLggr().With("txHash", txHash)

	lggr.Debug("RPC call: evmclient.Client#TransactionByHash")

	start := time.Now()
	if http != nil {
		tx, _, err = http.geth.TransactionByHash(ctx, txHash)
		err = r.wrapHTTP(err)
	} else {
		tx, _, err = ws.geth.TransactionByHash(ctx, txHash)
		err = r.wrapWS(err)
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "TransactionByHash",
		"receipt", tx,
	)

	return
}

func (r *rpcClient) HeaderByNumber(ctx context.Context, number *big.Int) (header *types.Header, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := r.newRqLggr().With("number", number)

	lggr.Debug("RPC call: evmclient.Client#HeaderByNumber")
	start := time.Now()
	if http != nil {
		header, err = http.geth.HeaderByNumber(ctx, number)
		err = r.wrapHTTP(err)
	} else {
		header, err = ws.geth.HeaderByNumber(ctx, number)
		err = r.wrapWS(err)
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "HeaderByNumber", "header", header)

	return
}

func (r *rpcClient) HeaderByHash(ctx context.Context, hash common.Hash) (header *types.Header, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := r.newRqLggr().With("hash", hash)

	lggr.Debug("RPC call: evmclient.Client#HeaderByHash")
	start := time.Now()
	if http != nil {
		header, err = http.geth.HeaderByHash(ctx, hash)
		err = r.wrapHTTP(err)
	} else {
		header, err = ws.geth.HeaderByHash(ctx, hash)
		err = r.wrapWS(err)
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "HeaderByHash",
		"header", header,
	)

	return
}

func (r *rpcClient) HeadByNumber(ctx context.Context, number *big.Int) (head *evmtypes.Head, err error) {
	hex := ToBlockNumArg(number)
	err = r.CallContext(ctx, &head, "eth_getBlockByNumber", hex, false)
	if err != nil {
		return nil, err
	}
	if head == nil {
		err = ethereum.NotFound
		return
	}
	head.EVMChainID = utils.NewBig(r.ConfiguredChainID())
	return
}

func (r *rpcClient) HeadByHash(ctx context.Context, hash common.Hash) (head *evmtypes.Head, err error) {
	err = r.CallContext(ctx, &head, "eth_getBlockByHash", hash.Hex(), false)
	if err != nil {
		return nil, err
	}
	if head == nil {
		err = ethereum.NotFound
		return
	}
	head.EVMChainID = utils.NewBig(r.ConfiguredChainID())
	return
}

func (r *rpcClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return err
	}
	defer cancel()
	lggr := r.newRqLggr().With("tx", tx)

	lggr.Debug("RPC call: evmclient.Client#SendTransaction")
	start := time.Now()
	if http != nil {
		err = r.wrapHTTP(http.geth.SendTransaction(ctx, tx))
	} else {
		err = r.wrapWS(ws.geth.SendTransaction(ctx, tx))
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "SendTransaction")

	return err
}

// PendingNonceAt returns one higher than the highest nonce from both mempool and mined transactions
func (r *rpcClient) PendingNonceAt(ctx context.Context, account common.Address) (nonce uint64, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return 0, err
	}
	defer cancel()
	lggr := r.newRqLggr().With("account", account)

	lggr.Debug("RPC call: evmclient.Client#PendingNonceAt")
	start := time.Now()
	if http != nil {
		nonce, err = http.geth.PendingNonceAt(ctx, account)
		err = r.wrapHTTP(err)
	} else {
		nonce, err = ws.geth.PendingNonceAt(ctx, account)
		err = r.wrapWS(err)
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "PendingNonceAt",
		"nonce", nonce,
	)

	return
}

// NonceAt is a bit of a misnomer. You might expect it to return the highest
// mined nonce at the given block number, but it actually returns the total
// transaction count which is the highest mined nonce + 1
func (r *rpcClient) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (nonce uint64, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return 0, err
	}
	defer cancel()
	lggr := r.newRqLggr().With("account", account, "blockNumber", blockNumber)

	lggr.Debug("RPC call: evmclient.Client#NonceAt")
	start := time.Now()
	if http != nil {
		nonce, err = http.geth.NonceAt(ctx, account, blockNumber)
		err = r.wrapHTTP(err)
	} else {
		nonce, err = ws.geth.NonceAt(ctx, account, blockNumber)
		err = r.wrapWS(err)
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "NonceAt",
		"nonce", nonce,
	)

	return
}

func (r *rpcClient) PendingCodeAt(ctx context.Context, account common.Address) (code []byte, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := r.newRqLggr().With("account", account)

	lggr.Debug("RPC call: evmclient.Client#PendingCodeAt")
	start := time.Now()
	if http != nil {
		code, err = http.geth.PendingCodeAt(ctx, account)
		err = r.wrapHTTP(err)
	} else {
		code, err = ws.geth.PendingCodeAt(ctx, account)
		err = r.wrapWS(err)
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "PendingCodeAt",
		"code", code,
	)

	return
}

func (r *rpcClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) (code []byte, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := r.newRqLggr().With("account", account, "blockNumber", blockNumber)

	lggr.Debug("RPC call: evmclient.Client#CodeAt")
	start := time.Now()
	if http != nil {
		code, err = http.geth.CodeAt(ctx, account, blockNumber)
		err = r.wrapHTTP(err)
	} else {
		code, err = ws.geth.CodeAt(ctx, account, blockNumber)
		err = r.wrapWS(err)
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "CodeAt",
		"code", code,
	)

	return
}

func (r *rpcClient) EstimateGas(ctx context.Context, c interface{}) (gas uint64, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return 0, err
	}
	defer cancel()
	call := c.(ethereum.CallMsg)
	lggr := r.newRqLggr().With("call", call)

	lggr.Debug("RPC call: evmclient.Client#EstimateGas")
	start := time.Now()
	if http != nil {
		gas, err = http.geth.EstimateGas(ctx, call)
		err = r.wrapHTTP(err)
	} else {
		gas, err = ws.geth.EstimateGas(ctx, call)
		err = r.wrapWS(err)
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "EstimateGas",
		"gas", gas,
	)

	return
}

func (r *rpcClient) SuggestGasPrice(ctx context.Context) (price *big.Int, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := r.newRqLggr()

	lggr.Debug("RPC call: evmclient.Client#SuggestGasPrice")
	start := time.Now()
	if http != nil {
		price, err = http.geth.SuggestGasPrice(ctx)
		err = r.wrapHTTP(err)
	} else {
		price, err = ws.geth.SuggestGasPrice(ctx)
		err = r.wrapWS(err)
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "SuggestGasPrice",
		"price", price,
	)

	return
}

func (r *rpcClient) CallContract(ctx context.Context, msg interface{}, blockNumber *big.Int) (val []byte, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := r.newRqLggr().With("callMsg", msg, "blockNumber", blockNumber)
	message := msg.(ethereum.CallMsg)

	lggr.Debug("RPC call: evmclient.Client#CallContract")
	start := time.Now()
	if http != nil {
		val, err = http.geth.CallContract(ctx, message, blockNumber)
		err = r.wrapHTTP(err)
	} else {
		val, err = ws.geth.CallContract(ctx, message, blockNumber)
		err = r.wrapWS(err)
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "CallContract",
		"val", val,
	)

	return

}

func (r *rpcClient) BlockByNumber(ctx context.Context, number *big.Int) (b *types.Block, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := r.newRqLggr().With("number", number)

	lggr.Debug("RPC call: evmclient.Client#BlockByNumber")
	start := time.Now()
	if http != nil {
		b, err = http.geth.BlockByNumber(ctx, number)
		err = r.wrapHTTP(err)
	} else {
		b, err = ws.geth.BlockByNumber(ctx, number)
		err = r.wrapWS(err)
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "BlockByNumber",
		"block", b,
	)

	return
}

func (r *rpcClient) BlockByHash(ctx context.Context, hash common.Hash) (b *types.Block, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := r.newRqLggr().With("hash", hash)

	lggr.Debug("RPC call: evmclient.Client#BlockByHash")
	start := time.Now()
	if http != nil {
		b, err = http.geth.BlockByHash(ctx, hash)
		err = r.wrapHTTP(err)
	} else {
		b, err = ws.geth.BlockByHash(ctx, hash)
		err = r.wrapWS(err)
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "BlockByHash",
		"block", b,
	)

	return
}

func (r *rpcClient) BlockNumber(ctx context.Context) (height uint64, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return 0, err
	}
	defer cancel()
	lggr := r.newRqLggr()

	lggr.Debug("RPC call: evmclient.Client#BlockNumber")
	start := time.Now()
	if http != nil {
		height, err = http.geth.BlockNumber(ctx)
		err = r.wrapHTTP(err)
	} else {
		height, err = ws.geth.BlockNumber(ctx)
		err = r.wrapWS(err)
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "BlockNumber",
		"height", height,
	)

	return
}

func (r *rpcClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (balance *big.Int, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := r.newRqLggr().With("account", account.Hex(), "blockNumber", blockNumber)

	lggr.Debug("RPC call: evmclient.Client#BalanceAt")
	start := time.Now()
	if http != nil {
		balance, err = http.geth.BalanceAt(ctx, account, blockNumber)
		err = r.wrapHTTP(err)
	} else {
		balance, err = ws.geth.BalanceAt(ctx, account, blockNumber)
		err = r.wrapWS(err)
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "BalanceAt",
		"balance", balance,
	)

	return
}

func (r *rpcClient) FilterEvents(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return r.FilterLogs(ctx, q)
}

func (r *rpcClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) (l []types.Log, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := r.newRqLggr().With("q", q)

	lggr.Debug("RPC call: evmclient.Client#FilterLogs")
	start := time.Now()
	if http != nil {
		l, err = http.geth.FilterLogs(ctx, q)
		err = r.wrapHTTP(err)
	} else {
		l, err = ws.geth.FilterLogs(ctx, q)
		err = r.wrapWS(err)
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "FilterLogs",
		"log", l,
	)

	return
}

func (r *rpcClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (sub ethereum.Subscription, err error) {
	ctx, cancel, ws, _, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := r.newRqLggr().With("q", q)

	lggr.Debug("RPC call: evmclient.Client#SubscribeFilterLogs")
	start := time.Now()
	sub, err = ws.geth.SubscribeFilterLogs(ctx, q, ch)
	if err == nil {
		r.registerSub(sub)
	}
	err = r.wrapWS(err)
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "SubscribeFilterLogs")

	return
}

func (r *rpcClient) SuggestGasTipCap(ctx context.Context) (tipCap *big.Int, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := r.newRqLggr()

	lggr.Debug("RPC call: evmclient.Client#SuggestGasTipCap")
	start := time.Now()
	if http != nil {
		tipCap, err = http.geth.SuggestGasTipCap(ctx)
		err = r.wrapHTTP(err)
	} else {
		tipCap, err = ws.geth.SuggestGasTipCap(ctx)
		err = r.wrapWS(err)
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "SuggestGasTipCap",
		"tipCap", tipCap,
	)

	return
}

func (r *rpcClient) ChainID() (chainID *big.Int, err error) { return r.chainID, nil }

func (r *rpcClient) ConfiguredChainID() *big.Int {
	chainID, err := r.ChainID()
	if err != nil {
		return chainID
	}
}

// newRqLggr generates a new logger with a unique request ID
func (r *rpcClient) newRqLggr() logger.Logger {
	return r.rpcLog.With(
		"requestID", uuid.New(),
	)
}

func (r *rpcClient) wrapWS(err error) error {
	err = wrap(err, fmt.Sprintf("primary websocket (%s)", r.ws.uri.Redacted()))
	return err
}

func (r *rpcClient) wrapHTTP(err error) error {
	err = wrap(err, fmt.Sprintf("primary http (%s)", r.http.uri.Redacted()))
	if err != nil {
		r.rpcLog.Debugw("Call failed", "err", err)
	} else {
		r.rpcLog.Trace("Call succeeded")
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

// makeLiveQueryCtxAndSafeGetClients wraps makeQueryCtx but returns error if node is not NodeStateAlive.
func (r *rpcClient) makeLiveQueryCtxAndSafeGetClients(parentCtx context.Context) (ctx context.Context, cancel context.CancelFunc, ws rawclient, http *rawclient, err error) {
	// Need to wrap in mutex because state transition can cancel and replace the
	// context
	r.stateMu.RLock()
	if r.state != NodeStateAlive {
		err = errors.Errorf("cannot execute RPC call on node with state: %s", r.state)
		r.stateMu.RUnlock()
		return
	}
	cancelCh := r.chStopInFlight
	ws = r.ws
	if r.http != nil {
		cp := *r.http
		http = &cp
	}
	r.stateMu.RUnlock()
	ctx, cancel = makeQueryCtx(parentCtx, cancelCh)
	return
}

func (r *rpcClient) makeQueryCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	return makeQueryCtx(ctx, r.getChStopInflight())
}

// getChStopInflight provides a convenience helper that mutex wraps a
// read to the chStopInFlight
func (r *rpcClient) getChStopInflight() chan struct{} {
	r.stateMu.RLock()
	defer r.stateMu.RUnlock()
	return r.chStopInFlight
}

func (r *rpcClient) getNodeMode() string {
	if r.http != nil {
		return "http"
	}
	return "websocket"
}

func (r *rpcClient) Name() string {
	return r.name
}

func (r *rpcClient) Order() int32 {
	return r.order
}

func Name(r *rpcClient) string {
	return r.name
}
