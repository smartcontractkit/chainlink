package standardcapability

import (
	"fmt"

	"github.com/google/uuid"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
)

const (
	EmptyStdCapConfig = "\"\""
)

func WorkerJobSpec(nodeID, name, command, config, oracleFactoryConfig string) *jobv1.ProposeJobRequest {
	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec: fmt.Sprintf(`
	type = "standardcapabilities"
	schemaVersion = 1
	externalJobID = "%s"
	name = "%s"
	forwardingAllowed = false
	command = "%s"
	config = %s
	%s
`,
			uuid.NewString(),
			name,
			command,
			config,
			oracleFactoryConfig),
	}
}
