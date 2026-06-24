package vaultutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

func TestNormalizeNamespace(t *testing.T) {
	assert.Equal(t, vaulttypes.DefaultNamespace, NormalizeNamespace(""))
	assert.Equal(t, "custom", NormalizeNamespace("custom"))
}

func TestNormalizeWorkflowOwnerAddress(t *testing.T) {
	t.Parallel()
	checksummed := "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
	require.Equal(t, checksummed, NormalizeWorkflowOwnerAddress("0xab5801a7d398351b8be11c439e05c5b3259aec9b"))
}

func TestCanonicalSecretIdentifier(t *testing.T) {
	t.Parallel()
	checksummed := "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
	id := CanonicalSecretIdentifier(&vaultcommon.SecretIdentifier{
		Key: "k", Owner: "0xab5801a7d398351b8be11c439e05c5b3259aec9b", Namespace: "",
	})
	require.Equal(t, checksummed, id.Owner)
	require.Equal(t, vaulttypes.DefaultNamespace, id.Namespace)
}
