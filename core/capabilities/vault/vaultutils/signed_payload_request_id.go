package vaultutils

import (
	"encoding/json"
	"errors"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

// SignedPayloadRequestID extracts the requestId field from a vault OCR signed response payload.
func SignedPayloadRequestID(method string, payload json.RawMessage) (string, error) {
	if len(payload) == 0 {
		return "", errors.New("signed payload is empty")
	}

	unmarshal := protojson.UnmarshalOptions{DiscardUnknown: true}
	switch method {
	case vaulttypes.MethodSecretsCreate:
		resp := &vaultcommon.CreateSecretsResponse{}
		if err := unmarshal.Unmarshal(payload, resp); err != nil {
			return "", fmt.Errorf("failed to unmarshal signed create secrets payload: %w", err)
		}
		return resp.RequestId, nil
	case vaulttypes.MethodSecretsUpdate:
		resp := &vaultcommon.UpdateSecretsResponse{}
		if err := unmarshal.Unmarshal(payload, resp); err != nil {
			return "", fmt.Errorf("failed to unmarshal signed update secrets payload: %w", err)
		}
		return resp.RequestId, nil
	case vaulttypes.MethodSecretsDelete:
		resp := &vaultcommon.DeleteSecretsResponse{}
		if err := unmarshal.Unmarshal(payload, resp); err != nil {
			return "", fmt.Errorf("failed to unmarshal signed delete secrets payload: %w", err)
		}
		return resp.RequestId, nil
	case vaulttypes.MethodSecretsList:
		resp := &vaultcommon.ListSecretIdentifiersResponse{}
		if err := unmarshal.Unmarshal(payload, resp); err != nil {
			return "", fmt.Errorf("failed to unmarshal signed list secret identifiers payload: %w", err)
		}
		return resp.RequestId, nil
	default:
		return "", fmt.Errorf("unsupported method for signed payload request id validation: %s", method)
	}
}
