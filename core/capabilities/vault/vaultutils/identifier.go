package vaultutils

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

// NormalizeNamespace returns DefaultNamespace when namespace is empty.
func NormalizeNamespace(namespace string) string {
	if namespace == "" {
		return vaulttypes.DefaultNamespace
	}
	return namespace
}

// NormalizeWorkflowOwnerAddress canonicalizes a workflow owner to EIP-55 checksum form.
// This matches read paths (workflow secrets fetch, confidential relay) that look up
// secrets using common.HexToAddress(owner).Hex().
func NormalizeWorkflowOwnerAddress(owner string) string {
	return common.HexToAddress(owner).Hex()
}

// KeyFor returns the storage key for a secret identifier using the owner string as provided.
// When owner canonicalization is enabled, callers must normalize identifiers before persistence.
func KeyFor(id *vaultcommon.SecretIdentifier) string {
	namespace := NormalizeNamespace(id.Namespace)
	return fmt.Sprintf("%s::%s::%s", id.Owner, namespace, id.Key)
}

// CanonicalSecretIdentifier returns a copy of id with namespace defaulting and owner
// normalized to EIP-55 checksum form for consistent storage keys and metadata.
func CanonicalSecretIdentifier(id *vaultcommon.SecretIdentifier) *vaultcommon.SecretIdentifier {
	if id == nil {
		return nil
	}
	namespace := NormalizeNamespace(id.Namespace)
	return &vaultcommon.SecretIdentifier{
		Key:       id.Key,
		Owner:     NormalizeWorkflowOwnerAddress(id.Owner),
		Namespace: namespace,
	}
}
