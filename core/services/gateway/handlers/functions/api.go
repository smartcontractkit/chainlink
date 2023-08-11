package functions

import "github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"

const (
	MethodSecretsSet  = "secrets_set"
	MethodSecretsList = "secrets_list"
)

type SecretsSetRequest struct {
	SlotID     uint   `json:"slot_id"`
	Version    uint64 `json:"version"`
	Expiration int64  `json:"expiration"`
	Payload    []byte `json:"payload"`
	Signature  []byte `json:"signature"`
}

// SecretsListRequest has empty payload

type SecretsResponseBase struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type SecretsSetResponse struct {
	SecretsResponseBase
}

type SecretsListResponse struct {
	SecretsResponseBase
	Rows []SecretsListRow `json:"rows,omitempty"`
}

type SecretsListRow struct {
	SlotID     uint   `json:"slot_id"`
	Version    uint64 `json:"version"`
	Expiration int64  `json:"expiration"`
}

// Gateway -> User response, which combines responses from several nodes
type CombinedSecretsResponse struct {
	SecretsResponseBase
	NodeResponses []*api.Message `json:"node_responses"`
}
