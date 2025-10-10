package vaultutils

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gibson042/canonicaljson-go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

func DigestForRequest(req jsonrpc.Request[json.RawMessage]) ([32]byte, error) {
	var seed proto.Message
	switch req.Method {
	case vaulttypes.MethodSecretsCreate:
		var createSecretsRequests vaultcommon.CreateSecretsRequest
		if err := json.Unmarshal(*req.Params, &createSecretsRequests); err != nil {
			return [32]byte{}, errors.New("error unmarshalling create secrets request: " + err.Error())
		}
		seed = &vaultcommon.CreateSecretsRequest{
			EncryptedSecrets: createSecretsRequests.EncryptedSecrets,
		}
	case vaulttypes.MethodSecretsUpdate:
		var updateSecretsRequests vaultcommon.UpdateSecretsRequest
		if err := json.Unmarshal(*req.Params, &updateSecretsRequests); err != nil {
			return [32]byte{}, errors.New("error unmarshalling update secrets request: " + err.Error())
		}
		seed = &vaultcommon.CreateSecretsRequest{
			EncryptedSecrets: updateSecretsRequests.EncryptedSecrets,
		}
	case vaulttypes.MethodSecretsList:
		var listSecretsRequests vaultcommon.ListSecretIdentifiersRequest
		if err := json.Unmarshal(*req.Params, &listSecretsRequests); err != nil {
			return [32]byte{}, errors.New("error unmarshalling list secrets request: " + err.Error())
		}
		seed = &vaultcommon.ListSecretIdentifiersRequest{
			Owner:     listSecretsRequests.Owner,
			Namespace: listSecretsRequests.Namespace,
		}
	case vaulttypes.MethodSecretsDelete:
		var deleteSecretsRequests vaultcommon.DeleteSecretsRequest
		if err := json.Unmarshal(*req.Params, &deleteSecretsRequests); err != nil {
			return [32]byte{}, errors.New("error unmarshalling delete secrets request: " + err.Error())
		}
		seed = &vaultcommon.DeleteSecretsRequest{
			Ids: deleteSecretsRequests.Ids,
		}
	default:
		return [32]byte{}, fmt.Errorf("unauthorized method: %s", req.Method)
	}
	canonicalSeed, err := ToCanonicalJSON(seed)
	if err != nil {
		return [32]byte{}, errors.New("error converting request to canonical JSON: " + err.Error())
	}
	return sha256.Sum256(canonicalSeed), nil
}

// ToCanonicalJSON converts a protobuf message to a stable, deterministic
// representation, including consistent sorting of keys and fields, and
// consistent spacing.
func ToCanonicalJSON(msg proto.Message) ([]byte, error) {
	jsonb, err := protojson.MarshalOptions{
		UseProtoNames:   false,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
	}.Marshal(msg)
	if err != nil {
		return nil, err
	}

	jsond := map[string]any{}
	err = json.Unmarshal(jsonb, &jsond)
	if err != nil {
		return nil, err
	}

	return canonicaljson.Marshal(jsond)
}
