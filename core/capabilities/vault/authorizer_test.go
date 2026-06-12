package vault_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	vault "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	vaultmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestAuthorizer_RejectsJWTBasedAuthWhenUnavailable(t *testing.T) {
	params, err := json.Marshal(vaultcommon.CreateSecretsRequest{})
	require.NoError(t, err)

	allowListBasedAuth := vaultmocks.NewAuthorizer(t)
	allowListBasedAuth.EXPECT().AuthorizeRequest(mock.Anything, mock.Anything).Maybe()

	a := vault.NewAuthorizer(allowListBasedAuth, nil, logger.TestLogger(t))

	authResult, err := a.AuthorizeRequest(t.Context(), jsonrpc.Request[json.RawMessage]{
		ID:     "1",
		Method: vaulttypes.MethodSecretsCreate,
		Params: (*json.RawMessage)(&params),
		Auth:   "jwt-token",
	})
	require.Nil(t, authResult)
	require.ErrorContains(t, err, "JWTBasedAuth is nil")
	allowListBasedAuth.AssertNotCalled(t, "AuthorizeRequest", mock.Anything, mock.Anything)
}

func TestAuthorizer_UsesJWTWhenGateEnabled(t *testing.T) {
	params, err := json.Marshal(vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{Id: &vaultcommon.SecretIdentifier{Owner: "0xworkflow", Namespace: "ns", Key: "k"}, EncryptedValue: "cipher"},
		},
	})
	require.NoError(t, err)

	req := jsonrpc.Request[json.RawMessage]{
		ID:     "1",
		Method: vaulttypes.MethodSecretsCreate,
		Params: (*json.RawMessage)(&params),
		Auth:   "jwt-token",
	}
	digest, err := req.Digest()
	require.NoError(t, err)

	jwtBasedAuth := vaultmocks.NewAuthorizer(t)
	jwtBasedAuth.EXPECT().AuthorizeRequest(mock.Anything, req).Return(vault.NewAuthResult("org-1", "0xworkflow", digest, time.Now().Add(time.Minute).Unix()), nil).Once()

	a := vault.NewAuthorizer(nil, jwtBasedAuth, logger.TestLogger(t))

	authResult, err := a.AuthorizeRequest(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, "org-1", authResult.OrgID())
	require.Equal(t, "0xworkflow", authResult.WorkflowOwner())
	require.Equal(t, "0xworkflow", authResult.AuthorizedOwner())
}

func TestAuthorizer_DelegatesDigestVerificationToJWTAuth(t *testing.T) {
	req := jsonrpc.Request[json.RawMessage]{
		ID:     "1",
		Method: vaulttypes.MethodPublicKeyGet,
		Auth:   "jwt-token",
	}

	jwtBasedAuth := vaultmocks.NewAuthorizer(t)
	jwtBasedAuth.EXPECT().AuthorizeRequest(mock.Anything, req).Return(vault.NewAuthResult("org-1", "", "wrong-digest", time.Now().Add(time.Minute).Unix()), nil).Once()

	a := vault.NewAuthorizer(nil, jwtBasedAuth, logger.TestLogger(t))

	authResult, err := a.AuthorizeRequest(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, "org-1", authResult.OrgID())
	require.Empty(t, authResult.WorkflowOwner())
	require.Empty(t, authResult.AuthorizedOwner())
}

func TestAuthorizer_RejectsJWTReplay(t *testing.T) {
	req := jsonrpc.Request[json.RawMessage]{
		ID:     "1",
		Method: vaulttypes.MethodPublicKeyGet,
		Auth:   "jwt-token",
	}
	digest, err := req.Digest()
	require.NoError(t, err)

	jwtBasedAuth := vaultmocks.NewAuthorizer(t)
	jwtBasedAuth.EXPECT().AuthorizeRequest(mock.Anything, req).Return(vault.NewAuthResult("org-1", "", digest, time.Now().Add(time.Minute).Unix()), nil).Twice()

	a := vault.NewAuthorizer(nil, jwtBasedAuth, logger.TestLogger(t))

	authResult, err := a.AuthorizeRequest(t.Context(), req)
	require.NoError(t, err)
	require.Empty(t, authResult.AuthorizedOwner())

	authResult, err = a.AuthorizeRequest(t.Context(), req)
	require.Nil(t, authResult)
	require.ErrorIs(t, err, vault.ErrRequestAlreadySeen)
}

func TestAuthorizer_RejectsAllowListBasedAuthReplay(t *testing.T) {
	allowListBasedAuth := vaultmocks.NewAuthorizer(t)
	// Use a method without secret identifiers so the owner-binding check is a no-op.
	req := jsonrpc.Request[json.RawMessage]{ID: "1", Method: vaulttypes.MethodPublicKeyGet}
	allowListBasedAuth.EXPECT().AuthorizeRequest(mock.Anything, req).Return(vault.NewAuthResult("", "0xabc", "digest-1", time.Now().Add(time.Minute).Unix()), nil).Twice()

	a := vault.NewAuthorizer(allowListBasedAuth, nil, logger.TestLogger(t))

	authResult, err := a.AuthorizeRequest(t.Context(), req)
	require.NoError(t, err)
	require.Empty(t, authResult.OrgID())
	require.Equal(t, "0xabc", authResult.WorkflowOwner())
	require.Equal(t, "0xabc", authResult.AuthorizedOwner())

	authResult, err = a.AuthorizeRequest(t.Context(), req)
	require.Nil(t, authResult)
	require.ErrorIs(t, err, vault.ErrRequestAlreadySeen)
}

func TestAuthorizer_PropagatesJWTValidationErrors(t *testing.T) {
	// JWT mock fails before owner binding; params are irrelevant here.
	req := jsonrpc.Request[json.RawMessage]{
		ID:     "1",
		Method: vaulttypes.MethodSecretsCreate,
		Auth:   "jwt-token",
	}

	jwtBasedAuth := vaultmocks.NewAuthorizer(t)
	jwtBasedAuth.EXPECT().AuthorizeRequest(mock.Anything, req).Return(nil, errors.New("bad token")).Once()

	a := vault.NewAuthorizer(nil, jwtBasedAuth, logger.TestLogger(t))

	authResult, err := a.AuthorizeRequest(t.Context(), req)
	require.Nil(t, authResult)
	require.ErrorContains(t, err, "bad token")
}

func TestAuthorizer_AllowListPath_RejectsCreateOwnerMismatch(t *testing.T) {
	params, err := json.Marshal(vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{Id: &vaultcommon.SecretIdentifier{Owner: "0xother", Namespace: "ns", Key: "k"}, EncryptedValue: "cipher"},
		},
	})
	require.NoError(t, err)

	req := jsonrpc.Request[json.RawMessage]{
		ID:     "1",
		Method: vaulttypes.MethodSecretsCreate,
		Params: (*json.RawMessage)(&params),
	}

	allowListBasedAuth := vaultmocks.NewAuthorizer(t)
	allowListBasedAuth.EXPECT().AuthorizeRequest(mock.Anything, req).Return(vault.NewAuthResult("", "0xauthorized", "digest-1", time.Now().Add(time.Minute).Unix()), nil).Once()

	a := vault.NewAuthorizer(allowListBasedAuth, nil, logger.TestLogger(t))

	authResult, err := a.AuthorizeRequest(t.Context(), req)
	require.Nil(t, authResult)
	require.ErrorContains(t, err, "encrypted secret owner at index 0 \"0xother\" does not match authorized workflow owner \"0xauthorized\"")
}

func TestAuthorizer_AllowListPath_RejectsUpdateOwnerMismatch(t *testing.T) {
	params, err := json.Marshal(vaultcommon.UpdateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{Id: &vaultcommon.SecretIdentifier{Owner: "0xother", Namespace: "ns", Key: "k"}, EncryptedValue: "cipher"},
		},
	})
	require.NoError(t, err)

	req := jsonrpc.Request[json.RawMessage]{
		ID:     "1",
		Method: vaulttypes.MethodSecretsUpdate,
		Params: (*json.RawMessage)(&params),
	}

	allowListBasedAuth := vaultmocks.NewAuthorizer(t)
	allowListBasedAuth.EXPECT().AuthorizeRequest(mock.Anything, req).Return(vault.NewAuthResult("", "0xauthorized", "digest-1", time.Now().Add(time.Minute).Unix()), nil).Once()

	a := vault.NewAuthorizer(allowListBasedAuth, nil, logger.TestLogger(t))

	authResult, err := a.AuthorizeRequest(t.Context(), req)
	require.Nil(t, authResult)
	require.ErrorContains(t, err, "encrypted secret owner at index 0 \"0xother\" does not match authorized workflow owner \"0xauthorized\"")
}

func TestAuthorizer_AllowListPath_RejectsDeleteOwnerMismatch(t *testing.T) {
	params, err := json.Marshal(vaultcommon.DeleteSecretsRequest{
		Ids: []*vaultcommon.SecretIdentifier{
			{Owner: "0xother", Namespace: "ns", Key: "k"},
		},
	})
	require.NoError(t, err)

	req := jsonrpc.Request[json.RawMessage]{
		ID:     "1",
		Method: vaulttypes.MethodSecretsDelete,
		Params: (*json.RawMessage)(&params),
	}

	allowListBasedAuth := vaultmocks.NewAuthorizer(t)
	allowListBasedAuth.EXPECT().AuthorizeRequest(mock.Anything, req).Return(vault.NewAuthResult("", "0xauthorized", "digest-1", time.Now().Add(time.Minute).Unix()), nil).Once()

	a := vault.NewAuthorizer(allowListBasedAuth, nil, logger.TestLogger(t))

	authResult, err := a.AuthorizeRequest(t.Context(), req)
	require.Nil(t, authResult)
	require.ErrorContains(t, err, "secret identifier owner at index 0 \"0xother\" does not match authorized workflow owner \"0xauthorized\"")
}

func TestAuthorizer_AllowListPath_RejectsListOwnerMismatch(t *testing.T) {
	params, err := json.Marshal(vaultcommon.ListSecretIdentifiersRequest{
		Owner:     "0xother",
		Namespace: "ns",
	})
	require.NoError(t, err)

	req := jsonrpc.Request[json.RawMessage]{
		ID:     "1",
		Method: vaulttypes.MethodSecretsList,
		Params: (*json.RawMessage)(&params),
	}

	allowListBasedAuth := vaultmocks.NewAuthorizer(t)
	allowListBasedAuth.EXPECT().AuthorizeRequest(mock.Anything, req).Return(vault.NewAuthResult("", "0xauthorized", "digest-1", time.Now().Add(time.Minute).Unix()), nil).Once()

	a := vault.NewAuthorizer(allowListBasedAuth, nil, logger.TestLogger(t))

	authResult, err := a.AuthorizeRequest(t.Context(), req)
	require.Nil(t, authResult)
	require.ErrorContains(t, err, "list secrets owner \"0xother\" does not match authorized workflow owner \"0xauthorized\"")
}

func TestAuthorizer_JWTPath_RejectsOwnerMismatch(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		buildParams func(mismatchedOwner string) json.RawMessage
		errContains string
	}{
		{
			name:   "create",
			method: vaulttypes.MethodSecretsCreate,
			buildParams: func(mismatchedOwner string) json.RawMessage {
				params, err := json.Marshal(vaultcommon.CreateSecretsRequest{
					EncryptedSecrets: []*vaultcommon.EncryptedSecret{
						{Id: &vaultcommon.SecretIdentifier{Owner: mismatchedOwner, Namespace: "ns", Key: "k"}, EncryptedValue: "cipher"},
					},
				})
				require.NoError(t, err)
				return params
			},
			errContains: "encrypted secret owner at index 0",
		},
		{
			name:   "update",
			method: vaulttypes.MethodSecretsUpdate,
			buildParams: func(mismatchedOwner string) json.RawMessage {
				params, err := json.Marshal(vaultcommon.UpdateSecretsRequest{
					EncryptedSecrets: []*vaultcommon.EncryptedSecret{
						{Id: &vaultcommon.SecretIdentifier{Owner: mismatchedOwner, Namespace: "ns", Key: "k"}, EncryptedValue: "cipher"},
					},
				})
				require.NoError(t, err)
				return params
			},
			errContains: "encrypted secret owner at index 0",
		},
		{
			name:   "delete",
			method: vaulttypes.MethodSecretsDelete,
			buildParams: func(mismatchedOwner string) json.RawMessage {
				params, err := json.Marshal(vaultcommon.DeleteSecretsRequest{
					Ids: []*vaultcommon.SecretIdentifier{
						{Owner: mismatchedOwner, Namespace: "ns", Key: "k"},
					},
				})
				require.NoError(t, err)
				return params
			},
			errContains: "secret identifier owner at index 0",
		},
		{
			name:   "list",
			method: vaulttypes.MethodSecretsList,
			buildParams: func(mismatchedOwner string) json.RawMessage {
				params, err := json.Marshal(vaultcommon.ListSecretIdentifiersRequest{
					Owner:     mismatchedOwner,
					Namespace: "ns",
				})
				require.NoError(t, err)
				return params
			},
			errContains: "list secrets owner",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mismatchedOwner := "0xother"
			authorizedOwner := "0xauthorized"
			params := tt.buildParams(mismatchedOwner)

			req := jsonrpc.Request[json.RawMessage]{
				ID:     "1",
				Method: tt.method,
				Params: &params,
				Auth:   "jwt-token",
			}

			jwtBasedAuth := vaultmocks.NewAuthorizer(t)
			jwtBasedAuth.EXPECT().AuthorizeRequest(mock.Anything, req).Return(vault.NewAuthResult("org-1", authorizedOwner, "digest-1", time.Now().Add(time.Minute).Unix()), nil).Once()

			a := vault.NewAuthorizer(nil, jwtBasedAuth, logger.TestLogger(t))

			authResult, err := a.AuthorizeRequest(t.Context(), req)
			require.Nil(t, authResult)
			require.ErrorContains(t, err, tt.errContains)
			require.ErrorContains(t, err, authorizedOwner)
		})
	}
}

func TestAuthorizer_RejectsOwnerBindingOnMalformedBatches(t *testing.T) {
	authorizedOwner := "0xauthorized"
	tests := []struct {
		name        string
		method      string
		buildParams func() json.RawMessage
		errContains string
	}{
		{
			name:   "create empty batch",
			method: vaulttypes.MethodSecretsCreate,
			buildParams: func() json.RawMessage {
				params, err := json.Marshal(vaultcommon.CreateSecretsRequest{EncryptedSecrets: []*vaultcommon.EncryptedSecret{}})
				require.NoError(t, err)
				return params
			},
			errContains: "request batch must contain at least 1 item",
		},
		{
			name:   "create nil secret id",
			method: vaulttypes.MethodSecretsCreate,
			buildParams: func() json.RawMessage {
				params, err := json.Marshal(vaultcommon.CreateSecretsRequest{
					EncryptedSecrets: []*vaultcommon.EncryptedSecret{
						{Id: nil, EncryptedValue: "ab"},
					},
				})
				require.NoError(t, err)
				return params
			},
			errContains: "secret ID must not be nil at index 0",
		},
		{
			name:   "update nil encrypted secret",
			method: vaulttypes.MethodSecretsUpdate,
			buildParams: func() json.RawMessage {
				params, err := json.Marshal(vaultcommon.UpdateSecretsRequest{
					EncryptedSecrets: []*vaultcommon.EncryptedSecret{nil},
				})
				require.NoError(t, err)
				return params
			},
			errContains: "encrypted secret must not be nil at index 0",
		},
		{
			name:   "delete nil secret identifier",
			method: vaulttypes.MethodSecretsDelete,
			buildParams: func() json.RawMessage {
				params, err := json.Marshal(vaultcommon.DeleteSecretsRequest{
					Ids: []*vaultcommon.SecretIdentifier{nil},
				})
				require.NoError(t, err)
				return params
			},
			errContains: "secret ID must not be nil at index 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := tt.buildParams()
			req := jsonrpc.Request[json.RawMessage]{
				ID:     "1",
				Method: tt.method,
				Params: &params,
			}

			allowListBasedAuth := vaultmocks.NewAuthorizer(t)
			allowListBasedAuth.EXPECT().AuthorizeRequest(mock.Anything, req).Return(vault.NewAuthResult("", authorizedOwner, "digest-1", time.Now().Add(time.Minute).Unix()), nil).Once()

			a := vault.NewAuthorizer(allowListBasedAuth, nil, logger.TestLogger(t))

			authResult, err := a.AuthorizeRequest(t.Context(), req)
			require.Nil(t, authResult)
			require.ErrorContains(t, err, tt.errContains)
		})
	}
}

func TestAuthorizer_RejectsOwnerBindingWhenParamsMissing(t *testing.T) {
	allowListBasedAuth := vaultmocks.NewAuthorizer(t)
	allowListBasedAuth.EXPECT().AuthorizeRequest(mock.Anything, mock.Anything).Return(vault.NewAuthResult("", "0xauthorized", "digest-1", time.Now().Add(time.Minute).Unix()), nil).Once()

	a := vault.NewAuthorizer(allowListBasedAuth, nil, logger.TestLogger(t))

	authResult, err := a.AuthorizeRequest(t.Context(), jsonrpc.Request[json.RawMessage]{
		ID:     "1",
		Method: vaulttypes.MethodSecretsCreate,
	})
	require.Nil(t, authResult)
	require.ErrorContains(t, err, "request params must not be nil")
}
