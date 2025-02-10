package headtracker

import (
	evmtypes "github.com/smartcontractkit/chainlink-integrations/evm/types"
)

type headsHeap struct {
	values []*evmtypes.Head
}

func (h *headsHeap) Len() int {
	return len(h.values)
}

func (h *headsHeap) Swap(i, j int) {
	h.values[i], h.values[j] = h.values[j], h.values[i]
}

func (h *headsHeap) Less(i, j int) bool {
	return h.values[i].Number < h.values[j].Number
}

func (h *headsHeap) Pop() any {
	n := len(h.values) - 1
	old := h.values[n]
	h.values[n] = nil
	h.values = h.values[:n]
	return old
}

func (h *headsHeap) Push(v any) {
	h.values = append(h.values, v.(*evmtypes.Head))
}

func (h *headsHeap) Peek() *evmtypes.Head {
	return h.values[0]
}
