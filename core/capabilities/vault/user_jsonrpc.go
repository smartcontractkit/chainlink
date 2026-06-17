package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"google.golang.org/protobuf/proto"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
)

// InvalidVaultParamsError marks structural validation failures that must surface as InvalidParamsError.
type InvalidVaultParamsError struct {
	Method string
	Err    error
}

func (e InvalidVaultParamsError) Error() string {
	if prefix := invalidVaultParamsErrorPrefix(e.Method); prefix != "" {
		return prefix + e.Err.Error()
	}
	return e.Err.Error()
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

// ValidateStructureBeforeAuth validates a user-facing vault JSON-RPC request before authorization.
// Params are left unchanged so JWT request digests match the bytes the client signed.
// When stripOwnerPrefix is true, owner-prefixed request IDs are removed from the envelope and
// params so node-side allowlist re-authorization can match the original digest.
func (r *RequestValidator) ValidateStructureBeforeAuth(
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
	return r.validateMethodParams(ctx, req.Method, req.ID, req.Params, opts)
}

// StampAuthorizedParams stamps the authorized request ID and namespace defaults on params.
func (r *RequestValidator) StampAuthorizedParams(req *jsonrpc.Request[json.RawMessage], requestID string) error {
	req.ID = requestID
	normalizedParams, err := stampUserJSONRPCParams(req.Method, req.Params, requestID)
	if err != nil {
		return err
	}
	req.Params = normalizedParams
	return nil
}

func (r *RequestValidator) validateMethodParams(
	ctx context.Context,
	method string,
	requestID string,
	params *json.RawMessage,
	opts UserJSONRPCValidationOptions,
) error {
	if params == nil {
		return InvalidVaultParamsError{Method: method, Err: errors.New("request params must not be nil")}
	}
	if !vaulttypes.IsGatewaySecretsMethod(method) {
		return fmt.Errorf("unsupported gateway secrets method for validation: %s", method)
	}

	desc, err := lookupGatewaySecretsMethodDescriptor(method)
	if err != nil {
		return err
	}

	err = vaultutils.InspectJSONRPCParams(params, desc.newRequest, func(parsed proto.Message) error {
		return desc.validate(ctx, r, opts, requestID, parsed)
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
	prefixedOwner, originalRequestID, ok := strings.Cut(requestID, vaulttypes.RequestIDSeparator)
	if !ok {
		return requestID, ""
	}
	return originalRequestID, prefixedOwner
}

func stampUserJSONRPCParams(method string, params *json.RawMessage, requestID string) (*json.RawMessage, error) {
	if params == nil || !vaulttypes.IsGatewaySecretsMethod(method) {
		return params, nil
	}

	desc, err := lookupGatewaySecretsMethodDescriptor(method)
	if err != nil {
		return nil, err
	}

	return vaultutils.TransformJSONRPCParams(params, desc.newRequest, func(parsed proto.Message) error {
		setUserJSONRPCRequestID(parsed, requestID)
		desc.applyDefaults(parsed)
		return nil
	})
}

func unstampOwnerPrefixedRequestIDInParams(req *jsonrpc.Request[json.RawMessage], requestID string) error {
	if req.Params == nil || !vaulttypes.IsGatewaySecretsMethod(req.Method) {
		return nil
	}

	desc, err := lookupGatewaySecretsMethodDescriptor(req.Method)
	if err != nil {
		return err
	}

	normalized, err := vaultutils.TransformJSONRPCParams(req.Params, desc.newRequest, func(parsed proto.Message) error {
		setUserJSONRPCRequestID(parsed, requestID)
		return nil
	})
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
