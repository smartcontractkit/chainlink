package vault

import (
	"encoding/json"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

// NormalizeRequestWithIdentity returns a copy of req with the supplied identity
// fields materialized into params for Vault secret-management methods. Requests
// with malformed params are returned unchanged so downstream parsing can surface
// the existing method-specific error paths.
func NormalizeRequestWithIdentity(req jsonrpc.Request[json.RawMessage], orgID, workflowOwner string) (jsonrpc.Request[json.RawMessage], error) {
	if req.Params == nil {
		return req, nil
	}

	rewrite := func(payload any) error {
		b, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		raw := json.RawMessage(b)
		req.Params = &raw
		return nil
	}

	switch req.Method {
	case vaulttypes.MethodSecretsCreate:
		parsed := &vaultcommon.CreateSecretsRequest{}
		if err := json.Unmarshal(*req.Params, parsed); err != nil {
			return req, nil
		}
		parsed.OrgId = orgID
		parsed.WorkflowOwner = workflowOwner
		return req, rewrite(parsed)
	case vaulttypes.MethodSecretsUpdate:
		parsed := &vaultcommon.UpdateSecretsRequest{}
		if err := json.Unmarshal(*req.Params, parsed); err != nil {
			return req, nil
		}
		parsed.OrgId = orgID
		parsed.WorkflowOwner = workflowOwner
		return req, rewrite(parsed)
	case vaulttypes.MethodSecretsDelete:
		parsed := &vaultcommon.DeleteSecretsRequest{}
		if err := json.Unmarshal(*req.Params, parsed); err != nil {
			return req, nil
		}
		parsed.OrgId = orgID
		parsed.WorkflowOwner = workflowOwner
		return req, rewrite(parsed)
	case vaulttypes.MethodSecretsList:
		parsed := &vaultcommon.ListSecretIdentifiersRequest{}
		if err := json.Unmarshal(*req.Params, parsed); err != nil {
			return req, nil
		}
		parsed.OrgId = orgID
		parsed.WorkflowOwner = workflowOwner
		return req, rewrite(parsed)
	default:
		return req, nil
	}
}

// StripRequestIdentity removes any org/workflow identity fields from Vault
// secret-management request params. This restores the request body shape used by
// the original client when the gateway had previously injected trusted identity
// fields for internal forwarding.
func StripRequestIdentity(req jsonrpc.Request[json.RawMessage]) (jsonrpc.Request[json.RawMessage], error) {
	return NormalizeRequestWithIdentity(req, "", "")
}
