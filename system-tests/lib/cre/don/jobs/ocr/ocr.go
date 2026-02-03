package ocr

import (
	"fmt"

	"github.com/google/uuid"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
)

func BootstrapJobSpec(nodeID string, name string, ocr3CapabilityAddress string, chainID uint64) *jobv1.ProposeJobRequest {
	uuid := uuid.NewString()

	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec: fmt.Sprintf(`
	type = "bootstrap"
	schemaVersion = 1
	externalJobID = "%s"
	name = "%s"
	contractID = "%s"
	contractConfigTrackerPollInterval = "1s"
	contractConfigConfirmations = 1
	relay = "evm"
	[relayConfig]
	chainID = %d
	providerType = "ocr3-capability"
`,
			uuid,
			"ocr3-bootstrap-"+name+fmt.Sprintf("-%d", chainID),
			ocr3CapabilityAddress,
			chainID),
	}
}
