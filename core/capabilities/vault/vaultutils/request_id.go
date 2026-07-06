package vaultutils

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
)

const subscriptionPhaseKey = "subscription"

// BuildWorkflowGetSecretsRequestID returns the pending-queue / OCR request ID for a
// workflow GetSecrets call. This is the only place this format should be defined.
//
// We need to generate sufficiently unique IDs accounting for two cases:
// 1. called during the subscription phase, in which case the executionID will be blank
// 2. called during execution, in which case it'll be present.
// The reference ID is unique per phase, so we need to differentiate when generating
// an ID.
func BuildWorkflowGetSecretsRequestID(md capabilities.RequestMetadata) string {
	phaseOrExecution := md.WorkflowExecutionID
	if phaseOrExecution == "" {
		phaseOrExecution = subscriptionPhaseKey
	}
	return fmt.Sprintf("%s::%s::%s", md.WorkflowID, phaseOrExecution, md.ReferenceID)
}
