package multinodeclient

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	NodeSelectionMode_HighestHead     = "HighestHead"
	NodeSelectionMode_RoundRobin      = "RoundRobin"
	NodeSelectionMode_TotalDifficulty = "TotalDifficulty"
	NodeSelectionMode_PriorityLevel   = "PriorityLevel"
)

type NodeSelector[
	CHAINID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK any,
	BLOCKHASH types.Hashable,
	TX any,
	TXHASH types.Hashable,
	EVENT any,
	EVENTOPS any, // event filter query options
	TXRECEIPT txmgrtypes.ChainReceipt[TXHASH, BLOCKHASH],
	FEE feetypes.Fee,
] interface {
	// Select returns a Node, or nil if none can be selected.
	// Implementation must be thread-safe.
	Select() Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]
	// Name returns the strategy name, e.g. "HighestHead" or "RoundRobin"
	Name() string
}

type Pool[
	CHAINID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK any,
	BLOCKHASH types.Hashable,
	TX any,
	TXHASH types.Hashable,
	EVENT any,
	EVENTOPS any, // event filter query options
	TXRECEIPT txmgrtypes.ChainReceipt[TXHASH, BLOCKHASH],
	FEE feetypes.Fee,
] interface {
	RPCClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]
}

type pool[
	CHAINID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK any,
	BLOCKHASH types.Hashable,
	TX any,
	TXHASH types.Hashable,
	EVENT any,
	EVENTOPS any,
	TXRECEIPT txmgrtypes.ChainReceipt[TXHASH, BLOCKHASH],
	FEE feetypes.Fee,
] struct {
	utils.StartStopOnce
	nodes               []Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]
	sendonlys           []Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]
	chainID             *big.Int
	chainType           config.ChainType
	logger              logger.Logger
	selectionMode       string
	noNewHeadsThreshold time.Duration
	nodeSelector        NodeSelector[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]

	activeMu   sync.RWMutex
	activeNode Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]

	chStop utils.StopChan
	wg     sync.WaitGroup
}

func NewPool[
	CHAINID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK any,
	BLOCKHASH types.Hashable,
	TX any,
	TXHASH types.Hashable,
	EVENT any,
	EVENTOPS any,
	TXRECEIPT txmgrtypes.ChainReceipt[TXHASH, BLOCKHASH],
	FEE feetypes.Fee,
](logger logger.Logger, selectionMode string, noNewHeadsTreshold time.Duration, nodes []Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE], sendonlys []Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE], chainID *big.Int, chainType config.ChainType,
) Pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE] {
	// Not sure if the below line is correct, but leaving it like that for now
	if &chainID == nil {
		panic("chainID is required")
	}

	nodeSelector := func() NodeSelector[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE] {
		switch selectionMode {
		case NodeSelectionMode_HighestHead:
			return NewHighestHeadNodeSelector[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE](nodes)
		case NodeSelectionMode_RoundRobin:
			return NewRoundRobinSelector[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE](nodes)
		case NodeSelectionMode_TotalDifficulty:
			return NewTotalDifficultyNodeSelector[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE](nodes)
		case NodeSelectionMode_PriorityLevel:
			return NewPriorityLevelNodeSelector[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE](nodes)
		default:
			panic(fmt.Sprintf("unsupported NodeSelectionMode: %s", selectionMode))
		}
	}()

	lggr := logger.Named("Pool").With("evmChainID", chainID.String())

	p := &pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]{
		nodes:               nodes,
		sendonlys:           sendonlys,
		chainID:             chainID,
		chainType:           chainType,
		logger:              lggr,
		selectionMode:       selectionMode,
		noNewHeadsThreshold: noNewHeadsTreshold,
		nodeSelector:        nodeSelector,
		chStop:              make(chan struct{}),
	}

	p.logger.Debugf("The pool is configured to use NodeSelectionMode: %s", selectionMode)

	return p
	// return &pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]{}
}

// selectNode returns the active Node, if it is still NodeStateAlive, otherwise it selects a new one from the NodeSelector.
func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) selectNode() (node Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) {
	p.activeMu.RLock()
	node = p.activeNode
	p.activeMu.RUnlock()
	if node != nil && node.State() == NodeStateAlive {
		return // still alive
	}

	// select a new one
	p.activeMu.Lock()
	defer p.activeMu.Unlock()
	node = p.activeNode
	if node != nil && node.State() == NodeStateAlive {
		return // another goroutine beat us here
	}

	p.activeNode = p.nodeSelector.Select()

	if p.activeNode == nil {
		p.logger.Criticalw("No live RPC nodes available", "NodeSelectionMode", p.nodeSelector.Name())
		errmsg := fmt.Errorf("no live nodes available for chain %s", p.chainID.String())
		p.SvcErrBuffer.Append(errmsg)
		return &erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]{errMsg: errmsg.Error()}
	}

	return p.activeNode
}

func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) BalanceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (*big.Int, error) {
	return p.selectNode().BalanceAt(ctx, account, blockNumber)
}

func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) BatchGetReceipts(
	ctx context.Context,
	attempts []txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
) (txReceipt []TXRECEIPT, txErr []error, err error) {
	return p.selectNode().BatchGetReceipts(ctx, attempts)
}

func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) BlockByHash(ctx context.Context, hash BLOCKHASH) (*BLOCK, error) {
	return p.selectNode().BlockByHash(ctx, hash)
}

func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) BlockByNumber(ctx context.Context, number *big.Int) (*BLOCK, error) {
	return p.selectNode().BlockByNumber(ctx, number)
}

func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return p.selectNode().CallContext(ctx, result, method, args...)
}

func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) CallContract(
	ctx context.Context,
	attempt txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
	blockNumber *big.Int,
) (rpcErr fmt.Stringer, extractErr error) {
	return p.selectNode().CallContract(ctx, attempt, blockNumber)
}

func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) ChainID() (CHAINID, error) {
	return p.selectNode().ChainID()
}

func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) ConfiguredChainID() CHAINID {
	return p.selectNode().ConfiguredChainID()
}

func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) FilterEvents(ctx context.Context, query EVENTOPS) ([]EVENT, error) {
	return p.selectNode().FilterEvents(ctx, query)
}

func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) LatestBlockHeight(ctx context.Context) (*big.Int, error) {
	return p.selectNode().LatestBlockHeight(ctx)
}

func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) PendingSequenceAt(ctx context.Context, addr ADDR) (SEQ, error) {
	return p.selectNode().PendingSequenceAt(ctx, addr)
}

func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) SendEmptyTransaction(
	ctx context.Context,
	newTxAttempt func(seq SEQ, feeLimit uint32, fee FEE, fromAddress ADDR) (attempt txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE], err error),
	seq SEQ,
	gasLimit uint32,
	fee FEE,
	fromAddress ADDR,
) (txhash string, err error) {
	return p.selectNode().SendEmptyTransaction(ctx, newTxAttempt, seq, gasLimit, fee, fromAddress)
}

func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) SendTransaction(ctx context.Context, tx *TX) error {
	return p.selectNode().SendTransaction(ctx, tx)
}

func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) SendTransactionReturnCode(
	ctx context.Context,
	tx txmgrtypes.Tx[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
	attempt txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
	lggr logger.Logger,
) (SendTxReturnCode, error) {
	return p.selectNode().SendTransactionReturnCode(ctx, tx, attempt, lggr)
}

func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) SequenceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (SEQ, error) {
	return p.selectNode().SequenceAt(ctx, account, blockNumber)
}

func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) SimulateTransaction(ctx context.Context, tx *TX) error {
	return p.selectNode().SimulateTransaction(ctx, tx)
}

func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) 	Subscribe(ctx context.Context, channel chan<- types.Head[BLOCKHASH], args ...interface{}) (types.Subscription, error) {
	return p.selectNode().Subscribe(ctx, channel, args)
}

func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) TokenBalance(ctx context.Context, account ADDR, tokenAddr ADDR) (*big.Int, error) {
	return p.selectNode().TokenBalance(ctx, account, tokenAddr)
}

func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) TransactionByHash(ctx context.Context, txHash TXHASH) (*TX, error) {
	return p.selectNode().TransactionByHash(ctx, txHash)
}

func (p *pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) TransactionReceipt(ctx context.Context, txHash TXHASH) (*TXRECEIPT, error) {
	return p.selectNode().TransactionReceipt(ctx, txHash)
}
