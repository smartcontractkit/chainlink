package secrets

import (
	"context"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils/chans"
	"github.com/smartcontractkit/chainlink/v2/core/utils/signalers"
)

var (
	defaultTickInterval = 12 * time.Second
	ContractName        = "WorkflowRegistry"
	ContractEventName   = "WorkflowForceUpdateSecretsRequestedV1"
)

type Syncer interface {
	Start(ctx context.Context) error
	Close() error
	SecretsFor(ctx context.Context, workflowOwner, workflowName string) (map[string]string, error)
}

type FetcherFunc func(ctx context.Context, url string) ([]byte, error)

type worker struct {
	services.StateMachine
	stopCh   services.StopChan
	eventsCh <-chan workflow_registry_wrapper.WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1
	timer    <-chan struct{}
	lggr     logger.Logger
	cr       types.ContractReader
	cfg      ContractEventPollerConfig
	orm      ORM
	gateway  FetcherFunc
	wg       sync.WaitGroup
}

func WithTimer(t <-chan struct{}) func(*worker) {
	return func(w *worker) {
		w.timer = t
	}
}

func New(
	lggr logger.Logger,
	startBlockNum uint64,
	qryCount uint64,
	contractAddr string,
	cr types.ContractReader,
	orm ORM,
	gateway FetcherFunc,
	opts ...func(*worker),
) *worker {
	w := &worker{
		lggr:    lggr,
		cr:      cr,
		orm:     orm,
		gateway: gateway,
		cfg: ContractEventPollerConfig{
			ContractName:      ContractName,
			ContractEventName: ContractEventName,
			ContractAddress:   contractAddr,
			QueryCount:        qryCount,
			StartBlockNum:     startBlockNum,
		},
		stopCh: make(services.StopChan),
	}

	if w.cfg.QueryCount == 0 {
		w.cfg.QueryCount = 20
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

func (w *worker) GetLogs() <-chan workflow_registry_wrapper.WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1 {
	return w.eventsCh
}

func (w *worker) Start(_ context.Context) error {
	return w.StartOnce("secrets_worker", func() error {
		ctx, cancel := w.stopCh.NewCtx()

		done, _ := w.syncForceUpdateSecretsEvents(ctx, w.cfg)

		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			defer cancel()
			<-done
		}()

		return nil
	})
}

func (w *worker) Close() error {
	return w.StopOnce("secrets_worker", func() error {
		close(w.stopCh)
		w.wg.Wait()
		return nil
	})
}

func (w *worker) SecretsFor(ctx context.Context, workflowOwner, workflowName string) (map[string]string, error) {
	return w.orm.SecretsFor(ctx, workflowOwner, workflowName)
}

func (w *worker) getTicker(ctx context.Context) <-chan struct{} {
	if w.timer == nil {
		w.timer = signalers.MakeTicker(ctx.Done(), defaultTickInterval)
	}
	return w.timer
}

// syncForceUpdateSecretsEvents synchronizes the force update secrets events from the contract
// to the local database.
func (w *worker) syncForceUpdateSecretsEvents(
	ctx context.Context,
	cfg ContractEventPollerConfig,
) (<-chan struct{}, <-chan error) {
	// Create the workers
	var (
		h                  = newForceUpdateSecretsHandler(w.orm, w.gateway)
		eventQueryWorker   = newQueryEventsWorker[workflow_registry_wrapper.WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1](w.getTicker(ctx), w.lggr, w.cr)
		eventHandlerWorker = newForceUpdateSecretsWorker(h, w.lggr)
	)

	// Start the workers
	var (
		doneQuerying, queryErrs, updateSecretsEvents = eventQueryWorker.Run(ctx, cfg)
		// Tee the update secrets events to the handler worker
		logs1, logs2              = chans.Tee(ctx.Done(), updateSecretsEvents)
		doneHandling, handlerErrs = eventHandlerWorker.Run(ctx, chans.Transform(ctx.Done(), newEvent, logs1))
	)

	// Set the events channel, which is a read copy of the update secrets events
	w.eventsCh = logs2

	// Merge all the done channels and error channels
	var (
		done  = chans.AllClosed(ctx.Done(), doneQuerying, doneHandling)
		errCh = chans.Merge(ctx.Done(), queryErrs, handlerErrs)
	)

	return done, errCh
}

type event struct {
	workflow_registry_wrapper.WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1
}

func newEvent(
	e workflow_registry_wrapper.WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1,
) URLGetter {
	return event{WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1: e}
}

func (e event) GetSecretsURL() string {
	return e.SecretsURL.Hex()
}
