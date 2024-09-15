package mercury

import (
	"context"
	"sync"
	"time"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

var (
	flushDeletesFrequency = time.Second
	pruneFrequency        = time.Hour
)

type PersistenceManager struct {
	lggr      logger.Logger
	orm       ORM
	serverURL string

	once   services.StateMachine
	stopCh services.StopChan
	wg     sync.WaitGroup

	deleteMu    sync.Mutex
	deleteQueue []*pb.TransmitRequest

	jobID int32

	maxTransmitQueueSize  int
	flushDeletesFrequency time.Duration
	pruneFrequency        time.Duration
}

func NewPersistenceManager(lggr logger.Logger, serverURL string, orm ORM, jobID int32, maxTransmitQueueSize int, flushDeletesFrequency, pruneFrequency time.Duration) *PersistenceManager {
	return &PersistenceManager{
		lggr:                  logger.Sugared(lggr).Named("MercuryPersistenceManager").With("serverURL", serverURL),
		orm:                   orm,
		serverURL:             serverURL,
		stopCh:                make(services.StopChan),
		jobID:                 jobID,
		maxTransmitQueueSize:  maxTransmitQueueSize,
		flushDeletesFrequency: flushDeletesFrequency,
		pruneFrequency:        pruneFrequency,
	}
}

func (pm *PersistenceManager) Start(ctx context.Context) error {
	return pm.once.StartOnce("MercuryPersistenceManager", func() error {
		pm.wg.Add(2)
		go pm.runFlushDeletesLoop()
		go pm.runPruneLoop()
		return nil
	})
}

func (pm *PersistenceManager) Close() error {
	return pm.once.StopOnce("MercuryPersistenceManager", func() error {
		close(pm.stopCh)
		pm.wg.Wait()
		return nil
	})
}

func (pm *PersistenceManager) Insert(ctx context.Context, req *pb.TransmitRequest, reportCtx ocrtypes.ReportContext) error {
	return pm.orm.InsertTransmitRequest(ctx, []string{pm.serverURL}, req, pm.jobID, reportCtx)
}

func (pm *PersistenceManager) Delete(ctx context.Context, req *pb.TransmitRequest) error {
	return pm.orm.DeleteTransmitRequests(ctx, pm.serverURL, []*pb.TransmitRequest{req})
}

func (pm *PersistenceManager) AsyncDelete(req *pb.TransmitRequest) {
	pm.addToDeleteQueue(req)
}

func (pm *PersistenceManager) Load(ctx context.Context) ([]*Transmission, error) {
	return pm.orm.GetTransmitRequests(ctx, pm.serverURL, pm.jobID)
}

func (pm *PersistenceManager) runFlushDeletesLoop() {
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
			queuedReqs := pm.resetDeleteQueue()
			if err := pm.orm.DeleteTransmitRequests(ctx, pm.serverURL, queuedReqs); err != nil {
				pm.lggr.Errorw("Failed to delete queued transmit requests", "err", err)
				pm.addToDeleteQueue(queuedReqs...)
			} else {
				pm.lggr.Debugw("Deleted queued transmit requests")
			}
		}
	}
}

func (pm *PersistenceManager) runPruneLoop() {
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
				if err := pm.orm.PruneTransmitRequests(ctx, pm.serverURL, pm.jobID, pm.maxTransmitQueueSize); err != nil {
					pm.lggr.Errorw("Failed to prune transmit requests table", "err", err)
				} else {
					pm.lggr.Debugw("Pruned transmit requests table")
				}
			}(ctx)
		}
	}
}

func (pm *PersistenceManager) addToDeleteQueue(reqs ...*pb.TransmitRequest) {
	pm.deleteMu.Lock()
	defer pm.deleteMu.Unlock()
	pm.deleteQueue = append(pm.deleteQueue, reqs...)
}

func (pm *PersistenceManager) resetDeleteQueue() []*pb.TransmitRequest {
	pm.deleteMu.Lock()
	defer pm.deleteMu.Unlock()
	queue := pm.deleteQueue
	pm.deleteQueue = nil
	return queue
}
