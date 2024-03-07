package cache

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
)

func Test_RMNStateCaching(t *testing.T) {
	ctx := tests.Context(t)
	lggr := logger.TestLogger(t)
	mockCommitStore := mocks.NewCommitStoreReader(t)
	mockOnRamp := mocks.NewOnRampReader(t)

	chainState := newChainHealthcheckWithCustomEviction(
		lggr,
		mockOnRamp,
		mockCommitStore,
		10*time.Hour,
		10*time.Hour,
	)

	// Chain is not cursed and healthy
	mockCommitStore.On("IsDown", ctx).Return(false, nil).Once()
	mockCommitStore.On("IsDestChainHealthy", ctx).Return(true, nil).Maybe()
	mockOnRamp.On("IsSourceCursed", ctx).Return(false, nil).Once()
	mockOnRamp.On("IsSourceChainHealthy", ctx).Return(true, nil).Maybe()
	healthy, err := chainState.IsHealthy(ctx, false)
	assert.NoError(t, err)
	assert.True(t, healthy)

	// Chain is cursed, but cache is stale
	mockCommitStore.On("IsDown", ctx).Return(true, nil).Once()
	mockOnRamp.On("IsSourceCursed", ctx).Return(true, nil).Once()
	healthy, err = chainState.IsHealthy(ctx, false)
	assert.NoError(t, err)
	assert.True(t, healthy)

	// Enforce cache refresh
	healthy, err = chainState.IsHealthy(ctx, true)
	assert.Nil(t, err)
	assert.False(t, healthy)

	// Chain is not cursed, but previous curse should be "sticky" even when force refreshing
	mockCommitStore.On("IsDown", ctx).Return(false, nil).Maybe()
	mockOnRamp.On("IsSourceCursed", ctx).Return(false, nil).Maybe()
	// Enforce cache refresh
	healthy, err = chainState.IsHealthy(ctx, true)
	assert.Nil(t, err)
	assert.False(t, healthy)
}

func Test_ChainStateIsCached(t *testing.T) {
	ctx := tests.Context(t)
	lggr := logger.TestLogger(t)
	mockCommitStore := mocks.NewCommitStoreReader(t)
	mockOnRamp := mocks.NewOnRampReader(t)

	chainState := newChainHealthcheckWithCustomEviction(
		lggr,
		mockOnRamp,
		mockCommitStore,
		10*time.Hour,
		10*time.Hour,
	)

	// Chain is not cursed and healthy
	mockCommitStore.On("IsDown", ctx).Return(false, nil).Maybe()
	mockCommitStore.On("IsDestChainHealthy", ctx).Return(true, nil).Once()
	mockOnRamp.On("IsSourceCursed", ctx).Return(false, nil).Maybe()
	mockOnRamp.On("IsSourceChainHealthy", ctx).Return(true, nil).Once()
	healthy, err := chainState.IsHealthy(ctx, false)
	assert.NoError(t, err)
	assert.True(t, healthy)

	// Chain is not healthy
	mockCommitStore.On("IsDestChainHealthy", ctx).Return(false, nil).Once()
	mockOnRamp.On("IsSourceChainHealthy", ctx).Return(false, nil).Once()
	healthy, err = chainState.IsHealthy(ctx, false)
	assert.NoError(t, err)
	assert.False(t, healthy)

	// Previous value is returned
	mockCommitStore.On("IsDestChainHealthy", ctx).Return(true, nil).Maybe()
	mockOnRamp.On("IsSourceChainHealthy", ctx).Return(true, nil).Maybe()
	healthy, err = chainState.IsHealthy(ctx, false)
	assert.NoError(t, err)
	assert.False(t, healthy)
}

func Test_ChainStateIsHealthy(t *testing.T) {
	ctx := tests.Context(t)

	testCases := []struct {
		name                 string
		commitStoreDown      bool
		commitStoreErr       error
		onRampCursed         bool
		onRampErr            error
		sourceChainUnhealthy bool
		sourceChainErr       error
		destChainUnhealthy   bool
		destChainErr         error

		expectedState bool
		expectedErr   bool
	}{
		{
			name:          "all components healthy",
			expectedState: true,
		},
		{
			name:            "CommitStore is down",
			commitStoreDown: true,
			expectedState:   false,
		},
		{
			name:           "CommitStore error",
			commitStoreErr: errors.New("commit store error"),
			expectedErr:    true,
		},
		{
			name:          "OnRamp is cursed",
			onRampCursed:  true,
			expectedState: false,
		},
		{
			name:        "OnRamp error",
			onRampErr:   errors.New("onramp error"),
			expectedErr: true,
		},
		{
			name:                 "Source chain is unhealthy",
			sourceChainUnhealthy: true,
			expectedState:        false,
		},
		{
			name:           "Source chain error",
			sourceChainErr: errors.New("source chain error"),
			expectedErr:    true,
		},
		{
			name:               "Destination chain is unhealthy",
			destChainUnhealthy: true,
			expectedState:      false,
		},
		{
			name:         "Destination chain error",
			destChainErr: errors.New("destination chain error"),
			expectedErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCommitStore := mocks.NewCommitStoreReader(t)
			mockOnRamp := mocks.NewOnRampReader(t)

			mockCommitStore.On("IsDown", ctx).Return(tc.commitStoreDown, tc.commitStoreErr).Maybe()
			mockCommitStore.On("IsDestChainHealthy", ctx).Return(!tc.destChainUnhealthy, tc.destChainErr).Maybe()
			mockOnRamp.On("IsSourceCursed", ctx).Return(tc.onRampCursed, tc.onRampErr).Maybe()
			mockOnRamp.On("IsSourceChainHealthy", ctx).Return(!tc.sourceChainUnhealthy, tc.sourceChainErr).Maybe()

			chainState := NewChainHealthcheck(logger.TestLogger(t), mockOnRamp, mockCommitStore)
			healthy, err := chainState.IsHealthy(ctx, false)

			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedState, healthy)
			}
		})
	}
}
