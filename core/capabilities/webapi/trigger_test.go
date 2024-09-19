package webapi

import (
	"encoding/json"
	"flag"
	"testing"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	registrymock "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	corelogger "github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	gcmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

const (
	workflowID1          = "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
	workflowExecutionID1 = "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeed"
	owner1               = "0x00000000000000000000000000000000000000aa"
)

type testHarness struct {
	registry  *registrymock.CapabilitiesRegistry
	connector *gcmocks.GatewayConnector
	lggr      logger.Logger
	config    TriggerConfig
	trigger   *triggerConnectorHandler
}

func setup(t *testing.T) testHarness {
	privateKey, _ := testutils.NewPrivateKeyAndAddress(t)
	registry := registrymock.NewCapabilitiesRegistry(t)
	connector := gcmocks.NewGatewayConnector(t)
	lggr := corelogger.TestLogger(t)
	config := TriggerConfig{
		RateLimiter: common.RateLimiterConfig{
			GlobalRPS:      100.0,
			GlobalBurst:    100,
			PerSenderRPS:   100.0,
			PerSenderBurst: 100,
		},
		AllowedSenders: []ethCommon.Address{ethCommon.HexToAddress("a")},
	}
	trigger, err := NewTrigger(config, registry, connector, privateKey, lggr)
	require.NoError(t, err)

	return testHarness{
		registry:  registry,
		connector: connector,
		lggr:      lggr,
		config:    config,
		trigger:   trigger,
	}
}

func gatewayRequest(t *testing.T) *api.Message {
	// TODO: are flags like this ok? this is how the upload_workflow test script does it
	privateKey := flag.String("private_key", "65456ffb8af4a2b93959256a8e04f6f2fe0943579fb3c9c3350593aabb89023f", "Private key to sign the message with")
	messageID := flag.String("id", "12345", "Request ID")
	methodName := flag.String("method", "web_trigger", "Method name")
	donID := flag.String("don_id", "workflow_don_1", "DON ID")

	flag.Parse()
	key, err := crypto.HexToECDSA(*privateKey)
	require.NoError(t, err)

	payload := `{
          trigger_id: "web-trigger@1.0.0",
          trigger_event_id: "action_1234567890",
          timestamp: 1234567890,
          topics: ["daily_price_update"],
					params: {
						bid: "101",
						ask: "102"
					}
        }
`
	payloadJSON := []byte(payload)
	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: *messageID,
			Method:    *methodName,
			DonId:     *donID,
			Payload:   json.RawMessage(payloadJSON),
		},
	}
	err = msg.Sign(key)
	require.NoError(t, err)

	return msg
}

func TestCapability_Execute(t *testing.T) {
	th := setup(t)
	ctx := testutils.Context(t)

	t.Run("happy case", func(t *testing.T) {
		triggerReq := capabilities.TriggerRegistrationRequest{
			Metadata: capabilities.RequestMetadata{
				WorkflowID:    workflowID1,
				WorkflowOwner: owner1,
			},
		}
		_, err := th.trigger.RegisterTrigger(ctx, triggerReq)
		require.NoError(t, err)

		gatewayRequest := gatewayRequest(t)

		th.connector.On("SendToGateway", mock.Anything, mock.Anything).Return(nil).Once()

		// TODO: verify SendToGateway called
		th.trigger.HandleGatewayMessage(ctx, "gateway1", gatewayRequest)

		// TODO: verify message sent to trigger channel
	})

	// TODO: allowedSenders fail
	// TODO: rateLimit fail
	// TODO: empty allowedSenders
	// TODO: missing required parameters
	// TODO: invalid message
	// TODO: other edge cases?  empty topics?
	// TODO: Test duplicate messages, ie PENDING returned.
	// TODO: Test message sent to multiple trigger channels
}

func TestRegisterUnregister(t *testing.T) {
	th := setup(t)
	ctx := testutils.Context(t)

	triggerReq := capabilities.TriggerRegistrationRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID:    workflowID1,
			WorkflowOwner: owner1,
		},
	}
	_, err := th.trigger.RegisterTrigger(ctx, triggerReq)
	require.NoError(t, err)

	err = th.trigger.UnregisterTrigger(ctx, triggerReq)
	require.NoError(t, err)
}
