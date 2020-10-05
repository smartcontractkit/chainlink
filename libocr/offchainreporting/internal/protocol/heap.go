package protocol

import (
	"container/heap"

	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)

type MinHeapTimeToContractReport struct {
	internal MinHeapTimeToContractReportInternal
}

func (h *MinHeapTimeToContractReport) Push(item MinHeapTimeToContractReportItem) {
	heap.Push(&h.internal, item)
}

func (h *MinHeapTimeToContractReport) Pop() MinHeapTimeToContractReportItem {
	return heap.Pop(&h.internal).(MinHeapTimeToContractReportItem)
}

func (h *MinHeapTimeToContractReport) Peek() MinHeapTimeToContractReportItem {
	return h.internal[0]
}

func (h *MinHeapTimeToContractReport) Len() int {
	return h.internal.Len()
}

type MinHeapTimeToContractReportItem struct {
	types.PendingTransmissionKey
	types.PendingTransmission
}

type MinHeapTimeToContractReportInternal []MinHeapTimeToContractReportItem

func (pq MinHeapTimeToContractReportInternal) Len() int { return len(pq) }

func (pq MinHeapTimeToContractReportInternal) Less(i, j int) bool {
	return pq[i].Time.Before(pq[j].Time)
}

func (pq MinHeapTimeToContractReportInternal) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *MinHeapTimeToContractReportInternal) Push(x interface{}) {
	item := x.(MinHeapTimeToContractReportItem)
	*pq = append(*pq, item)
}

func (pq *MinHeapTimeToContractReportInternal) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
