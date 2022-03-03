package evm_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
)

func Test_OCRContractTracker_LatestBlockHeight(t *testing.T) {
	t.Parallel()

	t.Run("on L2 chains, always returns 0", func(t *testing.T) {
		uni := newContractTrackerUni(t, evmtest.ChainOptimismMainnet(t))
		l, err := uni.configTracker.LatestBlockHeight(context.Background())
		require.NoError(t, err)

		assert.Equal(t, uint64(0), l)
	})

	t.Run("before first head incoming, looks up on-chain", func(t *testing.T) {
		uni := newContractTrackerUni(t)
		uni.ec.On("HeadByNumber", mock.AnythingOfType("*context.cancelCtx"), (*big.Int)(nil)).Return(&evmtypes.Head{Number: 42}, nil)

		l, err := uni.configTracker.LatestBlockHeight(context.Background())
		require.NoError(t, err)

		assert.Equal(t, uint64(42), l)
	})

	t.Run("Before first head incoming, on client error returns error", func(t *testing.T) {
		uni := newContractTrackerUni(t)
		uni.ec.On("HeadByNumber", mock.AnythingOfType("*context.cancelCtx"), (*big.Int)(nil)).Return(nil, nil).Once()

		_, err := uni.configTracker.LatestBlockHeight(context.Background())
		assert.EqualError(t, err, "got nil head")

		uni.ec.On("HeadByNumber", mock.AnythingOfType("*context.cancelCtx"), (*big.Int)(nil)).Return(nil, errors.New("bar")).Once()

		_, err = uni.configTracker.LatestBlockHeight(context.Background())
		assert.EqualError(t, err, "bar")

		uni.ec.AssertExpectations(t)
	})

	t.Run("after first head incoming, uses cached value", func(t *testing.T) {
		uni := newContractTrackerUni(t)

		uni.configTracker.OnNewLongestChain(context.Background(), &evmtypes.Head{Number: 42})

		l, err := uni.configTracker.LatestBlockHeight(context.Background())
		require.NoError(t, err)

		assert.Equal(t, uint64(42), l)
	})

	t.Run("if headbroadcaster has it, uses the given value on start", func(t *testing.T) {
		uni := newContractTrackerUni(t)

		uni.hb.On("Subscribe", uni.configTracker).Return(&evmtypes.Head{Number: 42}, func() {})
		require.NoError(t, uni.configTracker.Start())

		l, err := uni.configTracker.LatestBlockHeight(context.Background())
		require.NoError(t, err)

		assert.Equal(t, uint64(42), l)

		uni.hb.AssertExpectations(t)

		require.NoError(t, uni.configTracker.Close())
	})
}
