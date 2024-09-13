package mercurytransmitter

import (
	"context"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

var (
	flushDeletesFrequency = time.Second
	pruneFrequency        = time.Hour
)

// persistenceManager scopes an ORM to a single serverURL and handles cleanup
// and asynchronous deletion
type persistenceManager struct {
	lggr      logger.Logger
	orm       ORM
	serverURL string

	once   services.StateMachine
	stopCh services.StopChan
	wg     sync.WaitGroup

	deleteMu    sync.Mutex
	deleteQueue [][32]byte

	maxTransmitQueueSize  int
	flushDeletesFrequency time.Duration
	pruneFrequency        time.Duration
}

func NewPersistenceManager(lggr logger.Logger, orm ORM, serverURL string, maxTransmitQueueSize int, flushDeletesFrequency, pruneFrequency time.Duration) *persistenceManager {
	return &persistenceManager{
		orm:                   orm,
		lggr:                  logger.Sugared(lggr).Named("LLOPersistenceManager").With("serverURL", serverURL),
		serverURL:             serverURL,
		stopCh:                make(services.StopChan),
		maxTransmitQueueSize:  maxTransmitQueueSize,
		flushDeletesFrequency: flushDeletesFrequency,
		pruneFrequency:        pruneFrequency,
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
	return pm.orm.Get(ctx, pm.serverURL)
}

func (pm *persistenceManager) runFlushDeletesLoop() {
	defer pm.wg.Done()

	ctx, cancel := pm.stopCh.Ctx(context.Background())
	defer cancel()

	ticker := services.NewTicker(pm.flushDeletesFrequency)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			queuedTransmissionHashes := pm.resetDeleteQueue()
			if len(queuedTransmissionHashes) == 0 {
				continue
			}
			if err := pm.orm.Delete(ctx, queuedTransmissionHashes); err != nil {
				pm.lggr.Errorw("Failed to delete queued transmit requests", "err", err)
				pm.addToDeleteQueue(queuedTransmissionHashes...)
			} else {
				pm.lggr.Debugw("Deleted queued transmit requests")
			}
		}
	}
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
			return
		case <-ticker.C:
			func(ctx context.Context) {
				ctx, cancelPrune := context.WithTimeout(sqlutil.WithoutDefaultTimeout(ctx), time.Minute)
				defer cancelPrune()
				if err := pm.orm.Prune(ctx, pm.serverURL, pm.maxTransmitQueueSize); err != nil {
					pm.lggr.Errorw("Failed to prune transmit requests table", "err", err)
				} else {
					pm.lggr.Debugw("Pruned transmit requests table")
				}
			}(ctx)
		}
	}
}

func (pm *persistenceManager) addToDeleteQueue(hashes ...[32]byte) {
	pm.deleteMu.Lock()
	defer pm.deleteMu.Unlock()
	pm.deleteQueue = append(pm.deleteQueue, hashes...)
}

func (pm *persistenceManager) resetDeleteQueue() [][32]byte {
	pm.deleteMu.Lock()
	defer pm.deleteMu.Unlock()
	queue := pm.deleteQueue
	pm.deleteQueue = nil
	return queue
}
