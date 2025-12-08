package standardcapability

import (
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
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

func GetCommand(binaryPathOnHost string, provider infra.Provider) (string, error) {
	containerPath, pathErr := crecapabilities.DefaultContainerDirectory(provider.Type)
	if pathErr != nil {
		return "", errors.Wrapf(pathErr, "failed to get default container directory for infra type %s", provider.Type)
	}

	return filepath.Join(containerPath, filepath.Base(binaryPathOnHost)), nil
}
