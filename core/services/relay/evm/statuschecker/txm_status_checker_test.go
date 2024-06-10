package statuschecker

import (
	"context"
	"errors"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTxManager is a mock implementation of TxManager
type MockTxManager struct {
	mock.Mock
}

func (m *MockTxManager) GetTransactionStatus(ctx context.Context, transactionID string) (TransactionStatus, error) {
	args := m.Called(ctx, transactionID)
	return args.Get(0).(TransactionStatus), args.Error(1)
}

func Test_CheckMessageStatus(t *testing.T) {
	testutils.SkipShort(t, "")
	ctx := context.Background()
	mockTxManager := new(MockTxManager)
	checker := NewTransactionStatusChecker(mockTxManager)

	msgID := "test-message-id"

	// Define test cases
	testCases := []struct {
		name            string
		setupMock       func()
		expectedStatus  []TransactionStatus
		expectedCounter int
		expectedError   error
	}{
		{
			name: "No transactions found",
			setupMock: func() {
				mockTxManager.Mock = mock.Mock{}
				mockTxManager.On("GetTransactionStatus", ctx, "test-message-id-0").Return(Unknown, errors.New("failed to find transaction with IdempotencyKey test-message-id-0"))
			},
			expectedStatus:  nil,
			expectedCounter: -1,
			expectedError:   nil,
		},
		{
			name: "Single transaction found",
			setupMock: func() {
				mockTxManager.Mock = mock.Mock{}
				mockTxManager.On("GetTransactionStatus", ctx, "test-message-id-0").Return(Finalized, nil)
				mockTxManager.On("GetTransactionStatus", ctx, "test-message-id-1").Return(Unknown, errors.New("failed to find transaction with IdempotencyKey test-message-id-1"))
			},
			expectedStatus:  []TransactionStatus{Finalized},
			expectedCounter: 0,
			expectedError:   nil,
		},
		{
			name: "Multiple transactions found",
			setupMock: func() {
				mockTxManager.Mock = mock.Mock{}
				mockTxManager.On("GetTransactionStatus", ctx, "test-message-id-0").Return(Finalized, nil)
				mockTxManager.On("GetTransactionStatus", ctx, "test-message-id-1").Return(Failed, nil)
				mockTxManager.On("GetTransactionStatus", ctx, "test-message-id-2").Return(Unknown, errors.New("failed to find transaction with IdempotencyKey test-message-id-2"))
			},
			expectedStatus:  []TransactionStatus{Finalized, Failed},
			expectedCounter: 1,
			expectedError:   nil,
		},
		{
			name: "Error during transaction retrieval",
			setupMock: func() {
				mockTxManager.Mock = mock.Mock{}
				mockTxManager.On("GetTransactionStatus", ctx, "test-message-id-0").Return(Unknown, nil)
				mockTxManager.On("GetTransactionStatus", ctx, "test-message-id-1").Return(Unknown, errors.New("failed to find transaction with IdempotencyKey test-message-id-1"))
			},
			expectedStatus:  []TransactionStatus{Unknown},
			expectedCounter: 0,
			expectedError:   nil,
		},
		{
			name: "Unknown status with dummy error",
			setupMock: func() {
				mockTxManager.Mock = mock.Mock{}
				mockTxManager.On("GetTransactionStatus", ctx, "test-message-id-0").Return(Unknown, errors.New("dummy error"))
			},
			expectedStatus:  nil,
			expectedCounter: -1,
			expectedError:   errors.New("dummy error"),
		},
		{
			name: "Not unknown status with error",
			setupMock: func() {
				mockTxManager.Mock = mock.Mock{}
				mockTxManager.On("GetTransactionStatus", ctx, "test-message-id-0").Return(Finalized, errors.New("dummy error"))
			},
			expectedStatus:  nil,
			expectedCounter: -1,
			expectedError:   errors.New("dummy error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			statuses, counter, err := checker.CheckMessageStatus(ctx, msgID)
			assert.Equal(t, tc.expectedStatus, statuses)
			assert.Equal(t, tc.expectedCounter, counter)
			assert.Equal(t, tc.expectedError, err)
			mockTxManager.AssertExpectations(t)
		})
	}
}
