package functions_test

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"testing"
	"time"

	geth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/onsi/gomega"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/functions"
	sfmocks "github.com/smartcontractkit/chainlink/v2/core/services/functions/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	gwconnector "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	gcmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	fallowMocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/allowlist/mocks"
	fsubMocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/subscriptions/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
	s4mocks "github.com/smartcontractkit/chainlink/v2/core/services/s4/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newOffchainRequest(t *testing.T, sender []byte, ageSec uint64) (*api.Message, functions.RequestID) {
	requestId := make([]byte, 32)
	_, err := rand.Read(requestId)
	require.NoError(t, err)
	request := &functions.OffchainRequest{
		RequestId:         requestId,
		RequestInitiator:  sender,
		SubscriptionId:    1,
		SubscriptionOwner: sender,
		Timestamp:         uint64(time.Now().Unix()) - ageSec,
	}

	internalId := functions.InternalId(request.RequestInitiator, request.RequestId)
	req, err := json.Marshal(request)
	require.NoError(t, err)
	msg := &api.Message{
		Body: api.MessageBody{
			DonId:     "fun4",
			MessageId: "1",
			Method:    "heartbeat",
			Payload:   req,
		},
	}
	return msg, internalId
}

func TestFunctionsConnectorHandler(t *testing.T) {
	t.Parallel()

	logger := logger.TestLogger(t)
	privateKey, addr := testutils.NewPrivateKeyAndAddress(t)
	storage := s4mocks.NewStorage(t)
	connector := gcmocks.NewGatewayConnector(t)
	allowlist := fallowMocks.NewOnchainAllowlist(t)
	rateLimiter, err := hc.NewRateLimiter(hc.RateLimiterConfig{GlobalRPS: 100.0, GlobalBurst: 100, PerSenderRPS: 100.0, PerSenderBurst: 100})
	subscriptions := fsubMocks.NewOnchainSubscriptions(t)
	reportCh := make(chan *functions.OffchainResponse)
	offchainTransmitter := sfmocks.NewOffchainTransmitter(t)
	offchainTransmitter.On("ReportChannel", mock.Anything).Return(reportCh)
	listener := sfmocks.NewFunctionsListener(t)
	require.NoError(t, err)
	allowlist.On("Start", mock.Anything).Return(nil)
	allowlist.On("Close", mock.Anything).Return(nil)
	subscriptions.On("Start", mock.Anything).Return(nil)
	subscriptions.On("Close", mock.Anything).Return(nil)
	config := &config.PluginConfig{
		GatewayConnectorConfig: &gwconnector.ConnectorConfig{
			NodeAddress: addr.Hex(),
		},
		MinimumSubscriptionBalance: *assets.NewLinkFromJuels(100),
		RequestTimeoutSec:          1_000,
		AllowedHeartbeatInitiators: []string{crypto.PubkeyToAddress(privateKey.PublicKey).Hex()},
	}
	handler, err := functions.NewFunctionsConnectorHandler(config, privateKey, storage, allowlist, rateLimiter, subscriptions, listener, offchainTransmitter, logger)
	require.NoError(t, err)

	handler.SetConnector(connector)

	err = handler.Start(testutils.Context(t))
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, handler.Close())
	})

	t.Run("Sign", func(t *testing.T) {
		signature, err := handler.Sign([]byte("test"))
		require.NoError(t, err)

		signer, err := common.ExtractSigner(signature, []byte("test"))
		require.NoError(t, err)
		require.Equal(t, addr.Bytes(), signer)
	})

	t.Run("HandleGatewayMessage", func(t *testing.T) {
		t.Run("secrets_list", func(t *testing.T) {
			msg := api.Message{
				Body: api.MessageBody{
					DonId:     "fun4",
					MessageId: "1",
					Method:    "secrets_list",
					Sender:    addr.Hex(),
				},
			}
			require.NoError(t, msg.Sign(privateKey))

			ctx := testutils.Context(t)
			snapshot := []*s4.SnapshotRow{
				{SlotId: 1, Version: 1, Expiration: 1},
				{SlotId: 2, Version: 2, Expiration: 2},
			}
			storage.On("List", ctx, addr).Return(snapshot, nil).Once()
			allowlist.On("Allow", addr).Return(true).Once()
			connector.On("SendToGateway", ctx, "gw1", mock.Anything).Run(func(args mock.Arguments) {
				msg, ok := args[2].(*api.Message)
				require.True(t, ok)
				require.Equal(t, `{"success":true,"rows":[{"slot_id":1,"version":1,"expiration":1},{"slot_id":2,"version":2,"expiration":2}]}`, string(msg.Body.Payload))
			}).Return(nil).Once()

			handler.HandleGatewayMessage(ctx, "gw1", &msg)

			t.Run("orm error", func(t *testing.T) {
				storage.On("List", ctx, addr).Return(nil, errors.New("boom")).Once()
				allowlist.On("Allow", addr).Return(true).Once()
				connector.On("SendToGateway", ctx, "gw1", mock.Anything).Run(func(args mock.Arguments) {
					msg, ok := args[2].(*api.Message)
					require.True(t, ok)
					require.Equal(t, `{"success":false,"error_message":"Failed to list secrets: boom"}`, string(msg.Body.Payload))
				}).Return(nil).Once()

				handler.HandleGatewayMessage(ctx, "gw1", &msg)
			})

			t.Run("not allowed", func(t *testing.T) {
				allowlist.On("Allow", addr).Return(false).Once()
				handler.HandleGatewayMessage(ctx, "gw1", &msg)
			})
		})

		t.Run("secrets_set", func(t *testing.T) {
			ctx := testutils.Context(t)
			key := s4.Key{
				Address: addr,
				SlotId:  3,
				Version: 4,
			}
			record := s4.Record{
				Expiration: 5,
				Payload:    []byte("test"),
			}
			signature, err := s4.NewEnvelopeFromRecord(&key, &record).Sign(privateKey)
			signatureB64 := base64.StdEncoding.EncodeToString(signature)
			require.NoError(t, err)

			msg := api.Message{
				Body: api.MessageBody{
					DonId:     "fun4",
					MessageId: "1",
					Method:    "secrets_set",
					Sender:    addr.Hex(),
					Payload:   json.RawMessage(`{"slot_id":3,"version":4,"expiration":5,"payload":"dGVzdA==","signature":"` + signatureB64 + `"}`),
				},
			}
			require.NoError(t, msg.Sign(privateKey))

			storage.On("Put", ctx, &key, &record, signature).Return(nil).Once()
			allowlist.On("Allow", addr).Return(true).Once()
			subscriptions.On("GetMaxUserBalance", mock.Anything).Return(big.NewInt(100), nil).Once()
			connector.On("SendToGateway", ctx, "gw1", mock.Anything).Run(func(args mock.Arguments) {
				msg, ok := args[2].(*api.Message)
				require.True(t, ok)
				require.Equal(t, `{"success":true}`, string(msg.Body.Payload))
			}).Return(nil).Once()

			handler.HandleGatewayMessage(ctx, "gw1", &msg)

			t.Run("orm error", func(t *testing.T) {
				storage.On("Put", ctx, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("boom")).Once()
				allowlist.On("Allow", addr).Return(true).Once()
				subscriptions.On("GetMaxUserBalance", mock.Anything).Return(big.NewInt(100), nil).Once()
				connector.On("SendToGateway", ctx, "gw1", mock.Anything).Run(func(args mock.Arguments) {
					msg, ok := args[2].(*api.Message)
					require.True(t, ok)
					require.Equal(t, `{"success":false,"error_message":"Failed to set secret: boom"}`, string(msg.Body.Payload))
				}).Return(nil).Once()

				handler.HandleGatewayMessage(ctx, "gw1", &msg)
			})

			t.Run("missing signature", func(t *testing.T) {
				msg.Body.Payload = json.RawMessage(`{"slot_id":3,"version":4,"expiration":5,"payload":"dGVzdA=="}`)
				require.NoError(t, msg.Sign(privateKey))
				storage.On("Put", ctx, mock.Anything, mock.Anything, mock.Anything).Return(s4.ErrWrongSignature).Once()
				allowlist.On("Allow", addr).Return(true).Once()
				subscriptions.On("GetMaxUserBalance", mock.Anything).Return(big.NewInt(100), nil).Once()
				connector.On("SendToGateway", ctx, "gw1", mock.Anything).Run(func(args mock.Arguments) {
					msg, ok := args[2].(*api.Message)
					require.True(t, ok)
					require.Equal(t, `{"success":false,"error_message":"Failed to set secret: wrong signature"}`, string(msg.Body.Payload))
				}).Return(nil).Once()

				handler.HandleGatewayMessage(ctx, "gw1", &msg)
			})

			t.Run("malformed request", func(t *testing.T) {
				msg.Body.Payload = json.RawMessage(`{sdfgdfgoscsicosd:sdf:::sdf ::; xx}`)
				require.NoError(t, msg.Sign(privateKey))
				allowlist.On("Allow", addr).Return(true).Once()
				subscriptions.On("GetMaxUserBalance", mock.Anything).Return(big.NewInt(100), nil).Once()
				connector.On("SendToGateway", ctx, "gw1", mock.Anything).Run(func(args mock.Arguments) {
					msg, ok := args[2].(*api.Message)
					require.True(t, ok)
					require.Equal(t, `{"success":false,"error_message":"Bad request to set secret: invalid character 's' looking for beginning of object key string"}`, string(msg.Body.Payload))
				}).Return(nil).Once()

				handler.HandleGatewayMessage(ctx, "gw1", &msg)
			})

			t.Run("insufficient balance", func(t *testing.T) {
				allowlist.On("Allow", addr).Return(true).Once()
				subscriptions.On("GetMaxUserBalance", mock.Anything).Return(big.NewInt(0), nil).Once()
				connector.On("SendToGateway", ctx, "gw1", mock.Anything).Run(func(args mock.Arguments) {
					msg, ok := args[2].(*api.Message)
					require.True(t, ok)
					require.Equal(t, `{"success":false,"error_message":"user subscription has insufficient balance"}`, string(msg.Body.Payload))
				}).Return(nil).Once()

				handler.HandleGatewayMessage(ctx, "gw1", &msg)
			})
		})

		t.Run("unsupported method", func(t *testing.T) {
			msg := api.Message{
				Body: api.MessageBody{
					DonId:     "fun4",
					MessageId: "1",
					Method:    "foobar",
					Sender:    addr.Hex(),
					Payload:   []byte("whatever"),
				},
			}
			require.NoError(t, msg.Sign(privateKey))

			allowlist.On("Allow", addr).Return(true).Once()
			handler.HandleGatewayMessage(testutils.Context(t), "gw1", &msg)
		})
	})

	t.Run("heartbeat success", func(t *testing.T) {
		ctx := testutils.Context(t)
		msg, internalId := newOffchainRequest(t, addr.Bytes(), 0)
		require.NoError(t, msg.Sign(privateKey))

		// first call to trigger the request
		var response functions.HeartbeatResponse
		allowlist.On("Allow", addr).Return(true).Once()
		handlerCalled := make(chan struct{})
		listener.On("HandleOffchainRequest", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			handlerCalled <- struct{}{}
		}).Return(nil).Once()
		connector.On("SendToGateway", mock.Anything, "gw1", mock.Anything).Run(func(args mock.Arguments) {
			respMsg, ok := args[2].(*api.Message)
			require.True(t, ok)
			require.NoError(t, json.Unmarshal(respMsg.Body.Payload, &response))
			require.Equal(t, functions.RequestStatePending, response.Status)
		}).Return(nil).Once()
		handler.HandleGatewayMessage(ctx, "gw1", msg)
		<-handlerCalled

		// async response computation
		reportCh <- &functions.OffchainResponse{
			RequestId: internalId[:],
			Result:    []byte("ok!"),
		}
		reportCh <- &functions.OffchainResponse{} // sending second item to make sure the first one got processed

		// second call to collect the response
		allowlist.On("Allow", addr).Return(true).Once()
		connector.On("SendToGateway", mock.Anything, "gw1", mock.Anything).Run(func(args mock.Arguments) {
			respMsg, ok := args[2].(*api.Message)
			require.True(t, ok)
			require.NoError(t, json.Unmarshal(respMsg.Body.Payload, &response))
			require.Equal(t, functions.RequestStateComplete, response.Status)
		}).Return(nil).Once()
		handler.HandleGatewayMessage(ctx, "gw1", msg)
	})

	t.Run("heartbeat internal error", func(t *testing.T) {
		ctx := testutils.Context(t)
		msg, _ := newOffchainRequest(t, addr.Bytes(), 0)
		require.NoError(t, msg.Sign(privateKey))

		// first call to trigger the request
		var response functions.HeartbeatResponse
		allowlist.On("Allow", addr).Return(true).Once()
		handlerCalled := make(chan struct{})
		listener.On("HandleOffchainRequest", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			handlerCalled <- struct{}{}
		}).Return(errors.New("boom")).Once()
		connector.On("SendToGateway", mock.Anything, "gw1", mock.Anything).Return(nil).Once()
		handler.HandleGatewayMessage(ctx, "gw1", msg)
		<-handlerCalled

		// collect the response - should eventually result in an internal error
		gomega.NewGomegaWithT(t).Eventually(func() bool {
			returnedState := 0
			allowlist.On("Allow", addr).Return(true).Once()
			connector.On("SendToGateway", mock.Anything, "gw1", mock.Anything).Run(func(args mock.Arguments) {
				respMsg, ok := args[2].(*api.Message)
				require.True(t, ok)
				require.NoError(t, json.Unmarshal(respMsg.Body.Payload, &response))
				returnedState = response.Status
			}).Return(nil).Once()
			handler.HandleGatewayMessage(ctx, "gw1", msg)
			return returnedState == functions.RequestStateInternalError
		}, testutils.WaitTimeout(t), 50*time.Millisecond).Should(gomega.BeTrue())
	})

	t.Run("heartbeat sender address doesn't match", func(t *testing.T) {
		ctx := testutils.Context(t)
		msg, _ := newOffchainRequest(t, geth_common.BytesToAddress([]byte("0x1234")).Bytes(), 0)
		require.NoError(t, msg.Sign(privateKey))

		var response functions.HeartbeatResponse
		allowlist.On("Allow", addr).Return(true).Once()
		connector.On("SendToGateway", mock.Anything, "gw1", mock.Anything).Run(func(args mock.Arguments) {
			respMsg, ok := args[2].(*api.Message)
			require.True(t, ok)
			require.NoError(t, json.Unmarshal(respMsg.Body.Payload, &response))
			require.Equal(t, functions.RequestStateInternalError, response.Status)
		}).Return(nil).Once()
		handler.HandleGatewayMessage(ctx, "gw1", msg)
	})

	t.Run("heartbeat request too old", func(t *testing.T) {
		ctx := testutils.Context(t)
		msg, _ := newOffchainRequest(t, addr.Bytes(), 10_000)
		require.NoError(t, msg.Sign(privateKey))

		var response functions.HeartbeatResponse
		allowlist.On("Allow", addr).Return(true).Once()
		connector.On("SendToGateway", mock.Anything, "gw1", mock.Anything).Run(func(args mock.Arguments) {
			respMsg, ok := args[2].(*api.Message)
			require.True(t, ok)
			require.NoError(t, json.Unmarshal(respMsg.Body.Payload, &response))
			require.Equal(t, functions.RequestStateInternalError, response.Status)
		}).Return(nil).Once()
		handler.HandleGatewayMessage(ctx, "gw1", msg)
	})
}
