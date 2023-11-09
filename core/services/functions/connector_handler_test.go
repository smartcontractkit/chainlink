package functions_test

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	gcmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	gfmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
	s4mocks "github.com/smartcontractkit/chainlink/v2/core/services/s4/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFunctionsConnectorHandler(t *testing.T) {
	t.Parallel()

	logger := logger.TestLogger(t)
	privateKey, addr := testutils.NewPrivateKeyAndAddress(t)
	storage := s4mocks.NewStorage(t)
	connector := gcmocks.NewGatewayConnector(t)
	allowlist := gfmocks.NewOnchainAllowlist(t)
	rateLimiter, err := hc.NewRateLimiter(hc.RateLimiterConfig{GlobalRPS: 100.0, GlobalBurst: 100, PerSenderRPS: 100.0, PerSenderBurst: 100})
	subscriptions := gfmocks.NewOnchainSubscriptions(t)
	require.NoError(t, err)
	allowlist.On("Start", mock.Anything).Return(nil)
	allowlist.On("Close", mock.Anything).Return(nil)
	subscriptions.On("Start", mock.Anything).Return(nil)
	subscriptions.On("Close", mock.Anything).Return(nil)
	handler, err := functions.NewFunctionsConnectorHandler(addr.Hex(), privateKey, storage, allowlist, rateLimiter, subscriptions, *assets.NewLinkFromJuels(100), logger)
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
			subscriptions.On("GetMaxUserBalance", mock.Anything).Return(big.NewInt(100), nil).Once()
			connector.On("SendToGateway", ctx, "gw1", mock.Anything).Run(func(args mock.Arguments) {
				msg, ok := args[2].(*api.Message)
				require.True(t, ok)
				require.Equal(t, `{"success":true,"rows":[{"slot_id":1,"version":1,"expiration":1},{"slot_id":2,"version":2,"expiration":2}]}`, string(msg.Body.Payload))

			}).Return(nil).Once()

			handler.HandleGatewayMessage(ctx, "gw1", &msg)

			t.Run("orm error", func(t *testing.T) {
				storage.On("List", ctx, addr).Return(nil, errors.New("boom")).Once()
				allowlist.On("Allow", addr).Return(true).Once()
				subscriptions.On("GetMaxUserBalance", mock.Anything).Return(big.NewInt(100), nil).Once()
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
			subscriptions.On("GetMaxUserBalance", mock.Anything).Return(big.NewInt(100), nil).Once()
			handler.HandleGatewayMessage(testutils.Context(t), "gw1", &msg)
		})
	})
}
