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

func TestAuthorizer_SkipsOwnerBindingWhenParamsMissing(t *testing.T) {
	allowListBasedAuth := vaultmocks.NewAuthorizer(t)
	allowListBasedAuth.EXPECT().AuthorizeRequest(mock.Anything, mock.Anything).Return(vault.NewAuthResult("", "0xauthorized", "digest-1", time.Now().Add(time.Minute).Unix()), nil).Once()

	a := vault.NewAuthorizer(allowListBasedAuth, nil, logger.TestLogger(t))

	authResult, err := a.AuthorizeRequest(t.Context(), jsonrpc.Request[json.RawMessage]{
		ID:     "1",
		Method: vaulttypes.MethodSecretsCreate,
	})
	require.NoError(t, err)
	require.Equal(t, "0xauthorized", authResult.AuthorizedOwner())
}
