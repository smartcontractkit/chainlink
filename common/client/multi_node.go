package client

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

var (
	// PromMultiNodeRPCNodeStates reports current RPC node state
	PromMultiNodeRPCNodeStates = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "multi_node_states",
		Help: "The number of RPC nodes currently in the given state for the given chain",
	}, []string{"network", "chainId", "state"})
	// PromMultiNodeInvariantViolations reports violation of our assumptions
	PromMultiNodeInvariantViolations = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "multi_node_invariant_violations",
		Help: "The number of invariant violations",
	}, []string{"network", "chainId", "invariant"})
	ErroringNodeError = fmt.Errorf("no live nodes available")
)

// MultiNode is a generalized multi node client interface that includes methods to interact with different chains.
// It also handles multiple node RPC connections simultaneously.
type MultiNode[
	CHAIN_ID types.ID,
	RPC_CLIENT any,
] interface {
	Dial(ctx context.Context) error
	ChainID() CHAIN_ID
	// SelectRPC - returns the best healthy RPCClient
	SelectRPC() (RPC_CLIENT, error)
	// DoAll - calls `do` sequentially on all healthy RPCClients.
	// `do` can abort subsequent calls by returning `false`.
	// Returns error if `do` was not called or context returns an error.
	DoAll(ctx context.Context, do func(ctx context.Context, rpc RPC_CLIENT, isSendOnly bool) bool) error
	// NodeStates - returns RPCs' states
	NodeStates() map[string]NodeState
	Close() error
}

type multiNode[
	CHAIN_ID types.ID,
	RPC_CLIENT any,
] struct {
	services.StateMachine
	primaryNodes   []Node[CHAIN_ID, RPC_CLIENT]
	sendOnlyNodes  []SendOnlyNode[CHAIN_ID, RPC_CLIENT]
	chainID        CHAIN_ID
	lggr           logger.SugaredLogger
	selectionMode  string
	nodeSelector   NodeSelector[CHAIN_ID, RPC_CLIENT]
	leaseDuration  time.Duration
	leaseTicker    *time.Ticker
	chainFamily    string
	reportInterval time.Duration

	activeMu   sync.RWMutex
	activeNode Node[CHAIN_ID, RPC_CLIENT]

	chStop services.StopChan
	wg     sync.WaitGroup
}

func NewMultiNode[
	CHAIN_ID types.ID,
	RPC_CLIENT any,
](
	lggr logger.Logger,
	selectionMode string, // type of the "best" RPC selector (e.g HighestHead, RoundRobin, etc.)
	leaseDuration time.Duration, // defines interval on which new "best" RPC should be selected
	primaryNodes []Node[CHAIN_ID, RPC_CLIENT],
	sendOnlyNodes []SendOnlyNode[CHAIN_ID, RPC_CLIENT],
	chainID CHAIN_ID, // configured chain ID (used to verify that passed primaryNodes belong to the same chain)
	chainFamily string, // name of the chain family - used in the metrics
) MultiNode[CHAIN_ID, RPC_CLIENT] {
	nodeSelector := newNodeSelector(selectionMode, primaryNodes)
	// Prometheus' default interval is 15s, set this to under 7.5s to avoid
	// aliasing (see: https://en.wikipedia.org/wiki/Nyquist_frequency)
	const reportInterval = 6500 * time.Millisecond
	c := &multiNode[CHAIN_ID, RPC_CLIENT]{
		primaryNodes:   primaryNodes,
		sendOnlyNodes:  sendOnlyNodes,
		chainID:        chainID,
		lggr:           logger.Sugared(lggr).Named("MultiNode").With("chainID", chainID.String()),
		selectionMode:  selectionMode,
		nodeSelector:   nodeSelector,
		chStop:         make(services.StopChan),
		leaseDuration:  leaseDuration,
		chainFamily:    chainFamily,
		reportInterval: reportInterval,
	}

	c.lggr.Debugf("The MultiNode is configured to use NodeSelectionMode: %s", selectionMode)

	return c
}

func (c *multiNode[CHAIN_ID, RPC_CLIENT]) ChainID() CHAIN_ID {
	return c.chainID
}

func (c *multiNode[CHAIN_ID, RPC_CLIENT]) DoAll(ctx context.Context, do func(ctx context.Context, rpc RPC_CLIENT, isSendOnly bool) bool) error {
	callsCompleted := 0
	for _, n := range c.primaryNodes {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if n.State() != NodeStateAlive {
			continue
		}
		if do(ctx, n.RPC(), false) {
			callsCompleted++
		}
	}
	if callsCompleted == 0 {
		return fmt.Errorf("no calls were completed")
	}

	for _, n := range c.sendOnlyNodes {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if n.State() != NodeStateAlive {
			continue
		}
		do(ctx, n.RPC(), false)
	}
	return nil
}

func (c *multiNode[CHAIN_ID, RPC_CLIENT]) NodeStates() map[string]NodeState {
	states := map[string]NodeState{}
	for _, n := range c.primaryNodes {
		states[n.String()] = n.State()
	}
	for _, n := range c.sendOnlyNodes {
		states[n.String()] = n.State()
	}
	return states
}

// Dial starts every node in the pool
//
// Nodes handle their own redialing and runloops, so this function does not
// return any error if the nodes aren't available
func (c *multiNode[CHAIN_ID, RPC_CLIENT]) Dial(ctx context.Context) error {
	return c.StartOnce("MultiNode", func() (merr error) {
		if len(c.primaryNodes) == 0 {
			return fmt.Errorf("no available nodes for chain %s", c.chainID.String())
		}
		var ms services.MultiStart
		for _, n := range c.primaryNodes {
			if n.ConfiguredChainID().String() != c.chainID.String() {
				return ms.CloseBecause(fmt.Errorf("node %s has configured chain ID %s which does not match multinode configured chain ID of %s", n.String(), n.ConfiguredChainID().String(), c.chainID.String()))
			}
			/* TODO: Dmytro's PR on local finality handles this better.
			rawNode, ok := n.(*node[CHAIN_ID, *evmtypes.Head, RPC_CLIENT])
			if ok {
				// This is a bit hacky but it allows the node to be aware of
				// pool state and prevent certain state transitions that might
				// otherwise leave no primaryNodes available. It is better to have one
				// node in a degraded state than no primaryNodes at all.
				rawNode.nLiveNodes = c.nLiveNodes
			}
			*/
			// node will handle its own redialing and automatic recovery
			if err := ms.Start(ctx, n); err != nil {
				return err
			}
		}
		for _, s := range c.sendOnlyNodes {
			if s.ConfiguredChainID().String() != c.chainID.String() {
				return ms.CloseBecause(fmt.Errorf("sendonly node %s has configured chain ID %s which does not match multinode configured chain ID of %s", s.String(), s.ConfiguredChainID().String(), c.chainID.String()))
			}
			if err := ms.Start(ctx, s); err != nil {
				return err
			}
		}
		c.wg.Add(1)
		go c.runLoop()

		if c.leaseDuration.Seconds() > 0 && c.selectionMode != NodeSelectionModeRoundRobin {
			c.lggr.Infof("The MultiNode will switch to best node every %s", c.leaseDuration.String())
			c.wg.Add(1)
			go c.checkLeaseLoop()
		} else {
			c.lggr.Info("Best node switching is disabled")
		}

		return nil
	})
}

// Close tears down the MultiNode and closes all nodes
func (c *multiNode[CHAIN_ID, RPC_CLIENT]) Close() error {
	return c.StopOnce("MultiNode", func() error {
		close(c.chStop)
		c.wg.Wait()

		return services.CloseAll(services.MultiCloser(c.primaryNodes), services.MultiCloser(c.sendOnlyNodes))
	})
}

// SelectRPC returns an RPC of an active node. If there are no active nodes it returns an error.
// Call this method from your chain-specific client implementation to access any chain-specific rpc calls.
func (c *multiNode[CHAIN_ID, RPC_CLIENT]) SelectRPC() (rpc RPC_CLIENT, err error) {
	n, err := c.selectNode()
	if err != nil {
		return rpc, err
	}
	return n.RPC(), nil
}

// selectNode returns the active Node, if it is still NodeStateAlive, otherwise it selects a new one from the NodeSelector.
func (c *multiNode[CHAIN_ID, RPC_CLIENT]) selectNode() (node Node[CHAIN_ID, RPC_CLIENT], err error) {
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
		c.lggr.Criticalw("No live RPC nodes available", "NodeSelectionMode", c.nodeSelector.Name())
		errmsg := fmt.Errorf("no live nodes available for chain %s", c.chainID.String())
		c.SvcErrBuffer.Append(errmsg)
		err = ErroringNodeError
	}

	return c.activeNode, err
}

// nLiveNodes returns the number of currently alive nodes, as well as the highest block number and greatest total difficulty.
// totalDifficulty will be 0 if all nodes return nil.
func (c *multiNode[CHAIN_ID, RPC_CLIENT]) nLiveNodes() (nLiveNodes int, blockNumber int64, totalDifficulty *big.Int) {
	totalDifficulty = big.NewInt(0)
	for _, n := range c.primaryNodes {
		if s, chainInfo := n.StateAndLatest(); s == NodeStateAlive {
			nLiveNodes++
			if chainInfo.BlockNumber > blockNumber {
				blockNumber = chainInfo.BlockNumber
			}
			if chainInfo.BlockDifficulty != nil && chainInfo.BlockDifficulty.Cmp(totalDifficulty) > 0 {
				totalDifficulty = chainInfo.BlockDifficulty
			}
		}
	}
	return
}

func (c *multiNode[CHAIN_ID, RPC_CLIENT]) checkLease() {
	bestNode := c.nodeSelector.Select()
	for _, n := range c.primaryNodes {
		// Terminate client subscriptions. Services are responsible for reconnecting, which will be routed to the new
		// best node. Only terminate connections with more than 1 subscription to account for the aliveLoop subscription
		if n.State() == NodeStateAlive && n != bestNode {
			c.lggr.Infof("Switching to best node from %q to %q", n.String(), bestNode.String())
			n.UnsubscribeAll()
		}
	}

	c.activeMu.Lock()
	if bestNode != c.activeNode {
		c.activeNode = bestNode
	}
	c.activeMu.Unlock()
}

func (c *multiNode[CHAIN_ID, RPC_CLIENT]) checkLeaseLoop() {
	defer c.wg.Done()
	c.leaseTicker = time.NewTicker(c.leaseDuration)
	defer c.leaseTicker.Stop()

	for {
		select {
		case <-c.leaseTicker.C:
			c.checkLease()
		case <-c.chStop:
			return
		}
	}
}

func (c *multiNode[CHAIN_ID, RPC_CLIENT]) runLoop() {
	defer c.wg.Done()

	c.report()

	monitor := time.NewTicker(utils.WithJitter(c.reportInterval))
	defer monitor.Stop()

	for {
		select {
		case <-monitor.C:
			c.report()
		case <-c.chStop:
			return
		}
	}
}

func (c *multiNode[CHAIN_ID, RPC_CLIENT]) report() {
	type nodeWithState struct {
		Node  string
		State string
	}

	var total, dead int
	counts := make(map[NodeState]int)
	nodeStates := make([]nodeWithState, len(c.primaryNodes))
	for i, n := range c.primaryNodes {
		state := n.State()
		nodeStates[i] = nodeWithState{n.String(), state.String()}
		total++
		if state != NodeStateAlive {
			dead++
		}
		counts[state]++
	}
	for _, state := range allNodeStates {
		count := counts[state]
		PromMultiNodeRPCNodeStates.WithLabelValues(c.chainFamily, c.chainID.String(), state.String()).Set(float64(count))
	}

	live := total - dead
	c.lggr.Tracew(fmt.Sprintf("MultiNode state: %d/%d nodes are alive", live, total), "nodeStates", nodeStates)
	if total == dead {
		rerr := fmt.Errorf("no primary nodes available: 0/%d nodes are alive", total)
		c.lggr.Criticalw(rerr.Error(), "nodeStates", nodeStates)
		c.SvcErrBuffer.Append(rerr)
	} else if dead > 0 {
		c.lggr.Errorw(fmt.Sprintf("At least one primary node is dead: %d/%d nodes are alive", live, total), "nodeStates", nodeStates)
	}
}
