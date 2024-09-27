package webapicapabilities

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	gwcommon "github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"

	handlermocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
)

const (
	defaultSendChannelBufferSize = 1000
	privateKey1                  = "65456ffb8af4a2b93959256a8e04f6f2fe0943579fb3c9c3350593aabb89023f"
	privateKey2                  = "65456ffb8af4a2b93959256a8e04f6f2fe0943579fb3c9c3350593aabb89023e"
	triggerID1                   = "5"
	triggerID2                   = "6"
	workflowID1                  = "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
	workflowExecutionID1         = "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeed"
	owner1                       = "0x00000000000000000000000000000000000000aa"
	address1                     = "0x853d51d5d9935964267a5050aC53aa63ECA39bc5"
)

func setupHandler(t *testing.T) (*handler, *handlermocks.DON, []gwcommon.TestNode) {
	lggr := logger.TestLogger(t)
	don := handlermocks.NewDON(t)

	handlerConfig := HandlerConfig{
		MaxAllowedMessageAgeSec: 30,
	}
	cfgBytes, err := json.Marshal(handlerConfig)
	require.NoError(t, err)
	donConfig := &config.DONConfig{
		Members: []config.NodeConfig{},
		F:       1,
	}
	nodes := gwcommon.NewTestNodes(t, 2)
	for id, n := range nodes {
		donConfig.Members = append(donConfig.Members, config.NodeConfig{
			Name:    fmt.Sprintf("node_%d", id),
			Address: n.Address,
		})
	}

	handler, err := NewWorkflowHandler(json.RawMessage(cfgBytes), donConfig, don, lggr)
	require.NoError(t, err)
	return handler, don, nodes
}

func triggerRequest(t *testing.T, privateKey string, topics string, methodName string, timestamp string) *api.Message {
	messageID := "12345"
	if methodName == "" {
		methodName = MethodWebAPITrigger
	}
	if timestamp == "" {
		timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	}
	donID := "workflow_don_1"

	key, err := crypto.HexToECDSA(privateKey)
	require.NoError(t, err)

	payload := `{
         "trigger_id": "` + TriggerType + `",
          "trigger_event_id": "action_1234567890",
          "timestamp": ` + timestamp + `,
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

func requireNoChanMsg[T any](t *testing.T, ch <-chan T) {
	timedOut := false
	select {
	case <-ch:
	case <-time.After(100 * time.Millisecond):
		timedOut = true
	}
	require.True(t, timedOut)
}

func TestHandlerReceiveHTTPMessageFromClient(t *testing.T) {
	handler, don, _ := setupHandler(t)
	ctx := testutils.Context(t)
	msg := triggerRequest(t, privateKey1, `["daily_price_update"]`, "", "")

	t.Run("happy case", func(t *testing.T) {
		ch := make(chan handlers.UserCallbackPayload, defaultSendChannelBufferSize)

		// sends to 2 dons
		don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			require.Equal(t, msg, args.Get(2))
		}).Return(nil).Once()
		don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			require.Equal(t, msg, args.Get(2))
		}).Return(nil).Once()

		err := handler.HandleUserMessage(ctx, msg, ch)
		require.NoError(t, err)
		requireNoChanMsg(t, ch)

		err = handler.HandleNodeMessage(ctx, msg, "")
		require.NoError(t, err)

		resp := <-ch
		require.Equal(t, handlers.UserCallbackPayload{Msg: msg, ErrCode: api.NoError, ErrMsg: ""}, resp)
		_, open := <-ch
		require.Equal(t, open, false)
	})

	t.Run("sad case invalid method", func(t *testing.T) {
		invalidMsg := triggerRequest(t, privateKey1, `["daily_price_update"]`, "foo", "")
		ch := make(chan handlers.UserCallbackPayload, defaultSendChannelBufferSize)
		err := handler.HandleUserMessage(ctx, invalidMsg, ch)
		require.NoError(t, err)
		resp := <-ch
		require.Equal(t, handlers.UserCallbackPayload{Msg: invalidMsg, ErrCode: api.HandlerError, ErrMsg: "invalid method foo"}, resp)
		_, open := <-ch
		require.Equal(t, open, false)
	})

	t.Run("sad case stale message", func(t *testing.T) {
		invalidMsg := triggerRequest(t, privateKey1, `["daily_price_update"]`, "", "123456")
		ch := make(chan handlers.UserCallbackPayload, defaultSendChannelBufferSize)
		err := handler.HandleUserMessage(ctx, invalidMsg, ch)
		require.NoError(t, err)
		resp := <-ch
		require.Equal(t, handlers.UserCallbackPayload{Msg: invalidMsg, ErrCode: api.HandlerError, ErrMsg: "stale message"}, resp)
		_, open := <-ch
		require.Equal(t, open, false)
	})
	// TODO: Validate Senders and rate limit chck, pending question in trigger about where senders and rate limits are validated
}
