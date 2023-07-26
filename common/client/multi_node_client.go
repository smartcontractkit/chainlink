package client

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/pkg/errors"
	clienttypes "github.com/smartcontractkit/chainlink/v2/common/chains/client"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/smartcontractkit/chainlink/v2/common/types"
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

type MultiNodeClient[
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
	Dial(ctx context.Context) error
	Close()
	NodeStates() map[string]string
}

type multiNodeClient[
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
	chainID             CHAINID
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

func NewMultiNodeClient[
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
](logger logger.Logger, selectionMode string, noNewHeadsTreshold time.Duration, nodes []Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE], sendonlys []Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE], chainID CHAINID, chainType config.ChainType,
) MultiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE] {
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

	lggr := logger.Named("MultiNodeClient").With("evmChainID", chainID.String())

	c := &multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]{
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

	c.logger.Debugf("The MultiNodeClient is configured to use NodeSelectionMode: %s", selectionMode)

	return c
}

// selectNode returns the active Node, if it is still NodeStateAlive, otherwise it selects a new one from the NodeSelector.
func (c *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) selectNode() (node Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) {
	c.activeMu.RLock()
	node = c.activeNode
	c.activeMu.RUnlock()
	if node != nil && node.State() == NodeStateAlive {
		return // still alive
	}

	// select a new one
	c.activeMu.Lock()
	defer c.activeMu.Unlock()
	node = c.activeNode
	if node != nil && node.State() == NodeStateAlive {
		return // another goroutine beat us here
	}

	c.activeNode = c.nodeSelector.Select()

	if c.activeNode == nil {
		c.logger.Criticalw("No live RPC nodes available", "NodeSelectionMode", c.nodeSelector.Name())
		errmsg := fmt.Errorf("no live nodes available for chain %s", c.chainID.String())
		c.SvcErrBuffer.Append(errmsg)
		return &erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]{errMsg: errmsg.Error()}
	}

	return c.activeNode
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) Dial(ctx context.Context) error {
	return client.StartOnce("Pool", func() (merr error) {
		if len(client.nodes) == 0 {
			return errors.Errorf("no available nodes for chain %s", client.chainID.String())
		}
		var ms services.MultiStart
		for _, n := range client.nodes {
			chainID, err := n.ChainID()
			if err != nil {
				return errors.Errorf("Invalid ChainID")
			}
			if chainID.String() != client.chainID.String() {
				return ms.CloseBecause(errors.Errorf("node %s has chain ID %s which does not match pool chain ID of %s", n.String(), n.ChainID().String(), client.chainID.String()))
			}
			rawNode, ok := n.(*node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE])
			if ok {
				// This is a bit hacky but it allows the node to be aware of
				// pool state and prevent certain state transitions that might
				// otherwise leave no nodes available. It is better to have one
				// node in a degraded state than no nodes at all.
				rawNode.nLiveNodes = rawNode.nLiveNodes
			}
			// node will handle its own redialing and automatic recovery
			if err := ms.Start(ctx, n); err != nil {
				return err
			}
		}
		for _, s := range client.sendonlys {
			chainID, err := s.ChainID()
			if err != nil {
				return errors.Errorf("Invalid ChainID")
			}
			if chainID.String() != client.chainID.String() {
				return ms.CloseBecause(errors.Errorf("sendonly node %s has chain ID %s which does not match pool chain ID of %s", s.String(), s.ChainID().String(), client.chainID.String()))
			}
			if err := ms.Start(ctx, s); err != nil {
				return err
			}
		}
		client.wg.Add(1)
		go client.runLoop()

		return nil
	})
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) BalanceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (*big.Int, error) {
	return client.selectNode().BalanceAt(ctx, account, blockNumber)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) BatchCallContext(ctx context.Context, b []any) error {
	return client.selectNode().BatchCallContext(ctx, b)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) BatchCallContextAll(ctx context.Context, b []any) error {
	return client.selectNode().BatchCallContextAll(ctx, b)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) BatchGetReceipts(
	ctx context.Context,
	attempts []txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
) (txReceipt []TXRECEIPT, txErr []error, err error) {
	return client.selectNode().BatchGetReceipts(ctx, attempts)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) BlockByHash(ctx context.Context, hash BLOCKHASH) (*BLOCK, error) {
	return client.selectNode().BlockByHash(ctx, hash)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) BlockByNumber(ctx context.Context, number *big.Int) (*BLOCK, error) {
	return client.selectNode().BlockByNumber(ctx, number)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return client.selectNode().CallContext(ctx, result, method, args...)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) CallContract(
	ctx context.Context,
	attempt txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
	blockNumber *big.Int,
) (rpcErr fmt.Stringer, extractErr error) {
	return client.selectNode().CallContract(ctx, attempt, blockNumber)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) ChainID() (CHAINID, error) {
	return client.selectNode().ChainID()
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) CodeAt(ctx context.Context, account ADDR, blockNumber *big.Int) ([]byte, error) {
	return client.selectNode().CodeAt(ctx, account, blockNumber)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) ConfiguredChainID() CHAINID {
	return client.selectNode().ConfiguredChainID()
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) EstimateGas(ctx context.Context, call any) (gas uint64, err error) {
	return client.selectNode().EstimateGas(ctx, call)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) FilterEvents(ctx context.Context, query EVENTOPS) ([]EVENT, error) {
	return client.selectNode().FilterEvents(ctx, query)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) HeadByNumber(ctx context.Context, number *big.Int) (head *types.Head[BLOCKHASH], err error) {
	return client.selectNode().HeadByNumber(ctx, number)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) HeadByHash(ctx context.Context, hash BLOCKHASH) (head *types.Head[BLOCKHASH], err error) {
	return client.selectNode().HeadByHash(ctx, hash)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) IsL2() bool {
	return client.selectNode().IsL2()
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) LatestBlockHeight(ctx context.Context) (*big.Int, error) {
	return client.selectNode().LatestBlockHeight(ctx)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) LINKBalance(ctx context.Context, accountAddress ADDR, linkAddress ADDR) (*assets.Link, error) {
	return client.selectNode().LINKBalance(ctx, accountAddress, linkAddress)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) PendingSequenceAt(ctx context.Context, addr ADDR) (SEQ, error) {
	return client.selectNode().PendingSequenceAt(ctx, addr)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) SendEmptyTransaction(
	ctx context.Context,
	newTxAttempt func(seq SEQ, feeLimit uint32, fee FEE, fromAddress ADDR) (attempt txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE], err error),
	seq SEQ,
	gasLimit uint32,
	fee FEE,
	fromAddress ADDR,
) (txhash string, err error) {
	return client.selectNode().SendEmptyTransaction(ctx, newTxAttempt, seq, gasLimit, fee, fromAddress)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) SendTransaction(ctx context.Context, tx *TX) error {
	return client.selectNode().SendTransaction(ctx, tx)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) SendTransactionReturnCode(
	ctx context.Context,
	tx txmgrtypes.Tx[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
	attempt txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
	lggr logger.Logger,
) (clienttypes.SendTxReturnCode, error) {
	return client.selectNode().SendTransactionReturnCode(ctx, tx, attempt, lggr)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) SequenceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (SEQ, error) {
	return client.selectNode().SequenceAt(ctx, account, blockNumber)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) SimulateTransaction(ctx context.Context, tx *TX) error {
	return client.selectNode().SimulateTransaction(ctx, tx)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) Subscribe(ctx context.Context, channel chan<- types.Head[BLOCKHASH], args ...interface{}) (types.Subscription, error) {
	return client.selectNode().Subscribe(ctx, channel, args)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) TokenBalance(ctx context.Context, account ADDR, tokenAddr ADDR) (*big.Int, error) {
	return client.selectNode().TokenBalance(ctx, account, tokenAddr)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) TransactionByHash(ctx context.Context, txHash TXHASH) (*TX, error) {
	return client.selectNode().TransactionByHash(ctx, txHash)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) TransactionReceipt(ctx context.Context, txHash TXHASH) (*TXRECEIPT, error) {
	return client.selectNode().TransactionReceipt(ctx, txHash)
}
