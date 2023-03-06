package protocol

import (
	"container/heap"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// Type safe wrapper around MinHeapTimeToContractReportInternal
type MinHeapTimeToPendingTransmission struct {
	internal MinHeapTimeToPendingTransmissionInternal
}

func (h *MinHeapTimeToPendingTransmission) Push(item MinHeapTimeToPendingTransmissionItem) {
	heap.Push(&h.internal, item)
}

func (h *MinHeapTimeToPendingTransmission) Pop() MinHeapTimeToPendingTransmissionItem {
	return heap.Pop(&h.internal).(MinHeapTimeToPendingTransmissionItem)
}

func (h *MinHeapTimeToPendingTransmission) Peek() MinHeapTimeToPendingTransmissionItem {
	return h.internal[0]
}

func (h *MinHeapTimeToPendingTransmission) Len() int {
	return h.internal.Len()
}

type MinHeapTimeToPendingTransmissionItem struct {
	types.ReportTimestamp
	types.PendingTransmission
}

// Implements heap.Interface and uses interface{} all over the place.
type MinHeapTimeToPendingTransmissionInternal []MinHeapTimeToPendingTransmissionItem

func (pq MinHeapTimeToPendingTransmissionInternal) Len() int { return len(pq) }

func (pq MinHeapTimeToPendingTransmissionInternal) Less(i, j int) bool {
	return pq[i].Time.Before(pq[j].Time)
}

func (pq MinHeapTimeToPendingTransmissionInternal) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *MinHeapTimeToPendingTransmissionInternal) Push(x interface{}) {
	item := x.(MinHeapTimeToPendingTransmissionItem)
	*pq = append(*pq, item)
}

func (pq *MinHeapTimeToPendingTransmissionInternal) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
