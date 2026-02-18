package presenters

import (
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/workflowkey"
)

type WorkflowKeyResource struct {
	JAID
	PublicKey string `json:"publicKey"`
}

// GetName implements the api2go EntityNamer interface
func (WorkflowKeyResource) GetName() string {
	return "workflowKeys"
}

func NewWorkflowKeyResource(key workflowkey.Key) *WorkflowKeyResource {
	return &WorkflowKeyResource{
		JAID:      NewJAID(key.PublicKeyString()),
		PublicKey: key.PublicKeyString(),
	}
}

func NewWorkflowKeyResources(keys []workflowkey.Key) []WorkflowKeyResource {
	rs := []WorkflowKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewWorkflowKeyResource(key))
	}

	return rs
}
