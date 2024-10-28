package syncer

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const name = "WorkflowRegistrySyncer"

type WorkflowRegistrySyncer interface {
	services.Service
	Sync(ctx context.Context, isInitialSync bool) error
}

var _ WorkflowRegistrySyncer = (*workflowRegistry)(nil)

type workflowRegistry struct {
	services.StateMachine
	stopCh   services.StopChan
	updateCh chan any
	errCh    chan error
	tickerCh <-chan struct{}
	lggr     logger.Logger
	orm      WorkflowRegistryDS
	wg       sync.WaitGroup
	mu       sync.Mutex
}

func (w *workflowRegistry) Start(ctx context.Context) error {
	return w.StartOnce(w.Name(), func() error {
		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			w.syncLoop()
		}()

		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			w.updateStateLoop()
		}()

		return nil
	})
}

func (w *workflowRegistry) Sync(ctx context.Context, isInitialSync bool) error {
	return errors.New("not implemented")
}

func (w *workflowRegistry) Close() error {
	return w.StopOnce(w.Name(), func() error {
		close(w.stopCh)
		w.wg.Wait()
		return nil
	})
}

func (w *workflowRegistry) Ready() error {
	return nil
}

func (w *workflowRegistry) HealthReport() map[string]error {
	return nil
}

func (w *workflowRegistry) Name() string {
	return name
}

func (w *workflowRegistry) SecretsFor(workflowOwner, workflowName string) (map[string]string, error) {
	// TODO: actually get this from the right place.
	return map[string]string{}, nil
}

var (
	defaultTickInterval = 12 * time.Second
)

func NewWorkflowRegistry(
	lggr logger.Logger,
	orm WorkflowRegistryDS,
) *workflowRegistry {
	return &workflowRegistry{
		stopCh:   make(services.StopChan),
		updateCh: make(chan any),
		errCh:    make(chan error, 1),
		lggr:     lggr.Named(name),
		orm:      orm,
	}
}

func (w *workflowRegistry) syncLoop() {
	ctx, cancel := w.stopCh.NewCtx()
	defer cancel()

	ticker := w.getTicker(ctx)

	// Sync for a first time outside the loop; this means we'll check the cached state
	// or start a remote sync immediately once spinning up syncLoop, as by default a ticker will
	// fire for the first time at T+N, where N is the interval.
	//
	// TODO(mstreet3): metrics
	w.lggr.Debug("starting initial sync with remote registry")
	err := w.Sync(ctx, true)
	if err != nil {
		w.lggr.Errorw("failed to sync with remote registry", "error", err)
		w.sinkError(err)
	}

	for {
		select {
		case <-w.stopCh:
			return
		case <-ticker:
			w.lggr.Debug("starting regular sync with the remote registry")
			err := w.Sync(ctx, false)
			if err != nil {
				w.lggr.Errorw("failed to sync with remote registry", "error", err)
				w.sinkError(err)
			}
		}
	}
}

func (w *workflowRegistry) updateStateLoop() {
	ctx, cancel := w.stopCh.NewCtx()
	defer cancel()

	for {
		select {
		case <-w.stopCh:
			return
		case localRegistry, ok := <-w.updateCh:
			if !ok {
				// channel has been closed, terminating.
				return
			}
			if err := w.orm.AddLocalRegistry(ctx, localRegistry); err != nil {
				w.lggr.Errorw("failed to save state to local registry", "error", err)
				w.sinkError(err)
			}
		}
	}
}

// sinkError is a non-blocking attempt to sink errors to the registry syncer's error channel.
// If no one reads, the errors are dropped.
func (w *workflowRegistry) sinkError(err error) {
	select {
	case w.errCh <- err:
	default:
	}
	return
}

func (w *workflowRegistry) getTicker(ctx context.Context) <-chan struct{} {
	if w.tickerCh != nil {
		return w.tickerCh
	}

	w.tickerCh = makeTicker(ctx.Done(), defaultTickInterval)
	return w.tickerCh
}

func makeTicker(stop <-chan struct{}, d time.Duration) <-chan struct{} {
	ticker := make(chan struct{})
	internalTicker := time.NewTicker(d)

	go func() {
		defer close(ticker)
		defer internalTicker.Stop()

		for {
			select {
			case <-stop:
				return
			case <-internalTicker.C:
			}
		}
	}()

	return ticker
}
