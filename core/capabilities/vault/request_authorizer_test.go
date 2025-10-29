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

func TestRequestAuthorizer_CreateSecrets(t *testing.T) {
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

func TestRequestAuthorizer_UpdateSecrets(t *testing.T) {
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

func TestRequestAuthorizer_DeleteSecrets(t *testing.T) {
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

func TestRequestAuthorizer_ListSecrets(t *testing.T) {
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
	auth := NewRequestAuthorizer(lggr, mockSyncer)

	// Happy path
	digest, err := allowlistedRequest.Digest()
	require.NoError(t, err)
	digestBytes, err := hex.DecodeString(digest)
	require.NoError(t, err)
	expiry := uint64(time.Now().UTC().Unix() + 100) //nolint:gosec // it is a safe conversion
	allowlisted := []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{
		{
			RequestDigest:   [32]byte(digestBytes),
			Owner:           owner,
			ExpiryTimestamp: uint32(expiry), //nolint:gosec // it is a safe conversion
		},
	}
	mockSyncer.On("GetAllowlistedRequests", mock.Anything).Return(allowlisted)
	isAuthorized, gotOwner, err := auth.AuthorizeRequest(context.Background(), allowlistedRequest)
	require.True(t, isAuthorized, err)
	require.Equal(t, owner.Hex(), gotOwner)
	require.NoError(t, err)

	// Already authorized
	isAuthorized, _, err = auth.AuthorizeRequest(context.Background(), allowlistedRequest)
	require.False(t, isAuthorized)
	require.ErrorContains(t, err, "already authorized previously")

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
	isAuthorized, _, err = auth.AuthorizeRequest(context.Background(), allowlistedReqCopy)
	require.False(t, isAuthorized)
	require.ErrorContains(t, err, "authorization expired")

	isAuthorized, _, err = auth.AuthorizeRequest(context.Background(), notAllowlistedRequest)
	require.False(t, isAuthorized)
	require.ErrorContains(t, err, "not allowlisted")
}
