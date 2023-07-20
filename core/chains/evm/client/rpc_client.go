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
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/smartcontractkit/chainlink/v2/common/multinodeclient"
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
func NewRPCClient(nodeCfg config.NodePool, noNewHeadsThreshold time.Duration, lggr logger.Logger, wsuri url.URL, httpuri *url.URL, name string, id int32, chainID *big.Int, nodeOrder int32) multinodeclient.ChainRPCClient[
	*big.Int,
	evmtypes.Nonce,
	common.Address,
	evmtypes.Block,
	common.Hash,
	evmtypes.Transaction,
	common.Hash,
	evmtypes.Receipt,
	evmtypes.Log,
	ethereum.FilterQuery,
	*evmtypes.Receipt,
	evmtypes.EvmFee,
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

// newRqLggr generates a new logger with a unique request ID
func (r *rpcClient) newRqLggr(mode string) logger.Logger {
	return r.rpcLog.With(
		"requestID", uuid.New(),
		"mode", mode,
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

func (r *rpcClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (balance *big.Int, err error) {
	ctx, cancel, ws, http, err := r.makeLiveQueryCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	lggr := r.newRqLggr(switching(r)).With("account", account.Hex(), "blockNumber", blockNumber)

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

func Name(r *rpcClient) string {
	return r.name
}
