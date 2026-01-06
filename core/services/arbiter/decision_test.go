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
// It implements services.Service interface as required by the updated ShardConfigReader.
type mockShardConfigReader struct {
	desiredCount uint64
	err          error
}

func (m *mockShardConfigReader) Start(ctx context.Context) error {
	return nil
}

func (m *mockShardConfigReader) Close() error {
	return nil
}

func (m *mockShardConfigReader) Ready() error {
	return nil
}

func (m *mockShardConfigReader) HealthReport() map[string]error {
	return nil
}

func (m *mockShardConfigReader) Name() string {
	return "mockShardConfigReader"
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
			expectedResult: 10, // approved = onChainMax (since we just return on-chain value)
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
			name:           "desired exceeds limit",
			desiredCount:   15,
			onChainMax:     10,
			expectedResult: 10, // approved = onChainMax
			expectError:    false,
		},
		{
			name:           "on-chain limit zero - minimum 1 applied",
			desiredCount:   5,
			onChainMax:     0,
			expectedResult: 1, // minimum of 1 shard
			expectError:    false,
		},
		{
			name:           "small on-chain limit",
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

			engine := NewDecisionEngine(mockReader, logger.Sugared(lggr))

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

	t.Run("large on-chain limit", func(t *testing.T) {
		mockReader := &mockShardConfigReader{
			desiredCount: 1000,
		}
		engine := NewDecisionEngine(mockReader, logger.Sugared(lggr))

		result, err := engine.ComputeApprovedCount(context.Background(), 500)

		require.NoError(t, err)
		assert.Equal(t, 1000, result) // returns on-chain value
	})

	t.Run("exactly at on-chain limit", func(t *testing.T) {
		mockReader := &mockShardConfigReader{
			desiredCount: 7,
		}
		engine := NewDecisionEngine(mockReader, logger.Sugared(lggr))

		result, err := engine.ComputeApprovedCount(context.Background(), 7)

		require.NoError(t, err)
		assert.Equal(t, 7, result)
	})

	t.Run("on-chain limit is zero - minimum 1 applied", func(t *testing.T) {
		mockReader := &mockShardConfigReader{
			desiredCount: 0,
		}
		engine := NewDecisionEngine(mockReader, logger.Sugared(lggr))

		result, err := engine.ComputeApprovedCount(context.Background(), 5)

		require.NoError(t, err)
		// on-chain returns 0, but minimum is 1
		assert.Equal(t, 1, result)
	})
}

func TestDecisionEngine_ContextCancellation(t *testing.T) {
	lggr := logger.TestLogger(t)

	t.Run("context cancellation propagated", func(t *testing.T) {
		mockReader := &mockShardConfigReader{
			err: context.Canceled,
		}
		engine := NewDecisionEngine(mockReader, logger.Sugared(lggr))

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := engine.ComputeApprovedCount(ctx, 5)

		require.Error(t, err)
	})
}
