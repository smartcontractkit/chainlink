package mercurytransmitter

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	heap "github.com/esote/minmaxheap"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

type asyncDeleter interface {
	AsyncDelete(hash [32]byte)
	DonID() uint32
}

var _ services.Service = (*transmitQueue)(nil)

var promTransmitQueueLoad = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "llo",
	Subsystem: "mercurytransmitter",
	Name:      "transmit_queue_load",
	Help:      "Current count of items in the transmit queue",
},
	[]string{"donID", "serverURL", "capacity"},
)

// Prometheus' default interval is 15s, set this to under 7.5s to avoid
// aliasing (see: https://en.wikipedia.org/wiki/Nyquist_frequency)
const promInterval = 6500 * time.Millisecond

// TransmitQueue is the high-level package that everything outside of this file should be using
// It stores pending transmissions, yielding the latest (highest priority) first to the caller
type transmitQueue struct {
	services.StateMachine

	cond         sync.Cond
	lggr         logger.SugaredLogger
	asyncDeleter asyncDeleter
	mu           *sync.RWMutex

	pq     *priorityQueue
	maxlen int
	closed bool

	// monitor loop
	stopMonitor       func()
	transmitQueueLoad prometheus.Gauge
}

type TransmitQueue interface {
	services.Service

	BlockingPop() (t *Transmission)
	Push(t *Transmission) (ok bool)
	Init(ts []*Transmission) error
	IsEmpty() bool
}

// maxlen controls how many items will be stored in the queue
// 0 means unlimited - be careful, this can cause memory leaks
func NewTransmitQueue(lggr logger.Logger, serverURL string, maxlen int, asyncDeleter asyncDeleter) TransmitQueue {
	mu := new(sync.RWMutex)
	return &transmitQueue{
		services.StateMachine{},
		sync.Cond{L: mu},
		logger.Sugared(lggr).Named("TransmitQueue"),
		asyncDeleter,
		mu,
		nil, // pq needs to be initialized by calling tq.Init before use
		maxlen,
		false,
		nil,
		promTransmitQueueLoad.WithLabelValues(strconv.FormatUint(uint64(asyncDeleter.DonID()), 10), serverURL, strconv.FormatInt(int64(maxlen), 10)),
	}
}

func (tq *transmitQueue) Init(ts []*Transmission) error {
	if len(ts) > tq.maxlen {
		return fmt.Errorf("transmit queue is too small to hold %d transmissions", len(ts))
	}
	tq.lggr.Debugw("Initializing transmission queue", "nTransmissions", len(ts), "maxlen", tq.maxlen)
	pq := priorityQueue(ts)
	heap.Init(&pq) // ensure the heap is ordered
	tq.pq = &pq
	return nil
}

func (tq *transmitQueue) Push(t *Transmission) (ok bool) {
	tq.cond.L.Lock()
	defer tq.cond.L.Unlock()

	if tq.closed {
		return false
	}

	if tq.maxlen != 0 {
		for tq.pq.Len() >= tq.maxlen {
			// evict oldest entries to make room
			removed := heap.PopMax(tq.pq)
			lggr := tq.lggr
			if removed, ok := removed.(*Transmission); ok {
				hash := removed.Hash()
				lggr = lggr.With("transmissionHash", hex.EncodeToString(hash[:]))
				tq.asyncDeleter.AsyncDelete(hash)
			}
			lggr.Criticalw(fmt.Sprintf("Transmit queue is full; dropping oldest transmission (reached max length of %d)", tq.maxlen), "transmission", removed)
		}
	}

	heap.Push(tq.pq, t)
	tq.cond.Signal()

	return true
}

// BlockingPop will block until at least one item is in the heap, and then return it
// If the queue is closed, it will immediately return nil
func (tq *transmitQueue) BlockingPop() (t *Transmission) {
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

func (tq *transmitQueue) IsEmpty() bool {
	return tq.Len() == 0
}

func (tq *transmitQueue) Len() int {
	tq.cond.L.Lock()
	defer tq.cond.L.Unlock()

	sz := tq.pq.Len()
	tq.cond.Signal()
	return sz
}

func (tq *transmitQueue) Start(context.Context) error {
	return tq.StartOnce("TransmitQueue", func() error {
		t := services.NewTicker(promInterval)
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

func (tq *transmitQueue) Close() error {
	return tq.StopOnce("TransmitQueue", func() error {
		tq.cond.L.Lock()
		tq.closed = true
		tq.cond.L.Unlock()
		tq.cond.Broadcast()
		tq.stopMonitor()
		return nil
	})
}

func (tq *transmitQueue) monitorLoop(c <-chan time.Time, chStop <-chan struct{}, wg *sync.WaitGroup) {
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

func (tq *transmitQueue) report() {
	tq.mu.RLock()
	length := tq.pq.Len()
	tq.mu.RUnlock()
	tq.transmitQueueLoad.Set(float64(length))
}

func (tq *transmitQueue) Ready() error {
	return nil
}
func (tq *transmitQueue) Name() string { return tq.lggr.Name() }
func (tq *transmitQueue) HealthReport() map[string]error {
	report := map[string]error{tq.Name(): errors.Join(
		tq.status(),
	)}
	return report
}

func (tq *transmitQueue) status() (merr error) {
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
func (tq *transmitQueue) pop() *Transmission {
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
	// i.e. a later seqNr is "less" than an earlier one
	return pq[i].SeqNr > pq[j].SeqNr
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Pop() any {
	n := len(*pq)
	if n == 0 {
		return nil
	}
	old := *pq
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

func (pq *priorityQueue) Push(x any) {
	*pq = append(*pq, x.(*Transmission))
}
