package vault

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	syncerv2mocks "github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/v2/mocks"
)

// allowlistTestOwner is the payload/entry owner shared across allowlist auth tests.
// The owner-scoped lookup requires the request payload owner to match the owner on
// the allowlisted entry, so both must use this address.
var allowlistTestOwner = common.Address{1, 2, 3}

func TestAllowListBasedAuth_CreateSecrets(t *testing.T) {
	params, err := json.Marshal(vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id: &vaultcommon.SecretIdentifier{
					Key:       "a",
					Namespace: "b",
					Owner:     allowlistTestOwner.Hex(),
				},
				EncryptedValue: "encrypted-value",
			},
		},
	})
	allowListedReq := jsonrpc.Request[json.RawMessage]{
		ID:     "123",
		Method: vaulttypes.MethodSecretsCreate,
		Params: (*json.RawMessage)(&params),
	}
	require.NoError(t, err)
	notAllowedParams, err := json.Marshal(vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id: &vaultcommon.SecretIdentifier{
					Key:       "not allowed",
					Namespace: "b",
					Owner:     allowlistTestOwner.Hex(),
				},
				EncryptedValue: "encrypted-value",
			},
		},
	})
	require.NoError(t, err)
	notAllowListedReq := jsonrpc.Request[json.RawMessage]{
		ID:     "123",
		Method: vaulttypes.MethodSecretsCreate,
		Params: (*json.RawMessage)(&notAllowedParams),
	}

	require.NoError(t, err)
	testAuthForRequests(t, allowListedReq, notAllowListedReq)
}

func TestAllowListBasedAuth_UpdateSecrets(t *testing.T) {
	params, err := json.Marshal(vaultcommon.UpdateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id: &vaultcommon.SecretIdentifier{
					Key:       "a",
					Namespace: "b",
					Owner:     allowlistTestOwner.Hex(),
				},
				EncryptedValue: "encrypted-value",
			},
		},
	})
	allowListedReq := jsonrpc.Request[json.RawMessage]{
		ID:     "123",
		Method: vaulttypes.MethodSecretsUpdate,
		Params: (*json.RawMessage)(&params),
	}
	require.NoError(t, err)
	notAllowedParams, err := json.Marshal(vaultcommon.UpdateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id: &vaultcommon.SecretIdentifier{
					Key:       "not allowed",
					Namespace: "b",
					Owner:     allowlistTestOwner.Hex(),
				},
				EncryptedValue: "encrypted-value",
			},
		},
	})
	require.NoError(t, err)
	notAllowListedReq := jsonrpc.Request[json.RawMessage]{
		ID:     "123",
		Method: vaulttypes.MethodSecretsUpdate,
		Params: (*json.RawMessage)(&notAllowedParams),
	}
	require.NoError(t, err)
	testAuthForRequests(t, allowListedReq, notAllowListedReq)
}

func TestAllowListBasedAuth_DeleteSecrets(t *testing.T) {
	params, err := json.Marshal(vaultcommon.DeleteSecretsRequest{
		Ids: []*vaultcommon.SecretIdentifier{
			{
				Key:       "a",
				Namespace: "b",
				Owner:     allowlistTestOwner.Hex(),
			},
		},
	})
	allowListedReq := jsonrpc.Request[json.RawMessage]{
		ID:     "123",
		Method: vaulttypes.MethodSecretsDelete,
		Params: (*json.RawMessage)(&params),
	}
	require.NoError(t, err)
	notAllowedParams, err := json.Marshal(vaultcommon.DeleteSecretsRequest{
		Ids: []*vaultcommon.SecretIdentifier{
			{
				Key:       "not allowed",
				Namespace: "b",
				Owner:     allowlistTestOwner.Hex(),
			},
		},
	})
	require.NoError(t, err)
	notAllowListedReq := jsonrpc.Request[json.RawMessage]{
		ID:     "123",
		Method: vaulttypes.MethodSecretsDelete,
		Params: (*json.RawMessage)(&notAllowedParams),
	}
	require.NoError(t, err)
	testAuthForRequests(t, allowListedReq, notAllowListedReq)
}

func TestAllowListBasedAuth_ListSecrets(t *testing.T) {
	params, err := json.Marshal(vaultcommon.ListSecretIdentifiersRequest{
		Namespace: "b",
		Owner:     allowlistTestOwner.Hex(),
	})
	allowListedReq := jsonrpc.Request[json.RawMessage]{
		ID:     "123",
		Method: vaulttypes.MethodSecretsList,
		Params: (*json.RawMessage)(&params),
	}
	require.NoError(t, err)
	notAllowedParams, err := json.Marshal(vaultcommon.ListSecretIdentifiersRequest{
		Namespace: "not allowed",
		Owner:     allowlistTestOwner.Hex(),
	})
	require.NoError(t, err)
	notAllowListedReq := jsonrpc.Request[json.RawMessage]{
		ID:     "123",
		Method: vaulttypes.MethodSecretsList,
		Params: (*json.RawMessage)(&notAllowedParams),
	}
	require.NoError(t, err)
	testAuthForRequests(t, allowListedReq, notAllowListedReq)
}

func testAuthForRequests(t *testing.T, allowlistedRequest, notAllowlistedRequest jsonrpc.Request[json.RawMessage]) {
	lggr := logger.TestLogger(t)
	owner := allowlistTestOwner

	mockSyncer := syncerv2mocks.NewWorkflowRegistrySyncer(t)
	auth := NewAllowListBasedAuth(lggr, mockSyncer)
	auth.retryCount = 0
	auth.retryInterval = time.Millisecond

	// Happy path
	digest, err := allowlistedRequest.Digest()
	require.NoError(t, err)
	digestBytes, err := hex.DecodeString(digest)
	require.NoError(t, err)
	expiry := time.Now().UTC().Unix() + 100
	allowlisted := []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{
		{
			RequestDigest:   [32]byte(digestBytes),
			Owner:           owner,
			ExpiryTimestamp: uint32(expiry), //nolint:gosec // it is a safe conversion
		},
	}
	mockSyncer.On("GetAllowlistedRequests", mock.Anything).Return(allowlisted)
	authResult, err := auth.AuthorizeRequest(t.Context(), allowlistedRequest)
	require.NoError(t, err)
	require.Equal(t, owner.Hex(), authResult.AuthorizedOwner())
	require.Equal(t, expiry, authResult.ExpiresAt())
	require.NotEmpty(t, authResult.Digest())

	// Same request is still authorized here; replay protection lives in the generic Authorizer.
	authResult, err = auth.AuthorizeRequest(t.Context(), allowlistedRequest)
	require.NoError(t, err)
	require.Equal(t, owner.Hex(), authResult.AuthorizedOwner())

	// Expired request
	allowlistedReqCopy := allowlistedRequest
	allowlistedReqCopy.ID = "456"
	allowlistedReqCopyDigest, err := allowlistedReqCopy.Digest()
	require.NoError(t, err)
	allowlistedReqCopyDigestBytes, err := hex.DecodeString(allowlistedReqCopyDigest)
	require.NoError(t, err)
	allowlisted[0].RequestDigest = [32]byte(allowlistedReqCopyDigestBytes)
	allowlisted[0].ExpiryTimestamp = uint32(time.Now().UTC().Unix() - 1) //nolint:gosec // it is a safe conversion
	mockSyncer.On("GetAllowlistedRequests", mock.Anything).Return(allowlisted)
	authResult, err = auth.AuthorizeRequest(t.Context(), allowlistedReqCopy)
	require.Nil(t, authResult)
	require.ErrorContains(t, err, "authorization expired")

	authResult, err = auth.AuthorizeRequest(t.Context(), notAllowlistedRequest)
	require.Nil(t, authResult)
	require.ErrorContains(t, err, "not allowlisted")
}

func TestAllowListBasedAuth_RetriesUntilRequestIsAllowlisted(t *testing.T) {
	lggr := logger.TestLogger(t)
	owner := allowlistTestOwner
	req := makeListSecretsRequest(t, "123", "b", owner.Hex())

	digest, err := req.Digest()
	require.NoError(t, err)
	digestBytes, err := hex.DecodeString(digest)
	require.NoError(t, err)
	expiry := time.Now().UTC().Unix() + 100
	allowlisted := []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{
		{
			RequestDigest:   [32]byte(digestBytes),
			Owner:           owner,
			ExpiryTimestamp: uint32(expiry), //nolint:gosec // it is a safe conversion
		},
	}

	mockSyncer := syncerv2mocks.NewWorkflowRegistrySyncer(t)
	auth := NewAllowListBasedAuth(lggr, mockSyncer)
	auth.retryCount = 2
	auth.retryInterval = time.Millisecond

	mockSyncer.On("GetAllowlistedRequests", mock.Anything).Return([]workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{}).Once()
	mockSyncer.On("GetAllowlistedRequests", mock.Anything).Return([]workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{}).Once()
	mockSyncer.On("GetAllowlistedRequests", mock.Anything).Return(allowlisted).Once()

	authResult, err := auth.AuthorizeRequest(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, owner.Hex(), authResult.AuthorizedOwner())
	require.Equal(t, expiry, authResult.ExpiresAt())
}

func TestAllowListBasedAuth_FailsAfterAllowlistReadRetries(t *testing.T) {
	lggr := logger.TestLogger(t)
	req := makeListSecretsRequest(t, "123", "b", allowlistTestOwner.Hex())

	mockSyncer := syncerv2mocks.NewWorkflowRegistrySyncer(t)
	mockSyncer.On("GetAllowlistedRequests", mock.Anything).Return([]workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{}).Times(3)

	auth := NewAllowListBasedAuth(lggr, mockSyncer)
	auth.retryCount = 2
	auth.retryInterval = time.Millisecond

	authResult, err := auth.AuthorizeRequest(t.Context(), req)
	require.Nil(t, authResult)
	require.ErrorContains(t, err, "not allowlisted")
}

func TestAllowListBasedAuth_StopsRetriesWhenContextCanceled(t *testing.T) {
	lggr := logger.TestLogger(t)
	req := makeListSecretsRequest(t, "123", "b", allowlistTestOwner.Hex())

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	mockSyncer := syncerv2mocks.NewWorkflowRegistrySyncer(t)
	mockSyncer.On("GetAllowlistedRequests", mock.Anything).Return([]workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{}).Once()

	auth := NewAllowListBasedAuth(lggr, mockSyncer)
	auth.retryCount = 2
	auth.retryInterval = time.Second

	authResult, err := auth.AuthorizeRequest(ctx, req)
	require.Nil(t, authResult)
	require.ErrorIs(t, err, context.Canceled)
}

func TestAllowListBasedAuth_OwnerScopedLookup(t *testing.T) {
	t.Parallel()

	victim := allowlistTestOwner        // the owner declared in the request payload
	attacker := common.Address{9, 9, 9} // a different owner that registered the same digest

	req := makeListSecretsRequest(t, "123", "b", victim.Hex())
	digest, err := req.Digest()
	require.NoError(t, err)
	digestBytes, err := hex.DecodeString(digest)
	require.NoError(t, err)
	expiry := time.Now().UTC().Unix() + 100

	newAuth := func(t *testing.T, entries []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest) *allowListBasedAuth {
		mockSyncer := syncerv2mocks.NewWorkflowRegistrySyncer(t)
		mockSyncer.On("GetAllowlistedRequests", mock.Anything).Return(entries)
		auth := NewAllowListBasedAuth(logger.TestLogger(t), mockSyncer)
		auth.retryCount = 0
		auth.retryInterval = time.Millisecond
		return auth
	}

	t.Run("ignores an entry registered under a different owner", func(t *testing.T) {
		t.Parallel()
		// Only the attacker's entry exists for this digest; the victim's request
		// must not be authorized under the attacker's owner.
		auth := newAuth(t, []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{
			{RequestDigest: [32]byte(digestBytes), Owner: attacker, ExpiryTimestamp: uint32(expiry)}, //nolint:gosec // safe conversion
		})
		authResult, err := auth.AuthorizeRequest(t.Context(), req)
		require.Nil(t, authResult)
		require.ErrorContains(t, err, "not allowlisted")
	})

	t.Run("matches the victim's own entry even when a poison entry shares the digest and is listed first", func(t *testing.T) {
		t.Parallel()
		// The attacker's entry is intentionally first in the list; owner-scoped
		// matching must still resolve to the victim's entry.
		auth := newAuth(t, []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{
			{RequestDigest: [32]byte(digestBytes), Owner: attacker, ExpiryTimestamp: uint32(expiry)}, //nolint:gosec // safe conversion
			{RequestDigest: [32]byte(digestBytes), Owner: victim, ExpiryTimestamp: uint32(expiry)},   //nolint:gosec // safe conversion
		})
		authResult, err := auth.AuthorizeRequest(t.Context(), req)
		require.NoError(t, err)
		require.Equal(t, victim.Hex(), authResult.AuthorizedOwner())
	})
}

func makeListSecretsRequest(t *testing.T, id, namespace, owner string) jsonrpc.Request[json.RawMessage] {
	t.Helper()

	params, err := json.Marshal(vaultcommon.ListSecretIdentifiersRequest{
		Namespace: namespace,
		Owner:     owner,
	})
	require.NoError(t, err)

	return jsonrpc.Request[json.RawMessage]{
		ID:     id,
		Method: vaulttypes.MethodSecretsList,
		Params: (*json.RawMessage)(&params),
	}
}
