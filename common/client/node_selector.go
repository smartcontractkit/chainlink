package client

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

const (
	NodeSelectionModeHighestHead     = "HighestHead"
	NodeSelectionModeRoundRobin      = "RoundRobin"
	NodeSelectionModeTotalDifficulty = "TotalDifficulty"
	NodeSelectionModePriorityLevel   = "PriorityLevel"
)

//go:generate mockery --quiet --name NodeSelector --structname mockNodeSelector --filename "mock_node_selector_test.go" --inpackage --case=underscore
type NodeSelector[
	CHAIN_ID types.ID,
	RPC any,
] interface {
	// Select returns a Node, or nil if none can be selected.
	// Implementation must be thread-safe.
	Select() Node[CHAIN_ID, RPC]
	// Name returns the strategy name, e.g. "HighestHead" or "RoundRobin"
	Name() string
}

func newNodeSelector[
	CHAIN_ID types.ID,
	RPC any,
](selectionMode string, nodes []Node[CHAIN_ID, RPC]) NodeSelector[CHAIN_ID, RPC] {
	switch selectionMode {
	case NodeSelectionModeHighestHead:
		return NewHighestHeadNodeSelector[CHAIN_ID, RPC](nodes)
	case NodeSelectionModeRoundRobin:
		return NewRoundRobinSelector[CHAIN_ID, RPC](nodes)
	case NodeSelectionModeTotalDifficulty:
		return NewTotalDifficultyNodeSelector[CHAIN_ID, RPC](nodes)
	case NodeSelectionModePriorityLevel:
		return NewPriorityLevelNodeSelector[CHAIN_ID, RPC](nodes)
	default:
		panic(fmt.Sprintf("unsupported NodeSelectionMode: %s", selectionMode))
	}
}
