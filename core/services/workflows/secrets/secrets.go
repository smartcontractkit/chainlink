package syncer

import (
	"context"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils/signalers"
)

var (
	defaultTickInterval = 12 * time.Second
	ContractName        = "WorkflowRegistry"
	ContractEventName   = "WorkflowForceUpdateSecretsRequestedV1"
)

type FetcherFunc func(ctx context.Context, url string) ([]byte, error)

type workerProbes[T any] struct {
	Done      <-chan struct{}
	Err       <-chan error
	Heartbeat <-chan struct{}
	Logs      <-chan T
}

type worker struct {
	services.StateMachine
	stopCh  services.StopChan
	timer   <-chan struct{}
	lggr    logger.Logger
	cr      types.ContractReader
	cfg     ContractEventPollerConfig
	addr    string
	orm     ORM
	gateway FetcherFunc
	wg      sync.WaitGroup
}

func newSecretsWorker(
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

func (w *worker) Start(_ context.Context) error {
	return w.StartOnce("secrets_worker", func() error {
		ctx, cancel := w.stopCh.NewCtx()

		probes := w.fetchSecrets(ctx, w.cfg)

		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			defer cancel()
			<-probes.Done
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

func (w *worker) getTicker(ctx context.Context) <-chan struct{} {
	if w.timer == nil {
		w.timer = signalers.MakeTicker(ctx.Done(), defaultTickInterval)
	}
	return w.timer
}

func (w *worker) fetchSecrets(
	ctx context.Context,
	cfg ContractEventPollerConfig,
) workerProbes[workflow_registry_wrapper.WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1] {
	var (
		done          = make(chan struct{})
		errsCh        = make(chan error)
		wg            sync.WaitGroup
		ctxwc, cancel = context.WithCancel(ctx)
	)

	// Create the workers
	var (
		h                  = newForceUpdateSecretsHandler(w.orm, w.gateway)
		eventQueryWorker   = newQueryEventsWorker[workflow_registry_wrapper.WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1](w.getTicker(ctx), w.lggr, w.cr)
		eventHandlerWorker = newForceUpdateSecretsWorker(h, w.lggr)
	)

	// Start the workers
	var (
		doneQuerying, queryErrs, updateSecretsEvents = eventQueryWorker.Run(ctxwc, cfg)
		doneHandling, handlerErrs                    = eventHandlerWorker.Run(ctxwc, nil)
	)

	// Wait for the workers to finish
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-doneQuerying
		<-doneHandling
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range queryErrs {
			select {
			case <-ctxwc.Done():
				return
			case errsCh <- err:
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range handlerErrs {
			select {
			case <-ctxwc.Done():
				return
			case errsCh <- err:
			}
		}
	}()

	// Close channels when done
	go func() {
		defer close(done)
		defer close(errsCh)
		defer cancel()
		wg.Wait()
	}()

	return workerProbes[workflow_registry_wrapper.WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1]{
		Done: done,
		Err:  errsCh,
		Logs: updateSecretsEvents,
	}
}
