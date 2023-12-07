package client

import (
	"net/url"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type testNodeConfig struct {
	pollFailureThreshold uint32
	pollInterval         time.Duration
	selectionMode        string
	syncThreshold        uint32
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

type testNode struct {
	*node[types.ID, Head, NodeClient[types.ID, Head]]
}

type testNodeOpts struct {
	config              testNodeConfig
	noNewHeadsThreshold time.Duration
	lggr                logger.Logger
	wsuri               url.URL
	httpuri             *url.URL
	name                string
	id                  int32
	chainID             types.ID
	nodeOrder           int32
	rpc                 *mockNodeClient[types.ID, Head]
	chainFamily         string
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

	nodeI := NewNode[types.ID, Head, NodeClient[types.ID, Head]](opts.config, opts.noNewHeadsThreshold, opts.lggr,
		opts.wsuri, opts.httpuri, opts.name, opts.id, opts.chainID, opts.nodeOrder, opts.rpc, opts.chainFamily)

	return testNode{
		nodeI.(*node[types.ID, Head, NodeClient[types.ID, Head]]),
	}
}
