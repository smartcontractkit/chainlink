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

func TestAllowListBasedAuth_CreateSecrets(t *testing.T) {
	params, err := json.Marshal(vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id: &vaultcommon.SecretIdentifier{
					Key:       "a",
					Namespace: "b",
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
	})
	allowListedReq := jsonrpc.Request[json.RawMessage]{
		ID:     "123",
		Method: vaulttypes.MethodSecretsList,
		Params: (*json.RawMessage)(&params),
	}
	require.NoError(t, err)
	notAllowedParams, err := json.Marshal(vaultcommon.ListSecretIdentifiersRequest{
		Namespace: "not allowed",
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
	owner := common.Address{1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

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
	owner := common.Address{1, 2, 3}
	req := makeListSecretsRequest(t, "123", "b")

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
	req := makeListSecretsRequest(t, "123", "b")

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
	req := makeListSecretsRequest(t, "123", "b")

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

func makeListSecretsRequest(t *testing.T, id, namespace string) jsonrpc.Request[json.RawMessage] {
	t.Helper()

	params, err := json.Marshal(vaultcommon.ListSecretIdentifiersRequest{
		Namespace: namespace,
	})
	require.NoError(t, err)

	return jsonrpc.Request[json.RawMessage]{
		ID:     id,
		Method: vaulttypes.MethodSecretsList,
		Params: (*json.RawMessage)(&params),
	}
}
