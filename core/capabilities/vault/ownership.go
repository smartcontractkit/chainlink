package vault

import (
	"encoding/json"
	"fmt"
	"strings"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

// StampAuthorizedOwnerOnCreate overwrites secret identifier owners in a create request.
func StampAuthorizedOwnerOnCreate(request *vaultcommon.CreateSecretsRequest, authorizedOwner string) {
	if request == nil || authorizedOwner == "" {
		return
	}
	for _, encryptedSecret := range request.EncryptedSecrets {
		if encryptedSecret != nil && encryptedSecret.Id != nil {
			encryptedSecret.Id.Owner = authorizedOwner
		}
	}
}

// StampAuthorizedOwnerOnUpdate overwrites secret identifier owners in an update request.
func StampAuthorizedOwnerOnUpdate(request *vaultcommon.UpdateSecretsRequest, authorizedOwner string) {
	if request == nil || authorizedOwner == "" {
		return
	}
	for _, encryptedSecret := range request.EncryptedSecrets {
		if encryptedSecret != nil && encryptedSecret.Id != nil {
			encryptedSecret.Id.Owner = authorizedOwner
		}
	}
}

// StampAuthorizedOwnerOnDelete overwrites secret identifier owners in a delete request.
func StampAuthorizedOwnerOnDelete(request *vaultcommon.DeleteSecretsRequest, authorizedOwner string) {
	if request == nil || authorizedOwner == "" {
		return
	}
	for _, id := range request.Ids {
		if id != nil {
			id.Owner = authorizedOwner
		}
	}
}

// StampAuthorizedOwnerOnList overwrites the list request owner field.
func StampAuthorizedOwnerOnList(request *vaultcommon.ListSecretIdentifiersRequest, authorizedOwner string) {
	if request == nil || authorizedOwner == "" {
		return
	}
	request.Owner = authorizedOwner
}

// StampAuthorizedOwnerFromRequestID stamps owner fields using the owner prefix in a
// gateway-prefixed request ID (owner::userRequestID). No-op when the ID has no prefix.
func StampAuthorizedOwnerFromRequestID(requestID string, request any) error {
	authorizedOwner := AuthorizedOwnerFromRequestID(requestID)
	if authorizedOwner == "" {
		return nil
	}
	switch r := request.(type) {
	case *vaultcommon.CreateSecretsRequest:
		StampAuthorizedOwnerOnCreate(r, authorizedOwner)
	case *vaultcommon.UpdateSecretsRequest:
		StampAuthorizedOwnerOnUpdate(r, authorizedOwner)
	case *vaultcommon.DeleteSecretsRequest:
		StampAuthorizedOwnerOnDelete(r, authorizedOwner)
	case *vaultcommon.ListSecretIdentifiersRequest:
		StampAuthorizedOwnerOnList(r, authorizedOwner)
	default:
		return fmt.Errorf("unsupported request type %T for owner stamping", request)
	}
	return nil
}

// AuthorizedOwnerFromRequestID returns the owner prefix from a gateway request ID of the
// form owner::userRequestID.
func AuthorizedOwnerFromRequestID(prefixedID string) string {
	idx := strings.Index(prefixedID, vaulttypes.RequestIDSeparator)
	if idx == -1 {
		return ""
	}
	return prefixedID[:idx]
}

func rewriteRequestParams(req *jsonrpc.Request[json.RawMessage], payload any) error {
	params, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	raw := json.RawMessage(params)
	req.Params = &raw
	return nil
}
