package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
)

func Test_RMNStateCaching(t *testing.T) {
	ctx := tests.Context(t)
	lggr := logger.TestLogger(t)
	mockCommitStore := mocks.NewCommitStoreReader(t)
	mockOnRamp := mocks.NewOnRampReader(t)

	chainState := newChainHealthcheckWithCustomEviction(lggr, mockOnRamp, mockCommitStore, 10*time.Hour, 10*time.Hour)

	// Chain is not cursed and healthy
	mockCommitStore.On("IsDown", ctx).Return(false, nil).Once()
	mockCommitStore.On("IsDestChainHealthy", ctx).Return(true, nil).Maybe()
	mockOnRamp.On("IsSourceCursed", ctx).Return(false, nil).Once()
	mockOnRamp.On("IsSourceChainHealthy", ctx).Return(true, nil).Maybe()
	healthy, err := chainState.IsHealthy(ctx)
	assert.NoError(t, err)
	assert.True(t, healthy)

	// Chain is cursed, but cache is stale
	mockCommitStore.On("IsDown", ctx).Return(true, nil).Once()
	mockOnRamp.On("IsSourceCursed", ctx).Return(true, nil).Once()
	healthy, err = chainState.IsHealthy(ctx)
	assert.NoError(t, err)
	assert.True(t, healthy)

	// Enforce cache refresh
	_, err = chainState.refresh(ctx)
	assert.NoError(t, err)

	healthy, err = chainState.IsHealthy(ctx)
	assert.Nil(t, err)
	assert.False(t, healthy)

	// Chain is not cursed, but previous curse should be "sticky" even when force refreshing
	mockCommitStore.On("IsDown", ctx).Return(false, nil).Maybe()
	mockOnRamp.On("IsSourceCursed", ctx).Return(false, nil).Maybe()
	// Enforce cache refresh
	_, err = chainState.refresh(ctx)
	assert.NoError(t, err)

	healthy, err = chainState.IsHealthy(ctx)
	assert.Nil(t, err)
	assert.False(t, healthy)
}

func Test_ChainStateIsCached(t *testing.T) {
	ctx := tests.Context(t)
	lggr := logger.TestLogger(t)
	mockCommitStore := mocks.NewCommitStoreReader(t)
	mockOnRamp := mocks.NewOnRampReader(t)

	chainState := newChainHealthcheckWithCustomEviction(lggr, mockOnRamp, mockCommitStore, 10*time.Hour, 10*time.Hour)

	// Chain is not cursed and healthy
	mockCommitStore.On("IsDown", ctx).Return(false, nil).Maybe()
	mockCommitStore.On("IsDestChainHealthy", ctx).Return(true, nil).Once()
	mockOnRamp.On("IsSourceCursed", ctx).Return(false, nil).Maybe()
	mockOnRamp.On("IsSourceChainHealthy", ctx).Return(true, nil).Once()

	_, err := chainState.refresh(ctx)
	assert.NoError(t, err)

	healthy, err := chainState.IsHealthy(ctx)
	assert.NoError(t, err)
	assert.True(t, healthy)

	// Chain is not healthy
	mockCommitStore.On("IsDestChainHealthy", ctx).Return(false, nil).Once()
	mockOnRamp.On("IsSourceChainHealthy", ctx).Return(false, nil).Once()
	_, err = chainState.refresh(ctx)
	assert.NoError(t, err)

	healthy, err = chainState.IsHealthy(ctx)
	assert.NoError(t, err)
	assert.False(t, healthy)

	// Previous value is returned
	mockCommitStore.On("IsDestChainHealthy", ctx).Return(true, nil).Maybe()
	mockOnRamp.On("IsSourceChainHealthy", ctx).Return(true, nil).Maybe()

	_, err = chainState.refresh(ctx)
	assert.NoError(t, err)

	healthy, err = chainState.IsHealthy(ctx)
	assert.NoError(t, err)
	assert.False(t, healthy)
}

func Test_ChainStateIsHealthy(t *testing.T) {
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
			ctx := tests.Context(t)
			mockCommitStore := mocks.NewCommitStoreReader(t)
			mockOnRamp := mocks.NewOnRampReader(t)

			mockCommitStore.On("IsDown", ctx).Return(tc.commitStoreDown, tc.commitStoreErr).Maybe()
			mockCommitStore.On("IsDestChainHealthy", ctx).Return(!tc.destChainUnhealthy, tc.destChainErr).Maybe()
			mockOnRamp.On("IsSourceCursed", ctx).Return(tc.onRampCursed, tc.onRampErr).Maybe()
			mockOnRamp.On("IsSourceChainHealthy", ctx).Return(!tc.sourceChainUnhealthy, tc.sourceChainErr).Maybe()

			chainState := newChainHealthcheckWithCustomEviction(logger.TestLogger(t), mockOnRamp, mockCommitStore, 10*time.Hour, 10*time.Hour)

			healthy, err := chainState.IsHealthy(ctx)

			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedState, healthy)
			}
		})
	}
}

func Test_RefreshingInBackground(t *testing.T) {
	mockCommitStore := newCommitStoreWrapper(t, true, nil)
	mockCommitStore.CommitStoreReader.On("IsDestChainHealthy", mock.Anything).Return(true, nil).Maybe()

	mockOnRamp := newOnRampWrapper(t, true, nil)
	mockOnRamp.OnRampReader.On("IsSourceChainHealthy", mock.Anything).Return(true, nil).Maybe()

	chainState := newChainHealthcheckWithCustomEviction(
		logger.TestLogger(t),
		mockOnRamp,
		mockCommitStore,
		10*time.Microsecond,
		10*time.Microsecond,
	)
	require.NoError(t, chainState.Start(tests.Context(t)))

	// All healthy
	assertHealthy(t, chainState, true)

	// Commit store not healthy
	mockCommitStore.set(false, nil)
	assertHealthy(t, chainState, false)

	// Commit store error
	mockCommitStore.set(false, fmt.Errorf("commit store error"))
	assertError(t, chainState)

	// Commit store is back
	mockCommitStore.set(true, nil)
	assertHealthy(t, chainState, true)

	// OnRamp not healthy
	mockOnRamp.set(false, nil)
	assertHealthy(t, chainState, false)

	// OnRamp error
	mockOnRamp.set(false, fmt.Errorf("onramp error"))
	assertError(t, chainState)

	// All back in healthy state
	mockOnRamp.set(true, nil)
	assertHealthy(t, chainState, true)

	require.NoError(t, chainState.Close())
}

func assertHealthy(t *testing.T, ch *chainHealthcheck, expected bool) {
	assert.Eventually(t, func() bool {
		healthy, err := ch.IsHealthy(testutils.Context(t))
		return err == nil && healthy == expected
	}, testutils.WaitTimeout(t), testutils.TestInterval)
}

func assertError(t *testing.T, ch *chainHealthcheck) {
	assert.Eventually(t, func() bool {
		_, err := ch.IsHealthy(testutils.Context(t))
		return err != nil
	}, testutils.WaitTimeout(t), testutils.TestInterval)
}

type fakeStatusWrapper struct {
	*mocks.CommitStoreReader
	*mocks.OnRampReader

	healthy bool
	err     error
	mu      *sync.Mutex
}

func newCommitStoreWrapper(t *testing.T, healthy bool, err error) *fakeStatusWrapper {
	return &fakeStatusWrapper{
		CommitStoreReader: mocks.NewCommitStoreReader(t),
		healthy:           healthy,
		err:               err,
		mu:                new(sync.Mutex),
	}
}

func newOnRampWrapper(t *testing.T, healthy bool, err error) *fakeStatusWrapper {
	return &fakeStatusWrapper{
		OnRampReader: mocks.NewOnRampReader(t),
		healthy:      healthy,
		err:          err,
		mu:           new(sync.Mutex),
	}
}

func (f *fakeStatusWrapper) IsDown(context.Context) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return !f.healthy, f.err
}

func (f *fakeStatusWrapper) IsSourceCursed(context.Context) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return !f.healthy, f.err
}

func (f *fakeStatusWrapper) Close() error {
	return nil
}

func (f *fakeStatusWrapper) set(healthy bool, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.healthy = healthy
	f.err = err
}
