package vault_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	core_mocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	vaulttypesmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	connector_mocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
)

func TestGatewayHandler_HandleGatewayMessage(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := context.Background()

	tests := []struct {
		name          string
		setupMocks    func(*vaulttypesmocks.SecretsService, *connector_mocks.GatewayConnector)
		request       *jsonrpc.Request[json.RawMessage]
		expectedError bool
	}{
		{
			name: "success - create secrets",
			setupMocks: func(ss *vaulttypesmocks.SecretsService, gc *connector_mocks.GatewayConnector) {
				ss.EXPECT().CreateSecrets(mock.Anything, mock.MatchedBy(func(req *vaultcommon.CreateSecretsRequest) bool {
					return len(req.EncryptedSecrets) == 1 &&
						req.EncryptedSecrets[0].Id.Key == "test-secret"
				})).Return(&vaulttypes.Response{ID: "test-secret"}, nil)

				gc.On("SendToGateway", mock.Anything, "gateway-1", mock.MatchedBy(func(resp *jsonrpc.Response[json.RawMessage]) bool {
					return resp.Error == nil
				})).Return(nil)
			},
			request: &jsonrpc.Request[json.RawMessage]{
				Method: vaulttypes.MethodSecretsCreate,
				ID:     "1",
				Params: func() *json.RawMessage {
					params, _ := json.Marshal(vaultcommon.CreateSecretsRequest{
						RequestId: "test-request-id",
						EncryptedSecrets: []*vaultcommon.EncryptedSecret{
							{
								Id: &vaultcommon.SecretIdentifier{
									Key: "test-secret",
								},
								EncryptedValue: "encrypted-value",
							},
						},
					})
					raw := json.RawMessage(params)
					return &raw
				}(),
			},
			expectedError: false,
		},
		{
			name: "failure - service error",
			setupMocks: func(ss *vaulttypesmocks.SecretsService, gc *connector_mocks.GatewayConnector) {
				ss.EXPECT().CreateSecrets(mock.Anything, mock.Anything).Return(nil, errors.New("service error"))

				gc.On("SendToGateway", mock.Anything, "gateway-1", mock.MatchedBy(func(resp *jsonrpc.Response[json.RawMessage]) bool {
					return resp.Error != nil &&
						resp.Error.Code == api.ToJSONRPCErrorCode(api.FatalError)
				})).Return(nil)
			},
			request: &jsonrpc.Request[json.RawMessage]{
				Method: vaulttypes.MethodSecretsCreate,
				ID:     "1",
				Params: func() *json.RawMessage {
					params, _ := json.Marshal(vaultcommon.CreateSecretsRequest{
						RequestId: "test-request-id",
						EncryptedSecrets: []*vaultcommon.EncryptedSecret{
							{
								Id: &vaultcommon.SecretIdentifier{
									Key: "test-secret",
								},
								EncryptedValue: "encrypted-value",
							},
						},
					})
					raw := json.RawMessage(params)
					return &raw
				}(),
			},
			expectedError: false,
		},
		{
			name: "failure - invalid method",
			setupMocks: func(ss *vaulttypesmocks.SecretsService, gc *connector_mocks.GatewayConnector) {
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
			setupMocks: func(ss *vaulttypesmocks.SecretsService, gc *connector_mocks.GatewayConnector) {
				gc.On("SendToGateway", mock.Anything, "gateway-1", mock.MatchedBy(func(resp *jsonrpc.Response[json.RawMessage]) bool {
					return resp.Error != nil &&
						resp.Error.Code == api.ToJSONRPCErrorCode(api.UserMessageParseError)
				})).Return(nil)
			},
			request: &jsonrpc.Request[json.RawMessage]{
				Method: vaulttypes.MethodSecretsCreate,
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
			setupMocks: func(ss *vaulttypesmocks.SecretsService, gc *connector_mocks.GatewayConnector) {
				ss.EXPECT().DeleteSecrets(mock.Anything, mock.MatchedBy(func(req *vaultcommon.DeleteSecretsRequest) bool {
					return len(req.Ids) == 1 &&
						req.Ids[0].Key == "Foo" &&
						req.Ids[0].Namespace == "Bar" &&
						req.Ids[0].Owner == "Owner"
				})).Return(&vaulttypes.Response{ID: "test-secret"}, nil)

				gc.On("SendToGateway", mock.Anything, "gateway-1", mock.MatchedBy(func(resp *jsonrpc.Response[json.RawMessage]) bool {
					return resp.Error == nil
				})).Return(nil)
			},
			request: &jsonrpc.Request[json.RawMessage]{
				Method: vaulttypes.MethodSecretsDelete,
				ID:     "1",
				Params: func() *json.RawMessage {
					params, _ := json.Marshal(vaultcommon.DeleteSecretsRequest{
						RequestId: "test-secret",
						Ids: []*vaultcommon.SecretIdentifier{
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
			secretsService := vaulttypesmocks.NewSecretsService(t)
			gwConnector := connector_mocks.NewGatewayConnector(t)
			capRegistry := core_mocks.NewCapabilitiesRegistry(t)

			tt.setupMocks(secretsService, gwConnector)

			handler, err := vaultcap.NewGatewayHandler(capRegistry, secretsService, gwConnector, lggr)
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

	secretsService := vaulttypesmocks.NewSecretsService(t)
	gwConnector := connector_mocks.NewGatewayConnector(t)
	capRegistry := core_mocks.NewCapabilitiesRegistry(t)

	handler, err := vaultcap.NewGatewayHandler(capRegistry, secretsService, gwConnector, lggr)
	require.NoError(t, err)

	t.Run("start", func(t *testing.T) {
		gwConnector.On("AddHandler", mock.Anything, vaulttypes.GetSupportedMethods(lggr), handler).Return(nil).Once()
		err := handler.Start(ctx)
		require.NoError(t, err)
	})

	t.Run("close", func(t *testing.T) {
		gwConnector.On("RemoveHandler", mock.Anything, vaulttypes.GetSupportedMethods(lggr)).Return(nil).Once()
		err := handler.Close()
		require.NoError(t, err)
	})

	t.Run("id", func(t *testing.T) {
		id, err := handler.ID(ctx)
		require.NoError(t, err)
		assert.Equal(t, vaultcap.HandlerName, id)
	})
}
