package client

import (
	"net/url"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	clientMocks "github.com/smartcontractkit/chainlink/v2/common/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type testNode struct {
	*node[types.ID, Head, NodeClient[types.ID, Head]]
}

type testNodeOpts struct {
	config      clientMocks.NodeConfig
	chainConfig clientMocks.ChainConfig
	lggr        logger.Logger
	wsuri       url.URL
	httpuri     *url.URL
	name        string
	id          int32
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
