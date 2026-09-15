package vault

import (
	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
)

// buildRejectedOutcome constructs an Outcome that rejects every item in a request with the
// supplied error message. It is used by the include-invalid StateTransition path when an item
// receives F+1 error contributions, so the failed request still produces a per-item response
// rather than being silently dropped.
func buildRejectedOutcome(id string, payload proto.Message, requestType vaultcommon.RequestType, errMsg string) *vaultcommon.Outcome {
	o := &vaultcommon.Outcome{
		Id:          id,
		RequestType: requestType,
	}
	switch requestType {
	case vaultcommon.RequestType_GET_SECRETS:
		o.Response = &vaultcommon.Outcome_GetSecretsResponse{
			GetSecretsResponse: &vaultcommon.GetSecretsResponse{
				Responses: rejectedGetSecretsResponses(payload, errMsg),
			},
		}
	case vaultcommon.RequestType_CREATE_SECRETS:
		o.Response = &vaultcommon.Outcome_CreateSecretsResponse{
			CreateSecretsResponse: &vaultcommon.CreateSecretsResponse{
				Responses: rejectedCreateSecretsResponses(payload, errMsg),
			},
		}
	case vaultcommon.RequestType_UPDATE_SECRETS:
		o.Response = &vaultcommon.Outcome_UpdateSecretsResponse{
			UpdateSecretsResponse: &vaultcommon.UpdateSecretsResponse{
				Responses: rejectedUpdateSecretsResponses(payload, errMsg),
			},
		}
	case vaultcommon.RequestType_DELETE_SECRETS:
		o.Response = &vaultcommon.Outcome_DeleteSecretsResponse{
			DeleteSecretsResponse: &vaultcommon.DeleteSecretsResponse{
				Responses: rejectedDeleteSecretsResponses(payload, errMsg),
			},
		}
	case vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS:
		o.Response = &vaultcommon.Outcome_ListSecretIdentifiersResponse{
			ListSecretIdentifiersResponse: &vaultcommon.ListSecretIdentifiersResponse{
				Success: false,
				Error:   errMsg,
			},
		}
	case vaultcommon.RequestType_UNKNOWN:
	}
	return o
}

func rejectedGetSecretsResponses(payload proto.Message, errMsg string) []*vaultcommon.SecretResponse {
	req, ok := payload.(*vaultcommon.GetSecretsRequest)
	if !ok || len(req.Requests) == 0 {
		return []*vaultcommon.SecretResponse{{
			Result: &vaultcommon.SecretResponse_Error{Error: errMsg},
		}}
	}
	resps := make([]*vaultcommon.SecretResponse, len(req.Requests))
	for i, sr := range req.Requests {
		resps[i] = &vaultcommon.SecretResponse{
			Id: sr.Id,
			Result: &vaultcommon.SecretResponse_Error{
				Error: errMsg,
			},
		}
	}
	return resps
}

func rejectedCreateSecretsResponses(payload proto.Message, errMsg string) []*vaultcommon.CreateSecretResponse {
	req, ok := payload.(*vaultcommon.CreateSecretsRequest)
	if !ok || len(req.EncryptedSecrets) == 0 {
		return []*vaultcommon.CreateSecretResponse{{
			Success: false,
			Error:   errMsg,
		}}
	}
	resps := make([]*vaultcommon.CreateSecretResponse, len(req.EncryptedSecrets))
	for i, sr := range req.EncryptedSecrets {
		resps[i] = &vaultcommon.CreateSecretResponse{
			Id:      sr.Id,
			Success: false,
			Error:   errMsg,
		}
	}
	return resps
}

func rejectedUpdateSecretsResponses(payload proto.Message, errMsg string) []*vaultcommon.UpdateSecretResponse {
	req, ok := payload.(*vaultcommon.UpdateSecretsRequest)
	if !ok || len(req.EncryptedSecrets) == 0 {
		return []*vaultcommon.UpdateSecretResponse{{
			Success: false,
			Error:   errMsg,
		}}
	}
	resps := make([]*vaultcommon.UpdateSecretResponse, len(req.EncryptedSecrets))
	for i, sr := range req.EncryptedSecrets {
		resps[i] = &vaultcommon.UpdateSecretResponse{
			Id:      sr.Id,
			Success: false,
			Error:   errMsg,
		}
	}
	return resps
}

func rejectedDeleteSecretsResponses(payload proto.Message, errMsg string) []*vaultcommon.DeleteSecretResponse {
	req, ok := payload.(*vaultcommon.DeleteSecretsRequest)
	if !ok || len(req.Ids) == 0 {
		return []*vaultcommon.DeleteSecretResponse{{
			Success: false,
			Error:   errMsg,
		}}
	}
	resps := make([]*vaultcommon.DeleteSecretResponse, len(req.Ids))
	for i, id := range req.Ids {
		resps[i] = &vaultcommon.DeleteSecretResponse{
			Id:      id,
			Success: false,
			Error:   errMsg,
		}
	}
	return resps
}
