package vaultutils

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
)

const subscriptionPhaseKey = "subscription"

// BuildWorkflowGetSecretsRequestID returns the pending-queue / OCR request ID for a
// workflow GetSecrets call. This is the only place this format should be defined.
func BuildWorkflowGetSecretsRequestID(md capabilities.RequestMetadata) string {
	phaseOrExecution := md.WorkflowExecutionID
	if phaseOrExecution == "" {
		phaseOrExecution = subscriptionPhaseKey
	}
	return fmt.Sprintf("%s::%s::%s", md.WorkflowID, phaseOrExecution, md.ReferenceID)
}
