package client

import (
	"net/url"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	clientMocks "github.com/smartcontractkit/chainlink/v2/common/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type testNodeConfig struct {
	pollFailureThreshold       uint32
	pollInterval               time.Duration
	selectionMode              string
	syncThreshold              uint32
	nodeIsSyncingEnabled       bool
	enforceRepeatableRead      bool
	finalizedBlockPollInterval time.Duration
	deathDeclarationDelay      time.Duration
	newHeadsPollInterval       time.Duration
}

func (n testNodeConfig) NewHeadsPollInterval() time.Duration {
	return n.newHeadsPollInterval
}

func (n testNodeConfig) PollFailureThreshold() uint32 {
	return n.pollFailureThreshold
}

func (n testNodeConfig) PollInterval() time.Duration {
	return n.pollInterval
}

func (n testNodeConfig) SelectionMode() string {
	return n.selectionMode
}

func (n testNodeConfig) SyncThreshold() uint32 {
	return n.syncThreshold
}

func (n testNodeConfig) NodeIsSyncingEnabled() bool {
	return n.nodeIsSyncingEnabled
}

func (n testNodeConfig) FinalizedBlockPollInterval() time.Duration {
	return n.finalizedBlockPollInterval
}

func (n testNodeConfig) EnforceRepeatableRead() bool {
	return n.enforceRepeatableRead
}

func (n testNodeConfig) DeathDeclarationDelay() time.Duration {
	return n.deathDeclarationDelay
}

type testNode struct {
	*node[types.ID, Head, NodeClient[types.ID, Head]]
}

type testNodeOpts struct {
	config      testNodeConfig
	chainConfig clientMocks.ChainConfig
	lggr        logger.Logger
	wsuri       *url.URL
	httpuri     *url.URL
	name        string
	id          int
	chainID     types.ID
	nodeOrder   int32
	rpc         *mockNodeClient[types.ID, Head]
	chainFamily string
}

func newTestNode(t *testing.T, opts testNodeOpts) testNode {
	if opts.lggr == nil {
		opts.lggr = logger.Test(t)
	}

	if opts.name == "" {
		opts.name = "tes node"
	}

	if opts.chainFamily == "" {
		opts.chainFamily = "test node chain family"
	}

	if opts.chainID == nil {
		opts.chainID = types.RandomID()
	}

	if opts.id == 0 {
		opts.id = 42
	}

	nodeI := NewNode[types.ID, Head, NodeClient[types.ID, Head]](opts.config, opts.chainConfig, opts.lggr,
		opts.wsuri, opts.httpuri, opts.name, opts.id, opts.chainID, opts.nodeOrder, opts.rpc, opts.chainFamily)

	return testNode{
		nodeI.(*node[types.ID, Head, NodeClient[types.ID, Head]]),
	}
}
