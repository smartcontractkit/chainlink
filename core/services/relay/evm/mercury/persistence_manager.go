package mercury

import (
	"context"
	"sync"
	"time"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	flushDeletesFrequency = time.Second
	pruneFrequency        = time.Hour
)

type PersistenceManager struct {
	lggr logger.Logger
	orm  ORM

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

func NewPersistenceManager(lggr logger.Logger, orm ORM, jobID int32, maxTransmitQueueSize int, flushDeletesFrequency, pruneFrequency time.Duration) *PersistenceManager {
	return &PersistenceManager{
		lggr:                  lggr.Named("MercuryPersistenceManager"),
		orm:                   orm,
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
	return pm.orm.InsertTransmitRequest(req, pm.jobID, reportCtx, pg.WithParentCtx(ctx))
}

func (pm *PersistenceManager) Delete(ctx context.Context, req *pb.TransmitRequest) error {
	return pm.orm.DeleteTransmitRequests([]*pb.TransmitRequest{req}, pg.WithParentCtx(ctx))
}

func (pm *PersistenceManager) AsyncDelete(req *pb.TransmitRequest) {
	pm.addToDeleteQueue(req)
}

func (pm *PersistenceManager) Load(ctx context.Context) ([]*Transmission, error) {
	return pm.orm.GetTransmitRequests(pm.jobID, pg.WithParentCtx(ctx))
}

func (pm *PersistenceManager) runFlushDeletesLoop() {
	defer pm.wg.Done()

	ctx, cancel := pm.stopCh.Ctx(context.Background())
	defer cancel()

	ticker := time.NewTicker(utils.WithJitter(pm.flushDeletesFrequency))
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			queuedReqs := pm.resetDeleteQueue()
			if err := pm.orm.DeleteTransmitRequests(queuedReqs, pg.WithParentCtx(ctx)); err != nil {
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

	ctx, cancel := pm.stopCh.Ctx(context.Background())
	defer cancel()

	ticker := time.NewTicker(utils.WithJitter(pm.pruneFrequency))
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			if err := pm.orm.PruneTransmitRequests(pm.jobID, pm.maxTransmitQueueSize, pg.WithParentCtx(ctx), pg.WithLongQueryTimeout()); err != nil {
				pm.lggr.Errorw("Failed to prune transmit requests table", "err", err)
			} else {
				pm.lggr.Debugw("Pruned transmit requests table")
			}
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
