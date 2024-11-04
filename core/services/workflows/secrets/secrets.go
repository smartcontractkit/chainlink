package syncer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils/signalers"
)

var (
	defaultTickInterval = 12 * time.Second
	ContractName        = "WorkflowRegistry"
	ContractEventName   = "WorkflowForceUpdateSecretsRequestedV1"
)

type FetcherFunc func(ctx context.Context, url string) ([]byte, error)

type workerProbes struct {
	Done      <-chan struct{}
	Err       <-chan error
	Heartbeat <-chan struct{}
	Logs      <-chan any
}

type worker struct {
	services.StateMachine
	stopCh  services.StopChan
	timer   <-chan struct{}
	lggr    logger.Logger
	cr      types.ContractReader
	cfg     ContractEventPollerConfig
	addr    string
	orm     SecretsUpdater
	gateway FetcherFunc
	wg      sync.WaitGroup
}

func newSecretsWorker(
	lggr logger.Logger,
	startBlockNum uint64,
	qryCount uint64,
	contractAddr string,
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

		probes := fetchSecrets(ctx, w.getTicker(ctx), w.lggr, w.cfg, w.cr, w.orm, w.gateway)

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
	lggr logger.Logger,
	cfg ContractEventPollerConfig,
	cr types.ContractReader,
	_ SecretsUpdater,
	_ FetcherFunc,
) workerProbes {
	var (
		done          = make(chan struct{})
		logsCh        = make(chan any)
		errsCh        = make(chan error)
		hbCh          = make(chan struct{})
		ctxwc, cancel = context.WithCancel(ctx)

		cleanup = func() {
			defer close(done)
			defer close(logsCh)
			defer close(errsCh)
			cancel()
		}

		// pump heart before each unit of work
		beat = func() {
			select {
			case hbCh <- struct{}{}:
			default:
			}
		}

		sendErr = func(err error) {
			select {
			case errsCh <- err:
			case <-ctx.Done():
			default:
			}
		}

		sendLog = func(ld any) {
			select {
			case logsCh <- ld:
			case <-ctx.Done():
			}
		}
	)

	go func() {
		defer cleanup()

		// create query
		var (
			logs    []types.Sequence
			err     error
			logData values.Value

			boundContracts = []types.BoundContract{
				{
					Name:    cfg.ContractName,
					Address: cfg.ContractAddress,
				},
			}
			cursor       = ""
			limitAndSort = query.LimitAndSort{
				SortBy: []query.SortBy{query.NewSortByTimestamp(query.Asc)},
				Limit:  query.Limit{Count: cfg.QueryCount},
			}
		)

		err = cr.Bind(ctx, boundContracts)
		if err != nil {
			sendErr(err)
			return
		}

		for {
			select {
			case <-ctxwc.Done():
				return
			case _, open := <-timer:
				if !open {
					return
				}

				beat()

				if cursor != "" {
					limitAndSort.Limit = query.CursorLimit(cursor, query.CursorFollowing, cfg.QueryCount)
				}

				logs, err = cr.QueryKey(
					ctx,
					types.BoundContract{
						Name:    cfg.ContractName,
						Address: cfg.ContractAddress,
					},
					query.KeyFilter{
						Key: cfg.ContractEventName,
						Expressions: []query.Expression{
							query.Confidence(primitives.Finalized),
							query.Block(fmt.Sprintf("%d", cfg.StartBlockNum), primitives.Gte),
						},
					},
					limitAndSort,
					&logData,
				)

				if err != nil {
					lggr.Errorw("QueryKey failure", "err", err)
					sendErr(err)
					continue
				}

				// ChainReader QueryKey API provides logs including the cursor value and not
				// after the cursor value. If the response only consists of the log corresponding
				// to the cursor and no log after it, then we understand that there are no new
				// logs
				if len(logs) == 1 && logs[0].Cursor == cursor {
					lggr.Infow("No new logs since", "cursor", cursor)
					continue
				}

				for _, log := range logs {
					if log.Cursor == cursor {
						continue
					}

					vm, err := values.WrapMap(log.Data)
					if err != nil {
						sendErr(err)
						continue
					}

					var dm map[string]any
					err = vm.UnwrapTo(&dm)
					if err != nil {
						sendErr(err)
						continue
					}

					lggr.Debugf("Owner: %x, SecretsURL: %x, WorkflowNames: %v", dm["Owner"], dm["SecretsURL"], dm["WorkflowNames"])
					sendLog(dm)
					cursor = log.Cursor
				}
			}
		}
	}()

	return workerProbes{
		Done:      done,
		Logs:      logsCh,
		Err:       errsCh,
		Heartbeat: hbCh,
	}
}
