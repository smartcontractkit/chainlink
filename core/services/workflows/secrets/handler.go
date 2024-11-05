package syncer

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type Event interface {
	GetSecretsURL() string
}

type event struct {
	workflow_registry_wrapper.WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1
}

func newEvent(
	e workflow_registry_wrapper.WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1,
) Event {
	return event{WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1: e}
}

func (e event) GetSecretsURL() string {
	return e.SecretsURL.Hex()
}

type Handler interface {
	ForceUpdateSecrets(context.Context, Event) error
}

type handler struct {
	orm     ORM
	fetcher FetcherFunc
}

func newForceUpdateSecretsHandler(orm ORM, fetcher FetcherFunc) Handler {
	return &handler{orm: orm, fetcher: fetcher}
}

func (h *handler) ForceUpdateSecrets(
	ctx context.Context,
	event Event,
) error {
	// Fetch the secrets url from ORM
	url, err := h.orm.GetSecretsURL(ctx, event.GetSecretsURL())
	if err != nil {
		return err
	}

	// Fetch the contents of the secrets file from the url via the fetcher
	secrets, err := h.fetcher(ctx, url)
	if err != nil {
		return err
	}

	// Update the secrets in the ORM
	if _, err := h.orm.Update(ctx, url, string(secrets)); err != nil {
		return err
	}

	return nil
}

type forceUpdateSecretsWorker struct {
	h    Handler
	lggr logger.Logger
}

func newForceUpdateSecretsWorker(h Handler, lggr logger.Logger) forceUpdateSecretsWorker {
	return forceUpdateSecretsWorker{h: h, lggr: lggr}
}

func (w *forceUpdateSecretsWorker) Run(
	ctx context.Context,
	events <-chan workflow_registry_wrapper.WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1,
) (<-chan struct{}, <-chan error) {
	var (
		done = make(chan struct{})
		errs = make(chan error, 1)

		sendErr = func(err error) {
			select {
			case errs <- err:
			case <-ctx.Done():
			default:
			}
		}
	)

	go func() {
		defer close(done)
		defer close(errs)

		for {
			select {
			case <-ctx.Done():
				return
			case event, open := <-events:
				if !open {
					return
				}
				if err := w.h.ForceUpdateSecrets(ctx, newEvent(event)); err != nil {
					w.lggr.Errorf("error handling update secrets: %v", err)
					sendErr(err)
				}
			}
		}
	}()

	return done, errs
}

func newQueryEventsWorker[T any](
	timer <-chan struct{},
	lggr logger.Logger,
	cr ContractReader,
) queryEventsWorker[T] {
	return queryEventsWorker[T]{timer: timer, lggr: lggr, cr: cr}
}

type ContractReader interface {
	Bind(context.Context, []types.BoundContract) error
	QueryKey(
		context.Context,
		types.BoundContract,
		query.KeyFilter,
		query.LimitAndSort,
		any,
	) ([]types.Sequence, error)
}

type queryEventsWorker[T any] struct {
	timer <-chan struct{}
	lggr  logger.Logger
	cr    ContractReader
}

func (w *queryEventsWorker[T]) Run(
	ctx context.Context,
	cfg ContractEventPollerConfig,
) (
	<-chan struct{},
	<-chan error,
	<-chan T,
) {
	var (
		done   = make(chan struct{})
		errs   = make(chan error, 1)
		events = make(chan T, 1)

		sendErr = func(err error) {
			select {
			case errs <- err:
			case <-ctx.Done():
			default:
			}
		}

		sendLog = func(ld T) {
			select {
			case events <- ld:
			case <-ctx.Done():
			}
		}
	)

	go func() {
		defer close(done)
		defer close(errs)
		defer close(events)

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

		err = w.cr.Bind(ctx, boundContracts)
		if err != nil {
			sendErr(err)
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			case _, open := <-w.timer:
				if !open {
					return
				}

				if cursor != "" {
					limitAndSort.Limit = query.CursorLimit(cursor, query.CursorFollowing, cfg.QueryCount)
				}

				logs, err = w.cr.QueryKey(
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
					w.lggr.Errorw("QueryKey failure", "err", err)
					sendErr(err)
					continue
				}

				// ChainReader QueryKey API provides logs including the cursor value and not
				// after the cursor value. If the response only consists of the log corresponding
				// to the cursor and no log after it, then we understand that there are no new
				// logs
				if len(logs) == 1 && logs[0].Cursor == cursor {
					w.lggr.Infow("No new logs since", "cursor", cursor)
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

					var event T
					err = mapstructure.Decode(dm, &event)
					if err != nil {
						sendErr(err)
						continue
					}

					sendLog(event)
					cursor = log.Cursor
				}
			}
		}
	}()

	return done, errs, events
}
