package v2

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cosmos/gogoproto/proto"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

func DigestForRequest(req jsonrpc.Request[json.RawMessage]) ([32]byte, error) {
	var seed any
	switch req.Method {
	case vaulttypes.MethodSecretsCreate:
		var createSecretsRequests vaultcommon.CreateSecretsRequest
		if err := json.Unmarshal(*req.Params, &createSecretsRequests); err != nil {
			return [32]byte{}, errors.New("error unmarshalling create secrets request: " + err.Error())
		}
		seed = vaultcommon.CreateSecretsRequest{
			EncryptedSecrets: createSecretsRequests.EncryptedSecrets,
		}
	case vaulttypes.MethodSecretsUpdate:
		var updateSecretsRequests vaultcommon.UpdateSecretsRequest
		if err := json.Unmarshal(*req.Params, &updateSecretsRequests); err != nil {
			return [32]byte{}, errors.New("error unmarshalling update secrets request: " + err.Error())
		}
		seed = vaultcommon.UpdateSecretsRequest{
			EncryptedSecrets: updateSecretsRequests.EncryptedSecrets,
		}
	case vaulttypes.MethodSecretsList:
		var listSecretsRequests vaultcommon.ListSecretIdentifiersRequest
		if err := json.Unmarshal(*req.Params, &listSecretsRequests); err != nil {
			return [32]byte{}, errors.New("error unmarshalling list secrets request: " + err.Error())
		}
		seed = vaultcommon.ListSecretIdentifiersRequest{
			Owner:     listSecretsRequests.Owner,
			Namespace: listSecretsRequests.Namespace,
		}
	case vaulttypes.MethodSecretsDelete:
		var deleteSecretsRequests vaultcommon.DeleteSecretsRequest
		if err := json.Unmarshal(*req.Params, &deleteSecretsRequests); err != nil {
			return [32]byte{}, errors.New("error unmarshalling delete secrets request: " + err.Error())
		}
		seed = vaultcommon.DeleteSecretsRequest{
			Ids: deleteSecretsRequests.Ids,
		}
	default:
		return [32]byte{}, fmt.Errorf("unauthorized method: %s", req.Method)
	}

	return CalculateRequestDigest(seed), nil
}

// CalculateRequestDigest creates a SHA256 digest of the request for integrity verification
// This function is shared between client (JWT generation) and server (JWT validation)
func CalculateRequestDigest(req any) [32]byte {
	var data []byte
	if m, ok := req.(proto.Message); ok {
		// Use protobuf canonical serialization
		serialized, err := proto.Marshal(m)
		if err == nil {
			data = serialized
		} else {
			// fallback to string representation if marshal fails
			data = []byte(fmt.Sprintf("%v", req))
		}
	} else if s, ok := req.(fmt.Stringer); ok {
		data = []byte(s.String())
	} else {
		data = []byte(fmt.Sprintf("%v", req))
	}

	hash := sha256.Sum256(data)
	return hash
}
