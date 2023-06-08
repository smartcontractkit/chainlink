package mercury

import (
	"container/heap"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var _ services.ServiceCtx = (*TransmitQueue)(nil)

var transmitQueueLoad = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "mercury_transmit_queue_load",
	Help: "Percent of transmit queue capacity used",
},
	[]string{"feedID", "capacity"},
)

// Prometheus' default interval is 15s, set this to under 7.5s to avoid
// aliasing (see: https://en.wikipedia.org/wiki/Nyquist_frequency)
const promInterval = 6500 * time.Millisecond

// TransmitQueue is the high-level package that everything outside of this file should be using
// It stores pending transmissions, yielding the latest (highest priority) first to the caller
type TransmitQueue struct {
	utils.StartStopOnce

	cond sync.Cond
	lggr logger.Logger
	mu   *sync.RWMutex

	pq     *priorityQueue
	maxlen int
	closed bool

	// monitor loop
	stopMonitor       func()
	transmitQueueLoad prometheus.Gauge
}

type Transmission struct {
	Req       *pb.TransmitRequest    // the payload to transmit
	ReportCtx ocrtypes.ReportContext // contains priority information (latest epoch/round wins)

	// The index is needed by update and is maintained by the heap.Interface
	// methods
	// It should NOT be set manually
	index int // the index of the item in the heap
}

// maxlen controls how many items will be stored in the queue
// 0 means unlimited - be careful, this can cause memory leaks
func NewTransmitQueue(lggr logger.Logger, feedID string, maxlen int) *TransmitQueue {
	pq := new(priorityQueue)
	heap.Init(pq) // for completeness
	mu := new(sync.RWMutex)
	return &TransmitQueue{
		utils.StartStopOnce{}, sync.Cond{L: mu}, lggr.Named("TransmitQueue"), mu, pq, maxlen, false, nil,
		transmitQueueLoad.WithLabelValues(feedID, fmt.Sprintf("%d", maxlen)),
	}
}

func (tq *TransmitQueue) Push(req *pb.TransmitRequest, reportCtx ocrtypes.ReportContext) (ok bool) {
	tq.cond.L.Lock()
	defer tq.cond.L.Unlock()

	if tq.closed {
		return false
	}

	if tq.maxlen != 0 && tq.pq.Len() == tq.maxlen {
		// evict oldest entry to make room
		tq.lggr.Criticalf("Transmit queue is full; dropping oldest transmission (reached max length of %d)", tq.maxlen)
		heap.Remove(tq.pq, tq.pq.Len()-1)
	}

	heap.Push(tq.pq, &Transmission{req, reportCtx, -1})
	tq.cond.Signal()

	return true
}

// BlockingPop will block until at least one item is in the heap, and then return it
// If the queue is closed, it will immediately return nil
func (tq *TransmitQueue) BlockingPop() (t *Transmission) {
	tq.cond.L.Lock()
	defer tq.cond.L.Unlock()
	if tq.closed {
		return nil
	}
	for t = tq.pop(); t == nil; t = tq.pop() {
		tq.cond.Wait()
		if tq.closed {
			return nil
		}
	}
	return t
}

func (tq *TransmitQueue) IsEmpty() bool {
	tq.mu.RLock()
	defer tq.mu.RUnlock()
	return tq.pq.Len() == 0
}

func (tq *TransmitQueue) Start(context.Context) error {
	return tq.StartOnce("TransmitQueue", func() error {
		t := time.NewTicker(utils.WithJitter(promInterval))
		wg := new(sync.WaitGroup)
		chStop := make(chan struct{})
		tq.stopMonitor = func() {
			t.Stop()
			close(chStop)
			wg.Wait()
		}
		wg.Add(1)
		go tq.monitorLoop(t.C, chStop, wg)
		return nil
	})
}

func (tq *TransmitQueue) Close() error {
	return tq.StopOnce("TransmitQueue", func() error {
		tq.cond.L.Lock()
		tq.closed = true
		tq.cond.L.Unlock()
		tq.cond.Broadcast()
		tq.stopMonitor()
		return nil
	})
}

func (tq *TransmitQueue) monitorLoop(c <-chan time.Time, chStop <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-c:
			tq.report()
		case <-chStop:
			return
		}
	}
}

func (tq *TransmitQueue) report() {
	tq.mu.RLock()
	length := tq.pq.Len()
	tq.mu.RUnlock()
	tq.transmitQueueLoad.Set(float64(length))
}

func (tq *TransmitQueue) Ready() error {
	return nil
}
func (tq *TransmitQueue) Name() string { return tq.lggr.Name() }
func (tq *TransmitQueue) HealthReport() map[string]error {
	report := map[string]error{tq.Name(): errors.Join(
		tq.status(),
	)}
	return report
}

func (tq *TransmitQueue) status() (merr error) {
	tq.mu.RLock()
	length := tq.pq.Len()
	closed := tq.closed
	tq.mu.RUnlock()
	if tq.maxlen != 0 && length > (tq.maxlen/2) {
		merr = errors.Join(merr, fmt.Errorf("transmit priority queue is greater than 50%% full (%d/%d)", length, tq.maxlen))
	}
	if closed {
		merr = errors.New("transmit queue is closed")
	}
	return merr
}

// pop latest Transmission from the heap
// Not thread-safe
func (tq *TransmitQueue) pop() *Transmission {
	if tq.pq.Len() == 0 {
		return nil
	}
	return heap.Pop(tq.pq).(*Transmission)
}

// HEAP
// Adapted from https://pkg.go.dev/container/heap#example-package-PriorityQueue

// WARNING: None of these methods are thread-safe, caller must synchronize

var _ heap.Interface = &priorityQueue{}

type priorityQueue []*Transmission

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	// We want Pop to give us the latest round, so we use greater than here
	// i.e. a later epoch/round is "less" than an earlier one
	return pq[i].ReportCtx.ReportTimestamp.Epoch > pq[j].ReportCtx.ReportTimestamp.Epoch &&
		pq[i].ReportCtx.ReportTimestamp.Round > pq[j].ReportCtx.ReportTimestamp.Round
}

func (pq *priorityQueue) Pop() any {
	n := len(*pq)
	if n == 0 {
		return nil
	}
	old := *pq
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *priorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Transmission)
	item.index = n
	*pq = append(*pq, item)
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
