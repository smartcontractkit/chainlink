package vaultutils

import (
	"encoding/json"
	"errors"

	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

// InspectJSONRPCParams unmarshals params into a fresh proto from newRequest, runs fn, and does not remarshal.
func InspectJSONRPCParams(
	params *json.RawMessage,
	newRequest func() proto.Message,
	fn func(parsed proto.Message) error,
) error {
	if params == nil {
		return errors.New("request params must not be nil")
	}
	parsed := newRequest()
	if err := json.Unmarshal(*params, parsed); err != nil {
		return err
	}
	return fn(parsed)
}

// TransformJSONRPCParams unmarshals params, runs prepare, and remarshals.
func TransformJSONRPCParams(
	params *json.RawMessage,
	newRequest func() proto.Message,
	prepare func(parsed proto.Message) error,
) (*json.RawMessage, error) {
	if params == nil {
		return nil, errors.New("request params must not be nil")
	}
	parsed := newRequest()
	if err := json.Unmarshal(*params, parsed); err != nil {
		return nil, err
	}
	if err := prepare(parsed); err != nil {
		return nil, err
	}
	return MarshalUserJSONRPCParams(parsed)
}

// MarshalUserJSONRPCParams marshals a user secrets request payload to JSON-RPC params.
func MarshalUserJSONRPCParams(payload any) (*json.RawMessage, error) {
	params, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	raw := json.RawMessage(params)
	return &raw, nil
}

// ApplyEncryptedSecretNamespaceDefaults sets the default namespace on encrypted secret identifiers.
func ApplyEncryptedSecretNamespaceDefaults(encryptedSecrets []*vaultcommon.EncryptedSecret) {
	for _, secretItem := range encryptedSecrets {
		if secretItem != nil && secretItem.Id != nil && secretItem.Id.Namespace == "" {
			secretItem.Id.Namespace = vaulttypes.DefaultNamespace
		}
	}
}

// ApplySecretIdentifierNamespaceDefaults sets the default namespace on secret identifiers.
func ApplySecretIdentifierNamespaceDefaults(ids []*vaultcommon.SecretIdentifier) {
	for _, id := range ids {
		if id != nil && id.Namespace == "" {
			id.Namespace = vaulttypes.DefaultNamespace
		}
	}
}
