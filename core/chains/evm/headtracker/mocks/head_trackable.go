package mocks

import (
	"context"
	"sync/atomic"

	"github.com/smartcontractkit/chainlink-integrations/evm/types"
)

// MockHeadTrackable allows you to mock HeadTrackable
type MockHeadTrackable struct {
	onNewHeadCount atomic.Int32
}

// OnNewLongestChain increases the OnNewLongestChainCount count by one
func (m *MockHeadTrackable) OnNewLongestChain(context.Context, *types.Head) {
	m.onNewHeadCount.Add(1)
}

// OnNewLongestChainCount returns the count of new heads, safely.
func (m *MockHeadTrackable) OnNewLongestChainCount() int32 {
	return m.onNewHeadCount.Load()
}
