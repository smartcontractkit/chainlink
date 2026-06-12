package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

// InvalidVaultParamsError marks structural validation failures that must surface as InvalidParamsError.
type InvalidVaultParamsError struct {
	Method string
	Err    error
}

func (e InvalidVaultParamsError) Error() string {
	switch e.Method {
	case vaulttypes.MethodSecretsCreate:
		return fmt.Sprintf("failed to validate create secrets request: %s", e.Err.Error())
	case vaulttypes.MethodSecretsUpdate:
		return fmt.Sprintf("failed to validate update secrets request: %s", e.Err.Error())
	case vaulttypes.MethodSecretsDelete:
		return fmt.Sprintf("failed to validate delete secrets request: %s", e.Err.Error())
	case vaulttypes.MethodSecretsList:
		return fmt.Sprintf("failed to validate list secret identifiers request: %s", e.Err.Error())
	default:
		return e.Err.Error()
	}
}

func (e InvalidVaultParamsError) Unwrap() error {
	return e.Err
}

// IsInvalidVaultParamsError reports whether err is a pre-auth validation failure.
func IsInvalidVaultParamsError(err error) bool {
	var invalidParams InvalidVaultParamsError
	return errors.As(err, &invalidParams)
}

// UserJSONRPCValidationOptions configures pre-authorization JSON-RPC validation.
type UserJSONRPCValidationOptions struct {
	PublicKey           *tdh2easy.PublicKey
	SkipLabelValidation bool
}

// PrepareUserJSONRPCRequest validates a user-facing vault JSON-RPC request before authorization.
// Params are left unchanged so JWT request digests match the bytes the client signed.
// When stripOwnerPrefix is true, owner-prefixed request IDs are removed from the envelope and
// params so node-side allowlist re-authorization can match the original digest.
func (r *RequestValidator) PrepareUserJSONRPCRequest(
	ctx context.Context,
	req *jsonrpc.Request[json.RawMessage],
	opts UserJSONRPCValidationOptions,
	stripOwnerPrefix bool,
) error {
	if stripOwnerPrefix {
		if err := stripOwnerPrefixForNodeAuth(req); err != nil {
			return InvalidVaultParamsError{Method: req.Method, Err: err}
		}
	}
	return r.validateUserJSONRPCParamsStructure(ctx, req.Method, req.ID, req.Params, opts)
}

// FinalizeAuthorizedJSONRPCRequest stamps the authorized request ID and namespace defaults on params.
func (r *RequestValidator) FinalizeAuthorizedJSONRPCRequest(req *jsonrpc.Request[json.RawMessage], requestID string) error {
	req.ID = requestID
	normalizedParams, err := normalizeUserJSONRPCParams(req.Method, req.Params, requestID)
	if err != nil {
		return err
	}
	req.Params = normalizedParams
	return nil
}

func (r *RequestValidator) validateUserJSONRPCParamsStructure(
	ctx context.Context,
	method string,
	requestID string,
	params *json.RawMessage,
	opts UserJSONRPCValidationOptions,
) error {
	if params == nil {
		return InvalidVaultParamsError{Method: method, Err: errors.New("request params must not be nil")}
	}
	if !vaulttypes.IsUserSecretsMethod(method) {
		return fmt.Errorf("unsupported vault user method for validation: %s", method)
	}

	err := inspectUserJSONRPCParams(method, params, func(parsed any) error {
		switch method {
		case vaulttypes.MethodSecretsCreate:
			request := parsed.(*vaultcommon.CreateSecretsRequest)
			request.RequestId = coalesceRequestID(request.RequestId, requestID)
			return r.ValidateCreateSecretsRequest(ctx, opts.PublicKey, request, opts.SkipLabelValidation)
		case vaulttypes.MethodSecretsUpdate:
			request := parsed.(*vaultcommon.UpdateSecretsRequest)
			request.RequestId = coalesceRequestID(request.RequestId, requestID)
			return r.ValidateUpdateSecretsRequest(ctx, opts.PublicKey, request, opts.SkipLabelValidation)
		case vaulttypes.MethodSecretsDelete:
			request := parsed.(*vaultcommon.DeleteSecretsRequest)
			request.RequestId = coalesceRequestID(request.RequestId, requestID)
			return r.ValidateDeleteSecretsRequest(ctx, request)
		case vaulttypes.MethodSecretsList:
			request := parsed.(*vaultcommon.ListSecretIdentifiersRequest)
			request.RequestId = coalesceRequestID(request.RequestId, requestID)
			return r.ValidateListSecretIdentifiersRequest(ctx, request)
		default:
			return fmt.Errorf("unsupported vault user method for validation: %s", method)
		}
	})
	if err != nil {
		return InvalidVaultParamsError{Method: method, Err: err}
	}
	return nil
}

func stripOwnerPrefixForNodeAuth(req *jsonrpc.Request[json.RawMessage]) error {
	originalRequestID, _ := stripPrefixedVaultRequestID(req.ID)
	req.ID = originalRequestID
	return unstampOwnerPrefixedRequestIDInParams(req, originalRequestID)
}

func stripPrefixedVaultRequestID(requestID string) (originalRequestID, prefixedOwner string) {
	idx := strings.Index(requestID, vaulttypes.RequestIDSeparator)
	if idx == -1 {
		return requestID, ""
	}
	return requestID[idx+len(vaulttypes.RequestIDSeparator):], requestID[:idx]
}

func normalizeUserJSONRPCParams(method string, params *json.RawMessage, requestID string) (*json.RawMessage, error) {
	if params == nil || !vaulttypes.IsUserSecretsMethod(method) {
		return params, nil
	}

	return transformUserJSONRPCParams(method, params, requestID, func(parsed any) error {
		switch method {
		case vaulttypes.MethodSecretsCreate:
			applyEncryptedSecretNamespaceDefaults(parsed.(*vaultcommon.CreateSecretsRequest).EncryptedSecrets)
		case vaulttypes.MethodSecretsUpdate:
			applyEncryptedSecretNamespaceDefaults(parsed.(*vaultcommon.UpdateSecretsRequest).EncryptedSecrets)
		case vaulttypes.MethodSecretsDelete:
			applySecretIdentifierNamespaceDefaults(parsed.(*vaultcommon.DeleteSecretsRequest).Ids)
		case vaulttypes.MethodSecretsList:
			request := parsed.(*vaultcommon.ListSecretIdentifiersRequest)
			if request.Namespace == "" {
				request.Namespace = vaulttypes.DefaultNamespace
			}
		}
		return nil
	})
}

func unstampOwnerPrefixedRequestIDInParams(req *jsonrpc.Request[json.RawMessage], requestID string) error {
	if req.Params == nil || !vaulttypes.IsUserSecretsMethod(req.Method) {
		return nil
	}

	normalized, err := transformUserJSONRPCParams(req.Method, req.Params, requestID, func(any) error { return nil })
	if err != nil {
		return err
	}
	req.Params = normalized
	return nil
}

func coalesceRequestID(paramsRequestID, envelopeRequestID string) string {
	if paramsRequestID != "" {
		return paramsRequestID
	}
	return envelopeRequestID
}

func inspectUserJSONRPCParams(
	method string,
	params *json.RawMessage,
	afterUnmarshal func(parsed any) error,
) error {
	switch method {
	case vaulttypes.MethodSecretsCreate:
		parsed := &vaultcommon.CreateSecretsRequest{}
		if err := json.Unmarshal(*params, parsed); err != nil {
			return err
		}
		return afterUnmarshal(parsed)
	case vaulttypes.MethodSecretsUpdate:
		parsed := &vaultcommon.UpdateSecretsRequest{}
		if err := json.Unmarshal(*params, parsed); err != nil {
			return err
		}
		return afterUnmarshal(parsed)
	case vaulttypes.MethodSecretsDelete:
		parsed := &vaultcommon.DeleteSecretsRequest{}
		if err := json.Unmarshal(*params, parsed); err != nil {
			return err
		}
		return afterUnmarshal(parsed)
	case vaulttypes.MethodSecretsList:
		parsed := &vaultcommon.ListSecretIdentifiersRequest{}
		if err := json.Unmarshal(*params, parsed); err != nil {
			return err
		}
		return afterUnmarshal(parsed)
	default:
		return fmt.Errorf("unsupported vault user method for validation: %s", method)
	}
}

// transformUserJSONRPCParams unmarshals user secrets params, stamps requestID, runs afterUnmarshal, and re-marshals.
func transformUserJSONRPCParams(
	method string,
	params *json.RawMessage,
	requestID string,
	afterUnmarshal func(parsed any) error,
) (*json.RawMessage, error) {
	switch method {
	case vaulttypes.MethodSecretsCreate:
		parsed := &vaultcommon.CreateSecretsRequest{}
		if err := json.Unmarshal(*params, parsed); err != nil {
			return nil, err
		}
		parsed.RequestId = requestID
		if err := afterUnmarshal(parsed); err != nil {
			return nil, err
		}
		return marshalUserJSONRPCParams(parsed)
	case vaulttypes.MethodSecretsUpdate:
		parsed := &vaultcommon.UpdateSecretsRequest{}
		if err := json.Unmarshal(*params, parsed); err != nil {
			return nil, err
		}
		parsed.RequestId = requestID
		if err := afterUnmarshal(parsed); err != nil {
			return nil, err
		}
		return marshalUserJSONRPCParams(parsed)
	case vaulttypes.MethodSecretsDelete:
		parsed := &vaultcommon.DeleteSecretsRequest{}
		if err := json.Unmarshal(*params, parsed); err != nil {
			return nil, err
		}
		parsed.RequestId = requestID
		if err := afterUnmarshal(parsed); err != nil {
			return nil, err
		}
		return marshalUserJSONRPCParams(parsed)
	case vaulttypes.MethodSecretsList:
		parsed := &vaultcommon.ListSecretIdentifiersRequest{}
		if err := json.Unmarshal(*params, parsed); err != nil {
			return nil, err
		}
		parsed.RequestId = requestID
		if err := afterUnmarshal(parsed); err != nil {
			return nil, err
		}
		return marshalUserJSONRPCParams(parsed)
	default:
		return nil, fmt.Errorf("unsupported vault user method for validation: %s", method)
	}
}

func marshalUserJSONRPCParams(payload any) (*json.RawMessage, error) {
	params, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	raw := json.RawMessage(params)
	return &raw, nil
}

func applyEncryptedSecretNamespaceDefaults(encryptedSecrets []*vaultcommon.EncryptedSecret) {
	for _, secretItem := range encryptedSecrets {
		if secretItem != nil && secretItem.Id != nil && secretItem.Id.Namespace == "" {
			secretItem.Id.Namespace = vaulttypes.DefaultNamespace
		}
	}
}

func applySecretIdentifierNamespaceDefaults(ids []*vaultcommon.SecretIdentifier) {
	for _, id := range ids {
		if id != nil && id.Namespace == "" {
			id.Namespace = vaulttypes.DefaultNamespace
		}
	}
}
