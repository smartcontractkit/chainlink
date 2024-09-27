package trigger

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	registrymock "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	corelogger "github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	gcmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/webapicapabilities"
)

const (
	privateKey1          = "65456ffb8af4a2b93959256a8e04f6f2fe0943579fb3c9c3350593aabb89023f"
	privateKey2          = "65456ffb8af4a2b93959256a8e04f6f2fe0943579fb3c9c3350593aabb89023e"
	triggerID1           = "5"
	triggerID2           = "6"
	workflowID1          = "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
	workflowExecutionID1 = "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeed"
	owner1               = "0x00000000000000000000000000000000000000aa"
	address1             = "0x853d51d5d9935964267a5050aC53aa63ECA39bc5"
)

type testHarness struct {
	registry  *registrymock.CapabilitiesRegistry
	connector *gcmocks.GatewayConnector
	lggr      logger.Logger
	config    string
	trigger   *triggerConnectorHandler
}

func workflowTriggerConfig(_ testHarness, addresses []string, topics []string) (*values.Map, error) {
	var rateLimitConfig, err = values.NewMap(map[string]any{
		"GlobalRPS":      100.0,
		"GlobalBurst":    101,
		"PerSenderRPS":   102.0,
		"PerSenderBurst": 103,
	})
	if err != nil {
		return nil, err
	}

	triggerRegistrationConfig, err := values.NewMap(map[string]interface{}{
		"RateLimiter":    rateLimitConfig,
		"AllowedSenders": addresses,
		"AllowedTopics":  topics,
		"RequiredParams": []string{"bid", "ask"},
	})
	return triggerRegistrationConfig, err
}

func setup(t *testing.T) testHarness {
	registry := registrymock.NewCapabilitiesRegistry(t)
	connector := gcmocks.NewGatewayConnector(t)
	lggr := corelogger.TestLogger(t)
	config := ""

	trigger, err := NewTrigger(config, registry, connector, lggr)
	require.NoError(t, err)

	return testHarness{
		registry:  registry,
		connector: connector,
		lggr:      lggr,
		config:    config,
		trigger:   trigger,
	}
}

func gatewayRequest(t *testing.T, privateKey string, topics string, methodName string) *api.Message {
	messageID := "12345"
	if methodName == "" {
		methodName = webapicapabilities.MethodWebAPITrigger
	}
	donID := "workflow_don_1"

	key, err := crypto.HexToECDSA(privateKey)
	require.NoError(t, err)

	payload := `{
         "trigger_id": "` + TriggerType + `",
          "trigger_event_id": "action_1234567890",
          "timestamp": 1234567890,
          "topics": ` + topics + `,
					"params": {
						"bid": "101",
						"ask": "102"
					}
        }
`
	payloadJSON := []byte(payload)
	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: messageID,
			Method:    methodName,
			DonId:     donID,
			Payload:   json.RawMessage(payloadJSON),
		},
	}
	err = msg.Sign(key)
	require.NoError(t, err)
	return msg
}

func getResponseFromArg(arg interface{}) (webapicapabilities.TriggerResponsePayload, error) {
	var response webapicapabilities.TriggerResponsePayload
	err := json.Unmarshal((&(arg.(*api.Message)).Body).Payload, &response)
	return response, err
}

func requireNoChanMsg[T any](t *testing.T, ch <-chan T) {
	timedOut := false
	select {
	case <-ch:
	case <-time.After(100 * time.Millisecond):
		timedOut = true
	}
	require.True(t, timedOut)
}

func requireChanMsg[T capabilities.TriggerResponse](t *testing.T, ch <-chan capabilities.TriggerResponse) (capabilities.TriggerResponse, error) {
	timedOut := false
	select {
	case resp := <-ch:
		return resp, nil
	case <-time.After(100 * time.Millisecond):
		timedOut = true
	}
	require.False(t, timedOut)
	return capabilities.TriggerResponse{}, errors.New("channel timeout")
}

func TestTriggerExecute(t *testing.T) {
	th := setup(t)
	ctx := testutils.Context(t)
	ctx, cancelContext := context.WithDeadline(ctx, time.Now().Add(10*time.Second))
	Config, _ := workflowTriggerConfig(th, []string{address1}, []string{"daily_price_update", "ad_hoc_price_update"})
	triggerReq := capabilities.TriggerRegistrationRequest{
		TriggerID: triggerID1,
		Metadata: capabilities.RequestMetadata{
			WorkflowID:    workflowID1,
			WorkflowOwner: owner1,
		},
		Config: Config,
	}
	channel, err := th.trigger.RegisterTrigger(ctx, triggerReq)
	require.NoError(t, err)

	Config2, err := workflowTriggerConfig(th, []string{address1}, []string{"daily_price_update2", "ad_hoc_price_update"})
	require.NoError(t, err)

	triggerReq2 := capabilities.TriggerRegistrationRequest{
		TriggerID: triggerID2,
		Metadata: capabilities.RequestMetadata{
			WorkflowID:    workflowID1,
			WorkflowOwner: owner1,
		},
		Config: Config2,
	}
	channel2, err := th.trigger.RegisterTrigger(ctx, triggerReq2)
	require.NoError(t, err)

	t.Run("happy case single topic to single workflow", func(t *testing.T) {
		gatewayRequest := gatewayRequest(t, privateKey1, `["daily_price_update"]`, "")

		th.connector.On("SendToGateway", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			resp, _ := getResponseFromArg(args.Get(2))
			require.Equal(t, webapicapabilities.TriggerResponsePayload{Status: "ACCEPTED"}, resp)
		}).Return(nil).Once()

		th.trigger.HandleGatewayMessage(ctx, "gateway1", gatewayRequest)

		received, err := requireChanMsg(t, channel)
		require.Equal(t, received.Event.TriggerType, TriggerType)
		require.NoError(t, err)

		requireNoChanMsg(t, channel2)
		data := received.Event.Outputs
		var payload webapicapabilities.TriggerRequestPayload
		err = data.UnwrapTo(&payload)
		require.NoError(t, err)
		require.Equal(t, payload.Topics, []string{"daily_price_update"})
	})

	t.Run("happy case single different topic 2 workflows.", func(t *testing.T) {
		gatewayRequest := gatewayRequest(t, privateKey1, `["ad_hoc_price_update"]`, "")

		th.connector.On("SendToGateway", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			resp, _ := getResponseFromArg(args.Get(2))
			require.Equal(t, webapicapabilities.TriggerResponsePayload{Status: "ACCEPTED"}, resp)
		}).Return(nil).Once()

		th.trigger.HandleGatewayMessage(ctx, "gateway1", gatewayRequest)

		sent := <-channel
		require.Equal(t, sent.Event.TriggerType, TriggerType)
		data := sent.Event.Outputs
		var payload webapicapabilities.TriggerRequestPayload
		err := data.UnwrapTo(&payload)
		require.NoError(t, err)
		require.Equal(t, payload.Topics, []string{"ad_hoc_price_update"})

		sent2 := <-channel2
		require.Equal(t, sent2.Event.TriggerType, TriggerType)
		data2 := sent2.Event.Outputs
		var payload2 webapicapabilities.TriggerRequestPayload
		err2 := data2.UnwrapTo(&payload2)
		require.NoError(t, err2)
		require.Equal(t, payload2.Topics, []string{"ad_hoc_price_update"})
	})

	t.Run("sad case empty topic 2 workflows", func(t *testing.T) {
		gatewayRequest := gatewayRequest(t, privateKey1, `[]`, "")

		th.connector.On("SendToGateway", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			resp, _ := getResponseFromArg(args.Get(2))
			require.Equal(t, webapicapabilities.TriggerResponsePayload{Status: "ERROR", ErrorMessage: "No Matching Workflow Topics"}, resp)
		}).Return(nil).Once()

		th.trigger.HandleGatewayMessage(ctx, "gateway1", gatewayRequest)

		requireNoChanMsg(t, channel)
		requireNoChanMsg(t, channel2)
	})

	t.Run("sad case topic with no workflows", func(t *testing.T) {
		gatewayRequest := gatewayRequest(t, privateKey1, `["foo"]`, "")
		th.connector.On("SendToGateway", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			resp, _ := getResponseFromArg(args.Get(2))
			require.Equal(t, webapicapabilities.TriggerResponsePayload{Status: "ERROR", ErrorMessage: "No Matching Workflow Topics"}, resp)
		}).Return(nil).Once()

		th.trigger.HandleGatewayMessage(ctx, "gateway1", gatewayRequest)
		requireNoChanMsg(t, channel)
		requireNoChanMsg(t, channel2)
	})

	t.Run("sad case Not Allowed Sender", func(t *testing.T) {
		gatewayRequest := gatewayRequest(t, privateKey2, `["ad_hoc_price_update"]`, "")
		th.connector.On("SendToGateway", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			resp, _ := getResponseFromArg(args.Get(2))

			require.Equal(t, webapicapabilities.TriggerResponsePayload{Status: "ERROR", ErrorMessage: "Unauthorized Sender"}, resp)
		}).Return(nil).Once()

		th.trigger.HandleGatewayMessage(ctx, "gateway1", gatewayRequest)
		requireNoChanMsg(t, channel)
		requireNoChanMsg(t, channel2)
	})

	t.Run("sad case Invalid Method", func(t *testing.T) {
		gatewayRequest := gatewayRequest(t, privateKey2, `["ad_hoc_price_update"]`, "boo")
		th.connector.On("SendToGateway", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			resp, _ := getResponseFromArg(args.Get(2))
			require.Equal(t, webapicapabilities.TriggerResponsePayload{Status: "ERROR", ErrorMessage: "unsupported method boo"}, resp)
		}).Return(nil).Once()

		th.trigger.HandleGatewayMessage(ctx, "gateway1", gatewayRequest)
		requireNoChanMsg(t, channel)
		requireNoChanMsg(t, channel2)
	})

	err = th.trigger.UnregisterTrigger(ctx, triggerReq)
	require.NoError(t, err)
	err = th.trigger.UnregisterTrigger(ctx, triggerReq2)
	require.NoError(t, err)
	cancelContext()
}

func TestRegisterNoAllowedSenders(t *testing.T) {
	th := setup(t)
	ctx := testutils.Context(t)
	Config, _ := workflowTriggerConfig(th, []string{}, []string{"daily_price_update"})

	triggerReq := capabilities.TriggerRegistrationRequest{
		TriggerID: triggerID1,
		Metadata: capabilities.RequestMetadata{
			WorkflowID:    workflowID1,
			WorkflowOwner: owner1,
		},
		Config: Config,
	}
	_, err := th.trigger.RegisterTrigger(ctx, triggerReq)
	require.Error(t, err)

	gatewayRequest(t, privateKey1, `["daily_price_update"]`, "")
}

func TestRegisterUnregister(t *testing.T) {
	th := setup(t)
	ctx := testutils.Context(t)
	Config, err := workflowTriggerConfig(th, []string{address1}, []string{"daily_price_update"})
	require.NoError(t, err)

	triggerReq := capabilities.TriggerRegistrationRequest{
		TriggerID: triggerID1,
		Metadata: capabilities.RequestMetadata{
			WorkflowID:    workflowID1,
			WorkflowOwner: owner1,
		},
		Config: Config,
	}

	channel, err := th.trigger.RegisterTrigger(ctx, triggerReq)
	require.NoError(t, err)
	require.NotEmpty(t, th.trigger.registeredWorkflows[triggerID1])

	err = th.trigger.UnregisterTrigger(ctx, triggerReq)
	require.NoError(t, err)
	sent, open := <-channel
	require.Equal(t, open, false)
	require.Equal(t, capabilities.TriggerResponse(capabilities.TriggerResponse{Event: capabilities.TriggerEvent{TriggerType: "", ID: "", Outputs: (*values.Map)(nil)}, Err: error(nil)}), sent)
}
