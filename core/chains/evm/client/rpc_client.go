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

	clienttypes "github.com/smartcontractkit/chainlink/v2/common/chains/client"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type RPCClient interface {
	clienttypes.ChainRPCClient[
		*big.Int,
		common.Hash,
		*evmtypes.Head,
		ethereum.Subscription,
	]
	clienttypes.RPCClient[
		*big.Int,
		evmtypes.Nonce,
		common.Address,
		common.Hash,
		types.Transaction,
		common.Hash,
		types.Log,
		ethereum.FilterQuery,
		*evmtypes.Receipt,
		*assets.Wei,
		*evmtypes.Head,
		ethereum.Subscription,
	]
}

type rpcClient struct {
	utils.StartStopOnce
	lfcLog  logger.Logger
	rpcLog  logger.Logger
	name    string
	id      int32
	chainID *big.Int

	ws   rawclient
	http *rawclient

	stateMu sync.RWMutex // protects state* fields
	state   clienttypes.NodeState

	// Need to track subscriptions because closing the RPC does not (always?)
	// close the underlying subscription
	subs []ethereum.Subscription

	// chStopInFlight can be closed to immediately cancel all in-flight requests on
	// this rpcClient. Closing and replacing should be serialized through
	// stateMu since it can happen on state transitions as well as rpcClient Close.
	chStopInFlight chan struct{}
	// rpcClientCtx is the rpcClient lifetime's context
	rpcClientCtx context.Context
	// cancelRpcClientCtx cancels rpcClientCtx when stopping the rpcClient
	cancelRpcClientCtx context.CancelFunc
	// wg waits for subsidiary goroutines
	wg sync.WaitGroup
}

// NewRPCCLient returns a new *rpcClient as clienttypes.RPCClient
func NewRPCClient(lggr logger.Logger, wsuri url.URL, httpuri *url.URL, name string, id int32, chainID *big.Int) RPCClient {
	r := new(rpcClient)
	r.name = name
	r.id = id
	r.chainID = chainID
	r.ws.uri = wsuri
	if httpuri != nil {
		r.http = &rawclient{uri: *httpuri}
	}
	r.chStopInFlight = make(chan struct{})
	r.rpcClientCtx, r.cancelRpcClientCtx = context.WithCancel(context.Background())
	lggr = lggr.Named("Client").With(
		"clientTier", "primary",
		"clientName", name,
		"client", r.String(),
		"evmChainID", chainID,
	)
	r.lfcLog = lggr.Named("Lifecycle")
	r.rpcLog = lggr.Named("RPC")

	return r
}

// Not thread-safe
// Pure dial: does not mutate rpcClient "state" field.
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

		r.cancelRpcClientCtx()
		r.cancelInflightRequests()
		r.state = clienttypes.NodeStateClosed
		return nil
	})
}

func (r *rpcClient) SetState(state clienttypes.NodeState) {
	r.state = state
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
			r.name,                         // rpcClient name
			rpcDomain,                      // rpc domain
			"false",                        // is send only
			strconv.FormatBool(err == nil), // is successful
			callName,                       // rpc call name
		).
		Observe(float64(callDuration))
}

func (r *rpcClient) getRPCDomain() string {
	if r.http != nil {
		return r.http.uri.Host
	}
	return r.ws.uri.Host
}

// registerSub adds the sub to the rpcClient list
func (r *rpcClient) registerSub(sub ethereum.Subscription) {
	r.stateMu.Lock()
	defer r.stateMu.Unlock()
	r.subs = append(r.subs, sub)
}

// disconnectAll disconnects all clients connected to the rpcClient
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
	// Not sure if this is the best solution.
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

func (r *rpcClient) BlockByNumber(ctx context.Context, number *big.Int) (head *evmtypes.Head, err error) {
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

func (r *rpcClient) BlockByHash(ctx context.Context, hash common.Hash) (head *evmtypes.Head, err error) {
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

func (r *rpcClient) SendTransactionReturnCode(ctx context.Context, tx *types.Transaction) (clienttypes.SendTxReturnCode, error) {
	from, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
	if err != nil {
		return clienttypes.Fatal, err
	}
	err = r.SendTransaction(ctx, tx)

	return NewSendOnlyErrorReturnCode(err, r.rpcLog, tx, from, true)
}

func (r *rpcClient) SimulateTransaction(ctx context.Context, tx *types.Transaction) error {
	// todo: implement if used
	return errors.New("SimulateTransaction not implemented")
}

func (r *rpcClient) SendEmptyTransaction(
	ctx context.Context,
	newTxAttempt func(nonce evmtypes.Nonce, feeLimit uint32, fee evmtypes.EvmFee, fromAddress common.Address) (attempt any, err error),
	nonce evmtypes.Nonce,
	gasLimit uint32,
	fee evmtypes.EvmFee,
	fromAddress common.Address,
) (txhash string, err error) {
	// todo: implement if used
	return "", errors.New("SendEmptyTransaction not implemented")
}

// PendingSequenceAt returns one higher than the highest nonce from both mempool and mined transactions
func (r *rpcClient) PendingSequenceAt(ctx context.Context, account common.Address) (nonce evmtypes.Nonce, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return 0, err
	}
	defer cancel()
	lggr := r.newRqLggr().With("account", account)

	lggr.Debug("RPC call: evmclient.Client#PendingNonceAt")
	start := time.Now()
	var n uint64
	if http != nil {
		n, err = http.geth.PendingNonceAt(ctx, account)
		nonce = evmtypes.Nonce(int64(n))
		err = r.wrapHTTP(err)
	} else {
		n, err = ws.geth.PendingNonceAt(ctx, account)
		nonce = evmtypes.Nonce(int64(n))
		err = r.wrapWS(err)
	}
	duration := time.Since(start)

	r.logResult(lggr, err, duration, r.getRPCDomain(), "PendingNonceAt",
		"nonce", nonce,
	)

	return
}

// SequenceAt is a bit of a misnomer. You might expect it to return the highest
// mined nonce at the given block number, but it actually returns the total
// transaction count which is the highest mined nonce + 1
func (r *rpcClient) SequenceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (nonce evmtypes.Nonce, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)
	if err != nil {
		return 0, err
	}
	defer cancel()
	lggr := r.newRqLggr().With("account", account, "blockNumber", blockNumber)

	lggr.Debug("RPC call: evmclient.Client#NonceAt")
	start := time.Now()
	var n uint64
	if http != nil {
		n, err = http.geth.NonceAt(ctx, account, blockNumber)
		nonce = evmtypes.Nonce(int64(n))
		err = r.wrapHTTP(err)
	} else {
		n, err = ws.geth.NonceAt(ctx, account, blockNumber)
		nonce = evmtypes.Nonce(int64(n))
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

func (r *rpcClient) LatestBlockHeight(ctx context.Context) (*big.Int, error) {
	var height big.Int
	h, err := r.BlockNumber(ctx)
	return height.SetUint64(h), err
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

// TokenBalance returns the balance of the given address for the token contract address.
func (r *rpcClient) TokenBalance(ctx context.Context, address common.Address, contractAddress common.Address) (*big.Int, error) {
	result := ""
	numLinkBigInt := new(big.Int)
	functionSelector := evmtypes.HexToFunctionSelector(BALANCE_OF_ADDRESS_FUNCTION_SELECTOR) // balanceOf(address)
	data := utils.ConcatBytes(functionSelector.Bytes(), common.LeftPadBytes(address.Bytes(), utils.EVMWordByteLen))
	args := CallArgs{
		To:   contractAddress,
		Data: data,
	}
	err := r.CallContext(ctx, &result, "eth_call", args, "latest")
	if err != nil {
		return numLinkBigInt, err
	}
	numLinkBigInt.SetString(result, 0)
	return numLinkBigInt, nil
}

// LINKBalance returns the balance of LINK at the given address
func (r *rpcClient) LINKBalance(ctx context.Context, address common.Address, linkAddress common.Address) (*assets.Link, error) {
	balance, err := r.TokenBalance(ctx, address, linkAddress)
	if err != nil {
		return assets.NewLinkFromJuels(0), err
	}
	return (*assets.Link)(balance), nil
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

func (r *rpcClient) ClientVersion(ctx context.Context) (version string, err error) {
	err = r.CallContext(ctx, &version, "web3_clientVersion")
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

// Returns the ChainID according to the geth client. This is useful for functions like verify()
// the common node.
func (r *rpcClient) ClientChainID(ctx context.Context) (chainID *big.Int, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtxAndSafeGetClients(ctx)

	defer cancel()

	if http != nil {
		chainID, err = http.geth.ChainID(ctx)
		err = r.wrapHTTP(err)
	} else {
		chainID, err = ws.geth.ChainID(ctx)
		err = r.wrapWS(err)
	}
	return
}

func (r *rpcClient) ConfiguredChainID() *big.Int {
	chainID, _ := r.ChainID()
	return chainID
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

// makeLiveQueryCtxAndSafeGetClients wraps makeQueryCtx but returns error if rpcClient is not clienttypes.NodeStateAlive.
func (r *rpcClient) makeLiveQueryCtxAndSafeGetClients(parentCtx context.Context) (ctx context.Context, cancel context.CancelFunc, ws rawclient, http *rawclient, err error) {
	// Need to wrap in mutex because state transition can cancel and replace the
	// context
	r.stateMu.RLock()
	if r.state != clienttypes.NodeStateAlive {
		err = errors.Errorf("cannot execute RPC call on rpcClient with state: %s", r.state)
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

func (r *rpcClient) Name() string {
	return r.name
}

func Name(r *rpcClient) string {
	return r.name
}
