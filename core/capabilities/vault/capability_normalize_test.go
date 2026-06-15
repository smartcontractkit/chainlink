package vault

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

func TestCanonicalizeVaultRequest_CreateSecrets(t *testing.T) {
	t.Parallel()
	lowerOwner := "0x0001020304050607080900010203040506070809"
	req := &vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id: &vaultcommon.SecretIdentifier{
					Key:       "mykey",
					Owner:     lowerOwner,
					Namespace: "",
				},
			},
		},
	}

	canonicalizeVaultRequest(req)

	expectedOwner := common.HexToAddress(lowerOwner).Hex()
	require.Equal(t, expectedOwner, req.EncryptedSecrets[0].Id.Owner)
	require.Equal(t, "main", req.EncryptedSecrets[0].Id.Namespace)
}

func TestCanonicalizeRequestIfEnabled_RespectsGate(t *testing.T) {
	t.Parallel()
	lowerOwner := "0x0001020304050607080900010203040506070809"
	req := &vaultcommon.ListSecretIdentifiersRequest{
		Owner:     lowerOwner,
		Namespace: "",
	}

	disabledCap := &Capability{ownerAddressCanonicalizationEnabled: limits.NewGateLimiter(false)}
	unchanged := disabledCap.canonicalizeRequestIfEnabled(t.Context(), req)
	require.Same(t, req, unchanged)
	require.Equal(t, lowerOwner, req.Owner)

	enabledCap := &Capability{ownerAddressCanonicalizationEnabled: limits.NewGateLimiter(true)}
	enabledCap.canonicalizeRequestIfEnabled(t.Context(), req)
	require.Equal(t, common.HexToAddress(lowerOwner).Hex(), req.Owner)
}
