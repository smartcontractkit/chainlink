package syncer

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Handler(t *testing.T) {
	lggr := logger.TestLogger(t)
	t.Run("success", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		giveURL := "https://original-url.com"
		giveEvent := WorkflowRegistryEvent{
			EventType: ForceUpdateSecretsEvent,
			Data: map[string]any{
				"SecretsURL": giveURL,
			},
		}

		fetcher := func(_ context.Context, _ string) ([]byte, error) {
			return []byte("contents"), nil
		}
		mockORM.EXPECT().Update(matches.AnyContext, giveURL, "contents").Return(int64(1), nil)
		h := newEventHandler(lggr, mockORM, fetcher)
		err := h.Handle(ctx, giveEvent)
		require.NoError(t, err)
	})

	t.Run("fails with unsupported event type", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)

		giveEvent := WorkflowRegistryEvent{}
		fetcher := func(_ context.Context, _ string) ([]byte, error) {
			return []byte("contents"), nil
		}

		h := newEventHandler(lggr, mockORM, fetcher)
		err := h.Handle(ctx, giveEvent)
		require.Error(t, err)
		require.Contains(t, err.Error(), "event type unsupported")
	})

	t.Run("fails to get secrets url", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		h := newEventHandler(lggr, mockORM, nil)
		err := h.Handle(ctx, WorkflowRegistryEvent{
			EventType: ForceUpdateSecretsEvent,
			Data: map[string]any{
				"SecretsURL": assert.AnError,
			},
		})
		require.Error(t, err)
		require.ErrorContains(t, err, assert.AnError.Error())
	})

	t.Run("fails to fetch contents", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		giveURL := "http://example.com"

		fetcher := func(_ context.Context, _ string) ([]byte, error) {
			return nil, assert.AnError
		}

		h := newEventHandler(lggr, mockORM, fetcher)
		err := h.Handle(ctx, WorkflowRegistryEvent{
			EventType: ForceUpdateSecretsEvent,
			Data: map[string]any{
				"SecretsURL": giveURL,
			},
		})
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("fails to update secrets", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		giveURL := "http://example.com"

		fetcher := func(_ context.Context, _ string) ([]byte, error) {
			return []byte("contents"), nil
		}
		mockORM.EXPECT().Update(matches.AnyContext, giveURL, "contents").Return(0, assert.AnError)
		h := newEventHandler(lggr, mockORM, fetcher)
		err := h.Handle(ctx, WorkflowRegistryEvent{
			EventType: ForceUpdateSecretsEvent,
			Data: map[string]any{
				"SecretsURL": giveURL,
			},
		})
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})
}

func Test_getURLHash(t *testing.T) {
	giveURL := "http://example.com"

	giveEvent := WorkflowRegistryEvent{
		Data: map[string]any{
			"SecretsURL": giveURL,
		},
	}
	gotURL, err := getSecretsURL(giveEvent)
	require.NoError(t, err)

	assert.Equal(t, giveURL, gotURL)
}
