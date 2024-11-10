package syncer

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/triggers/logevent/logeventcap"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/capabilities/workflows/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockEvent struct {
	SecretsURLHash string
}

func (e mockEvent) GetURLHash() string {
	return e.SecretsURLHash
}

func Test_Handler(t *testing.T) {
	lggr := logger.TestLogger(t)
	t.Run("success", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		giveURL := "https://original-url.com"
		hash := common.Keccak256Hash([]byte(giveURL))

		giveEvent := WorkflowRegistryEvent{
			Output: logeventcap.Output{
				Data: map[string]any{
					"SecretsURL": []byte(hash),
				},
			},
		}

		mockORM.EXPECT().GetSecretsURL(matches.AnyContext, hash).Return(giveURL, nil)
		fetcher := func(_ context.Context, _ string) ([]byte, error) {
			return []byte("contents"), nil
		}
		mockORM.EXPECT().Update(matches.AnyContext, giveURL, "contents").Return(int64(1), nil)
		h := newForceUpdateSecretsHandler(lggr, mockORM, fetcher)
		err := h.Handle(ctx, giveEvent)
		require.NoError(t, err)
	})

	t.Run("fails to get secrets", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		giveURL := "http://example.com"
		hash := common.Keccak256Hash([]byte(giveURL))
		mockORM.EXPECT().GetSecretsURL(matches.AnyContext, hash).Return("", assert.AnError)
		h := newForceUpdateSecretsHandler(lggr, mockORM, nil)
		err := h.Handle(ctx, WorkflowRegistryEvent{
			Output: logeventcap.Output{
				Data: map[string]any{
					"SecretsURL": []byte(hash),
				},
			},
		})
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("fails to fetch contents", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		giveURL := "http://example.com"
		hash := common.Keccak256Hash([]byte(giveURL))
		mockORM.EXPECT().GetSecretsURL(matches.AnyContext, hash).Return("", assert.AnError)
		fetcher := func(_ context.Context, _ string) ([]byte, error) {
			return nil, assert.AnError
		}

		h := newForceUpdateSecretsHandler(lggr, mockORM, fetcher)
		err := h.Handle(ctx, WorkflowRegistryEvent{
			Output: logeventcap.Output{
				Data: map[string]any{
					"SecretsURL": []byte(hash),
				},
			},
		})
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("fails to update secrets", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		ctx := testutils.Context(t)
		giveURL := "http://example.com"
		hash := common.Keccak256Hash([]byte(giveURL))
		mockORM.EXPECT().GetSecretsURL(matches.AnyContext, hash).Return(giveURL, nil)
		fetcher := func(_ context.Context, _ string) ([]byte, error) {
			return []byte("contents"), nil
		}
		mockORM.EXPECT().Update(matches.AnyContext, giveURL, "contents").Return(0, assert.AnError)
		h := newForceUpdateSecretsHandler(lggr, mockORM, fetcher)
		err := h.Handle(ctx, WorkflowRegistryEvent{
			Output: logeventcap.Output{
				Data: map[string]any{
					"SecretsURL": []byte(hash),
				},
			},
		})
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})
}

func Test_getURLHash(t *testing.T) {
	giveURL := "http://example.com"
	hash := common.Keccak256Hash([]byte(giveURL))
	giveEvent := WorkflowRegistryEvent{
		Output: logeventcap.Output{
			Data: map[string]any{
				"SecretsURL": []byte(hash),
			},
		},
	}
	gotHash, err := getURLHash(giveEvent)
	require.NoError(t, err)

	assert.Equal(t, hash, gotHash)
}
