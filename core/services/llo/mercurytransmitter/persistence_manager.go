package mercurytransmitter

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/jpillora/backoff"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

const (
	// DeleteQueueMaxSize is a sanity limit to avoid unbounded memory consumption.
	// We should never get anywhere close to this under normal operation.
	DeleteQueueMaxSize = 1_000_000
	// DeleteBatchSize is the max number of transmission records to delete
	// in one query. Setting this larger may reduce overall total transaction
	// load on the DB at the expense of blocking inserts for longer.
	DeleteBatchSize = 1_000
	// FlushDeletesFrequency controls how often we wake up to check if there
	// are records in the delete queue, and if so, attempt to drain the queue
	// and delete them all.
	FlushDeletesFrequency = 15 * time.Second

	// PruneFrequency controls how often we wake up to check to see if the
	// transmissions table has exceeded its allowed size, and if so, truncate
	// it. This should already be automatically handled by the transmission
	// queue calling AsyncDelete, but it's here anyway for safety.
	PruneFrequency = 1 * time.Hour
	// PruneBatchSize is the max number of transmission records to delete in
	// one query when pruning the table.
	PruneBatchSize = 10_000

	// OvertimeDeleteTimeout is the maximum time we will spend trying to delete
	// queued transmissions after exit signal before giving up and logging an
	// error.
	OvertimeDeleteTimeout = 2 * time.Second
)

var (
	promTransmitQueueDeleteErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "llo",
		Subsystem: "mercurytransmitter",
		Name:      "transmit_queue_delete_error_count",
		Help:      "Running count of DB errors when trying to delete an item from the queue DB",
	},
		[]string{"donID", "serverURL"},
	)
)

// persistenceManager scopes an ORM to a single serverURL and handles cleanup
// and asynchronous deletion
type persistenceManager struct {
	lggr      logger.Logger
	orm       ORM
	serverURL string
	donID     uint32

	once   services.StateMachine
	stopCh services.StopChan
	wg     sync.WaitGroup

	deleteMu    sync.Mutex
	deleteQueue [][32]byte

	maxTransmitQueueSize  int
	flushDeletesFrequency time.Duration
	pruneFrequency        time.Duration
	maxAge                time.Duration

	transmitQueueDeleteErrorCount prometheus.Counter
}

func NewPersistenceManager(lggr logger.Logger, orm ORM, serverURL string, maxTransmitQueueSize int, flushDeletesFrequency, pruneFrequency, maxAge time.Duration) *persistenceManager {
	return &persistenceManager{
		logger.Sugared(lggr).Named("LLOPersistenceManager"),
		orm,
		serverURL,
		orm.DonID(),
		services.StateMachine{},
		make(services.StopChan),
		sync.WaitGroup{},
		sync.Mutex{},
		nil,
		maxTransmitQueueSize,
		flushDeletesFrequency,
		pruneFrequency,
		maxAge,
		promTransmitQueueDeleteErrorCount.WithLabelValues(strconv.Itoa(int(orm.DonID())), serverURL),
	}
}

func (pm *persistenceManager) Start(ctx context.Context) error {
	return pm.once.StartOnce("LLOMercuryPersistenceManager", func() error {
		pm.wg.Add(2)
		go pm.runFlushDeletesLoop()
		go pm.runPruneLoop()
		return nil
	})
}

func (pm *persistenceManager) Close() error {
	return pm.once.StopOnce("LLOMercuryPersistenceManager", func() error {
		close(pm.stopCh)
		pm.wg.Wait()
		return nil
	})
}

func (pm *persistenceManager) DonID() uint32 {
	return pm.orm.DonID()
}

func (pm *persistenceManager) AsyncDelete(hash [32]byte) {
	pm.addToDeleteQueue(hash)
}

func (pm *persistenceManager) Load(ctx context.Context) ([]*Transmission, error) {
	return pm.orm.Get(ctx, pm.serverURL, pm.maxTransmitQueueSize, pm.maxAge)
}

func (pm *persistenceManager) runFlushDeletesLoop() {
	defer pm.wg.Done()

	ctx, cancel := pm.stopCh.NewCtx()
	defer cancel()

	ticker := services.TickerConfig{
		// Don't prune right away, wait some time for the application to settle
		// down first
		Initial:   services.DefaultJitter.Apply(pm.flushDeletesFrequency),
		JitterPct: services.DefaultJitter,
	}.NewTicker(pm.flushDeletesFrequency)
	defer ticker.Stop()
	for {
		select {
		case <-pm.stopCh:
			q := pm.resetDeleteQueue()
			if len(q) > 0 {
				// make a final effort to clear the database that goes into
				// overtime
				overtimeCtx, cancel := context.WithTimeout(context.Background(), OvertimeDeleteTimeout)
				pm.deleteTransmissions(overtimeCtx, q, DeleteBatchSize)
				cancel()
				if n := pm.lenDeleteQueue(); n > 0 {
					pm.lggr.Errorw("Exiting with undeleted transmissions", "n", n)
				}
			}
			return
		case <-ticker.C:
			queuedTransmissionHashes := pm.resetDeleteQueue()
			if len(queuedTransmissionHashes) == 0 {
				continue
			}
			pm.deleteTransmissions(ctx, queuedTransmissionHashes, DeleteBatchSize)
		}
	}
}

// deleteTransmissions blocks until transmissions are deleted or context is canceled
// it auto-retries on errors
func (pm *persistenceManager) deleteTransmissions(ctx context.Context, hashes [][32]byte, batchSize int) {
	// Exponential backoff for very rarely occurring errors (DB disconnect etc)
	b := backoff.Backoff{
		Min:    10 * time.Millisecond,
		Max:    1 * time.Second,
		Factor: 2,
		Jitter: true,
	}

	for i := 0; i < len(hashes); i += batchSize { // batch deletes to avoid large transactions
		end := i + batchSize
		if end > len(hashes) {
			end = len(hashes)
		}
		deleteBatch := hashes[i:end]
		for {
			if err := pm.orm.Delete(ctx, deleteBatch); err != nil {
				pm.lggr.Errorw("Failed to delete queued transmit requests", "err", err)
				pm.transmitQueueDeleteErrorCount.Inc()
				select {
				case <-time.After(b.Duration()):
					// Wait a backoff duration before trying to delete again
					continue
				case <-ctx.Done():
					// put undeleted items back on the queue and exit
					pm.addToDeleteQueue(hashes[i:]...)
					return
				}
			}
			break
		}
	}
	pm.lggr.Debugw("Flushed delete queue", "nDeleted", len(hashes))
}

func (pm *persistenceManager) runPruneLoop() {
	defer pm.wg.Done()

	ctx, cancel := pm.stopCh.NewCtx()
	defer cancel()

	ticker := services.NewTicker(pm.pruneFrequency)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			overtimeCtx, cancel := context.WithTimeout(context.Background(), OvertimeDeleteTimeout)
			n, err := pm.orm.Prune(overtimeCtx, pm.serverURL, pm.maxTransmitQueueSize, PruneBatchSize)
			cancel()
			if err != nil {
				pm.lggr.Errorw("Failed to truncate transmit requests table on close", "err", err)
			} else if n > 0 {
				pm.lggr.Debugw("Truncated transmit requests table on close", "nDeleted", n)
			}
			return
		case <-ticker.C:
			n, err := pm.orm.Prune(ctx, pm.serverURL, pm.maxTransmitQueueSize, PruneBatchSize)
			if err != nil {
				pm.lggr.Errorw("Failed to prune transmit requests table", "err", err)
				continue
			}
			if n > 0 {
				pm.lggr.Debugw("Pruned transmit requests table", "nDeleted", n)
			}
		}
	}
}

func (pm *persistenceManager) addToDeleteQueue(hashes ...[32]byte) {
	pm.deleteMu.Lock()
	defer pm.deleteMu.Unlock()
	pm.deleteQueue = append(pm.deleteQueue, hashes...)
	if len(pm.deleteQueue) > DeleteQueueMaxSize {
		// NOTE: This could only happen if inserts are succeeding while deletes are
		// failing (or not fast enough) which would be very strange
		pm.lggr.Errorw("Delete queue is full; dropping transmissions", "hashes", hashes, "n", len(pm.deleteQueue))
		pm.deleteQueue = pm.deleteQueue[:DeleteQueueMaxSize]
	}
}

func (pm *persistenceManager) resetDeleteQueue() [][32]byte {
	pm.deleteMu.Lock()
	defer pm.deleteMu.Unlock()
	queue := pm.deleteQueue
	pm.deleteQueue = nil
	return queue
}

func (pm *persistenceManager) lenDeleteQueue() int {
	pm.deleteMu.Lock()
	defer pm.deleteMu.Unlock()
	return len(pm.deleteQueue)
}
