package target

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	registrymock "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	gcmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

const (
	workflowID1          = "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
	workflowID2          = "44f129ea13948d1c4eaa2bbc0e72319266364cba12b789174732b2f72b57088d"
	workflowExecutionID1 = "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeed"
	owner1               = "0x00000000000000000000000000000000000000aa"
)

var defaultConfig = webapi.ServiceConfig{
	RateLimiter: common.RateLimiterConfig{
		GlobalRPS:      100.0,
		GlobalBurst:    100,
		PerSenderRPS:   100.0,
		PerSenderBurst: 100,
	},
}

type testHarness struct {
	registry         *registrymock.CapabilitiesRegistry
	connector        *gcmocks.GatewayConnector
	lggr             logger.Logger
	config           webapi.ServiceConfig
	connectorHandler *webapi.OutgoingConnectorHandler
	capability       *Capability
}

func setup(t *testing.T, config webapi.ServiceConfig) testHarness {
	registry := registrymock.NewCapabilitiesRegistry(t)
	connector := gcmocks.NewGatewayConnector(t)
	lggr := logger.Test(t)
	connector.EXPECT().GatewayIDs().Return([]string{"gateway1"})
	connectorHandler, err := webapi.NewOutgoingConnectorHandler(connector, config, ghcapabilities.MethodWebAPITarget, lggr)
	require.NoError(t, err)

	capability, err := NewCapability(config, registry, connectorHandler, lggr)
	require.NoError(t, err)

	return testHarness{
		registry:         registry,
		connector:        connector,
		lggr:             lggr,
		config:           config,
		connectorHandler: connectorHandler,
		capability:       capability,
	}
}

func emptyWfConfig(t *testing.T) *values.Map {
	wfConfig, err := values.NewMap(map[string]interface{}{})
	require.NoError(t, err)
	return wfConfig
}

func inputsAndConfig(t *testing.T) (*values.Map, *values.Map) {
	data := map[string]string{
		"key": "value",
	}
	jsonData, err := json.Marshal(data)
	require.NoError(t, err)
	encoded := base64.StdEncoding.EncodeToString(jsonData)
	targetInput := map[string]any{
		"url":     "http://example.com",
		"method":  "POST",
		"headers": map[string]string{"Content-Type": "application/json"},
		"body":    encoded,
	}
	inputs, err := values.NewMap(targetInput)
	require.NoError(t, err)
	wfConfig, err := values.NewMap(map[string]interface{}{
		"timeoutMs": 1000,
		"schedule":  webapi.SingleNode,
	})
	require.NoError(t, err)
	return inputs, wfConfig
}

func capabilityRequest(t *testing.T) capabilities.CapabilityRequest {
	inputs, wfConfig := inputsAndConfig(t)

	return capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID:          workflowID1,
			WorkflowExecutionID: workflowExecutionID1,
		},
		Inputs: inputs,
		Config: wfConfig,
	}
}

func gatewayResponse(t *testing.T, msgID string) *api.Message {
	headers := map[string]string{"Content-Type": "application/json"}
	body := []byte("response body")
	responsePayload, err := json.Marshal(ghcapabilities.Response{
		StatusCode:     200,
		Headers:        headers,
		Body:           body,
		ExecutionError: false,
	})
	require.NoError(t, err)
	return &api.Message{
		Body: api.MessageBody{
			MessageId: msgID,
			Method:    ghcapabilities.MethodWebAPITarget,
			Payload:   responsePayload,
		},
	}
}

func TestRegisterUnregister(t *testing.T) {
	th := setup(t, defaultConfig)
	ctx := testutils.Context(t)

	regReq := capabilities.RegisterToWorkflowRequest{
		Metadata: capabilities.RegistrationMetadata{
			WorkflowID:    workflowID1,
			WorkflowOwner: owner1,
		},
	}
	err := th.capability.RegisterToWorkflow(ctx, regReq)
	require.NoError(t, err)

	t.Run("happy case", func(t *testing.T) {
		err = th.capability.UnregisterFromWorkflow(ctx, capabilities.UnregisterFromWorkflowRequest{
			Metadata: capabilities.RegistrationMetadata{
				WorkflowID:    workflowID1,
				WorkflowOwner: owner1,
			},
		})
		require.NoError(t, err)
	})

	t.Run("unregister non-existent workflow no error", func(t *testing.T) {
		err = th.capability.UnregisterFromWorkflow(ctx, capabilities.UnregisterFromWorkflowRequest{
			Metadata: capabilities.RegistrationMetadata{
				WorkflowID:    workflowID2,
				WorkflowOwner: owner1,
			},
		})
		require.NoError(t, err)
	})

	t.Run("reregister idempotent", func(t *testing.T) {
		regReq := capabilities.RegisterToWorkflowRequest{
			Metadata: capabilities.RegistrationMetadata{
				WorkflowID:    workflowID1,
				WorkflowOwner: owner1,
			},
		}
		err := th.capability.RegisterToWorkflow(ctx, regReq)
		require.NoError(t, err)
	})
}

func TestCapability_Execute(t *testing.T) {
	th := setup(t, defaultConfig)
	ctx := testutils.Context(t)
	th.connector.EXPECT().DonID().Return("donID")
	th.connector.EXPECT().GatewayIDs().Return([]string{"gateway2", "gateway1"})

	t.Run("happy case", func(t *testing.T) {
		regReq := capabilities.RegisterToWorkflowRequest{
			Metadata: capabilities.RegistrationMetadata{
				WorkflowID:    workflowID1,
				WorkflowOwner: owner1,
			},
		}
		err := th.capability.RegisterToWorkflow(ctx, regReq)
		require.NoError(t, err)

		req := capabilityRequest(t)
		msgID, err := getMessageID(req)
		require.NoError(t, err)

		gatewayResp := gatewayResponse(t, msgID)
		th.connector.EXPECT().AwaitConnection(mock.Anything, "gateway1").Return(nil)
		th.connector.On("SignAndSendToGateway", mock.Anything, "gateway1", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			th.connectorHandler.HandleGatewayMessage(ctx, "gateway1", gatewayResp)
		}).Once()

		resp, err := th.capability.Execute(ctx, req)
		require.NoError(t, err)
		verifyResp(t, resp)
	})

	t.Run("context cancelled while waiting for gateway response", func(t *testing.T) {
		regReq := capabilities.RegisterToWorkflowRequest{
			Metadata: capabilities.RegistrationMetadata{
				WorkflowID:    workflowID1,
				WorkflowOwner: owner1,
			},
		}
		err := th.capability.RegisterToWorkflow(ctx, regReq)
		require.NoError(t, err)

		req := capabilityRequest(t)
		_, err = getMessageID(req)
		require.NoError(t, err)

		newCtx, cancel := context.WithCancel(ctx)
		th.connector.On("SignAndSendToGateway", mock.Anything, "gateway1", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			cancel()
		}).Once()

		_, err = th.capability.Execute(newCtx, req)
		require.Error(t, err)
		require.Contains(t, err.Error(), "context canceled")
	})

	t.Run("invalid workflow ID during execute", func(t *testing.T) {
		regReq := capabilities.RegisterToWorkflowRequest{
			Metadata: capabilities.RegistrationMetadata{
				WorkflowID:    workflowID1,
				WorkflowOwner: owner1,
			},
		}
		err := th.capability.RegisterToWorkflow(ctx, regReq)
		require.NoError(t, err)

		inputs, wfConfig := inputsAndConfig(t)
		req := capabilities.CapabilityRequest{
			Metadata: capabilities.RequestMetadata{
				WorkflowID:          "invalid",
				WorkflowExecutionID: workflowExecutionID1,
			},
			Inputs: inputs,
			Config: wfConfig,
		}

		_, err = th.capability.Execute(ctx, req)
		require.Error(t, err)
		require.ErrorContains(t, err, "workflow ID is invalid")
	})

	t.Run("invalid workflow execution ID during execute", func(t *testing.T) {
		regReq := capabilities.RegisterToWorkflowRequest{
			Metadata: capabilities.RegistrationMetadata{
				WorkflowID:    workflowID1,
				WorkflowOwner: owner1,
			},
		}
		err := th.capability.RegisterToWorkflow(ctx, regReq)
		require.NoError(t, err)

		inputs, wfConfig := inputsAndConfig(t)
		req := capabilities.CapabilityRequest{
			Metadata: capabilities.RequestMetadata{
				WorkflowID:          workflowID1,
				WorkflowExecutionID: "invalid",
			},
			Inputs: inputs,
			Config: wfConfig,
		}

		_, err = th.capability.Execute(ctx, req)
		require.Error(t, err)
		require.ErrorContains(t, err, "workflow execution ID is invalid")
	})

	t.Run("unsupported delivery mode", func(t *testing.T) {
		regReq := capabilities.RegisterToWorkflowRequest{
			Metadata: capabilities.RegistrationMetadata{
				WorkflowID:    workflowID1,
				WorkflowOwner: owner1,
			},
		}
		err := th.capability.RegisterToWorkflow(ctx, regReq)
		require.NoError(t, err)

		targetInput := map[string]any{
			"url":     "http://example.com",
			"method":  "POST",
			"headers": map[string]string{"Content-Type": "application/json"},
		}
		inputs, err := values.NewMap(targetInput)

		require.NoError(t, err)
		wfConfig, err := values.NewMap(map[string]interface{}{
			"timeoutMs":    1000,
			"deliveryMode": "invalid",
		})
		require.NoError(t, err)

		req := capabilities.CapabilityRequest{
			Metadata: capabilities.RequestMetadata{
				WorkflowID:          workflowID1,
				WorkflowExecutionID: workflowExecutionID1,
			},
			Inputs: inputs,
			Config: wfConfig,
		}

		_, err = th.capability.Execute(ctx, req)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unsupported delivery mode")
	})

	t.Run("gateway connector error", func(t *testing.T) {
		regReq := capabilities.RegisterToWorkflowRequest{
			Metadata: capabilities.RegistrationMetadata{
				WorkflowID:    workflowID1,
				WorkflowOwner: owner1,
			},
		}
		err := th.capability.RegisterToWorkflow(ctx, regReq)
		require.NoError(t, err)

		req := capabilityRequest(t)
		require.NoError(t, err)

		th.connector.EXPECT().SignAndSendToGateway(mock.Anything, "gateway1", mock.Anything).Return(errors.New("gateway error")).Once()
		_, err = th.capability.Execute(ctx, req)
		require.Error(t, err)
		require.Contains(t, err.Error(), "gateway error")
	})

	t.Run("empty workflow config", func(t *testing.T) {
		regReq := capabilities.RegisterToWorkflowRequest{
			Metadata: capabilities.RegistrationMetadata{
				WorkflowID:    workflowID1,
				WorkflowOwner: owner1,
			},
		}
		err := th.capability.RegisterToWorkflow(ctx, regReq)
		require.NoError(t, err)

		inputs, _ := inputsAndConfig(t)
		req := capabilities.CapabilityRequest{
			Metadata: capabilities.RequestMetadata{
				WorkflowID:          workflowID1,
				WorkflowExecutionID: workflowExecutionID1,
			},
			Inputs: inputs,
			Config: emptyWfConfig(t),
		}

		msgID, err := getMessageID(req)
		require.NoError(t, err)
		gatewayResp := gatewayResponse(t, msgID)
		th.connector.On("SignAndSendToGateway", mock.Anything, "gateway1", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			th.connectorHandler.HandleGatewayMessage(ctx, "gateway1", gatewayResp)
		}).Once()

		resp, err := th.capability.Execute(ctx, req)
		require.NoError(t, err)
		verifyResp(t, resp)
	})
}

func verifyResp(t *testing.T, resp capabilities.CapabilityResponse) {
	var values map[string]any
	err := resp.Value.UnwrapTo(&values)
	require.NoError(t, err)
	statusCode, ok := values["statusCode"].(int64)
	require.True(t, ok)
	require.Equal(t, int64(200), statusCode)
	respBody, ok := values["body"].([]byte)
	require.True(t, ok)
	require.Equal(t, "response body", string(respBody))
}
