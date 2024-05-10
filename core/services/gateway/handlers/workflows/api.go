package workflows

import "github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"

const (
	Create = "create"
	Delete = "delete"
	List   = "list"
	Get    = "get"
	Update = "update"

	Commit = "commit"
	Abort  = "abort"
)

type ResponseBase struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type CreateRequest struct {
	TOML string `json:"toml"`
}

type CreateResponse struct {
	ResponseBase
}

type DeleteRequest struct {
	// ID string `json:"id"`
	//	ExternalJobID string `json:"external_job_id"`
	Name string `json:"name"`
}

type DeleteResponse struct {
	ResponseBase
}

// TODO deduplicate with functions/api.go
type CombinedResponse struct {
	ResponseBase
	NodeResponses []*api.Message `json:"node_responses"`
}
