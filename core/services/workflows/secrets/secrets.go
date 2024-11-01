package syncer

import (
	"context"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils/signalers"
)

var defaultTickInterval = 12 * time.Second

type UpdateSecretsCommand struct {
	SecretsURL    string
	Contents      []byte
	WorkflowName  string
	WorkflowOwner string
}

type Updater[T any] interface {
	Update(ctx context.Context, cmd T) (int64, error)
}

type SecretsUpdater = Updater[UpdateSecretsCommand]

type FetcherFunc func(ctx context.Context, url string) ([]byte, error)

type workerProbes struct {
	Done     <-chan struct{}
	Err      <-chan error
	Hearbeat <-chan struct{}
	Logs     <-chan any
}

func newWorkerProbes() workerProbes {
	return workerProbes{
		Done:     make(chan struct{}),
		Err:      make(chan error),
		Hearbeat: make(chan struct{}),
		Logs:     make(chan any),
	}
}

type worker struct {
	services.StateMachine
	stopCh  services.StopChan
	timer   <-chan struct{}
	lggr    logger.Logger
	cr      types.ContractReader
	orm     SecretsUpdater
	gateway FetcherFunc
	wg      sync.WaitGroup
}

func newSecretsWorker(
	lggr logger.Logger,
	cr types.ContractReader,
	orm SecretsUpdater,
	gateway FetcherFunc,
	opts ...func(*worker),
) *worker {
	w := &worker{
		lggr:    lggr,
		cr:      cr,
		orm:     orm,
		gateway: gateway,
		stopCh:  make(services.StopChan),
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func (w *worker) Start(_ context.Context) error {
	return w.StartOnce("secrets_worker", func() error {
		ctx, cancel := w.stopCh.NewCtx()
		probes := fetchSecrets(ctx, w.getTicker(ctx), w.cr, w.orm, w.gateway)

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

func fetchSecrets(
	ctx context.Context,
	timer <-chan struct{},
	cr types.ContractReader,
	orm SecretsUpdater,
	gateway FetcherFunc,
) workerProbes {
	probes := newWorkerProbes()

	go func() {

	}()

	return probes
}
