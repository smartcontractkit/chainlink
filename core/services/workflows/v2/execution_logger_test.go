package v2_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"

	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/executionlog"
)

// mockExecutionHelper is a simple mock implementation of host.ExecutionHelper
type mockExecutionHelper struct {
	mock.Mock
}

var _ host.ExecutionHelper = (*mockExecutionHelper)(nil)

func (m *mockExecutionHelper) CallCapability(ctx context.Context, request *sdkpb.CapabilityRequest) (*sdkpb.CapabilityResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdkpb.CapabilityResponse), args.Error(1)
}

func (m *mockExecutionHelper) GetSecrets(ctx context.Context, request *sdkpb.GetSecretsRequest) ([]*sdkpb.SecretResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*sdkpb.SecretResponse), args.Error(1)
}

func (m *mockExecutionHelper) GetWorkflowExecutionID() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockExecutionHelper) GetNodeTime() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *mockExecutionHelper) GetDONTime() (time.Time, error) {
	args := m.Called()
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *mockExecutionHelper) EmitUserLog(log string) error {
	args := m.Called(log)
	return args.Error(0)
}

func TestExecutionLogger_CallCapability(t *testing.T) {
	t.Parallel()

	t.Run("delegates to ExecutionHelper and calls OnExecutionEvent on success", func(t *testing.T) {
		mockHelper := new(mockExecutionHelper)
		var capturedEvent *executionlog.ExecutionEvent

		logger := v2.ExecutionCallbackHelper{
			ExecutionHelper: mockHelper,
			OnExecutionEvent: func(event *executionlog.ExecutionEvent) {
				capturedEvent = event
			},
		}

		ctx := context.Background()
		request := &sdkpb.CapabilityRequest{Id: "test-capability"}
		response := &sdkpb.CapabilityResponse{}

		mockHelper.On("CallCapability", ctx, request).Return(response, nil)

		result, err := logger.CallCapability(ctx, request)

		require.NoError(t, err)
		assert.Equal(t, response, result)
		require.NotNil(t, capturedEvent)
		capEvent := capturedEvent.GetCapabilityEvent()
		require.NotNil(t, capEvent)
		assert.Equal(t, request, capEvent.Request)
		assert.Equal(t, response, capEvent.Response)
		assert.Empty(t, capEvent.Error)
	})

	t.Run("delegates to ExecutionHelper and calls OnExecutionEvent on error", func(t *testing.T) {
		mockHelper := new(mockExecutionHelper)
		var capturedEvent *executionlog.ExecutionEvent

		logger := v2.ExecutionCallbackHelper{
			ExecutionHelper: mockHelper,
			OnExecutionEvent: func(event *executionlog.ExecutionEvent) {
				capturedEvent = event
			},
		}

		ctx := context.Background()
		request := &sdkpb.CapabilityRequest{Id: "test-capability"}
		expectedErr := errors.New("capability error")

		mockHelper.On("CallCapability", ctx, request).Return(nil, expectedErr)

		result, err := logger.CallCapability(ctx, request)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, result)
		require.NotNil(t, capturedEvent)
		capEvent := capturedEvent.GetCapabilityEvent()
		require.NotNil(t, capEvent)
		assert.Equal(t, request, capEvent.Request)
		assert.Nil(t, capEvent.Response)
		assert.Equal(t, expectedErr.Error(), capEvent.Error)
	})
}

func TestExecutionLogger_GetSecrets(t *testing.T) {
	t.Parallel()

	t.Run("delegates to ExecutionHelper and calls OnExecutionEvent on success", func(t *testing.T) {
		mockHelper := new(mockExecutionHelper)
		var capturedEvent *executionlog.ExecutionEvent

		logger := v2.ExecutionCallbackHelper{
			ExecutionHelper: mockHelper,
			OnExecutionEvent: func(event *executionlog.ExecutionEvent) {
				capturedEvent = event
			},
		}

		ctx := context.Background()
		request := &sdkpb.GetSecretsRequest{}
		response := []*sdkpb.SecretResponse{
			{},
		}

		mockHelper.On("GetSecrets", ctx, request).Return(response, nil)

		result, err := logger.GetSecrets(ctx, request)

		require.NoError(t, err)
		assert.Equal(t, response, result)
		require.NotNil(t, capturedEvent)
		secretsEvent := capturedEvent.GetSecretsEvent()
		require.NotNil(t, secretsEvent)
		assert.Equal(t, request, secretsEvent.Request)
		assert.Equal(t, response[0], secretsEvent.Response)
		assert.Empty(t, secretsEvent.Error)
	})

	t.Run("delegates to ExecutionHelper and calls OnExecutionEvent on error", func(t *testing.T) {
		mockHelper := new(mockExecutionHelper)
		var capturedEvent *executionlog.ExecutionEvent

		logger := v2.ExecutionCallbackHelper{
			ExecutionHelper: mockHelper,
			OnExecutionEvent: func(event *executionlog.ExecutionEvent) {
				capturedEvent = event
			},
		}

		ctx := context.Background()
		request := &sdkpb.GetSecretsRequest{}
		expectedErr := errors.New("secrets error")

		mockHelper.On("GetSecrets", ctx, request).Return(nil, expectedErr)

		result, err := logger.GetSecrets(ctx, request)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, result)
		require.NotNil(t, capturedEvent)
		secretsEvent := capturedEvent.GetSecretsEvent()
		require.NotNil(t, secretsEvent)
		assert.Equal(t, request, secretsEvent.Request)
		assert.Nil(t, secretsEvent.Response)
		assert.Equal(t, expectedErr.Error(), secretsEvent.Error)
	})
}

func TestExecutionLogger_GetWorkflowExecutionID(t *testing.T) {
	t.Parallel()

	t.Run("delegates to ExecutionHelper and calls OnExecutionEvent", func(t *testing.T) {
		mockHelper := new(mockExecutionHelper)
		var capturedEvent *executionlog.ExecutionEvent

		logger := v2.ExecutionCallbackHelper{
			ExecutionHelper: mockHelper,
			OnExecutionEvent: func(event *executionlog.ExecutionEvent) {
				capturedEvent = event
			},
		}

		expectedID := "workflow-execution-id-123"
		mockHelper.On("GetWorkflowExecutionID").Return(expectedID)

		result := logger.GetWorkflowExecutionID()

		assert.Equal(t, expectedID, result)
		require.NotNil(t, capturedEvent)
		assert.Equal(t, expectedID, capturedEvent.GetWorkflowIdEvent())
	})
}

func TestExecutionLogger_GetNodeTime(t *testing.T) {
	t.Parallel()

	t.Run("delegates to ExecutionHelper and calls OnExecutionEvent", func(t *testing.T) {
		mockHelper := new(mockExecutionHelper)
		var capturedEvent *executionlog.ExecutionEvent

		logger := v2.ExecutionCallbackHelper{
			ExecutionHelper: mockHelper,
			OnExecutionEvent: func(event *executionlog.ExecutionEvent) {
				capturedEvent = event
			},
		}

		expectedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		mockHelper.On("GetNodeTime").Return(expectedTime)

		result := logger.GetNodeTime()

		assert.Equal(t, expectedTime, result)
		require.NotNil(t, capturedEvent)
		nodeTimeEvent := capturedEvent.GetNodeTimeEvent()
		require.NotNil(t, nodeTimeEvent)
		require.NotNil(t, nodeTimeEvent.Response)
		assert.Equal(t, expectedTime, nodeTimeEvent.Response.AsTime())
	})
}

func TestExecutionLogger_GetDONTime(t *testing.T) {
	t.Parallel()

	t.Run("delegates to ExecutionHelper and calls OnExecutionEvent on success", func(t *testing.T) {
		mockHelper := new(mockExecutionHelper)
		var capturedEvent *executionlog.ExecutionEvent

		logger := v2.ExecutionCallbackHelper{
			ExecutionHelper: mockHelper,
			OnExecutionEvent: func(event *executionlog.ExecutionEvent) {
				capturedEvent = event
			},
		}

		expectedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		mockHelper.On("GetDONTime").Return(expectedTime, nil)

		result, err := logger.GetDONTime()

		require.NoError(t, err)
		assert.Equal(t, expectedTime, result)
		require.NotNil(t, capturedEvent)
		donTimeEvent := capturedEvent.GetDonTimeEvent()
		require.NotNil(t, donTimeEvent)
		require.NotNil(t, donTimeEvent.Response)
		assert.Equal(t, expectedTime, donTimeEvent.Response.AsTime())
		assert.Empty(t, donTimeEvent.Error)
	})

	t.Run("delegates to ExecutionHelper and calls OnExecutionEvent on error", func(t *testing.T) {
		mockHelper := new(mockExecutionHelper)
		var capturedEvent *executionlog.ExecutionEvent

		logger := v2.ExecutionCallbackHelper{
			ExecutionHelper: mockHelper,
			OnExecutionEvent: func(event *executionlog.ExecutionEvent) {
				capturedEvent = event
			},
		}

		expectedErr := errors.New("don time error")
		mockHelper.On("GetDONTime").Return(time.Time{}, expectedErr)

		result, err := logger.GetDONTime()

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.True(t, result.IsZero())
		require.NotNil(t, capturedEvent)
		donTimeEvent := capturedEvent.GetDonTimeEvent()
		require.NotNil(t, donTimeEvent)
		assert.Nil(t, donTimeEvent.Response)
		assert.Equal(t, expectedErr.Error(), donTimeEvent.Error)
	})
}

func TestExecutionLogger_EmitUserLog(t *testing.T) {
	t.Parallel()

	t.Run("delegates to ExecutionHelper and calls OnExecutionEvent on success", func(t *testing.T) {
		mockHelper := new(mockExecutionHelper)
		var capturedEvent *executionlog.ExecutionEvent

		logger := v2.ExecutionCallbackHelper{
			ExecutionHelper: mockHelper,
			OnExecutionEvent: func(event *executionlog.ExecutionEvent) {
				capturedEvent = event
			},
		}

		logMessage := "test log message"
		mockHelper.On("EmitUserLog", logMessage).Return(nil)

		err := logger.EmitUserLog(logMessage)

		require.NoError(t, err)
		require.NotNil(t, capturedEvent)
		emitLogEvent := capturedEvent.GetEmitLogEvent()
		require.NotNil(t, emitLogEvent)
		assert.Equal(t, logMessage, emitLogEvent.Log)
		assert.Empty(t, emitLogEvent.Error)
	})

	t.Run("delegates to ExecutionHelper and calls OnExecutionEvent on error", func(t *testing.T) {
		mockHelper := new(mockExecutionHelper)
		var capturedEvent *executionlog.ExecutionEvent

		logger := v2.ExecutionCallbackHelper{
			ExecutionHelper: mockHelper,
			OnExecutionEvent: func(event *executionlog.ExecutionEvent) {
				capturedEvent = event
			},
		}

		logMessage := "test log message"
		expectedErr := errors.New("emit log error")
		mockHelper.On("EmitUserLog", logMessage).Return(expectedErr)

		err := logger.EmitUserLog(logMessage)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		require.NotNil(t, capturedEvent)
		emitLogEvent := capturedEvent.GetEmitLogEvent()
		require.NotNil(t, emitLogEvent)
		assert.Equal(t, logMessage, emitLogEvent.Log)
		assert.Equal(t, expectedErr.Error(), emitLogEvent.Error)
	})
}

func TestExecutionLogger_OnExecutionEventNil(t *testing.T) {
	t.Parallel()

	t.Run("does not panic when OnExecutionEvent is nil", func(t *testing.T) {
		mockHelper := new(mockExecutionHelper)

		logger := v2.ExecutionCallbackHelper{
			ExecutionHelper:  mockHelper,
			OnExecutionEvent: nil,
		}

		ctx := context.Background()
		request := &sdkpb.CapabilityRequest{Id: "test"}
		response := &sdkpb.CapabilityResponse{}
		mockHelper.On("CallCapability", ctx, request).Return(response, nil)

		result, err := logger.CallCapability(ctx, request)

		require.NoError(t, err)
		assert.Equal(t, response, result)
	})
}
