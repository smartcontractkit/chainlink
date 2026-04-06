package vault_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	vaultcapmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	vaulttypesmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	connector_mocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
)

func TestGatewayHandler_HandleGatewayMessage(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := t.Context()
	authResult := func(orgID, workflowOwner string) *vaultcap.AuthResult {
		digestOwner := workflowOwner
		if orgID != "" {
			digestOwner = orgID
		}
		return vaultcap.NewAuthResult(orgID, workflowOwner, "digest-"+digestOwner, time.Now().Add(time.Minute).Unix())
	}

	tests := []struct {
		name          string
		setupMocks    func(*vaulttypesmocks.SecretsService, *connector_mocks.GatewayConnector, *vaultcapmocks.Authorizer)
		request       *jsonrpc.Request[json.RawMessage]
		expectedError bool
	}{
		{
			name: "success - create secrets",
			setupMocks: func(ss *vaulttypesmocks.SecretsService, gc *connector_mocks.GatewayConnector, ra *vaultcapmocks.Authorizer) {
				ra.EXPECT().AuthorizeRequest(mock.Anything, mock.MatchedBy(func(req jsonrpc.Request[json.RawMessage]) bool {
					return req.Method == vaulttypes.MethodSecretsCreate && req.ID == "1"
				})).Return(authResult("", "0xabc"), nil)
				ss.EXPECT().CreateSecrets(mock.Anything, mock.MatchedBy(func(req *vaultcommon.CreateSecretsRequest) bool {
					return len(req.EncryptedSecrets) == 1 &&
						req.EncryptedSecrets[0].Id.Key == "test-secret" &&
						req.EncryptedSecrets[0].Id.Owner == "0xAbC" &&
						req.RequestId == "0xabc"+vaulttypes.RequestIDSeparator+"1" &&
						req.OrgId == "" &&
						req.WorkflowOwner == "0xabc"
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
									Key:   "test-secret",
									Owner: "0xAbC",
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
			name: "success - create secrets propagates jwt auth identity",
			setupMocks: func(ss *vaulttypesmocks.SecretsService, gc *connector_mocks.GatewayConnector, ra *vaultcapmocks.Authorizer) {
				ra.EXPECT().AuthorizeRequest(mock.Anything, mock.MatchedBy(func(req jsonrpc.Request[json.RawMessage]) bool {
					return req.Method == vaulttypes.MethodSecretsCreate && req.ID == "1"
				})).Return(authResult("org-1", "0xworkflow"), nil)
				ss.EXPECT().CreateSecrets(mock.Anything, mock.MatchedBy(func(req *vaultcommon.CreateSecretsRequest) bool {
					return len(req.EncryptedSecrets) == 1 &&
						req.EncryptedSecrets[0].Id.Key == "test-secret" &&
						req.EncryptedSecrets[0].Id.Owner == "org-1" &&
						req.RequestId == "org-1"+vaulttypes.RequestIDSeparator+"1" &&
						req.OrgId == "org-1" &&
						req.WorkflowOwner == "0xworkflow"
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
						EncryptedSecrets: []*vaultcommon.EncryptedSecret{
							{
								Id: &vaultcommon.SecretIdentifier{
									Key:   "test-secret",
									Owner: "org-1",
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
			setupMocks: func(ss *vaulttypesmocks.SecretsService, gc *connector_mocks.GatewayConnector, ra *vaultcapmocks.Authorizer) {
				ra.EXPECT().AuthorizeRequest(mock.Anything, mock.Anything).Return(authResult("", "0xabc"), nil)
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
									Key:   "test-secret",
									Owner: "0xAbC",
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
			setupMocks: func(ss *vaulttypesmocks.SecretsService, gc *connector_mocks.GatewayConnector, ra *vaultcapmocks.Authorizer) {
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
			setupMocks: func(ss *vaulttypesmocks.SecretsService, gc *connector_mocks.GatewayConnector, ra *vaultcapmocks.Authorizer) {
				gc.On("SendToGateway", mock.Anything, "gateway-1", mock.MatchedBy(func(resp *jsonrpc.Response[json.RawMessage]) bool {
					return resp.Error != nil &&
						resp.Error.Code == api.ToJSONRPCErrorCode(api.HandlerError)
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
			setupMocks: func(ss *vaulttypesmocks.SecretsService, gc *connector_mocks.GatewayConnector, ra *vaultcapmocks.Authorizer) {
				ra.EXPECT().AuthorizeRequest(mock.Anything, mock.MatchedBy(func(req jsonrpc.Request[json.RawMessage]) bool {
					return req.Method == vaulttypes.MethodSecretsDelete && req.ID == "1"
				})).Return(authResult("", "0xabc"), nil)
				ss.EXPECT().DeleteSecrets(mock.Anything, mock.MatchedBy(func(req *vaultcommon.DeleteSecretsRequest) bool {
					return len(req.Ids) == 1 &&
						req.Ids[0].Key == "Foo" &&
						req.Ids[0].Namespace == "Bar" &&
						req.Ids[0].Owner == "0xAbC" &&
						req.RequestId == "0xabc"+vaulttypes.RequestIDSeparator+"1" &&
						req.OrgId == "" &&
						req.WorkflowOwner == "0xabc"
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
								Owner:     "0xAbC",
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
			name: "success - update secrets propagates jwt auth identity",
			setupMocks: func(ss *vaulttypesmocks.SecretsService, gc *connector_mocks.GatewayConnector, ra *vaultcapmocks.Authorizer) {
				ra.EXPECT().AuthorizeRequest(mock.Anything, mock.MatchedBy(func(req jsonrpc.Request[json.RawMessage]) bool {
					return req.Method == vaulttypes.MethodSecretsUpdate && req.ID == "1"
				})).Return(authResult("org-1", "0xworkflow"), nil)
				ss.EXPECT().UpdateSecrets(mock.Anything, mock.MatchedBy(func(req *vaultcommon.UpdateSecretsRequest) bool {
					return len(req.EncryptedSecrets) == 1 &&
						req.EncryptedSecrets[0].Id.Key == "updated-secret" &&
						req.EncryptedSecrets[0].Id.Owner == "org-1" &&
						req.RequestId == "org-1"+vaulttypes.RequestIDSeparator+"1" &&
						req.OrgId == "org-1" &&
						req.WorkflowOwner == "0xworkflow"
				})).Return(&vaulttypes.Response{ID: "updated-secret"}, nil)

				gc.On("SendToGateway", mock.Anything, "gateway-1", mock.MatchedBy(func(resp *jsonrpc.Response[json.RawMessage]) bool {
					return resp.Error == nil
				})).Return(nil)
			},
			request: &jsonrpc.Request[json.RawMessage]{
				Method: vaulttypes.MethodSecretsUpdate,
				ID:     "1",
				Params: func() *json.RawMessage {
					params, _ := json.Marshal(vaultcommon.UpdateSecretsRequest{
						EncryptedSecrets: []*vaultcommon.EncryptedSecret{
							{
								Id: &vaultcommon.SecretIdentifier{
									Key:   "updated-secret",
									Owner: "org-1",
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
			name: "success - delete secrets propagates jwt auth identity",
			setupMocks: func(ss *vaulttypesmocks.SecretsService, gc *connector_mocks.GatewayConnector, ra *vaultcapmocks.Authorizer) {
				ra.EXPECT().AuthorizeRequest(mock.Anything, mock.MatchedBy(func(req jsonrpc.Request[json.RawMessage]) bool {
					return req.Method == vaulttypes.MethodSecretsDelete && req.ID == "1"
				})).Return(authResult("org-1", "0xworkflow"), nil)
				ss.EXPECT().DeleteSecrets(mock.Anything, mock.MatchedBy(func(req *vaultcommon.DeleteSecretsRequest) bool {
					return len(req.Ids) == 1 &&
						req.Ids[0].Key == "Foo" &&
						req.Ids[0].Namespace == "Bar" &&
						req.Ids[0].Owner == "org-1" &&
						req.RequestId == "org-1"+vaulttypes.RequestIDSeparator+"1" &&
						req.OrgId == "org-1" &&
						req.WorkflowOwner == "0xworkflow"
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
						Ids: []*vaultcommon.SecretIdentifier{
							{
								Key:       "Foo",
								Namespace: "Bar",
								Owner:     "org-1",
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
			name: "success - list secrets propagates jwt auth identity",
			setupMocks: func(ss *vaulttypesmocks.SecretsService, gc *connector_mocks.GatewayConnector, ra *vaultcapmocks.Authorizer) {
				ra.EXPECT().AuthorizeRequest(mock.Anything, mock.MatchedBy(func(req jsonrpc.Request[json.RawMessage]) bool {
					return req.Method == vaulttypes.MethodSecretsList && req.ID == "1"
				})).Return(authResult("org-1", "0xworkflow"), nil)
				ss.EXPECT().ListSecretIdentifiers(mock.Anything, mock.MatchedBy(func(req *vaultcommon.ListSecretIdentifiersRequest) bool {
					return req.RequestId == "org-1"+vaulttypes.RequestIDSeparator+"1" &&
						req.Owner == "org-1" &&
						req.Namespace == "ns" &&
						req.OrgId == "org-1" &&
						req.WorkflowOwner == "0xworkflow"
				})).Return(&vaulttypes.Response{ID: "test-secret"}, nil)

				gc.On("SendToGateway", mock.Anything, "gateway-1", mock.MatchedBy(func(resp *jsonrpc.Response[json.RawMessage]) bool {
					return resp.Error == nil
				})).Return(nil)
			},
			request: &jsonrpc.Request[json.RawMessage]{
				Method: vaulttypes.MethodSecretsList,
				ID:     "1",
				Params: func() *json.RawMessage {
					params, _ := json.Marshal(vaultcommon.ListSecretIdentifiersRequest{
						Owner:     "user-supplied-owner",
						Namespace: "ns",
					})
					raw := json.RawMessage(params)
					return &raw
				}(),
			},
			expectedError: false,
		},
		{
			name: "failure - unauthorized request",
			setupMocks: func(ss *vaulttypesmocks.SecretsService, gc *connector_mocks.GatewayConnector, ra *vaultcapmocks.Authorizer) {
				ra.EXPECT().AuthorizeRequest(mock.Anything, mock.Anything).Return(nil, errors.New("not allowlisted"))
				gc.On("SendToGateway", mock.Anything, "gateway-1", mock.MatchedBy(func(resp *jsonrpc.Response[json.RawMessage]) bool {
					return resp.Error != nil &&
						resp.Error.Code == api.ToJSONRPCErrorCode(api.HandlerError) &&
						resp.Error.Message == "request not authorized: not allowlisted"
				})).Return(nil)
			},
			request: &jsonrpc.Request[json.RawMessage]{
				Method: vaulttypes.MethodSecretsCreate,
				ID:     "1",
				Params: func() *json.RawMessage {
					params, _ := json.Marshal(vaultcommon.CreateSecretsRequest{
						EncryptedSecrets: []*vaultcommon.EncryptedSecret{
							{
								Id: &vaultcommon.SecretIdentifier{
									Key:   "test-secret",
									Owner: "0xAbC",
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
			name: "success - strips owner prefix from forwarded request before authorization",
			setupMocks: func(ss *vaulttypesmocks.SecretsService, gc *connector_mocks.GatewayConnector, ra *vaultcapmocks.Authorizer) {
				ra.EXPECT().AuthorizeRequest(mock.Anything, mock.MatchedBy(func(req jsonrpc.Request[json.RawMessage]) bool {
					if req.Method != vaulttypes.MethodSecretsCreate || req.ID != "1" || req.Params == nil {
						return false
					}

					var parsed vaultcommon.CreateSecretsRequest
					if err := json.Unmarshal(*req.Params, &parsed); err != nil {
						return false
					}

					return parsed.RequestId == "1" &&
						len(parsed.EncryptedSecrets) == 1 &&
						parsed.EncryptedSecrets[0].Id != nil &&
						parsed.EncryptedSecrets[0].Id.Owner == "0xAbC"
				})).Return(authResult("", "0xabc"), nil)
				ss.EXPECT().CreateSecrets(mock.Anything, mock.MatchedBy(func(req *vaultcommon.CreateSecretsRequest) bool {
					return req.RequestId == "0xabc"+vaulttypes.RequestIDSeparator+"1" &&
						req.OrgId == "" &&
						req.WorkflowOwner == "0xabc"
				})).Return(&vaulttypes.Response{ID: "test-secret"}, nil)

				gc.On("SendToGateway", mock.Anything, "gateway-1", mock.MatchedBy(func(resp *jsonrpc.Response[json.RawMessage]) bool {
					return resp.Error == nil
				})).Return(nil)
			},
			request: &jsonrpc.Request[json.RawMessage]{
				Method: vaulttypes.MethodSecretsCreate,
				ID:     "0xAbC" + vaulttypes.RequestIDSeparator + "1",
				Params: func() *json.RawMessage {
					params, _ := json.Marshal(vaultcommon.CreateSecretsRequest{
						RequestId: "0xAbC" + vaulttypes.RequestIDSeparator + "1",
						EncryptedSecrets: []*vaultcommon.EncryptedSecret{
							{
								Id: &vaultcommon.SecretIdentifier{
									Key:   "test-secret",
									Owner: "0xAbC",
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
			name: "failure - owner mismatch against authorized owner",
			setupMocks: func(ss *vaulttypesmocks.SecretsService, gc *connector_mocks.GatewayConnector, ra *vaultcapmocks.Authorizer) {
				ra.EXPECT().AuthorizeRequest(mock.Anything, mock.Anything).Return(authResult("", "0xdef"), nil)
				gc.On("SendToGateway", mock.Anything, "gateway-1", mock.MatchedBy(func(resp *jsonrpc.Response[json.RawMessage]) bool {
					return resp.Error != nil &&
						resp.Error.Code == api.ToJSONRPCErrorCode(api.FatalError) &&
						resp.Error.Message == `secret ID owner "0xabc" does not match authorized owner "0xdef" at index 0`
				})).Return(nil)
			},
			request: &jsonrpc.Request[json.RawMessage]{
				Method: vaulttypes.MethodSecretsCreate,
				ID:     "1",
				Params: func() *json.RawMessage {
					params, _ := json.Marshal(vaultcommon.CreateSecretsRequest{
						EncryptedSecrets: []*vaultcommon.EncryptedSecret{
							{
								Id: &vaultcommon.SecretIdentifier{
									Key:   "test-secret",
									Owner: "0xabc",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secretsService := vaulttypesmocks.NewSecretsService(t)
			gwConnector := connector_mocks.NewGatewayConnector(t)
			allowListBasedAuth := vaultcapmocks.NewAuthorizer(t)
			limitsFactory := limits.Factory{Settings: cresettings.DefaultGetter}
			jwtBasedAuth, err := vaultcap.NewJWTBasedAuth(vaultcap.JWTBasedAuthConfig{}, limitsFactory, lggr, vaultcap.WithDisabledJWTBasedAuth())
			require.NoError(t, err)

			tt.setupMocks(secretsService, gwConnector, allowListBasedAuth)

			handler, err := vaultcap.NewGatewayHandler(
				secretsService,
				gwConnector,
				nil,
				lggr,
				limitsFactory,
				vaultcap.WithAuthorizer(vaultcap.NewAuthorizer(allowListBasedAuth, jwtBasedAuth, lggr)),
			)
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
	ctx := t.Context()

	secretsService := vaulttypesmocks.NewSecretsService(t)
	gwConnector := connector_mocks.NewGatewayConnector(t)
	allowListBasedAuth := vaultcapmocks.NewAuthorizer(t)
	limitsFactory := limits.Factory{Settings: cresettings.DefaultGetter}
	jwtBasedAuth, err := vaultcap.NewJWTBasedAuth(vaultcap.JWTBasedAuthConfig{}, limitsFactory, lggr, vaultcap.WithDisabledJWTBasedAuth())
	require.NoError(t, err)

	handler, err := vaultcap.NewGatewayHandler(
		secretsService,
		gwConnector,
		nil,
		lggr,
		limitsFactory,
		vaultcap.WithAuthorizer(vaultcap.NewAuthorizer(allowListBasedAuth, jwtBasedAuth, lggr)),
	)
	require.NoError(t, err)

	t.Run("start", func(t *testing.T) {
		gwConnector.On("AddHandler", mock.Anything, vaulttypes.Methods, handler).Return(nil).Once()
		err := handler.Start(ctx)
		require.NoError(t, err)
	})

	t.Run("close", func(t *testing.T) {
		gwConnector.On("RemoveHandler", mock.Anything, vaulttypes.Methods).Return(nil).Once()
		err := handler.Close()
		require.NoError(t, err)
	})

	t.Run("id", func(t *testing.T) {
		id, err := handler.ID(ctx)
		require.NoError(t, err)
		assert.Equal(t, vaultcap.HandlerName, id)
	})
}
