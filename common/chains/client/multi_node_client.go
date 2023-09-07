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
	// PromMultiNodeClientRPCNodeStates reports current RPC node state
	PromMultiNodeClientRPCNodeStates = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_rpc_node_states",
		Help: "The number of RPC nodes currently in the given state for the given chain",
	}, []string{"network", "chainId", "state"})
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
	RPC_CLIENT NodeClient[CHAIN_ID, HEAD],
] interface {
	// Select returns a Node, or nil if none can be selected.
	// Implementation must be thread-safe.
	Select() Node[CHAIN_ID, HEAD, RPC_CLIENT]
	// Name returns the strategy name, e.g. "HighestHead" or "RoundRobin"
	Name() string
}

type MultiNodeClient[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC_CLIENT NodeClient[CHAIN_ID, HEAD],
	TX any,
] interface {
	Dial(context.Context) error
	Close() error
	NodeStates() map[string]string
	SelectNode() Node[CHAIN_ID, HEAD, RPC_CLIENT]
	NodesAsSendOnlys() []SendOnlyNode[CHAIN_ID, RPC_CLIENT]
	WrapSendOnlyTransaction(ctx context.Context, lggr logger.Logger, tx TX, n SendOnlyNode[CHAIN_ID, RPC_CLIENT],
		f func(ctx context.Context, lggr logger.Logger, tx TX, n SendOnlyNode[CHAIN_ID, RPC_CLIENT]))

	runLoop()
	nLiveNodes() (int, int64, *utils.Big)
	report()
}

func ContextWithDefaultTimeout() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), QueryTimeout)
}

type multiNodeClient[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC_CLIENT NodeClient[CHAIN_ID, HEAD],
	TX any,
] struct {
	utils.StartStopOnce
	nodes               []Node[CHAIN_ID, HEAD, RPC_CLIENT]
	sendonlys           []SendOnlyNode[CHAIN_ID, RPC_CLIENT]
	chainID             CHAIN_ID
	logger              logger.Logger
	selectionMode       string
	noNewHeadsThreshold time.Duration
	nodeSelector        NodeSelector[CHAIN_ID, HEAD, RPC_CLIENT]
	chainFamily         string

	activeMu   sync.RWMutex
	activeNode Node[CHAIN_ID, HEAD, RPC_CLIENT]

	chStop utils.StopChan
	wg     sync.WaitGroup
}

func NewMultiNodeClient[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC_CLIENT NodeClient[CHAIN_ID, HEAD],
	TX any,
](
	logger logger.Logger,
	selectionMode string,
	noNewHeadsThreshold time.Duration,
	nodes []Node[CHAIN_ID, HEAD, RPC_CLIENT],
	sendonlys []SendOnlyNode[CHAIN_ID, RPC_CLIENT],
	chainID CHAIN_ID,
	chainFamily string,
) MultiNodeClient[CHAIN_ID, HEAD, RPC_CLIENT, TX] {
	nodeSelector := func() NodeSelector[CHAIN_ID, HEAD, RPC_CLIENT] {
		switch selectionMode {
		case NodeSelectionMode_HighestHead:
			return NewHighestHeadNodeSelector[CHAIN_ID, HEAD, RPC_CLIENT](nodes)
		case NodeSelectionMode_RoundRobin:
			return NewRoundRobinSelector[CHAIN_ID, HEAD, RPC_CLIENT](nodes)
		case NodeSelectionMode_TotalDifficulty:
			return NewTotalDifficultyNodeSelector[CHAIN_ID, HEAD, RPC_CLIENT](nodes)
		case NodeSelectionMode_PriorityLevel:
			return NewPriorityLevelNodeSelector[CHAIN_ID, HEAD, RPC_CLIENT](nodes)
		default:
			panic(fmt.Sprintf("unsupported NodeSelectionMode: %s", selectionMode))
		}
	}()

	lggr := logger.Named("MultiNodeClient").With("chainID", chainID.String())

	c := &multiNodeClient[CHAIN_ID, HEAD, RPC_CLIENT, TX]{
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

	c.logger.Debugf("The MultiNodeClient is configured to use NodeSelectionMode: %s", selectionMode)

	return c
}

// selectNode returns the active Node, if it is still NodeStateAlive, otherwise it selects a new one from the NodeSelector.
func (c *multiNodeClient[CHAIN_ID, HEAD, RPC_CLIENT, TX]) SelectNode() (node Node[CHAIN_ID, HEAD, RPC_CLIENT]) {
	c.activeMu.RLock()
	node = c.activeNode
	c.activeMu.RUnlock()
	if node != nil && node.State() == nodeStateAlive {
		return // still alive
	}

	// select a new one
	c.activeMu.Lock()
	defer c.activeMu.Unlock()
	node = c.activeNode
	if node != nil && node.State() == nodeStateAlive {
		return // another goroutine beat us here
	}

	c.activeNode = c.nodeSelector.Select()

	if c.activeNode == nil {
		c.logger.Criticalw("No live RPC nodes available", "NodeSelectionMode", c.nodeSelector.Name())
		errmsg := fmt.Errorf("no live nodes available for chain %s", c.chainID.String())
		c.SvcErrBuffer.Append(errmsg)
		return &erroringNode[CHAIN_ID, HEAD, RPC_CLIENT]{errMsg: errmsg.Error()}
	}

	return c.activeNode
}

func (c *multiNodeClient[CHAIN_ID, HEAD, RPC_CLIENT, TX]) Dial(ctx context.Context) error {
	return c.StartOnce("MultiNodeClient", func() (merr error) {
		if len(c.nodes) == 0 {
			return errors.Errorf("no available nodes for chain %s", c.chainID.String())
		}
		var ms services.MultiStart
		for _, n := range c.nodes {
			if n.ConfiguredChainID().String() != c.chainID.String() {
				return ms.CloseBecause(errors.Errorf("node %s has chain ID %s which does not match client chain ID of %s", n.String(), n.ConfiguredChainID().String(), c.chainID.String()))
			}
			rawNode, ok := n.(*node[CHAIN_ID, HEAD, RPC_CLIENT])
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
				return ms.CloseBecause(errors.Errorf("sendonly node %s has chain ID %s which does not match client chain ID of %s", s.String(), s.ConfiguredChainID().String(), c.chainID.String()))
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

func (c *multiNodeClient[CHAIN_ID, HEAD, RPC_CLIENT, TX]) NodesAsSendOnlys() (nodes []SendOnlyNode[CHAIN_ID, RPC_CLIENT]) {
	for _, n := range c.nodes {
		nodes = append(nodes, n)
	}
	nodes = append(nodes, c.sendonlys...)
	return
}

// Close tears down the pool and closes all nodes
func (c *multiNodeClient[CHAIN_ID, HEAD, RPC_CLIENT, TX]) Close() error {
	return c.StopOnce("MultiNodeClient", func() error {
		close(c.chStop)
		c.wg.Wait()

		return services.CloseAll(services.MultiCloser(c.nodes), services.MultiCloser(c.sendonlys))
	})
}

func (c *multiNodeClient[CHAIN_ID, HEAD, RPC_CLIENT, TX]) NodeStates() (states map[string]string) {
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
func (c *multiNodeClient[CHAIN_ID, HEAD, RPC_CLIENT, TX]) nLiveNodes() (nLiveNodes int, blockNumber int64, totalDifficulty *utils.Big) {
	totalDifficulty = utils.NewBigI(0)
	for _, n := range c.nodes {
		if s, num, td := n.StateAndLatest(); s == nodeStateAlive {
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

func (c *multiNodeClient[CHAIN_ID, HEAD, RPC_CLIENT, TX]) runLoop() {
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

func (c *multiNodeClient[CHAIN_ID, HEAD, RPC_CLIENT, TX]) report() {
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
		if state != nodeStateAlive {
			dead++
		}
		counts[state]++
	}
	for _, state := range allNodeStates {
		count := counts[state]
		PromMultiNodeClientRPCNodeStates.WithLabelValues(c.chainFamily, c.chainID.String(), state.String()).Set(float64(count))
	}

	live := total - dead
	c.logger.Tracew(fmt.Sprintf("MultiNodeClient state: %d/%d nodes are alive", live, total), "nodeStates", nodeStates)
	if total == dead {
		rerr := fmt.Errorf("no primary nodes available: 0/%d nodes are alive", total)
		c.logger.Criticalw(rerr.Error(), "nodeStates", nodeStates)
		c.SvcErrBuffer.Append(rerr)
	} else if dead > 0 {
		c.logger.Errorw(fmt.Sprintf("At least one primary node is dead: %d/%d nodes are alive", live, total), "nodeStates", nodeStates)
	}
}

func (c *multiNodeClient[CHAIN_ID, HEAD, RPC_CLIENT, TX]) WrapSendOnlyTransaction(ctx context.Context, lggr logger.Logger, tx TX, n SendOnlyNode[CHAIN_ID, RPC_CLIENT],
	f func(ctx context.Context, lggr logger.Logger, tx TX, n SendOnlyNode[CHAIN_ID, RPC_CLIENT]),
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
		c.logger.Debug("Cannot send transaction on sendonly node; multinodeclient is stopped", "node", n.String())
	}
}
