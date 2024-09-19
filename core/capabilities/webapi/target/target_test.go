package target

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	registrymock "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	gcmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/webcapabilities"
)

const (
	workflowID1          = "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
	workflowExecutionID1 = "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeed"
	owner1               = "0x00000000000000000000000000000000000000aa"
)

type testHarness struct {
	registry         *registrymock.CapabilitiesRegistry
	connector        *gcmocks.GatewayConnector
	lggr             logger.Logger
	config           Config
	connectorHandler *ConnectorHandler
	capability       *Capability
}

func setup(t *testing.T) testHarness {
	registry := registrymock.NewCapabilitiesRegistry(t)
	connector := gcmocks.NewGatewayConnector(t)
	lggr := logger.Test(t)
	config := Config{
		RateLimiter: common.RateLimiterConfig{
			GlobalRPS:      100.0,
			GlobalBurst:    100,
			PerSenderRPS:   100.0,
			PerSenderBurst: 100,
		},
	}
	connectorHandler, err := NewConnectorHandler(connector, config, lggr)
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

func capabilityRequest(t *testing.T) capabilities.CapabilityRequest {
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
		"schedule":  SingleNode,
	})
	require.NoError(t, err)

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
	responsePayload, err := json.Marshal(webcapabilities.TargetResponsePayload{
		StatusCode: 200,
		Headers:    headers,
		Body:       body,
		Success:    true,
	})
	require.NoError(t, err)
	return &api.Message{
		Body: api.MessageBody{
			MessageId: msgID,
			Method:    webcapabilities.MethodWebAPITarget,
			Payload:   responsePayload,
		},
	}
}

func TestRegisterUnregister(t *testing.T) {
	th := setup(t)
	ctx := testutils.Context(t)

	regReq := capabilities.RegisterToWorkflowRequest{
		Metadata: capabilities.RegistrationMetadata{
			WorkflowID:    workflowID1,
			WorkflowOwner: owner1,
		},
	}
	err := th.capability.RegisterToWorkflow(ctx, regReq)
	require.NoError(t, err)

	err = th.capability.UnregisterFromWorkflow(ctx, capabilities.UnregisterFromWorkflowRequest{
		Metadata: capabilities.RegistrationMetadata{
			WorkflowID:    workflowID1,
			WorkflowOwner: owner1,
		},
	})
	require.NoError(t, err)
}

func TestCapability_Execute(t *testing.T) {
	th := setup(t)
	ctx := testutils.Context(t)
	th.connector.On("DonID").Return("donID")
	th.connector.On("GatewayIDs").Return([]string{"gateway2", "gateway1"})

	t.Run("unregistered workflow", func(t *testing.T) {
		req := capabilityRequest(t)
		_, err := th.capability.Execute(ctx, req)
		require.Error(t, err)
		require.Contains(t, err.Error(), "workflow is not registered")
	})

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

		th.connector.On("SignAndSendToGateway", mock.Anything, "gateway1", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			th.connectorHandler.HandleGatewayMessage(ctx, "gateway1", gatewayResp)
		}).Once()

		resp, err := th.capability.Execute(ctx, req)
		require.NoError(t, err)
		var values map[string]any
		err = resp.Value.UnwrapTo(&values)
		require.NoError(t, err)
		fmt.Printf("values %+v", values)
		statusCode, ok := values["statusCode"].(int64)
		require.True(t, ok)
		require.Equal(t, int64(200), statusCode)
		respBody, ok := values["body"].(string)
		require.True(t, ok)
		require.Equal(t, "response body", respBody)
	})

	t.Run("unsupported schedule", func(t *testing.T) {
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
			"timeoutMs": 1000,
			"schedule":  "invalid",
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
		require.Contains(t, err.Error(), "unsupported schedule")
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

		th.connector.On("SignAndSendToGateway", mock.Anything, "gateway1", mock.Anything).Return(errors.New("gateway error")).Once()
		_, err = th.capability.Execute(ctx, req)
		require.Error(t, err)
		require.Contains(t, err.Error(), "gateway error")
	})
}
