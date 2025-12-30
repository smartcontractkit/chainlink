package arbiter

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// mockShardConfigReader is a mock implementation of ShardConfigReader for testing.
type mockShardConfigReader struct {
	desiredCount uint64
	err          error
}

func (m *mockShardConfigReader) GetDesiredShardCount(ctx context.Context) (uint64, error) {
	return m.desiredCount, m.err
}

func TestDecisionEngine_ComputeApprovedCount(t *testing.T) {
	tests := []struct {
		name           string
		desiredCount   int
		onChainMax     uint64
		shardConfigErr error
		expectedResult int
		expectError    bool
	}{
		{
			name:           "desired under limit",
			desiredCount:   5,
			onChainMax:     10,
			expectedResult: 5,
			expectError:    false,
		},
		{
			name:           "desired equals limit",
			desiredCount:   10,
			onChainMax:     10,
			expectedResult: 10,
			expectError:    false,
		},
		{
			name:           "desired exceeds limit - capped",
			desiredCount:   15,
			onChainMax:     10,
			expectedResult: 10,
			expectError:    false,
		},
		{
			name:           "desired zero - minimum 1",
			desiredCount:   0,
			onChainMax:     10,
			expectedResult: 1,
			expectError:    false,
		},
		{
			name:           "negative desired - minimum 1",
			desiredCount:   -5,
			onChainMax:     10,
			expectedResult: 1,
			expectError:    false,
		},
		{
			name:           "small on-chain limit caps result",
			desiredCount:   5,
			onChainMax:     3,
			expectedResult: 3,
			expectError:    false,
		},
		{
			name:           "on-chain limit of 1",
			desiredCount:   100,
			onChainMax:     1,
			expectedResult: 1,
			expectError:    false,
		},
		{
			name:           "shard config error",
			desiredCount:   5,
			onChainMax:     10,
			shardConfigErr: errors.New("contract read failed"),
			expectedResult: 0,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lggr := logger.TestLogger(t)

			mockReader := &mockShardConfigReader{
				desiredCount: tc.onChainMax,
				err:          tc.shardConfigErr,
			}

			engine := NewDecisionEngine(mockReader, lggr)

			result, err := engine.ComputeApprovedCount(context.Background(), tc.desiredCount)

			if tc.expectError {
				require.Error(t, err)
				assert.Equal(t, 0, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestDecisionEngine_ComputeApprovedCount_EdgeCases(t *testing.T) {
	lggr := logger.TestLogger(t)

	t.Run("large desired count capped to large on-chain limit", func(t *testing.T) {
		mockReader := &mockShardConfigReader{
			desiredCount: 1000,
		}
		engine := NewDecisionEngine(mockReader, lggr)

		result, err := engine.ComputeApprovedCount(context.Background(), 500)

		require.NoError(t, err)
		assert.Equal(t, 500, result)
	})

	t.Run("exactly at on-chain limit", func(t *testing.T) {
		mockReader := &mockShardConfigReader{
			desiredCount: 7,
		}
		engine := NewDecisionEngine(mockReader, lggr)

		result, err := engine.ComputeApprovedCount(context.Background(), 7)

		require.NoError(t, err)
		assert.Equal(t, 7, result)
	})

	t.Run("on-chain limit is zero - minimum 1 applied", func(t *testing.T) {
		mockReader := &mockShardConfigReader{
			desiredCount: 0,
		}
		engine := NewDecisionEngine(mockReader, lggr)

		result, err := engine.ComputeApprovedCount(context.Background(), 5)

		require.NoError(t, err)
		// approved = min(5, 0) = 0, but minimum is 1
		assert.Equal(t, 1, result)
	})
}

func TestDecisionEngine_ContextCancellation(t *testing.T) {
	lggr := logger.TestLogger(t)

	t.Run("context cancellation propagated", func(t *testing.T) {
		mockReader := &mockShardConfigReader{
			err: context.Canceled,
		}
		engine := NewDecisionEngine(mockReader, lggr)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := engine.ComputeApprovedCount(ctx, 5)

		require.Error(t, err)
	})
}
