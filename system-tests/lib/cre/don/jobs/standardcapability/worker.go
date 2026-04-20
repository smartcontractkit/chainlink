package standardcapability

import (
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/pkg/errors"

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

// DefaultCapabilitiesDir is where capability binaries live in the container (set during image build).
const DefaultCapabilitiesDir = "/usr/local/bin"

func GetCommand(binaryName string) (string, error) {
	if binaryName == "" {
		return "", errors.New("binary_name is required for capability config; set it in capability_defaults.toml or nodesets.capability_configs for plugin-based capabilities")
	}
	return filepath.Join(DefaultCapabilitiesDir, binaryName), nil
}
