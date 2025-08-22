package vault_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	core_mocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	vaultCap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	connector_mocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	vault_api "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/vault"
	pluginsvault "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/vault"
)

func TestGatewayHandler_HandleGatewayMessage(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := context.Background()

	tests := []struct {
		name          string
		setupMocks    func(*mocks.SecretsService, *connector_mocks.GatewayConnector)
		request       *jsonrpc.Request[json.RawMessage]
		expectedError bool
	}{
		{
			name: "success - create secrets",
			setupMocks: func(ss *mocks.SecretsService, gc *connector_mocks.GatewayConnector) {
				ss.EXPECT().CreateSecrets(mock.Anything, mock.MatchedBy(func(req *vault.CreateSecretsRequest) bool {
					return len(req.EncryptedSecrets) == 1 &&
						req.EncryptedSecrets[0].Id.Key == "test-secret"
				})).Return(&pluginsvault.Response{ID: "test-secret"}, nil)

				gc.On("SendToGateway", mock.Anything, "gateway-1", mock.MatchedBy(func(resp *jsonrpc.Response[json.RawMessage]) bool {
					return resp.Error == nil
				})).Return(nil)
			},
			request: &jsonrpc.Request[json.RawMessage]{
				Method: vault_api.MethodSecretsCreate,
				ID:     "1",
				Params: func() *json.RawMessage {
					params, _ := json.Marshal(vault_api.SecretsCreateRequest{
						ID:    "test-secret",
						Value: "encrypted-value",
					})
					raw := json.RawMessage(params)
					return &raw
				}(),
			},
			expectedError: false,
		},
		{
			name: "failure - service error",
			setupMocks: func(ss *mocks.SecretsService, gc *connector_mocks.GatewayConnector) {
				ss.EXPECT().CreateSecrets(mock.Anything, mock.Anything).Return(nil, errors.New("service error"))

				gc.On("SendToGateway", mock.Anything, "gateway-1", mock.MatchedBy(func(resp *jsonrpc.Response[json.RawMessage]) bool {
					return resp.Error != nil &&
						resp.Error.Code == api.ToJSONRPCErrorCode(api.FatalError)
				})).Return(nil)
			},
			request: &jsonrpc.Request[json.RawMessage]{
				Method: vault_api.MethodSecretsCreate,
				ID:     "1",
				Params: func() *json.RawMessage {
					params, _ := json.Marshal(vault_api.SecretsCreateRequest{
						ID:    "test-secret",
						Value: "encrypted-value",
					})
					raw := json.RawMessage(params)
					return &raw
				}(),
			},
			expectedError: false,
		},
		{
			name: "failure - invalid method",
			setupMocks: func(ss *mocks.SecretsService, gc *connector_mocks.GatewayConnector) {
				gc.On("SendToGateway", mock.Anything, "gateway-1", mock.MatchedBy(func(resp *jsonrpc.Response[json.RawMessage]) bool {
					return resp.Error != nil &&
						resp.Error.Code == api.ToJSONRPCErrorCode(api.UnsupportedMethodError)
				})).Return(nil)
			},
			request: &jsonrpc.Request[json.RawMessage]{
				Method: "invalid_method",
				ID:     "1",
			},
			expectedError: false,
		},
		{
			name: "failure - invalid request params",
			setupMocks: func(ss *mocks.SecretsService, gc *connector_mocks.GatewayConnector) {
				gc.On("SendToGateway", mock.Anything, "gateway-1", mock.MatchedBy(func(resp *jsonrpc.Response[json.RawMessage]) bool {
					return resp.Error != nil &&
						resp.Error.Code == api.ToJSONRPCErrorCode(api.UserMessageParseError)
				})).Return(nil)
			},
			request: &jsonrpc.Request[json.RawMessage]{
				Method: vault_api.MethodSecretsCreate,
				ID:     "1",
				Params: func() *json.RawMessage {
					raw := json.RawMessage([]byte(`{invalid json`))
					return &raw
				}(),
			},
			expectedError: false,
		},
		{
			name: "success - delete secrets",
			setupMocks: func(ss *mocks.SecretsService, gc *connector_mocks.GatewayConnector) {
				ss.EXPECT().DeleteSecrets(mock.Anything, mock.MatchedBy(func(req *vault.DeleteSecretsRequest) bool {
					return len(req.Ids) == 1 &&
						req.Ids[0].Key == "Foo" &&
						req.Ids[0].Namespace == "Bar" &&
						req.Ids[0].Owner == "Owner"
				})).Return(&pluginsvault.Response{ID: "test-secret"}, nil)

				gc.On("SendToGateway", mock.Anything, "gateway-1", mock.MatchedBy(func(resp *jsonrpc.Response[json.RawMessage]) bool {
					return resp.Error == nil
				})).Return(nil)
			},
			request: &jsonrpc.Request[json.RawMessage]{
				Method: vault_api.MethodSecretsDelete,
				ID:     "1",
				Params: func() *json.RawMessage {
					params, _ := json.Marshal(vault.DeleteSecretsRequest{
						RequestId: "test-secret",
						Ids: []*vault.SecretIdentifier{
							{

								Key:       "Foo",
								Namespace: "Bar",
								Owner:     "Owner",
							},
						},
					})
					raw := json.RawMessage(params)
					return &raw
				}(),
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secretsService := mocks.NewSecretsService(t)
			gwConnector := connector_mocks.NewGatewayConnector(t)
			capRegistry := core_mocks.NewCapabilitiesRegistry(t)

			tt.setupMocks(secretsService, gwConnector)

			handler, err := vaultCap.NewGatewayHandler(capRegistry, secretsService, gwConnector, lggr)
			require.NoError(t, err)

			err = handler.HandleGatewayMessage(ctx, "gateway-1", tt.request)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGatewayHandler_Lifecycle(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := context.Background()

	secretsService := mocks.NewSecretsService(t)
	gwConnector := connector_mocks.NewGatewayConnector(t)
	capRegistry := core_mocks.NewCapabilitiesRegistry(t)

	handler, err := vaultCap.NewGatewayHandler(capRegistry, secretsService, gwConnector, lggr)
	require.NoError(t, err)

	t.Run("start", func(t *testing.T) {
		err := handler.Start(ctx)
		require.NoError(t, err)
	})

	t.Run("close", func(t *testing.T) {
		err := handler.Close()
		require.NoError(t, err)
	})

	t.Run("id", func(t *testing.T) {
		id, err := handler.ID(ctx)
		require.NoError(t, err)
		assert.Equal(t, vaultCap.HandlerName, id)
	})
}
