package infra

import (
	"fmt"

	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

// Unfortunately, we need to construct some of these URLs before any environment is created, because they are used
// in CL node configs. This introduces a coupling between Helm charts used by CRIB and Docker container names used by CTFv2.
func InternalHost(nodeIndex int, nodeType cretypes.CapabilityFlag, donName string, infraDetails types.InfraInput) string {
	if infraDetails.InfraType == types.CRIB {
		if nodeType == cretypes.BootstrapNode {
			return fmt.Sprintf("%s-bt-%d", donName, nodeIndex)
		}
		return fmt.Sprintf("%s-%d", donName, nodeIndex)
	}

	return fmt.Sprintf("%s-node%d", donName, nodeIndex)
}

func InternalGatewayHost(nodeIndex int, nodeType cretypes.CapabilityFlag, donName string, infraDetails types.InfraInput) string {
	if infraDetails.InfraType == types.CRIB {
		host := fmt.Sprintf("%s-%d", donName, nodeIndex)
		if nodeType == cretypes.BootstrapNode {
			host = fmt.Sprintf("%s-bt-%d", donName, nodeIndex)
		}
		host += "-gtwnode"

		return host
	}

	return fmt.Sprintf("%s-node%d", donName, nodeIndex)
}

func ExternalGatewayHost(nodeIndex int, nodeType cretypes.CapabilityFlag, donName string, infraDetails types.InfraInput) string {
	if infraDetails.InfraType == types.CRIB {
		return infraDetails.CRIB.Namespace + "-gateway.main.stage.cldev.sh"
	}

	return "localhost"
}

func ExternalGatewayPort(infraDetails types.InfraInput) int {
	if infraDetails.InfraType == types.CRIB {
		return 80
	}

	return 5002
}
