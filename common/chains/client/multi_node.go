package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

var (
	// PromMultiNodeRPCNodeStates reports current RPC node state
	PromMultiNodeRPCNodeStates = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_rpc_node_states",
		Help: "The number of RPC nodes currently in the given state for the given chain",
	}, []string{"network", "chainId", "state"})
	ErroringNodeError = fmt.Errorf("no live nodes available")
)

const (
	NodeSelectionMode_HighestHead     = "HighestHead"
	NodeSelectionMode_RoundRobin      = "RoundRobin"
	NodeSelectionMode_TotalDifficulty = "TotalDifficulty"
	NodeSelectionMode_PriorityLevel   = "PriorityLevel"
)

type NodeSelector[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC NodeClient[CHAIN_ID, HEAD],
] interface {
	// Select returns a Node, or nil if none can be selected.
	// Implementation must be thread-safe.
	Select() Node[CHAIN_ID, HEAD, RPC]
	// Name returns the strategy name, e.g. "HighestHead" or "RoundRobin"
	Name() string
}

type MultiNode[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC NodeClient[CHAIN_ID, HEAD],
	TX any,
] interface {
	Dial(context.Context) error
	Close() error
	NodeStates() map[string]string
	SelectNode() (Node[CHAIN_ID, HEAD, RPC], error)
	NodesAsSendOnlys() []SendOnlyNode[CHAIN_ID, RPC]
	WrapSendOnlyTransaction(ctx context.Context, lggr logger.Logger, tx TX, n SendOnlyNode[CHAIN_ID, RPC],
		f func(ctx context.Context, lggr logger.Logger, tx TX, n SendOnlyNode[CHAIN_ID, RPC]))
}

func ContextWithDefaultTimeout() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), QueryTimeout)
}

type multiNode[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC NodeClient[CHAIN_ID, HEAD],
	TX any,
] struct {
	utils.StartStopOnce
	nodes               []Node[CHAIN_ID, HEAD, RPC]
	sendonlys           []SendOnlyNode[CHAIN_ID, RPC]
	chainID             CHAIN_ID
	logger              logger.Logger
	selectionMode       string
	noNewHeadsThreshold time.Duration
	nodeSelector        NodeSelector[CHAIN_ID, HEAD, RPC]
	chainFamily         string

	activeMu   sync.RWMutex
	activeNode Node[CHAIN_ID, HEAD, RPC]

	chStop utils.StopChan
	wg     sync.WaitGroup
}

func NewMultiNode[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC NodeClient[CHAIN_ID, HEAD],
	TX any,
](
	logger logger.Logger,
	selectionMode string,
	noNewHeadsThreshold time.Duration,
	nodes []Node[CHAIN_ID, HEAD, RPC],
	sendonlys []SendOnlyNode[CHAIN_ID, RPC],
	chainID CHAIN_ID,
	chainFamily string,
) MultiNode[CHAIN_ID, HEAD, RPC, TX] {
	nodeSelector := func() NodeSelector[CHAIN_ID, HEAD, RPC] {
		switch selectionMode {
		case NodeSelectionMode_HighestHead:
			return NewHighestHeadNodeSelector[CHAIN_ID, HEAD, RPC](nodes)
		case NodeSelectionMode_RoundRobin:
			return NewRoundRobinSelector[CHAIN_ID, HEAD, RPC](nodes)
		case NodeSelectionMode_TotalDifficulty:
			return NewTotalDifficultyNodeSelector[CHAIN_ID, HEAD, RPC](nodes)
		case NodeSelectionMode_PriorityLevel:
			return NewPriorityLevelNodeSelector[CHAIN_ID, HEAD, RPC](nodes)
		default:
			panic(fmt.Sprintf("unsupported NodeSelectionMode: %s", selectionMode))
		}
	}()

	lggr := logger.Named("MultiNode").With("chainID", chainID.String())

	c := &multiNode[CHAIN_ID, HEAD, RPC, TX]{
		nodes:               nodes,
		sendonlys:           sendonlys,
		chainID:             chainID,
		logger:              lggr,
		selectionMode:       selectionMode,
		noNewHeadsThreshold: noNewHeadsThreshold,
		nodeSelector:        nodeSelector,
		chStop:              make(chan struct{}),
		chainFamily:         chainFamily,
	}

	c.logger.Debugf("The MultiNode is configured to use NodeSelectionMode: %s", selectionMode)

	return c
}

// SelectNode returns the active Node, if it is still NodeStateAlive, otherwise it selects a new one from the NodeSelector.
func (c *multiNode[CHAIN_ID, HEAD, RPC, TX]) SelectNode() (node Node[CHAIN_ID, HEAD, RPC], err error) {
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
		err = ErroringNodeError
	}

	return c.activeNode, err
}

func (c *multiNode[CHAIN_ID, HEAD, RPC, TX]) Dial(ctx context.Context) error {
	return c.StartOnce("MultiNode", func() (merr error) {
		if len(c.nodes) == 0 {
			return errors.Errorf("no available nodes for chain %s", c.chainID.String())
		}
		var ms services.MultiStart
		for _, n := range c.nodes {
			if n.ConfiguredChainID().String() != c.chainID.String() {
				return ms.CloseBecause(errors.Errorf("node %s has configured chain ID %s which does not match multinode configured chain ID of %s", n.String(), n.ConfiguredChainID().String(), c.chainID.String()))
			}
			rawNode, ok := n.(*node[CHAIN_ID, HEAD, RPC])
			if ok {
				// This is a bit hacky but it allows the node to be aware of
				// client / pool state and prevent certain state transitions that might
				// otherwise leave no nodes available. It is better to have one
				// node in a degraded state than no nodes at all.
				rawNode.nLiveNodes = c.nLiveNodes
			}
			// node will handle its own redialing and automatic recovery
			if err := ms.Start(ctx, n); err != nil {
				return err
			}
		}
		for _, s := range c.sendonlys {
			if s.ConfiguredChainID().String() != c.chainID.String() {
				return ms.CloseBecause(errors.Errorf("sendonly node %s has configured chain ID %s which does not match multinode configured chain ID of %s", s.String(), s.ConfiguredChainID().String(), c.chainID.String()))
			}
			if err := ms.Start(ctx, s); err != nil {
				return err
			}
		}
		c.wg.Add(1)
		go c.runLoop()

		return nil
	})
}

func (c *multiNode[CHAIN_ID, HEAD, RPC, TX]) NodesAsSendOnlys() (nodes []SendOnlyNode[CHAIN_ID, RPC]) {
	for _, n := range c.nodes {
		nodes = append(nodes, n)
	}
	nodes = append(nodes, c.sendonlys...)
	return
}

// Close tears down the pool and closes all nodes
func (c *multiNode[CHAIN_ID, HEAD, RPC, TX]) Close() error {
	return c.StopOnce("MultiNode", func() error {
		close(c.chStop)
		c.wg.Wait()

		return services.CloseAll(services.MultiCloser(c.nodes), services.MultiCloser(c.sendonlys))
	})
}

func (c *multiNode[CHAIN_ID, HEAD, RPC, TX]) NodeStates() (states map[string]string) {
	states = make(map[string]string)
	for _, n := range c.nodes {
		states[n.Name()] = n.State().String()
	}
	for _, s := range c.sendonlys {
		states[s.Name()] = s.State().String()
	}
	return
}

// nLiveNodes returns the number of currently alive nodes, as well as the highest block number and greatest total difficulty.
// totalDifficulty will be 0 if all nodes return nil.
func (c *multiNode[CHAIN_ID, HEAD, RPC, TX]) nLiveNodes() (nLiveNodes int, blockNumber int64, totalDifficulty *utils.Big) {
	totalDifficulty = utils.NewBigI(0)
	for _, n := range c.nodes {
		if s, num, td := n.StateAndLatest(); s == NodeStateAlive {
			nLiveNodes++
			if num > blockNumber {
				blockNumber = num
			}
			if td != nil && td.Cmp(totalDifficulty) > 0 {
				totalDifficulty = td
			}
		}
	}
	return
}

func (c *multiNode[CHAIN_ID, HEAD, RPC, TX]) runLoop() {
	defer c.wg.Done()

	c.report()

	// Prometheus' default interval is 15s, set this to under 7.5s to avoid
	// aliasing (see: https://en.wikipedia.org/wiki/Nyquist_frequency)
	reportInterval := 6500 * time.Millisecond
	monitor := time.NewTicker(utils.WithJitter(reportInterval))
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

func (c *multiNode[CHAIN_ID, HEAD, RPC, TX]) report() {
	type nodeWithState struct {
		Node  string
		State string
	}

	var total, dead int
	counts := make(map[NodeState]int)
	nodeStates := make([]nodeWithState, len(c.nodes))
	for i, n := range c.nodes {
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
	c.logger.Tracew(fmt.Sprintf("MultiNode state: %d/%d nodes are alive", live, total), "nodeStates", nodeStates)
	if total == dead {
		rerr := fmt.Errorf("no primary nodes available: 0/%d nodes are alive", total)
		c.logger.Criticalw(rerr.Error(), "nodeStates", nodeStates)
		c.SvcErrBuffer.Append(rerr)
	} else if dead > 0 {
		c.logger.Errorw(fmt.Sprintf("At least one primary node is dead: %d/%d nodes are alive", live, total), "nodeStates", nodeStates)
	}
}

func (c *multiNode[CHAIN_ID, HEAD, RPC, TX]) WrapSendOnlyTransaction(ctx context.Context, lggr logger.Logger, tx TX, n SendOnlyNode[CHAIN_ID, RPC],
	f func(ctx context.Context, lggr logger.Logger, tx TX, n SendOnlyNode[CHAIN_ID, RPC]),
) {
	ok := c.IfNotStopped(func() {
		// Must wrap inside IfNotStopped to avoid waitgroup racing with Close
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			f(ctx, lggr, tx, n)
		}()
	})
	if !ok {
		c.logger.Debug("Cannot send transaction on sendonly node; multinode is stopped", "node", n.String())
	}
}
