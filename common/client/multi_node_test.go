package client

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type multiNodeRPCClient RPCClient[types.ID, types.Head[Hashable]]

type testMultiNode struct {
	*MultiNode[types.ID, multiNodeRPCClient]
}

type multiNodeOpts struct {
	logger                logger.Logger
	selectionMode         string
	leaseDuration         time.Duration
	nodes                 []Node[types.ID, multiNodeRPCClient]
	sendonlys             []SendOnlyNode[types.ID, multiNodeRPCClient]
	chainID               types.ID
	chainFamily           string
	deathDeclarationDelay time.Duration
}

func newTestMultiNode(t *testing.T, opts multiNodeOpts) testMultiNode {
	if opts.logger == nil {
		opts.logger = logger.Test(t)
	}

	result := NewMultiNode[types.ID, multiNodeRPCClient](
		opts.logger, opts.selectionMode, opts.leaseDuration, opts.nodes, opts.sendonlys, opts.chainID, opts.chainFamily, opts.deathDeclarationDelay)
	return testMultiNode{
		result,
	}
}

func newHealthyNode(t *testing.T, chainID types.ID) *mockNode[types.ID, multiNodeRPCClient] {
	return newNodeWithState(t, chainID, nodeStateAlive)
}

func newNodeWithState(t *testing.T, chainID types.ID, state nodeState) *mockNode[types.ID, multiNodeRPCClient] {
	node := newMockNode[types.ID, multiNodeRPCClient](t)
	node.On("ConfiguredChainID").Return(chainID).Once()
	node.On("Start", mock.Anything).Return(nil).Once()
	node.On("Close").Return(nil).Once()
	node.On("String").Return(fmt.Sprintf("healthy_node_%d", rand.Int())).Maybe()
	node.On("SetPoolChainInfoProvider", mock.Anything).Once()
	node.On("State").Return(state).Maybe()
	return node
}

func TestMultiNode_Dial(t *testing.T) {
	t.Parallel()

	newMockNode := newMockNode[types.ID, multiNodeRPCClient]
	newMockSendOnlyNode := newMockSendOnlyNode[types.ID, multiNodeRPCClient]

	t.Run("Fails without nodes", func(t *testing.T) {
		t.Parallel()
		mn := newTestMultiNode(t, multiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       types.RandomID(),
		})
		err := mn.Start(tests.Context(t))
		assert.EqualError(t, err, fmt.Sprintf("no available nodes for chain %s", mn.chainID.String()))
	})
	t.Run("Fails with wrong node's chainID", func(t *testing.T) {
		t.Parallel()
		node := newMockNode(t)
		multiNodeChainID := types.NewIDFromInt(10)
		nodeChainID := types.NewIDFromInt(11)
		node.On("ConfiguredChainID").Return(nodeChainID).Twice()
		const nodeName = "nodeName"
		node.On("String").Return(nodeName).Once()
		mn := newTestMultiNode(t, multiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       multiNodeChainID,
			nodes:         []Node[types.ID, multiNodeRPCClient]{node},
		})
		err := mn.Start(tests.Context(t))
		assert.EqualError(t, err, fmt.Sprintf("node %s has configured chain ID %s which does not match multinode configured chain ID of %s", nodeName, nodeChainID, mn.chainID))
	})
	t.Run("Fails if node fails", func(t *testing.T) {
		t.Parallel()
		node := newMockNode(t)
		chainID := types.RandomID()
		node.On("ConfiguredChainID").Return(chainID).Once()
		expectedError := errors.New("failed to start node")
		node.On("Start", mock.Anything).Return(expectedError).Once()
		node.On("SetPoolChainInfoProvider", mock.Anything).Once()
		mn := newTestMultiNode(t, multiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       chainID,
			nodes:         []Node[types.ID, multiNodeRPCClient]{node},
		})
		err := mn.Start(tests.Context(t))
		assert.EqualError(t, err, expectedError.Error())
	})

	t.Run("Closes started nodes on failure", func(t *testing.T) {
		t.Parallel()
		chainID := types.RandomID()
		node1 := newHealthyNode(t, chainID)
		node2 := newMockNode(t)
		node2.On("ConfiguredChainID").Return(chainID).Once()
		expectedError := errors.New("failed to start node")
		node2.On("Start", mock.Anything).Return(expectedError).Once()
		node2.On("SetPoolChainInfoProvider", mock.Anything).Once()

		mn := newTestMultiNode(t, multiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       chainID,
			nodes:         []Node[types.ID, multiNodeRPCClient]{node1, node2},
		})
		err := mn.Start(tests.Context(t))
		assert.EqualError(t, err, expectedError.Error())
	})
	t.Run("Fails with wrong send only node's chainID", func(t *testing.T) {
		t.Parallel()
		multiNodeChainID := types.NewIDFromInt(10)
		node := newHealthyNode(t, multiNodeChainID)
		sendOnly := newMockSendOnlyNode(t)
		sendOnlyChainID := types.NewIDFromInt(11)
		sendOnly.On("ConfiguredChainID").Return(sendOnlyChainID).Twice()
		const sendOnlyName = "sendOnlyNodeName"
		sendOnly.On("String").Return(sendOnlyName).Once()

		mn := newTestMultiNode(t, multiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       multiNodeChainID,
			nodes:         []Node[types.ID, multiNodeRPCClient]{node},
			sendonlys:     []SendOnlyNode[types.ID, multiNodeRPCClient]{sendOnly},
		})
		err := mn.Start(tests.Context(t))
		assert.EqualError(t, err, fmt.Sprintf("sendonly node %s has configured chain ID %s which does not match multinode configured chain ID of %s", sendOnlyName, sendOnlyChainID, mn.chainID))
	})

	newHealthySendOnly := func(t *testing.T, chainID types.ID) *mockSendOnlyNode[types.ID, multiNodeRPCClient] {
		node := newMockSendOnlyNode(t)
		node.On("ConfiguredChainID").Return(chainID).Once()
		node.On("Start", mock.Anything).Return(nil).Once()
		node.On("Close").Return(nil).Once()
		return node
	}
	t.Run("Fails on send only node failure", func(t *testing.T) {
		t.Parallel()
		chainID := types.NewIDFromInt(10)
		node := newHealthyNode(t, chainID)
		sendOnly1 := newHealthySendOnly(t, chainID)
		sendOnly2 := newMockSendOnlyNode(t)
		sendOnly2.On("ConfiguredChainID").Return(chainID).Once()
		expectedError := errors.New("failed to start send only node")
		sendOnly2.On("Start", mock.Anything).Return(expectedError).Once()

		mn := newTestMultiNode(t, multiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       chainID,
			nodes:         []Node[types.ID, multiNodeRPCClient]{node},
			sendonlys:     []SendOnlyNode[types.ID, multiNodeRPCClient]{sendOnly1, sendOnly2},
		})
		err := mn.Start(tests.Context(t))
		assert.EqualError(t, err, expectedError.Error())
	})
	t.Run("Starts successfully with healthy nodes", func(t *testing.T) {
		t.Parallel()
		chainID := types.NewIDFromInt(10)
		node := newHealthyNode(t, chainID)
		mn := newTestMultiNode(t, multiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       chainID,
			nodes:         []Node[types.ID, multiNodeRPCClient]{node},
			sendonlys:     []SendOnlyNode[types.ID, multiNodeRPCClient]{newHealthySendOnly(t, chainID)},
		})
		defer func() { assert.NoError(t, mn.Close()) }()
		err := mn.Start(tests.Context(t))
		require.NoError(t, err)
		selectedNode, err := mn.selectNode()
		require.NoError(t, err)
		assert.Equal(t, node, selectedNode)
	})
}

func TestMultiNode_Report(t *testing.T) {
	t.Parallel()
	t.Run("Dial starts periodical reporting", func(t *testing.T) {
		t.Parallel()
		chainID := types.RandomID()
		node1 := newHealthyNode(t, chainID)
		node2 := newNodeWithState(t, chainID, nodeStateOutOfSync)
		lggr, observedLogs := logger.TestObserved(t, zap.WarnLevel)
		mn := newTestMultiNode(t, multiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       chainID,
			nodes:         []Node[types.ID, multiNodeRPCClient]{node1, node2},
			logger:        lggr,
		})
		mn.reportInterval = tests.TestInterval
		mn.deathDeclarationDelay = tests.TestInterval
		defer func() { assert.NoError(t, mn.Close()) }()
		err := mn.Start(tests.Context(t))
		require.NoError(t, err)
		tests.AssertLogCountEventually(t, observedLogs, "At least one primary node is dead: 1/2 nodes are alive", 2)
	})
	t.Run("Report critical error on all node failure", func(t *testing.T) {
		t.Parallel()
		chainID := types.RandomID()
		node := newNodeWithState(t, chainID, nodeStateOutOfSync)
		lggr, observedLogs := logger.TestObserved(t, zap.WarnLevel)
		mn := newTestMultiNode(t, multiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       chainID,
			nodes:         []Node[types.ID, multiNodeRPCClient]{node},
			logger:        lggr,
		})
		mn.reportInterval = tests.TestInterval
		mn.deathDeclarationDelay = tests.TestInterval
		defer func() { assert.NoError(t, mn.Close()) }()
		err := mn.Start(tests.Context(t))
		require.NoError(t, err)
		tests.AssertLogCountEventually(t, observedLogs, "no primary nodes available: 0/1 nodes are alive", 2)
		err = mn.Healthy()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no primary nodes available: 0/1 nodes are alive")
	})
}

func TestMultiNode_CheckLease(t *testing.T) {
	t.Parallel()
	t.Run("Round robin disables lease check", func(t *testing.T) {
		t.Parallel()
		chainID := types.RandomID()
		node := newHealthyNode(t, chainID)
		lggr, observedLogs := logger.TestObserved(t, zap.InfoLevel)
		mn := newTestMultiNode(t, multiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       chainID,
			logger:        lggr,
			nodes:         []Node[types.ID, multiNodeRPCClient]{node},
		})
		defer func() { assert.NoError(t, mn.Close()) }()
		err := mn.Start(tests.Context(t))
		require.NoError(t, err)
		tests.RequireLogMessage(t, observedLogs, "Best node switching is disabled")
	})
	t.Run("Misconfigured lease check period won't start", func(t *testing.T) {
		t.Parallel()
		chainID := types.RandomID()
		node := newHealthyNode(t, chainID)
		lggr, observedLogs := logger.TestObserved(t, zap.InfoLevel)
		mn := newTestMultiNode(t, multiNodeOpts{
			selectionMode: NodeSelectionModeHighestHead,
			chainID:       chainID,
			logger:        lggr,
			nodes:         []Node[types.ID, multiNodeRPCClient]{node},
			leaseDuration: 0,
		})
		defer func() { assert.NoError(t, mn.Close()) }()
		err := mn.Start(tests.Context(t))
		require.NoError(t, err)
		tests.RequireLogMessage(t, observedLogs, "Best node switching is disabled")
	})
	t.Run("Lease check updates active node", func(t *testing.T) {
		t.Parallel()
		chainID := types.RandomID()
		node := newHealthyNode(t, chainID)
		node.On("UnsubscribeAllExceptAliveLoop")
		bestNode := newHealthyNode(t, chainID)
		nodeSelector := newMockNodeSelector[types.ID, multiNodeRPCClient](t)
		nodeSelector.On("Select").Return(bestNode)
		lggr, observedLogs := logger.TestObserved(t, zap.InfoLevel)
		mn := newTestMultiNode(t, multiNodeOpts{
			selectionMode: NodeSelectionModeHighestHead,
			chainID:       chainID,
			logger:        lggr,
			nodes:         []Node[types.ID, multiNodeRPCClient]{node, bestNode},
			leaseDuration: tests.TestInterval,
		})
		defer func() { assert.NoError(t, mn.Close()) }()
		mn.nodeSelector = nodeSelector
		err := mn.Start(tests.Context(t))
		require.NoError(t, err)
		tests.AssertLogEventually(t, observedLogs, fmt.Sprintf("Switching to best node from %q to %q", node.String(), bestNode.String()))
		tests.AssertEventually(t, func() bool {
			mn.activeMu.RLock()
			active := mn.activeNode
			mn.activeMu.RUnlock()
			return bestNode == active
		})
	})
	t.Run("NodeStates returns proper states", func(t *testing.T) {
		t.Parallel()
		chainID := types.NewIDFromInt(10)
		nodes := map[string]nodeState{
			"node_1": nodeStateAlive,
			"node_2": nodeStateUnreachable,
			"node_3": nodeStateDialed,
		}

		opts := multiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       chainID,
		}

		expectedResult := map[string]string{}
		for name, state := range nodes {
			node := newMockNode[types.ID, multiNodeRPCClient](t)
			node.On("State").Return(state).Once()
			node.On("String").Return(name).Once()
			opts.nodes = append(opts.nodes, node)

			sendOnly := newMockSendOnlyNode[types.ID, multiNodeRPCClient](t)
			sendOnlyName := "send_only_" + name
			sendOnly.On("State").Return(state).Once()
			sendOnly.On("String").Return(sendOnlyName).Once()
			opts.sendonlys = append(opts.sendonlys, sendOnly)

			expectedResult[name] = state.String()
			expectedResult[sendOnlyName] = state.String()
		}

		mn := newTestMultiNode(t, opts)
		states := mn.NodeStates()
		assert.Equal(t, expectedResult, states)
	})
}

func TestMultiNode_selectNode(t *testing.T) {
	t.Parallel()
	t.Run("Returns same node, if it's still healthy", func(t *testing.T) {
		t.Parallel()
		chainID := types.RandomID()
		node1 := newMockNode[types.ID, multiNodeRPCClient](t)
		node1.On("State").Return(nodeStateAlive).Once()
		node1.On("String").Return("node1").Maybe()
		node2 := newMockNode[types.ID, multiNodeRPCClient](t)
		node2.On("String").Return("node2").Maybe()
		mn := newTestMultiNode(t, multiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       chainID,
			nodes:         []Node[types.ID, multiNodeRPCClient]{node1, node2},
		})
		nodeSelector := newMockNodeSelector[types.ID, multiNodeRPCClient](t)
		nodeSelector.On("Select").Return(node1).Once()
		mn.nodeSelector = nodeSelector
		prevActiveNode, err := mn.selectNode()
		require.NoError(t, err)
		require.Equal(t, node1.String(), prevActiveNode.String())
		newActiveNode, err := mn.selectNode()
		require.NoError(t, err)
		require.Equal(t, prevActiveNode.String(), newActiveNode.String())
	})
	t.Run("Updates node if active is not healthy", func(t *testing.T) {
		t.Parallel()
		chainID := types.RandomID()
		oldBest := newMockNode[types.ID, multiNodeRPCClient](t)
		oldBest.On("String").Return("oldBest").Maybe()
		oldBest.On("UnsubscribeAllExceptAliveLoop")
		newBest := newMockNode[types.ID, multiNodeRPCClient](t)
		newBest.On("String").Return("newBest").Maybe()
		mn := newTestMultiNode(t, multiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       chainID,
			nodes:         []Node[types.ID, multiNodeRPCClient]{oldBest, newBest},
		})
		nodeSelector := newMockNodeSelector[types.ID, multiNodeRPCClient](t)
		nodeSelector.On("Select").Return(oldBest).Once()
		mn.nodeSelector = nodeSelector
		activeNode, err := mn.selectNode()
		require.NoError(t, err)
		require.Equal(t, oldBest.String(), activeNode.String())
		// old best died, so we should replace it
		oldBest.On("State").Return(nodeStateOutOfSync).Twice()
		nodeSelector.On("Select").Return(newBest).Once()
		newActiveNode, err := mn.selectNode()
		require.NoError(t, err)
		require.Equal(t, newBest.String(), newActiveNode.String())
	})
	t.Run("No active nodes - reports critical error", func(t *testing.T) {
		t.Parallel()
		chainID := types.RandomID()
		lggr, observedLogs := logger.TestObserved(t, zap.InfoLevel)
		mn := newTestMultiNode(t, multiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       chainID,
			logger:        lggr,
		})
		nodeSelector := newMockNodeSelector[types.ID, multiNodeRPCClient](t)
		nodeSelector.On("Select").Return(nil).Once()
		nodeSelector.On("Name").Return("MockedNodeSelector").Once()
		mn.nodeSelector = nodeSelector
		node, err := mn.selectNode()
		require.EqualError(t, err, ErroringNodeError.Error())
		require.Nil(t, node)
		tests.RequireLogMessage(t, observedLogs, "No live RPC nodes available")
	})
}

func TestMultiNode_ChainInfo(t *testing.T) {
	t.Parallel()
	type nodeParams struct {
		LatestChainInfo         ChainInfo
		HighestUserObservations ChainInfo
		State                   nodeState
	}
	testCases := []struct {
		Name                            string
		ExpectedNLiveNodes              int
		ExpectedLatestChainInfo         ChainInfo
		ExpectedHighestUserObservations ChainInfo
		NodeParams                      []nodeParams
	}{
		{
			Name: "no nodes",
			ExpectedLatestChainInfo: ChainInfo{
				TotalDifficulty: big.NewInt(0),
			},
			ExpectedHighestUserObservations: ChainInfo{
				TotalDifficulty: big.NewInt(0),
			},
		},
		{
			Name:               "Best node is not healthy",
			ExpectedNLiveNodes: 3,
			ExpectedLatestChainInfo: ChainInfo{
				BlockNumber:          20,
				FinalizedBlockNumber: 10,
				TotalDifficulty:      big.NewInt(10),
			},
			ExpectedHighestUserObservations: ChainInfo{
				BlockNumber:          1005,
				FinalizedBlockNumber: 995,
				TotalDifficulty:      big.NewInt(2005),
			},
			NodeParams: []nodeParams{
				{
					State: nodeStateOutOfSync,
					LatestChainInfo: ChainInfo{
						BlockNumber:          1000,
						FinalizedBlockNumber: 990,
						TotalDifficulty:      big.NewInt(2000),
					},
					HighestUserObservations: ChainInfo{
						BlockNumber:          1005,
						FinalizedBlockNumber: 995,
						TotalDifficulty:      big.NewInt(2005),
					},
				},
				{
					State: nodeStateAlive,
					LatestChainInfo: ChainInfo{
						BlockNumber:          20,
						FinalizedBlockNumber: 10,
						TotalDifficulty:      big.NewInt(9),
					},
					HighestUserObservations: ChainInfo{
						BlockNumber:          25,
						FinalizedBlockNumber: 15,
						TotalDifficulty:      big.NewInt(14),
					},
				},
				{
					State: nodeStateAlive,
					LatestChainInfo: ChainInfo{
						BlockNumber:          19,
						FinalizedBlockNumber: 9,
						TotalDifficulty:      big.NewInt(10),
					},
					HighestUserObservations: ChainInfo{
						BlockNumber:          24,
						FinalizedBlockNumber: 14,
						TotalDifficulty:      big.NewInt(15),
					},
				},
				{
					State: nodeStateAlive,
					LatestChainInfo: ChainInfo{
						BlockNumber:          11,
						FinalizedBlockNumber: 1,
						TotalDifficulty:      nil,
					},
					HighestUserObservations: ChainInfo{
						BlockNumber:          16,
						FinalizedBlockNumber: 6,
						TotalDifficulty:      nil,
					},
				},
			},
		},
	}

	chainID := types.RandomID()
	mn := newTestMultiNode(t, multiNodeOpts{
		selectionMode: NodeSelectionModeRoundRobin,
		chainID:       chainID,
	})
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			for _, params := range tc.NodeParams {
				node := newMockNode[types.ID, multiNodeRPCClient](t)
				mn.primaryNodes = append(mn.primaryNodes, node)
				node.On("StateAndLatest").Return(params.State, params.LatestChainInfo)
				node.On("HighestUserObservations").Return(params.HighestUserObservations)
			}

			nNodes, latestChainInfo := mn.LatestChainInfo()
			assert.Equal(t, tc.ExpectedNLiveNodes, nNodes)
			assert.Equal(t, tc.ExpectedLatestChainInfo, latestChainInfo)

			highestChainInfo := mn.HighestUserObservations()
			assert.Equal(t, tc.ExpectedHighestUserObservations, highestChainInfo)
		})
	}
}
