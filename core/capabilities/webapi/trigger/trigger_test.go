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
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	registrymock "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi/webapicap"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	gcmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
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
	address2             = "0x853d51d5d9935964267a5050aC53aa63ECA39bc6"
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
	lggr := logger.Test(t)
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

func gatewayRequest(t *testing.T, privateKey string, topics []string, methodName string) *jsonrpc.Request[json.RawMessage] {
	messageID := "12345"
	if methodName == "" {
		methodName = ghcapabilities.MethodWebAPITrigger
	}
	donID := "workflow_don_1"

	key, err := crypto.HexToECDSA(privateKey)
	require.NoError(t, err)

	payload := webapicap.TriggerRequestPayload{
		TriggerId:      TriggerType,
		TriggerEventId: "action_1234567890",
		Timestamp:      1234567890,
		Topics:         topics,
		Params: webapicap.TriggerRequestPayloadParams{
			"bid": "100",
			"ask": "101",
		},
	}

	payloadJSON, err := json.Marshal(payload)
	require.NoError(t, err)
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
	req, err := hc.ValidatedRequestFromMessage(msg)
	require.NoError(t, err)
	return req
}

func getResponseFromArg(arg interface{}) (ghcapabilities.TriggerResponsePayload, error) {
	resp := arg.(*jsonrpc.Response[json.RawMessage])
	var msg api.Message
	err := json.Unmarshal(*resp.Result, &msg)
	if err != nil {
		return ghcapabilities.TriggerResponsePayload{}, err
	}
	var payload ghcapabilities.TriggerResponsePayload
	err = json.Unmarshal(msg.Body.Payload, &payload)
	if err != nil {
		return ghcapabilities.TriggerResponsePayload{}, err
	}
	return payload, nil
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
		gatewayRequest := gatewayRequest(t, privateKey1, []string{"daily_price_update"}, "")

		th.connector.EXPECT().SignMessage(mock.Anything, mock.Anything).Return([]byte("signature"), nil).Once()
		th.connector.On("SendToGateway", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			resp, err2 := getResponseFromArg(args.Get(2))
			require.NoError(t, err2)
			require.Equal(t, ghcapabilities.TriggerResponsePayload{Status: "ACCEPTED"}, resp)
		}).Return(nil).Once()

		th.trigger.HandleGatewayMessage(ctx, "gateway1", gatewayRequest)

		received, chanErr := requireChanMsg(t, channel)
		require.Equal(t, TriggerType, received.Event.TriggerType)
		require.NoError(t, chanErr)

		requireNoChanMsg(t, channel2)
		data := received.Event.Outputs
		var payload webapicap.TriggerRequestPayload
		unwrapErr := data.UnwrapTo(&payload)
		require.NoError(t, unwrapErr)
		require.Equal(t, []string{"daily_price_update"}, payload.Topics)
	})

	t.Run("happy case single different topic 2 workflows.", func(t *testing.T) {
		gatewayRequest := gatewayRequest(t, privateKey1, []string{"ad_hoc_price_update"}, "")

		th.connector.EXPECT().SignMessage(mock.Anything, mock.Anything).Return([]byte("signature"), nil).Once()
		th.connector.On("SendToGateway", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			resp, err2 := getResponseFromArg(args.Get(2))
			require.NoError(t, err2)
			require.Equal(t, ghcapabilities.TriggerResponsePayload{Status: "ACCEPTED"}, resp)
		}).Return(nil).Once()

		th.trigger.HandleGatewayMessage(ctx, "gateway1", gatewayRequest)

		sent := <-channel
		require.Equal(t, TriggerType, sent.Event.TriggerType)
		data := sent.Event.Outputs
		var payload webapicap.TriggerRequestPayload
		unwrapErr := data.UnwrapTo(&payload)
		require.NoError(t, unwrapErr)
		require.Equal(t, []string{"ad_hoc_price_update"}, payload.Topics)

		sent2 := <-channel2
		require.Equal(t, TriggerType, sent2.Event.TriggerType)
		data2 := sent2.Event.Outputs
		var payload2 webapicap.TriggerRequestPayload
		err2 := data2.UnwrapTo(&payload2)
		require.NoError(t, err2)
		require.Equal(t, []string{"ad_hoc_price_update"}, payload2.Topics)
	})

	t.Run("sad case empty topic 2 workflows", func(t *testing.T) {
		gatewayRequest := gatewayRequest(t, privateKey1, nil, "")

		th.connector.EXPECT().SignMessage(mock.Anything, mock.Anything).Return([]byte("signature"), nil).Once()
		th.connector.On("SendToGateway", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			resp, err2 := getResponseFromArg(args.Get(2))
			require.NoError(t, err2)
			require.Equal(t, ghcapabilities.TriggerResponsePayload{Status: "ERROR", ErrorMessage: "empty Workflow Topics"}, resp)
		}).Return(nil).Once()

		th.trigger.HandleGatewayMessage(ctx, "gateway1", gatewayRequest)

		requireNoChanMsg(t, channel)
		requireNoChanMsg(t, channel2)
	})

	t.Run("sad case topic with no workflows", func(t *testing.T) {
		gatewayRequest := gatewayRequest(t, privateKey1, []string{"foo"}, "")
		th.connector.EXPECT().SignMessage(mock.Anything, mock.Anything).Return([]byte("signature"), nil).Once()
		th.connector.On("SendToGateway", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			resp, err2 := getResponseFromArg(args.Get(2))
			require.NoError(t, err2)
			require.Equal(t, ghcapabilities.TriggerResponsePayload{Status: "ERROR", ErrorMessage: "no Matching Workflow Topics"}, resp)
		}).Return(nil).Once()

		th.trigger.HandleGatewayMessage(ctx, "gateway1", gatewayRequest)
		requireNoChanMsg(t, channel)
		requireNoChanMsg(t, channel2)
	})

	t.Run("sad case Not Allowed Sender", func(t *testing.T) {
		gatewayRequest := gatewayRequest(t, privateKey2, []string{"ad_hoc_price_update"}, "")
		th.connector.EXPECT().SignMessage(mock.Anything, mock.Anything).Return([]byte("signature"), nil).Once()
		th.connector.On("SendToGateway", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			resp, err2 := getResponseFromArg(args.Get(2))
			require.NoError(t, err2)

			require.Equal(t, ghcapabilities.TriggerResponsePayload{Status: "ERROR", ErrorMessage: "unauthorized Sender 0x2dAC9f74Ee66e2D55ea1B8BE284caFedE048dB3A, messageID 12345"}, resp)
		}).Return(nil).Once()

		th.trigger.HandleGatewayMessage(ctx, "gateway1", gatewayRequest)
		requireNoChanMsg(t, channel)
		requireNoChanMsg(t, channel2)
	})

	t.Run("sad case Invalid Method", func(t *testing.T) {
		gatewayRequest := gatewayRequest(t, privateKey2, []string{"ad_hoc_price_update"}, "boo")
		th.connector.EXPECT().SignMessage(mock.Anything, mock.Anything).Return([]byte("signature"), nil).Once()
		th.connector.On("SendToGateway", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			resp, err2 := getResponseFromArg(args.Get(2))
			require.NoError(t, err2)
			require.Equal(t, ghcapabilities.TriggerResponsePayload{Status: "ERROR", ErrorMessage: "unsupported method boo"}, resp)
		}).Return(nil).Once()

		th.trigger.HandleGatewayMessage(ctx, "gateway1", gatewayRequest)
		requireNoChanMsg(t, channel)
		requireNoChanMsg(t, channel2)
	})

	t.Run("invalid message validation", func(t *testing.T) {
		// request with missing params
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "id",
			Method:  ghcapabilities.MethodWebAPITrigger,
		}
		th.connector.AssertNotCalled(t, "SignMessage")
		th.connector.AssertNotCalled(t, "SendToGateway")
		err = th.trigger.HandleGatewayMessage(ctx, "gateway1", req)
		require.NoError(t, err)
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

	gatewayRequest(t, privateKey1, []string{"daily_price_update"}, "")
}

func TestTriggerExecute2WorkflowsSameTopicDifferentAllowLists(t *testing.T) {
	th := setup(t)
	ctx := testutils.Context(t)
	ctx, cancelContext := context.WithDeadline(ctx, time.Now().Add(10*time.Second))
	Config, _ := workflowTriggerConfig(th, []string{address2}, []string{"daily_price_update"})
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

	Config2, err := workflowTriggerConfig(th, []string{address1}, []string{"daily_price_update"})
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
		gatewayRequest := gatewayRequest(t, privateKey1, []string{"daily_price_update"}, "")

		th.connector.EXPECT().SignMessage(mock.Anything, mock.Anything).Return([]byte("signature"), nil).Once()
		th.connector.On("SendToGateway", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			resp, err2 := getResponseFromArg(args.Get(2))
			require.NoError(t, err2)
			require.Equal(t, ghcapabilities.TriggerResponsePayload{Status: "ACCEPTED"}, resp)
		}).Return(nil).Once()

		th.trigger.HandleGatewayMessage(ctx, "gateway1", gatewayRequest)

		requireNoChanMsg(t, channel)
		received, chanErr := requireChanMsg(t, channel2)
		require.Equal(t, TriggerType, received.Event.TriggerType)
		require.NoError(t, chanErr)
		data := received.Event.Outputs
		var payload webapicap.TriggerRequestPayload
		unwrapErr := data.UnwrapTo(&payload)
		require.NoError(t, unwrapErr)
		require.Equal(t, []string{"daily_price_update"}, payload.Topics)
	})
	err = th.trigger.UnregisterTrigger(ctx, triggerReq)
	require.NoError(t, err)
	err = th.trigger.UnregisterTrigger(ctx, triggerReq2)
	require.NoError(t, err)
	cancelContext()
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
	_, open := <-channel
	require.False(t, open)
}
