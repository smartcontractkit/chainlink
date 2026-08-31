package vaultutils

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

func TestSignedPayloadRequestID(t *testing.T) {
	t.Parallel()

	payload, err := ToCanonicalJSON(&vaultcommon.CreateSecretsResponse{
		RequestId: "owner::req-1",
		Responses: []*vaultcommon.CreateSecretResponse{},
	})
	require.NoError(t, err)

	got, err := SignedPayloadRequestID(vaulttypes.MethodSecretsCreate, json.RawMessage(payload))
	require.NoError(t, err)
	require.Equal(t, "owner::req-1", got)

	got, err = SignedPayloadRequestID(vaulttypes.MethodSecretsCreate, json.RawMessage(`{"responses":[]}`))
	require.NoError(t, err)
	require.Empty(t, got)
}
