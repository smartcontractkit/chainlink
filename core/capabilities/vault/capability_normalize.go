package vault

import (
	"context"

	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
)

func (s *Capability) canonicalizeRequestIfEnabled(ctx context.Context, request proto.Message) proto.Message {
	if s.ownerAddressCanonicalizationEnabled == nil || !ownerAddressCanonicalizationEnabled(ctx, s.ownerAddressCanonicalizationEnabled) {
		return request
	}
	return canonicalizeVaultRequest(request)
}

func ownerAddressCanonicalizationEnabled(ctx context.Context, gate limits.GateLimiter) bool {
	if gate == nil {
		return false
	}
	enabled, err := gate.Limit(ctx)
	return err == nil && enabled
}

func canonicalizeVaultRequest(request proto.Message) proto.Message {
	switch r := request.(type) {
	case *vaultcommon.CreateSecretsRequest:
		for _, encryptedSecret := range r.EncryptedSecrets {
			if encryptedSecret != nil && encryptedSecret.Id != nil {
				encryptedSecret.Id = vaultutils.CanonicalSecretIdentifier(encryptedSecret.Id)
			}
		}
	case *vaultcommon.UpdateSecretsRequest:
		for _, encryptedSecret := range r.EncryptedSecrets {
			if encryptedSecret != nil && encryptedSecret.Id != nil {
				encryptedSecret.Id = vaultutils.CanonicalSecretIdentifier(encryptedSecret.Id)
			}
		}
	case *vaultcommon.DeleteSecretsRequest:
		for i, id := range r.Ids {
			r.Ids[i] = vaultutils.CanonicalSecretIdentifier(id)
		}
	case *vaultcommon.ListSecretIdentifiersRequest:
		r.Owner = vaultutils.NormalizeWorkflowOwnerAddress(r.Owner)
		r.Namespace = vaultutils.NormalizeNamespace(r.Namespace)
	case *vaultcommon.GetSecretsRequest:
		for _, secretRequest := range r.Requests {
			if secretRequest != nil && secretRequest.Id != nil {
				secretRequest.Id = vaultutils.CanonicalSecretIdentifier(secretRequest.Id)
			}
		}
	}
	return request
}
