package secrets

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	types "github.com/smartcontractkit/chainlink-common/pkg/types"
	query "github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/secrets/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockEvent struct {
	SecretsURL string
}

func (e mockEvent) GetSecretsURL() string {
	return e.SecretsURL
}

func Test_Handler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		giveURL := "http://example.com"
		hash := fmt.Sprintf("%x", keccak256Hash(giveURL))
		mockORM.EXPECT().GetSecretsURL(matches.AnyContext, hash).Return(giveURL, nil)
		fetcher := func(_ context.Context, _ string) ([]byte, error) {
			return []byte("contents"), nil
		}
		mockORM.EXPECT().Update(matches.AnyContext, giveURL, "contents").Return(int64(1), nil)
		h := newForceUpdateSecretsHandler(mockORM, fetcher)
		err := h.ForceUpdateSecrets(ctx, mockEvent{SecretsURL: hash})
		require.NoError(t, err)
	})

	t.Run("fails to get secrets", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		giveURL := "http://example.com"
		hash := fmt.Sprintf("%x", keccak256Hash(giveURL))
		mockORM.EXPECT().GetSecretsURL(matches.AnyContext, hash).Return("", assert.AnError)
		h := newForceUpdateSecretsHandler(mockORM, nil)
		err := h.ForceUpdateSecrets(ctx, mockEvent{SecretsURL: hash})
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("fails to fetch contents", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		giveURL := "http://example.com"
		hash := fmt.Sprintf("%x", keccak256Hash(giveURL))
		mockORM.EXPECT().GetSecretsURL(matches.AnyContext, hash).Return("", assert.AnError)
		fetcher := func(_ context.Context, _ string) ([]byte, error) {
			return nil, assert.AnError
		}

		h := newForceUpdateSecretsHandler(mockORM, fetcher)
		err := h.ForceUpdateSecrets(ctx, mockEvent{SecretsURL: hash})
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("fails to update secrets", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		giveURL := "http://example.com"
		hash := fmt.Sprintf("%x", keccak256Hash(giveURL))
		mockORM.EXPECT().GetSecretsURL(matches.AnyContext, hash).Return(giveURL, nil)
		fetcher := func(_ context.Context, _ string) ([]byte, error) {
			return []byte("contents"), nil
		}
		mockORM.EXPECT().Update(matches.AnyContext, giveURL, "contents").Return(0, assert.AnError)
		h := newForceUpdateSecretsHandler(mockORM, fetcher)
		err := h.ForceUpdateSecrets(ctx, mockEvent{SecretsURL: hash})
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})
}

func Test_HandlerWorker(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var (
			lggr        = logger.TestLogger(t)
			mockHandler = NewMockHandler[URLGetter](t)
			ctx, cancel = context.WithCancel(testutils.Context(t))
			worker      = newForceUpdateSecretsWorker(mockHandler, lggr)
			events      = make(chan URLGetter)
			giveEvent   = newEvent(workflow_registry_wrapper.WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1{})
		)

		mockHandler.EXPECT().ForceUpdateSecrets(matches.AnyContext, giveEvent).Return(nil)

		done, _ := worker.Run(ctx, events)

		events <- giveEvent

		cancel()
		<-done
	})

	t.Run("failure to handle event", func(t *testing.T) {
		var (
			lggr        = logger.TestLogger(t)
			mockHandler = NewMockHandler[URLGetter](t)
			ctx, cancel = context.WithCancel(testutils.Context(t))
			worker      = newForceUpdateSecretsWorker(mockHandler, lggr)
			events      = make(chan URLGetter)
			giveEvent   = newEvent(workflow_registry_wrapper.WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1{})
		)

		mockHandler.EXPECT().ForceUpdateSecrets(matches.AnyContext, giveEvent).Return(assert.AnError)

		done, errsCh := worker.Run(ctx, events)

		events <- giveEvent

		err := <-errsCh
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)

		cancel()
		<-done
	})
}

func Test_QueryEventsHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var (
			lggr        = logger.TestLogger(t)
			ctx, cancel = context.WithCancel(testutils.Context(t))
			reader      = NewMockContractReader(t)
			timer       = make(chan struct{})
			giveCfg     = ContractEventPollerConfig{
				ContractName:      "MockContract",
				ContractAddress:   "0xdeadbeef",
				ContractEventName: "MockEvent",
				StartBlockNum:     0,
				QueryCount:        1,
			}
			giveLog = types.Sequence{
				Data:   mockEvent{SecretsURL: "http://example.com"},
				Cursor: "cursor",
			}
		)

		worker := newQueryEventsWorker[mockEvent](timer, lggr, reader)

		reader.EXPECT().Bind(matches.AnyContext, []types.BoundContract{
			{
				Name:    giveCfg.ContractName,
				Address: giveCfg.ContractAddress,
			},
		}).Return(nil)

		reader.EXPECT().QueryKey(
			matches.AnyContext,
			types.BoundContract{
				Name:    giveCfg.ContractName,
				Address: giveCfg.ContractAddress,
			},
			query.KeyFilter{
				Key: giveCfg.ContractEventName,
				Expressions: []query.Expression{
					query.Confidence(primitives.Finalized),
					query.Block(strconv.FormatUint(giveCfg.StartBlockNum, 10), primitives.Gte),
				},
			},
			query.LimitAndSort{
				SortBy: []query.SortBy{query.NewSortByTimestamp(query.Asc)},
				Limit:  query.Limit{Count: giveCfg.QueryCount},
			},
			new(values.Value),
		).Return([]types.Sequence{
			giveLog,
		}, nil)

		done, _, logs := worker.Run(ctx, giveCfg)

		timer <- struct{}{}
		gotLog := <-logs
		cancel()
		<-done

		require.Equal(t, "http://example.com", gotLog.GetSecretsURL())
	})
}
