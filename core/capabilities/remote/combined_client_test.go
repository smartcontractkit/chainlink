package remote_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestCombinedClient_Info(t *testing.T) {
	ctx := testutils.Context(t)
	info := createTestCapabilityInfo("test-capability-info", commoncap.CapabilityTypeAction)
	client := remote.NewCombinedClient(info)

	returnedInfo, err := client.Info(ctx)
	require.NoError(t, err)
	assert.Equal(t, info, returnedInfo)
}

func TestCombinedClient_RegisterTrigger_Success(t *testing.T) {
	ctx := testutils.Context(t)
	info := createTestCapabilityInfo("test-trigger", commoncap.CapabilityTypeTrigger)
	client := remote.NewCombinedClient(info)

	mockTrigger := &mocks.TriggerCapability{}
	method := "test-method"

	client.SetTriggerSubscriber(method, mockTrigger)

	request := commoncap.TriggerRegistrationRequest{
		TriggerID: "test-trigger-id",
		Method:    method,
	}

	responseChan := make(<-chan commoncap.TriggerResponse, 1)
	mockTrigger.On("RegisterTrigger", ctx, request).Return(responseChan, nil)

	result, err := client.RegisterTrigger(ctx, request)
	require.NoError(t, err)
	assert.Equal(t, responseChan, result)
}

func TestCombinedClient_RegisterTrigger_MethodNotDefined(t *testing.T) {
	ctx := testutils.Context(t)
	info := createTestCapabilityInfo("test-trigger", commoncap.CapabilityTypeTrigger)
	client := remote.NewCombinedClient(info)

	request := commoncap.TriggerRegistrationRequest{
		TriggerID: "test-trigger-id",
		Method:    "undefined-method",
	}

	result, err := client.RegisterTrigger(ctx, request)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "method undefined-method not defined")
}

func TestCombinedClient_RegisterTrigger_ErrorFromSubscriber(t *testing.T) {
	ctx := testutils.Context(t)
	info := createTestCapabilityInfo("test-trigger", commoncap.CapabilityTypeTrigger)
	client := remote.NewCombinedClient(info)

	mockTrigger := &mocks.TriggerCapability{}
	method := "test-method"

	client.SetTriggerSubscriber(method, mockTrigger)

	request := commoncap.TriggerRegistrationRequest{
		TriggerID: "test-trigger-id",
		Method:    method,
	}

	expectedError := errors.New("registration failed")
	mockTrigger.On("RegisterTrigger", ctx, request).Return((<-chan commoncap.TriggerResponse)(nil), expectedError)

	result, err := client.RegisterTrigger(ctx, request)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

func TestCombinedClient_UnregisterTrigger_Success(t *testing.T) {
	ctx := testutils.Context(t)
	info := createTestCapabilityInfo("test-trigger", commoncap.CapabilityTypeTrigger)
	client := remote.NewCombinedClient(info)

	mockTrigger := &mocks.TriggerCapability{}
	method := "test-method"

	client.SetTriggerSubscriber(method, mockTrigger)

	request := commoncap.TriggerRegistrationRequest{
		TriggerID: "test-trigger-id",
		Method:    method,
	}

	mockTrigger.On("UnregisterTrigger", ctx, request).Return(nil)
	err := client.UnregisterTrigger(ctx, request)
	require.NoError(t, err)
}

func TestCombinedClient_UnregisterTrigger_MethodNotDefined(t *testing.T) {
	ctx := testutils.Context(t)
	info := createTestCapabilityInfo("test-trigger", commoncap.CapabilityTypeTrigger)
	client := remote.NewCombinedClient(info)

	request := commoncap.TriggerRegistrationRequest{
		TriggerID: "test-trigger-id",
		Method:    "undefined-method",
	}

	err := client.UnregisterTrigger(ctx, request)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "method undefined-method not defined")
}

func TestCombinedClient_UnregisterTrigger_ErrorFromSubscriber(t *testing.T) {
	ctx := testutils.Context(t)
	info := createTestCapabilityInfo("test-trigger", commoncap.CapabilityTypeTrigger)
	client := remote.NewCombinedClient(info)

	mockTrigger := &mocks.TriggerCapability{}
	method := "test-method"

	client.SetTriggerSubscriber(method, mockTrigger)

	request := commoncap.TriggerRegistrationRequest{
		TriggerID: "test-trigger-id",
		Method:    method,
	}

	expectedError := errors.New("unregistration failed")
	mockTrigger.On("UnregisterTrigger", ctx, request).Return(expectedError)

	err := client.UnregisterTrigger(ctx, request)
	require.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestCombinedClient_RegisterToWorkflow_NotSupported(t *testing.T) {
	ctx := testutils.Context(t)
	info := createTestCapabilityInfo("test-capability", commoncap.CapabilityTypeTrigger)
	client := remote.NewCombinedClient(info)

	request := commoncap.RegisterToWorkflowRequest{
		Metadata: commoncap.RegistrationMetadata{
			WorkflowID: "test-workflow",
		},
	}

	err := client.RegisterToWorkflow(ctx, request)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "RegisterToWorkflow is not supported by remote capabilities")
}

func TestCombinedClient_UnregisterFromWorkflow_NotSupported(t *testing.T) {
	ctx := testutils.Context(t)
	info := createTestCapabilityInfo("test-capability", commoncap.CapabilityTypeTrigger)
	client := remote.NewCombinedClient(info)

	request := commoncap.UnregisterFromWorkflowRequest{
		Metadata: commoncap.RegistrationMetadata{
			WorkflowID: "test-workflow",
		},
	}

	err := client.UnregisterFromWorkflow(ctx, request)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "UnregisterFromWorkflow is not supported by remote capabilities")
}

func TestCombinedClient_Execute_Success(t *testing.T) {
	ctx := testutils.Context(t)
	info := createTestCapabilityInfo("test-executable", commoncap.CapabilityTypeAction)
	client := remote.NewCombinedClient(info)

	mockExecutable := &mocks.ExecutableCapability{}
	method := "test-execute-method"

	client.SetExecutableClient(method, mockExecutable)

	request := commoncap.CapabilityRequest{
		Method: method,
		Config: nil,
		Inputs: nil,
		Metadata: commoncap.RequestMetadata{
			WorkflowID:          "test-workflow",
			WorkflowExecutionID: "test-execution",
		},
	}

	expectedResponse := commoncap.CapabilityResponse{
		Value: nil,
	}

	mockExecutable.On("Execute", ctx, request).Return(expectedResponse, nil)

	result, err := client.Execute(ctx, request)
	require.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
}

func TestCombinedClient_Execute_MethodNotDefined(t *testing.T) {
	ctx := testutils.Context(t)
	info := createTestCapabilityInfo("test-executable", commoncap.CapabilityTypeAction)
	client := remote.NewCombinedClient(info)

	request := commoncap.CapabilityRequest{
		Method: "undefined-method",
		Config: nil,
		Inputs: nil,
		Metadata: commoncap.RequestMetadata{
			WorkflowID:          "test-workflow",
			WorkflowExecutionID: "test-execution",
		},
	}

	result, err := client.Execute(ctx, request)
	require.Error(t, err)
	assert.Equal(t, commoncap.CapabilityResponse{}, result)
	assert.Contains(t, err.Error(), "method undefined-method not defined")
}

func TestCombinedClient_Execute_ErrorFromExecutable(t *testing.T) {
	ctx := testutils.Context(t)
	info := createTestCapabilityInfo("test-executable", commoncap.CapabilityTypeAction)
	client := remote.NewCombinedClient(info)

	// Create mock executable capability
	mockExecutable := &mocks.ExecutableCapability{}
	method := "test-execute-method"

	client.SetExecutableClient(method, mockExecutable)

	request := commoncap.CapabilityRequest{
		Method: method,
		Config: nil,
		Inputs: nil,
		Metadata: commoncap.RequestMetadata{
			WorkflowID:          "test-workflow",
			WorkflowExecutionID: "test-execution",
		},
	}

	expectedError := errors.New("execution failed")
	mockExecutable.On("Execute", ctx, request).Return(commoncap.CapabilityResponse{}, expectedError)

	result, err := client.Execute(ctx, request)
	require.Error(t, err)
	assert.Equal(t, commoncap.CapabilityResponse{}, result)
	assert.Equal(t, expectedError, err)
}

func TestCombinedClient_SetTriggerSubscriber(t *testing.T) {
	info := createTestCapabilityInfo("test-capability", commoncap.CapabilityTypeTrigger)

	client := remote.NewCombinedClient(info)
	mockTrigger := &mocks.TriggerCapability{}
	method := "test-method"

	client.SetTriggerSubscriber(method, mockTrigger)

	ctx := testutils.Context(t)
	request := commoncap.TriggerRegistrationRequest{
		TriggerID: "test-trigger-id",
		Method:    method,
	}

	responseChan := make(chan commoncap.TriggerResponse, 1)
	mockTrigger.On("RegisterTrigger", ctx, request).Return((<-chan commoncap.TriggerResponse)(responseChan), nil)

	_, err := client.RegisterTrigger(ctx, request)
	require.NoError(t, err)
}

func TestCombinedClient_SetExecutableClient(t *testing.T) {
	info := createTestCapabilityInfo("test-capability", commoncap.CapabilityTypeAction)

	client := remote.NewCombinedClient(info)
	mockExecutable := &mocks.ExecutableCapability{}
	method := "test-method"

	client.SetExecutableClient(method, mockExecutable)

	ctx := testutils.Context(t)
	request := commoncap.CapabilityRequest{
		Method: method,
		Config: nil,
		Inputs: nil,
		Metadata: commoncap.RequestMetadata{
			WorkflowID:          "test-workflow",
			WorkflowExecutionID: "test-execution",
		},
	}

	expectedResponse := commoncap.CapabilityResponse{
		Value: nil,
	}

	mockExecutable.On("Execute", ctx, request).Return(expectedResponse, nil)

	_, err := client.Execute(ctx, request)
	require.NoError(t, err)
}

func TestCombinedClient_MultipleMethodsAndCapabilities(t *testing.T) {
	ctx := testutils.Context(t)
	info := createTestCapabilityInfo("test-multi-capability", commoncap.CapabilityTypeAction)
	client := remote.NewCombinedClient(info)

	// Add multiple trigger subscribers
	mockTrigger1 := &mocks.TriggerCapability{}
	mockTrigger2 := &mocks.TriggerCapability{}
	triggerMethod1 := "trigger-method-1"
	triggerMethod2 := "trigger-method-2"

	client.SetTriggerSubscriber(triggerMethod1, mockTrigger1)
	client.SetTriggerSubscriber(triggerMethod2, mockTrigger2)

	client.SetTriggerSubscriber(triggerMethod1, mockTrigger1)
	client.SetTriggerSubscriber(triggerMethod2, mockTrigger2)

	// Add multiple executable clients
	mockExecutable1 := &mocks.ExecutableCapability{}
	mockExecutable2 := &mocks.ExecutableCapability{}
	execMethod1 := "exec-method-1"
	execMethod2 := "exec-method-2"

	client.SetExecutableClient(execMethod1, mockExecutable1)
	client.SetExecutableClient(execMethod2, mockExecutable2)

	// Test trigger method 1
	triggerRequest1 := commoncap.TriggerRegistrationRequest{
		TriggerID: "trigger-1",
		Method:    triggerMethod1,
	}
	responseChan1 := make(chan commoncap.TriggerResponse, 1)
	mockTrigger1.On("RegisterTrigger", ctx, triggerRequest1).Return((<-chan commoncap.TriggerResponse)(responseChan1), nil)

	_, err := client.RegisterTrigger(ctx, triggerRequest1)
	require.NoError(t, err)

	// Test trigger method 2
	triggerRequest2 := commoncap.TriggerRegistrationRequest{
		TriggerID: "trigger-2",
		Method:    triggerMethod2,
	}
	responseChan2 := make(chan commoncap.TriggerResponse, 1)
	mockTrigger2.On("RegisterTrigger", ctx, triggerRequest2).Return((<-chan commoncap.TriggerResponse)(responseChan2), nil)

	_, err = client.RegisterTrigger(ctx, triggerRequest2)
	require.NoError(t, err)

	// Test executable method 1
	execRequest1 := commoncap.CapabilityRequest{
		Method: execMethod1,
		Metadata: commoncap.RequestMetadata{
			WorkflowID:          "workflow-1",
			WorkflowExecutionID: "execution-1",
		},
	}
	execResponse1 := commoncap.CapabilityResponse{Value: nil}
	mockExecutable1.On("Execute", ctx, execRequest1).Return(execResponse1, nil)

	_, err = client.Execute(ctx, execRequest1)
	require.NoError(t, err)

	// Test executable method 2
	execRequest2 := commoncap.CapabilityRequest{
		Method: execMethod2,
		Metadata: commoncap.RequestMetadata{
			WorkflowID:          "workflow-2",
			WorkflowExecutionID: "execution-2",
		},
	}
	execResponse2 := commoncap.CapabilityResponse{Value: nil}
	mockExecutable2.On("Execute", ctx, execRequest2).Return(execResponse2, nil)

	_, err = client.Execute(ctx, execRequest2)
	require.NoError(t, err)

	// Assert all expectations
	mockTrigger1.AssertExpectations(t)
	mockTrigger2.AssertExpectations(t)
	mockExecutable1.AssertExpectations(t)
	mockExecutable2.AssertExpectations(t)
}

func createTestCapabilityInfo(id string, capType commoncap.CapabilityType) commoncap.CapabilityInfo {
	return commoncap.CapabilityInfo{
		ID:             id,
		CapabilityType: capType,
		Description:    "Test capability",
		DON:            nil,
		IsLocal:        false,
	}
}
