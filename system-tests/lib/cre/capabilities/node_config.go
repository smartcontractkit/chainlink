package capabilities

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func DefaultContainerDirectory(infraType infra.Type) (string, error) {
	switch infraType {
	case infra.Docker:
		// needs to match what CTFv2 uses by default, we should define a constant there and import it here
		return clnode.DefaultCapabilitiesDir, nil
	case infra.Kubernetes:
		// For Kubernetes, capabilities are already in the container image at /usr/local/bin
		return clnode.DefaultCapabilitiesDir, nil
	default:
		return "", fmt.Errorf("unknown infra type: %s", infraType)
	}
}
